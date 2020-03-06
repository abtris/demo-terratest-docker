package test

import (
	"fmt"
	"strings"
	"testing"
	"time"

	http_helper "github.com/gruntwork-io/terratest/modules/http-helper"
	"github.com/gruntwork-io/terratest/modules/k8s"
	"github.com/gruntwork-io/terratest/modules/random"
)

func TestKubernetesHelloWorldExample(t *testing.T) {

	kubeResourcePath := "../hello-world-kubernetes/deployment.yaml"

	// unique namespace
	namespaceName := fmt.Sprintf("test-%s", strings.ToLower(random.UniqueId()))

	options := k8s.NewKubectlOptions("", "", namespaceName)

	k8s.CreateNamespace(t, options, namespaceName)
	defer k8s.DeleteNamespace(t, options, namespaceName)

	defer k8s.KubectlDelete(t, options, kubeResourcePath)
	k8s.KubectlApply(t, options, kubeResourcePath)
	k8s.WaitUntilServiceAvailable(t, options, "hello-world-service", 10, 3*time.Second)
	service := k8s.GetService(t, options, "hello-world-service")
	url := fmt.Sprintf("http://%s", k8s.GetServiceEndpoint(t, options, service, 5000))
	http_helper.HttpGetWithRetry(t, url, nil, 200, "Hello world!", 30, 5*time.Second)
}
