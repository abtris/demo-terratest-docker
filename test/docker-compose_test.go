package test

import (
	"crypto/tls"
	"fmt"
	"testing"
	"time"

	"github.com/gruntwork-io/terratest/modules/docker"
	http_helper "github.com/gruntwork-io/terratest/modules/http-helper"
	"github.com/gruntwork-io/terratest/modules/random"
)

// An example of how to test the Packer template in examples/packer-docker-example completely locally using Terratest
// and Docker.
func TestPackerDockerExampleLocal(t *testing.T) {
	t.Parallel()

	tag := "go-webapp"
	buildOptions := &docker.BuildOptions{
		Tags: []string{tag},
	}

	docker.Build(t, "../hello-world-docker-compose", buildOptions)

	serverPort := 80
	expectedServerText := fmt.Sprintf("Hello, %s!", random.UniqueId())

	dockerOptions := &docker.Options{
		WorkingDir: "../hello-world-docker-compose",

		// Configure the port the web app will listen on and the text it will return using environment variables
		EnvVars: map[string]string{
			"SERVER_TEXT": expectedServerText,
		},
	}

	// website::tag::6::Make sure to shut down the Docker container at the end of the test.
	defer docker.RunDockerCompose(t, dockerOptions, "down")

	// website::tag::4::Run Docker Compose to fire up the web app. We run it in the background (-d) so it doesn't block this test.
	docker.RunDockerCompose(t, dockerOptions, "up", "-d")

	// It can take a few seconds for the Docker container boot up, so retry a few times
	maxRetries := 5
	timeBetweenRetries := 5 * time.Second
	url := fmt.Sprintf("http://localhost:%d/hello", serverPort)

	// Setup a TLS configuration to submit with the helper, a blank struct is acceptable
	tlsConfig := tls.Config{}

	// website::tag::5::Verify that we get back a 200 OK with the expected text
	http_helper.HttpGetWithRetry(t, url, &tlsConfig, 200, expectedServerText, maxRetries, timeBetweenRetries)
}
