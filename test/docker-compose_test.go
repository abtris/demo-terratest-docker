package test

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/gruntwork-io/terratest/modules/docker"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/magiconair/properties/assert"
)

func TestDockerComposeLocal(t *testing.T) {

	tag := "go-webapp"
	buildOptions := &docker.BuildOptions{
		Tags: []string{tag},
	}

	docker.Build(t, "../hello-world-docker-compose", buildOptions)

	tagNginx := "local/nginx"
	buildOptionsNginx := &docker.BuildOptions{
		Tags: []string{tagNginx},
	}

	docker.Build(t, "../nginx", buildOptionsNginx)

	serverPort := 80
	randomSuffix := random.UniqueId()
	expectedServerText := fmt.Sprintf("Hello, %s!", randomSuffix)

	dockerOptions := &docker.Options{
		WorkingDir: "../hello-world-docker-compose",

		EnvVars: map[string]string{
			"SERVER_TEXT":  expectedServerText,
			"SERVER_PORT":  strconv.Itoa(serverPort),
			"randomSuffix": randomSuffix,
		},
	}

	defer docker.RunDockerCompose(t, dockerOptions, "down")

	docker.RunDockerCompose(t, dockerOptions, "up", "-d")

	// https://circleci.com/docs/2.0/building-docker-images/
	// docker run --network container:my-app appropriate/curl --retry 10 --retry-connrefused http://localhost:8080
	opts := &docker.RunOptions{
		Command:      []string{"--retry", "5", "--retry-connrefused", "-s", "http://production_nginx:80/hello"},
		OtherOptions: []string{"--network", "testdockercomposelocal_teststack-network"},
	}
	tag = "appropriate/curl"
	output := docker.Run(t, tag, opts)
	assert.Equal(t, expectedServerText, output)

}
