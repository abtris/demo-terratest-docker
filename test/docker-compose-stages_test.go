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
	test_structure "github.com/gruntwork-io/terratest/modules/test-structure"
)

func TestDockerComposeWithStagesLocal(t *testing.T) {
	t.Parallel()

	workingDir := "../hello-world-docker-compose-stages"

	test_structure.RunTestStage(t, "build_docker_image", func() {
		buildImage(t, workingDir)
	})
	test_structure.RunTestStage(t, "run_docker_compose", func() {
		runCompose(t, workingDir)
	})
}

func buildImage(t *testing.T, workingDir string) {
	tag := "go-webapp"
	buildOptions := &docker.BuildOptions{
		Tags: []string{tag},
	}
	docker.Build(t, workingDir, buildOptions)
}

func runCompose(t *testing.T, workingDir string) {
	serverPort := 88
	randomSuffix := random.UniqueId()
	expectedServerText := fmt.Sprintf("Hello, %s!", randomSuffix)

	dockerOptions := &docker.Options{
		WorkingDir: workingDir,

		EnvVars: map[string]string{
			"SERVER_TEXT":  expectedServerText,
			"randomSuffix": randomSuffix,
			"SERVER_PORT":  strconv.Itoa(serverPort),
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
