package test

import (
	"crypto/tls"
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/gruntwork-io/terratest/modules/docker"
	http_helper "github.com/gruntwork-io/terratest/modules/http-helper"
	"github.com/gruntwork-io/terratest/modules/random"
)

func TestDockerComposeLocal(t *testing.T) {
	// t.Parallel()

	tag := "go-webapp"
	buildOptions := &docker.BuildOptions{
		Tags: []string{tag},
	}

	docker.Build(t, "../hello-world-docker-compose", buildOptions)

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

	maxRetries := 10
	timeBetweenRetries := 5 * time.Second
	url := fmt.Sprintf("http://localhost:%d/hello", serverPort)

	tlsConfig := tls.Config{}

	http_helper.HttpGetWithRetry(t, url, &tlsConfig, 200, expectedServerText, maxRetries, timeBetweenRetries)
}
