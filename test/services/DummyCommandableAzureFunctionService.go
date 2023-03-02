package services_test

import (
	"context"

	azureserv "github.com/pip-services3-gox/pip-services3-azure-gox/services"
	crefer "github.com/pip-services3-gox/pip-services3-commons-gox/refer"
)

type DummyCommandableAzureFunctionService struct {
	*azureserv.CommandableAzureFunctionService
}

func NewDummyCommandableAzureFunctionService() *DummyCommandableAzureFunctionService {
	c := DummyCommandableAzureFunctionService{}
	c.CommandableAzureFunctionService = azureserv.NewCommandableAzureFunctionService("dummies")
	c.DependencyResolver.Put(context.Background(), "controller", crefer.NewDescriptor("pip-services-dummies", "controller", "default", "*", "*"))
	return &c
}
