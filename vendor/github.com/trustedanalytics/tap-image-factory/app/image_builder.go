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
	"archive/tar"
	"bytes"
	"errors"
	"github.com/docker/engine-api/client"
	"github.com/docker/engine-api/types"
	"github.com/trustedanalytics/tap-catalog/models"
	"golang.org/x/net/context"
	"io"
)

const (
	artifactFileName = "artifact.tar"
)

type DockerClient struct {
	cli *client.Client
}

type ImageBuilder interface {
	CreateImage(artifact io.Reader, baseImage, imageId string) error
	buildImage(buildContext io.Reader, imageId string) error
	TagImage(imageId, tag string) error
	PushImage(tag string) error
}

func NewDockerClient() (*DockerClient, error) {
	cli, err := client.NewEnvClient()
	if err != nil {
		return nil, errors.New("Couldn't create docker client: " + err.Error())
	}
	dockerClient := DockerClient{cli}
	return &dockerClient, nil
}

func (d *DockerClient) CreateImage(artifact io.Reader, imageType models.ImageType, tag string) error {
	dockerfile, err := createDockerfile(imageType)
	if err != nil {
		return err
	}
	buildContext, err := createBuildContext(artifact, dockerfile)
	if err != nil {
		return err
	}
	return d.buildImage(buildContext, tag)
}

func (d *DockerClient) buildImage(buildContext io.Reader, tag string) error {
	buildOptions := types.ImageBuildOptions{Tags: []string{tag}}
	response, err := d.cli.ImageBuild(context.Background(), buildContext, buildOptions)
	if err != nil {
		return errors.New("Couldn't build docker image: " + err.Error())
	}
	responseStr, _ := StreamToString(response.Body)
	logger.Info(responseStr)
	if err != nil {
		return errors.New("Couldn't build docker image: " + err.Error())
	}
	return nil
}

func (d *DockerClient) TagImage(imageId, tag string) error {
	return d.cli.ImageTag(context.Background(), imageId, tag)
}

func (d *DockerClient) PushImage(tag string) error {
	//RegisterAuth is necessary even if no auth in registry - random value.
	//github issue: https://github.com/docker/docker/issues/10983
	_, err := d.cli.ImagePush(context.Background(), tag, types.ImagePushOptions{RegistryAuth: "cmFuZG9tX3ZhbHVl"})
	return err
}

func baseImageFromType(imageType models.ImageType) (string, error) {
	baseImage, exists := ImagesMap[imageType]
	if !exists {
		return "", errors.New("No such type - base image not detected.")
	}
	return GetHubAddressWithoutProtocol() + "/" + baseImage, nil
}

func createDockerfile(imageType models.ImageType) (io.Reader, error) {
	baseImage, err := baseImageFromType(imageType)
	if err != nil {
		return nil, err
	}
	buf := new(bytes.Buffer)
	buf.WriteString("FROM " + baseImage + "\n")
	buf.WriteString("ADD " + artifactFileName + " /root\n")
	buf.WriteString("ENV PORT 80\n")
	buf.WriteString("EXPOSE $PORT\n")
	buf.WriteString("WORKDIR /root\n")
	buf.WriteString("CMD [\"/root/run.sh\"]")
	return buf, nil
}

func createBuildContext(imageArtifact io.Reader, dockerfile io.Reader) (io.Reader, error) {
	imageArtifactBytes, err := StreamToByte(imageArtifact)
	if err != nil {
		return nil, err
	}
	dockerfileBytes, err := StreamToByte(dockerfile)
	if err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)
	tw := tar.NewWriter(buf)

	err = writeToTar(artifactFileName, imageArtifactBytes, tw)
	if err != nil {
		return nil, err
	}
	err = writeToTar("Dockerfile", dockerfileBytes, tw)
	if err != nil {
		return nil, err
	}
	tw.Close()

	return buf, nil
}

func writeToTar(fileName string, fileContents []byte, tw *tar.Writer) error {
	hdr := createHeader(fileName, int64(len(fileContents)))
	err := tw.WriteHeader(hdr)
	if err != nil {
		return errors.New("Couldn't write tar header for file: " + hdr.Name + ". Error: " + err.Error())
	}
	_, err = tw.Write(fileContents)
	if err != nil {
		return errors.New("Couldn't write tar body for file: " + hdr.Name + ". Error: " + err.Error())
	}
	return nil
}

func createHeader(fileName string, fileSize int64) *tar.Header {
	return &tar.Header{
		Name: fileName,
		Mode: 0700,
		Size: fileSize,
	}
}
