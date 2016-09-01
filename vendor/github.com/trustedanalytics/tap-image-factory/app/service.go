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
	"os"

	blobStoreApi "github.com/trustedanalytics/tap-blob-store/client"
	catalogApi "github.com/trustedanalytics/tap-catalog/client"
)

type BlobStoreApi interface {
	GetBlob(blobId string) ([]byte, error)
	GetImageBlob(imageId string) ([]byte, error)
	DeleteBlob(blobId string) error
	DeleteImageBlob(imageId string) error
}

func GetCatalogConnector() (*catalogApi.TapCatalogApiConnector, error) {
	address := GetCatalogAddressWithoutProtocol()
	if os.Getenv("CATALOG_SSL_CERT_FILE_LOCATION") != "" {
		return catalogApi.NewTapCatalogApiWithSSLAndBasicAuth(
			"https://"+address,
			os.Getenv("CATALOG_USER"),
			os.Getenv("CATALOG_PASS"),
			os.Getenv("CATALOG_SSL_CERT_FILE_LOCATION"),
			os.Getenv("CATALOG_SSL_KEY_FILE_LOCATION"),
			os.Getenv("CATALOG_SSL_CA_FILE_LOCATION"),
		)
	} else {
		return catalogApi.NewTapCatalogApiWithBasicAuth(
			"http://"+address,
			os.Getenv("CATALOG_USER"),
			os.Getenv("CATALOG_PASS"),
		)
	}
}

func GetBlobStoreConnector() (*blobStoreApi.TapBlobStoreApiConnector, error) {
	address := GetBlobStoreAddress()
	if os.Getenv("BLOB_STORE_SSL_CERT_FILE_LOCATION") != "" {
		return blobStoreApi.NewTapBlobStoreApiWithSSLAndBasicAuth(
			"https://"+address,
			os.Getenv("BLOB_STORE_USER"),
			os.Getenv("BLOB_STORE_PASS"),
			os.Getenv("BLOB_STORE_SSL_CERT_FILE_LOCATION"),
			os.Getenv("BLOB_STORE_SSL_KEY_FILE_LOCATION"),
			os.Getenv("BLOB_STORE_SSL_CA_FILE_LOCATION"),
		)
	} else {
		return blobStoreApi.NewTapBlobStoreApiWithBasicAuth(
			"http://"+address,
			os.Getenv("BLOB_STORE_USER"),
			os.Getenv("BLOB_STORE_PASS"),
		)
	}
}
