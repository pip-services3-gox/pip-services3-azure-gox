package clients_test

import (
	"context"
	"os"
	"testing"

	cconf "github.com/pip-services3-gox/pip-services3-commons-gox/config"
	cref "github.com/pip-services3-gox/pip-services3-commons-gox/refer"
)

func TestDummyCommandableAzureFunctionClient(t *testing.T) {

	appName := os.Getenv("AZURE_FUNCTION_APP_NAME")
	functionName := os.Getenv("AZURE_FUNCTION_NAME")
	protocol := os.Getenv("AZURE_FUNCTION_PROTOCOL")
	authCode := os.Getenv("AZURE_FUNCTION_AUTH_CODE")
	uri := os.Getenv("AZURE_FUNCTION_URI")

	// if uri == "" {
	// 	uri = "http://localhost:7071/api/HttpTrigger1" // uncomment for local test
	// }

	if uri == "" && (appName == "" || functionName == "" || protocol == "" || authCode == "") {
		t.Skip("No credentials set, skip TestDummyCommandableAzureFunctionClient")
	}

	restConfig := cconf.NewConfigParamsFromTuples(
		"connection.uri", uri,
		"connection.protocol", protocol,
		"connection.app_name", appName,
		"connection.function_name", functionName,
		"credential.auth_code", authCode,
	)

	client := NewDummyCommandableAzureFunctionClient()
	fixture := NewDummyClientFixture(client)

	client.Configure(context.Background(), restConfig)
	client.SetReferences(context.Background(), cref.NewEmptyReferences())
	client.Open(context.Background(), "")
	defer client.Close(context.Background(), "DummyCommandableAzureFunctionClient")

	t.Run("DummyCommandableAzureFunctionClient.CrudOperations", fixture.TestCrudOperations)
}
