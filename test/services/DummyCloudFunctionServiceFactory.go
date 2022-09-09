package services_test

import (
	cref "github.com/pip-services3-gox/pip-services3-commons-gox/refer"
	cbuild "github.com/pip-services3-gox/pip-services3-components-gox/build"
)

type DummyAzureFunctionServiceFactory struct {
	cbuild.Factory
	Descriptor                *cref.Descriptor
	ControllerDescriptor      *cref.Descriptor
	AzureServiceDescriptor    *cref.Descriptor
	CmdAzureServiceDescriptor *cref.Descriptor
}

func NewDummyAzureFunctionServiceFactory() *DummyAzureFunctionServiceFactory {

	c := DummyAzureFunctionServiceFactory{
		Factory:                   *cbuild.NewFactory(),
		Descriptor:                cref.NewDescriptor("pip-services-dummies", "factory", "default", "default", "1.0"),
		AzureServiceDescriptor:    cref.NewDescriptor("pip-services-dummies", "service", "azure-function", "*", "1.0"),
		CmdAzureServiceDescriptor: cref.NewDescriptor("pip-services-dummies", "service", "commandable-azure-function", "*", "1.0"),
	}

	c.RegisterType(c.AzureServiceDescriptor, NewDummyAzureFunctionService)
	c.RegisterType(c.CmdAzureServiceDescriptor, NewDummyCommandableAzureFunctionService)
	return &c
}
