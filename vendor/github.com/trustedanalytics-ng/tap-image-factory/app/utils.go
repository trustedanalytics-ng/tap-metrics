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
	"io"
	"os"
	"strings"
)

func GetImageWithHubAddressWithoutProtocol(image string) string {
	address := os.Getenv("HUB_ADDRESS")
	split := strings.SplitN(address, "://", 2)
	if len(split) == 2 {
		return split[1] + "/" + image
	}
	return address + "/" + image
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
