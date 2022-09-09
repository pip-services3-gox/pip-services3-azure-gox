package services

import (
	"context"
	"net/http"

	azureutil "github.com/pip-services3-gox/pip-services3-azure-gox/utils"
	ccomand "github.com/pip-services3-gox/pip-services3-commons-gox/commands"
	crun "github.com/pip-services3-gox/pip-services3-commons-gox/run"
	rpcserv "github.com/pip-services3-gox/pip-services3-rpc-gox/services"
)

// Abstract service that receives commands via Azure Function protocol
// to operations automatically generated for commands defined in ccomand.ICommandable components.
// Each command is exposed as invoke method that receives command name and parameters.
//
// Commandable services require only 3 lines of code to implement a robust external
// Azure Function-based remote interface.
//
// This service is intended to work inside Azure Function container that
// exploses registered actions externally.
//
// 	Configuration parameters:
//		- dependencies:
//			- controller:            override for Controller dependency
// 	References
//		- *:logger:*:*:1.0			(optional) ILogger components to pass log messages
//		- *:counters:*:*:1.0		(optional) ICounters components to pass collected measurements
//
// see AzureFunctionService
//
// 	Example:
//		type MyCommandableAzureFunctionService struct {
//			*azuresrv.CommandableAzureFunctionService
//		}
//
//		func NewMyCommandableAzureFunctionService() *MyCommandableAzureFunctionService {
//			c := MyCommandableAzureFunctionService{}
//			c.CommandableAzureFunctionService = azuresrv.NewCommandableAzureFunctionService("mydata")
//			c.DependencyResolver.Put(context.Background(), "controller", crefer.NewDescriptor("mygroup", "controller", "default", "*", "*"))
//			return &c
//		}
//
//		...
//
//		service := NewMyCommandableAzureFunctionService()
//		service.SetReferences(crefer.NewReferencesFromTuples(
//			crefer.NewDescriptor("mygroup","controller","default","default","1.0"), controller,
//		))
//		service.Open(ctx, "123")
//		fmt.Println("The Azure Function service is running")
//
type CommandableAzureFunctionService struct {
	*AzureFunctionService
	commandSet *ccomand.CommandSet
}

// Creates a new instance of the service.
// Parameters:
// 		- name 	a service name.
func NewCommandableAzureFunctionService(name string) *CommandableAzureFunctionService {
	c := CommandableAzureFunctionService{}
	c.AzureFunctionService = InheritAzureFunctionService(&c, name)

	return &c
}

// Returns body from Azure Function request.
// This method can be overloaded in child classes
// Parameters:
//		- req	Azure Function request
// Returns Parameters from request
func (c *CommandableAzureFunctionService) GetParameters(req *http.Request) *crun.Parameters {
	return azureutil.AzureFunctionRequestHelper.GetParameters(req)
}

// Registers all actions in Azure Function.
func (c *CommandableAzureFunctionService) Register() {
	resCtrl, depErr := c.DependencyResolver.GetOneRequired("controller")
	if depErr != nil {
		panic(depErr)
	}

	controller, ok := resCtrl.(ccomand.ICommandable)
	if !ok {
		c.Logger.Error(context.Background(), "CommandableHttpService", nil, "Can't cast Controller to ICommandable")
		return
	}

	c.commandSet = controller.GetCommandSet()
	commands := c.commandSet.Commands()

	for index := 0; index < len(commands); index++ {
		command := commands[index]
		name := command.Name()

		c.RegisterAction(name, nil, func(w http.ResponseWriter, r *http.Request) {
			correlationId := c.GetCorrelationId(r)
			args := c.GetParameters(r)
			args.Remove("correlation_id")

			timing := c.Instrument(r.Context(), correlationId, name)
			execRes, execErr := command.Execute(r.Context(), correlationId, args)
			timing.EndTiming(r.Context(), execErr)
			rpcserv.HttpResponseSender.SendResult(w, r, execRes, execErr)
		})
	}
}
