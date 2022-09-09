package main

import (
	"context"
	"log"
	"net/http"
	"os"

	aserv "github.com/pip-services3-gox/pip-services3-azure-gox/test/services"

	cconf "github.com/pip-services3-gox/pip-services3-commons-gox/config"
)

func main() {
	// create container
	config := cconf.NewConfigParamsFromTuples(
		"logger.descriptor", "pip-services:logger:console:default:1.0",
		"service.descriptor", "pip-services-dummies:service:azure-function:default:1.0",
		// "service.descriptor", "pip-services-dummies:service:commandable-azure-function:default:1.0",
	)

	ctx := context.Background()

	funcContainer := aserv.NewDummyAzureFunction()
	funcContainer.Configure(ctx, config)
	funcContainer.Open(ctx, "handler.main")

	handler := funcContainer.GetHandler()

	// run server
	listenAddr := ":8080"
	if val, ok := os.LookupEnv("FUNCTIONS_CUSTOMHANDLER_PORT"); ok {
		listenAddr = ":" + val
	}

	http.HandleFunc("/api/HttpTrigger1", handler)
	log.Printf("About to listen on %s. Go to https://127.0.0.1%s/", listenAddr, listenAddr)
	log.Fatal(http.ListenAndServe(listenAddr, nil))
}
