package containers_test

import (
	"context"

	azurecont "github.com/pip-services3-gox/pip-services3-azure-gox/containers"
	tbuild "github.com/pip-services3-gox/pip-services3-azure-gox/test/build"
	crefer "github.com/pip-services3-gox/pip-services3-commons-gox/refer"
)

type DummyCommandableAzureFunction struct {
	*azurecont.CommandableAzureFunction
}

func NewDummyCommandableAzureFunction() *DummyCommandableAzureFunction {
	c := DummyCommandableAzureFunction{}
	c.CommandableAzureFunction = azurecont.NewCommandableAzureFunctionWithParams("dummy", "Dummy commandable azure function")
	c.DependencyResolver.Put(context.Background(), "controller", crefer.NewDescriptor("pip-services-dummies", "controller", "default", "*", "*"))

	c.AddFactory(tbuild.NewDummyFactory())

	return &c
}
