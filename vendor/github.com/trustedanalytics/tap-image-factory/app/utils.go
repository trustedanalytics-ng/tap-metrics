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
	"bytes"
	"errors"
	"fmt"
	"github.com/trustedanalytics/tap-catalog/models"
	"io"
	"os"
	"strings"
)

type ImageGetResponse struct {
	ImageId    string `json:"id"`
	Type       string `json:"type"`
	State      string `json:"state"`
	AuditTrail models.AuditTrail
}

type ImageStatePutRequest struct {
	State string `json:"state"`
}

type BuildImagePostRequest struct {
	ImageId string `json:"id"`
}

var (
	ImagesMap = map[models.ImageType]string{
		"JAVA":   "tap-base-java:java8-jessie",
		"GO":     "tap-base-binary:binary-jessie",
		"NODEJS": "tap-base-node:node4.4-jessie",
		"PYTHON": "tap-base-python:python2.7-jessie",
	}
)

func GetQueueConnectionString() string {
	host := os.Getenv("QUEUE_HOST")
	port := os.Getenv("QUEUE_PORT")
	user := os.Getenv("QUEUE_USER")
	pass := os.Getenv("QUEUE_PASS")
	return fmt.Sprintf("amqp://%v:%v@%v:%v/", user, pass, host, port)
}

func GetQueueName() string {
	return os.Getenv("QUEUE_NAME")
}

func GetCatalogAddressWithoutProtocol() string {
	return fmt.Sprintf("%v:%v", os.Getenv("CATALOG_HOST"), os.Getenv("CATALOG_PORT"))
}

func GetBlobStoreAddress() string {
	return fmt.Sprintf("%v:%v", os.Getenv("BLOB_STORE_HOST"), os.Getenv("BLOB_STORE_PORT"))
}

func GetHubAddressWithoutProtocol() string {
	address := os.Getenv("HUB_ADDRESS")
	split := strings.SplitN(address, "://", 2)
	if len(split) == 2 {
		return split[1]
	}
	return address
}

func GetDockerApiVersion() string {
	return os.Getenv("DOCKER_API_VERSION")
}

func GetDockerHostAddress() string {
	return os.Getenv("DOCKER_HOST")
}

func StreamToByte(stream io.Reader) ([]byte, error) {
	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(stream)
	if err != nil {
		return nil, errors.New("Could not read stream into byte array: " + err.Error())
	}
	return buf.Bytes(), nil
}

func StreamToString(stream io.Reader) (string, error) {
	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(stream)
	if err != nil {
		return "", errors.New("Could not read stream into string: " + err.Error())
	}
	return buf.String(), nil
}
