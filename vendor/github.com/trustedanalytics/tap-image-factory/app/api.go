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
	"net/http"

	"github.com/gocraft/web"
	"github.com/op/go-logging"

	blobStoreApi "github.com/trustedanalytics/tap-blob-store/client"
	catalogApi "github.com/trustedanalytics/tap-catalog/client"
	commonLogger "github.com/trustedanalytics/tap-go-common/logger"
	"github.com/trustedanalytics/tap-go-common/util"
	"github.com/trustedanalytics/tap-image-factory/models"
)

var logger = initLogger()

func initLogger() *logging.Logger {
	logger, err := commonLogger.InitLogger("app")
	if err != nil {
		panic(err)
	}
	return logger
}

type Context struct {
	BlobStoreConnector     *blobStoreApi.TapBlobStoreApiConnector
	TapCatalogApiConnector catalogApi.TapCatalogApi
	DockerConnector        *DockerClient
}

var ctx *Context

func SetupContext() *Context {
	if ctx == nil {
		ctx = &Context{}
		tapBlobStoreConnector, err := GetBlobStoreConnector()
		if err != nil {
			logger.Panic(err)
		}
		ctx.BlobStoreConnector = tapBlobStoreConnector

		tapCatalogApiConnector, err := GetCatalogConnector()
		if err != nil {
			logger.Panic(err)
		}
		ctx.TapCatalogApiConnector = tapCatalogApiConnector

		dockerClient, err := NewDockerClient()
		if err != nil {
			logger.Panic(err)
		}
		ctx.DockerConnector = dockerClient
	}
	return ctx
}

func (c *Context) BuildImage(rw web.ResponseWriter, req *web.Request) {
	buildRequest := models.BuildImagePostRequest{}
	if err := util.ReadJson(req, &buildRequest); err != nil {
		util.Respond400(rw, err)
		return
	}

	go func() {
		if err := BuildAndPushImage(buildRequest.ImageId); err != nil {
			logger.Error("Building image error:", err)
		}
	}()
	util.WriteJson(rw, "", http.StatusAccepted)
}
