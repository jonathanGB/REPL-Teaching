package run

import (
	"archive/tar"
	"bytes"
	"fmt"
	"github.com/docker/docker/client"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"io/ioutil"
	"io"
	"context"
	"time"
)

func createTar(content []byte, extension string) (io.Reader, error) {
	// Create a buffer to write our archive to.
	buf := new(bytes.Buffer)
	tw := tar.NewWriter(buf)

	hdr := &tar.Header{
		Name: fmt.Sprintf("running.%s", extension),
		Mode: 0444,
		Size: int64(len(content)),
	}

	if err := tw.WriteHeader(hdr); err != nil {
		return nil, err
	}

	if _, err := tw.Write(content); err != nil {
		return nil, err
	}

	if err := tw.Close(); err != nil {
		return nil, err
	}

	return buf, nil
}

func extractTar(r io.Reader) (tarContent io.ReadWriter, err error) {
	// Open the tar archive for reading.
	tr := tar.NewReader(r)

	// Iterate through the files in the archive.
	_, err = tr.Next()
	if err != nil {
		return nil, err
	}

	tarContent = new(bytes.Buffer)
	if _, err = io.Copy(tarContent, tr); err != nil {
		return nil, err
	}

	return
}

func runQuery(script []byte, extension string) ([]byte, error) {
	scriptTarReader, err := createTar(script, extension)
	if err != nil {
		return nil, err
	}

	cli, err := client.NewEnvClient()
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15 * time.Second)
	defer cancel()

	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: fmt.Sprintf("run-%s", extension),
	}, nil, nil, "")
	if err != nil {
		return nil, fmt.Errorf("creating %v", err)
	}

	if err = cli.CopyToContainer(ctx, resp.ID, "/runs/", scriptTarReader, types.CopyToContainerOptions{}); err != nil {
		return nil, fmt.Errorf("copyto %v", err)
	}

	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		return nil, fmt.Errorf("starting %v", err)
	}

  if _, err = cli.ContainerWait(ctx, resp.ID); err != nil {
  	return nil, err
  }

	tarReader, _, err := cli.CopyFromContainer(ctx, resp.ID, "/runs/out")
	if err != nil {
		return nil, err
	}

	r, err := extractTar(tarReader)
	if err != nil {
		return nil, err
	}

	allContent, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	firstLineMarker := bytes.Index(allContent, []byte("\n"))
	content := allContent[firstLineMarker + 1:]

	if firstLineMarker > 0 && string(allContent[:firstLineMarker]) == "ok" {
		return content, nil
	}

	return content, fmt.Errorf("Erreur dans le script")
}
