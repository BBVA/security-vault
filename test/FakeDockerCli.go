package test

import (
	"io"
	"golang.org/x/net/context"
	"github.com/docker/engine-api/types"
)

type EventsTestMetrics struct {
	error error
	readCloser io.ReadCloser
	MethodCallMetrics
}

type CopyToContainerMetrics struct {
	error error
	MethodCallMetrics
}

type FakeDockerCli struct {
	events EventsTestMetrics
	copyToContainer CopyToContainerMetrics
}

func (f* FakeDockerCli) Events(ctx context.Context, options types.EventsOptions) (io.ReadCloser, error){
	f.events.Call()
	return f.events.readCloser,f.events.error
}

func (f* FakeDockerCli) CopyToContainer(ctx context.Context, container, path string, content io.Reader, options types.CopyToContainerOptions) error {
	f.copyToContainer.Call()
	return f.copyToContainer.error
}