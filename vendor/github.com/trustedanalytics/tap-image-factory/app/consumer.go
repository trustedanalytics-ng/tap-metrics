/**
 * Copyright (c) 2016 Intel Corporation
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *    http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package app

import (
	"sync"

	"github.com/streadway/amqp"

	"github.com/trustedanalytics/tap-go-common/queue"
	"github.com/trustedanalytics/tap-go-common/util"
	"github.com/trustedanalytics/tap-image-factory/models"
)

const (
	maxSimultaneousRoutines = 1
)

func StartConsumer(waitGroup *sync.WaitGroup) {
	waitGroup.Add(1)

	ch, conn := queue.GetConnectionChannel()
	queue.CreateExchangeWithQueueByRoutingKeys(ch, models.IMAGE_FACTORY_QUEUE_NAME, []string{models.IMAGE_FACTORY_IMAGE_ROUTING_KEY})
	queue.ConsumeMessages(ch, handleMessage, models.IMAGE_FACTORY_QUEUE_NAME, maxSimultaneousRoutines)

	defer conn.Close()
	defer ch.Close()

	logger.Info("Consuming stopped. Queue:", models.IMAGE_FACTORY_QUEUE_NAME)
	waitGroup.Done()
}

func handleMessage(msg amqp.Delivery) {
	buildImageRequest := models.BuildImagePostRequest{}
	err := util.ReadJsonFromByte(msg.Body, &buildImageRequest)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	if err := BuildAndPushImage(buildImageRequest.ImageId); err != nil {
		logger.Error("Building image error:", err)
	}
	err = msg.Ack(false)
	if err != nil {
		logger.Error("ACK error:", err)
	}
}
