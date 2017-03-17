package run

import (
	"github.com/docker/docker/client"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"context"
)

type DockerPool struct {

}

func getPool(num int) *DockerPool {
	// TODO: create pool of containers

}

func (dp *DockerPool) Query(script []byte, extension string) []byte {

}
