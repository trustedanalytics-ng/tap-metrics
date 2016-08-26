package app

import (
	"bytes"
	"fmt"
	"log"
	"net/http"

	"github.com/streadway/amqp"
	"github.com/trustedanalytics/tapng-go-common/util"
)

func failReceiverOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}

func StartConsumer(ctx Context) {
	logger.Info("Connecting to: " + GetQueueConnectionString())
	conn, err := amqp.Dial(GetQueueConnectionString())
	failReceiverOnError(err, "Failed to connect to Queue on address: "+GetQueueConnectionString())
	defer conn.Close()

	ch, err := conn.Channel()
	failReceiverOnError(err, "Failed to open a channel")
	defer ch.Close()

	/* DELETE THIS AFTER MONITOR IS IMPLEMENTED */
	err = ch.ExchangeDeclare(
		"tap.image-factory", // name
		"direct",            // type
		true,                // durable
		false,               // auto-deleted
		false,               // internal
		false,               // no-wait
		nil,                 // arguments
	)
	failReceiverOnError(err, "Failed to declare an exchange")

	q, err := ch.QueueDeclare(
		GetQueueName(), // name
		true,           // durable
		false,          // delete when usused
		true,           // exclusive
		false,          // no-wait
		nil,            // arguments
	)
	failReceiverOnError(err, "Failed to declare a queue")

	err = ch.QueueBind(
		q.Name,              // queue name
		"tap.image-factory", // routing key
		"tap.image-factory", // exchange
		false,               // no-wait
		nil,                 // arguments
	)
	failReceiverOnError(err, "Failed to bind a queue")
	/* DELETE THIS AFTER MONITOR IS IMPLEMENTED */

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer - empty means generate unique id
		true,   // auto-ack
		true,   // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failReceiverOnError(err, "Failed to register a consumer")

	go func() {
		for m := range msgs {
			handleMessage(ctx, m)
		}
	}()

	forever := make(chan bool)
	<-forever
}

func handleMessage(c Context, msg amqp.Delivery) {
	msg_json := BuildImagePostRequest{}
	err := util.ReadJsonFromByte(msg.Body, &msg_json)
	if err != nil {
		logger.Error(err.Error())
		return
	}
	c.updateImageWithState(msg_json.ImageId, "BUILDING")
	imgDetails, _, err := c.TapCatalogApiConnector.GetImage(msg_json.ImageId)
	if err != nil {
		c.updateImageWithState(msg_json.ImageId, "ERROR")
		logger.Error(err.Error())
		return
	}

	buffer := bytes.Buffer{}
	err = c.BlobStoreConnector.GetBlob(imgDetails.Id, &buffer)
	if err != nil {
		c.updateImageWithState(msg_json.ImageId, "ERROR")
		logger.Error(err.Error())
		return
	}
	tag := GetHubAddressWithoutProtocol() + "/" + imgDetails.Id

	err = c.DockerConnector.CreateImage(bytes.NewReader(buffer.Bytes()), imgDetails.Type, tag)
	if err != nil {
		c.updateImageWithState(msg_json.ImageId, "ERROR")
		logger.Error(err.Error())
		return
	}
	err = c.DockerConnector.PushImage(tag)
	if err != nil {
		c.updateImageWithState(msg_json.ImageId, "ERROR")
		logger.Error(err.Error())
		return
	}
	c.updateImageWithState(msg_json.ImageId, "READY")
	status, err := c.BlobStoreConnector.DeleteBlob(imgDetails.Id)
	if err != nil {
		logger.Error(err.Error())
		return
	}
	if status != http.StatusNoContent {
		logger.Warning("Blob removal failed. Actual status %v", status)
	}
}
