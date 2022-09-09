package services_test

import (
	azuresrv "github.com/pip-services3-gox/pip-services3-azure-gox/containers"
	tbuild "github.com/pip-services3-gox/pip-services3-azure-gox/test/build"
)

type DummyAzureFunction struct {
	*azuresrv.AzureFunction
}

func NewDummyAzureFunction() *DummyAzureFunction {
	c := DummyAzureFunction{AzureFunction: azuresrv.NewAzureFunctionWithParams("dummy", "Dummy azure function")}
	c.AddFactory(tbuild.NewDummyFactory())
	c.AddFactory(NewDummyAzureFunctionServiceFactory())

	return &c
}
