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
	"encoding/json"
	"github.com/gocraft/web"
	"net/http"

	blobStoreApi "github.com/trustedanalytics/tapng-blob-store/client"
	catalogApi "github.com/trustedanalytics/tapng-catalog/client"
	"github.com/trustedanalytics/tapng-catalog/models"
	"github.com/trustedanalytics/tapng-go-common/util"
	"github.com/trustedanalytics/tapng-image-factory/logger"
)

var (
	logger = logger_wrapper.InitLogger("main")
)

type Context struct {
	BlobStoreConnector     *blobStoreApi.TapBlobStoreApiConnector
	TapCatalogApiConnector *catalogApi.TapCatalogApiConnector
	DockerConnector        *DockerClient
}

func (c *Context) SetupContext() {
	tapBlobStoreConnector, err := GetBlobStoreConnector()
	if err != nil {
		logger.Panic(err)
	}
	c.BlobStoreConnector = tapBlobStoreConnector

	tapCatalogApiConnector, err := GetCatalogConnector()
	if err != nil {
		logger.Panic(err)
	}
	c.TapCatalogApiConnector = tapCatalogApiConnector

	dockerClient, err := NewDockerClient()
	if err != nil {
		logger.Panic(err)
	}
	c.DockerConnector = dockerClient
}
func (c *Context) BuildImage(rw web.ResponseWriter, req *web.Request) {
	req_json := BuildImagePostRequest{}
	err := util.ReadJson(req, &req_json)
	if err != nil {
		logger.Error(err.Error())
		rw.WriteHeader(400)
		return
	}
	imgDetails, _, err := c.TapCatalogApiConnector.GetImage(req_json.ImageId)
	if err != nil {
		logger.Error(err.Error())
		rw.WriteHeader(500)
		return
	}

	buffer := bytes.Buffer{}
	err = c.BlobStoreConnector.GetBlob(imgDetails.Id, &buffer)
	if err != nil {
		logger.Error(err.Error())
		rw.WriteHeader(500)
		return
	}
	marshalledValue, _ := json.Marshal("BUILDING")
	patches := []models.Patch{{Operation: models.OperationUpdate, Field: "State", Value: marshalledValue}}
	c.TapCatalogApiConnector.UpdateImage(imgDetails.Id, patches)
	if err != nil {
		logger.Error(err.Error())
		rw.WriteHeader(500)
		return
	}

	tag := GetHubAddressWithoutProtocol() + "/" + imgDetails.Id

	err = c.DockerConnector.CreateImage(bytes.NewReader(buffer.Bytes()), imgDetails.Type, tag)
	if err != nil {
		logger.Error(err.Error())
		rw.WriteHeader(500)
		return
	}
	err = c.DockerConnector.PushImage(tag)
	if err != nil {
		logger.Error(err.Error())
		rw.WriteHeader(500)
		return
	}
	marshalledValue, _ = json.Marshal("READY")
	patches = []models.Patch{{Operation: models.OperationUpdate, Field: "State", Value: marshalledValue}}
	c.TapCatalogApiConnector.UpdateImage(imgDetails.Id, patches)
	if err != nil {
		logger.Error(err.Error())
		rw.WriteHeader(500)
		return
	}
	status, err := c.BlobStoreConnector.DeleteBlob(imgDetails.Id)
	if err != nil {
		logger.Error(err.Error())
		rw.Write([]byte(err.Error()))
		rw.WriteHeader(500)
		return
	}
	if status != http.StatusNoContent {
		logger.Warning("Blob removal failed. Actual status %v", status)
	}

	rw.WriteHeader(201)
}

func (c *Context) updateImageWithState(imageId, state string) {
	marshalledValue, _ := json.Marshal(state)
	patches := []models.Patch{{Operation: models.OperationUpdate, Field: "State", Value: marshalledValue}}
	c.TapCatalogApiConnector.UpdateImage(imageId, patches)
}
