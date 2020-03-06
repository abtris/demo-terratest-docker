package test

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/gruntwork-io/terratest/modules/docker"
	"github.com/gruntwork-io/terratest/modules/random"
	test_structure "github.com/gruntwork-io/terratest/modules/test-structure"
	"github.com/magiconair/properties/assert"
)

func TestDockerComposeWithStagesLocal(t *testing.T) {

	workingDir := "../hello-world-docker-compose-stages"

	test_structure.RunTestStage(t, "build_docker_image", func() {
		buildImage(t, "go-webapp", "../demowebapp")
		buildImage(t, "local/nginx", "../nginx")
	})
	test_structure.RunTestStage(t, "run_docker_compose", func() {
		runCompose(t, workingDir)
	})
}

func buildImage(t *testing.T, tag string, workingDir string) {
	buildOptions := &docker.BuildOptions{
		Tags: []string{tag},
	}
	docker.Build(t, workingDir, buildOptions)
}

func runCompose(t *testing.T, workingDir string) {
	serverPort := 80
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

	// https://circleci.com/docs/2.0/building-docker-images/
	// docker run --network container:my-app appropriate/curl --retry 10 --retry-connrefused http://localhost:8080
	opts := &docker.RunOptions{
		Command:      []string{"--retry", "5", "--retry-connrefused", "-s", "http://production_nginx:80/hello"},
		OtherOptions: []string{"--network", "testdockercomposewithstageslocal_teststack-network"},
	}
	tag := "appropriate/curl"
	output := docker.Run(t, tag, opts)
	assert.Equal(t, expectedServerText, output)

}
