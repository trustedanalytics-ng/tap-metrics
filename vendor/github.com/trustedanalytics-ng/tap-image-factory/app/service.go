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
	blobStoreApi "github.com/trustedanalytics-ng/tap-blob-store/client"
	catalogApi "github.com/trustedanalytics-ng/tap-catalog/client"
	"github.com/trustedanalytics-ng/tap-go-common/util"
)

type BlobStoreApi interface {
	GetBlob(blobId string) ([]byte, error)
	GetImageBlob(imageId string) ([]byte, error)
	DeleteBlob(blobId string) error
	DeleteImageBlob(imageId string) error
}

func GetCatalogConnector() (catalogApi.TapCatalogApi, error) {
	address, username, password, err := util.GetConnectionParametersFromEnv("CATALOG")
	if err != nil {
		panic(err.Error())
	}
	return catalogApi.NewTapCatalogApiWithBasicAuth("https://"+address, username, password)
}

func GetBlobStoreConnector() (*blobStoreApi.TapBlobStoreApiConnector, error) {
	address, username, password, err := util.GetConnectionParametersFromEnv("BLOB_STORE")
	if err != nil {
		panic(err.Error())
	}
	return blobStoreApi.NewTapBlobStoreApiWithBasicAuth("https://"+address, username, password)
}
