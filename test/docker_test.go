package test

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/docker"
	"github.com/magiconair/properties/assert"
)

func TestDockerHelloWorldExample(t *testing.T) {
	tag := "teststack/demo-terratest-docker"
	buildOptions := &docker.BuildOptions{
		Tags: []string{tag},
	}

	docker.Build(t, "../hello-world", buildOptions)

	opts := &docker.RunOptions{Command: []string{"cat", "/test.txt"}}
	output := docker.Run(t, tag, opts)
	assert.Equal(t, "Hello, World!", output)
}
