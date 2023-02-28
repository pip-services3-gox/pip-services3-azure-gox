package containers

import (
	"context"
	"net/http"

	azureutil "github.com/pip-services3-gox/pip-services3-azure-gox/utils"
	ccomand "github.com/pip-services3-gox/pip-services3-commons-gox/commands"
	crun "github.com/pip-services3-gox/pip-services3-commons-gox/run"
	rpcserv "github.com/pip-services3-gox/pip-services3-rpc-gox/services"
)

// Abstract Azure Function function, that acts as a container to instantiate and run components
// and expose them via external entry point. All actions are automatically generated for commands
// defined in ICommandable components. Each command is exposed as an action defined by "cmd" parameter.
//
// Container configuration for this Azure Function is stored in "./config/config.yml" file.
// But this path can be overridden by <code>CONFIG_PATH</code> environment variable.
//
//	References
//		- *:logger:*:*:1.0							(optional) ILogger components to pass log messages
//		- *:counters:*:*:1.0						(optional) ICounters components to pass collected measurements
//		- *:service:azurefunc:*:1.0       		(optional) IAzureFunctionService services to handle action requests
//		- *:service:commandable-azurefunc:*:1.0	(optional) IAzureFunctionService services to handle action requests
//
//	Example:
//		type MyAzureFunction struct {
//			*containers.CommandableAzureFunction
//			controller IMyController
//		}
//
//		func NewMyAzureFunction() *MyAzureFunction {
//			c := MyAzureFunction{}
//			c.AzureFunction = containers.NewCommandableAzureFunctionWithParams("mygroup", "MyGroup AzureFunction")
//
//			return &c
//		}
//
//		...
//
//		AzureFunction := NewMyAzureFunction()
//		AzureFunction.Run(ctx)
//		fmt.Println("MyAzureFunction is started")
//
// Deprecated: This component has been deprecated. Use AzureFunctionService instead.
type CommandableAzureFunction struct {
	*AzureFunction
}

// Creates a new instance of this Azure Function.
func NewCommandableAzureFunction() *CommandableAzureFunction {
	c := CommandableAzureFunction{}
	c.AzureFunction = InheritAzureFunction(&c)
	return &c
}

// Creates a new instance of this Azure Function.
// Parameters:
//		- name	(optional) a container name (accessible via ContextInfo)
//		- description	(optional) a container description (accessible via ContextInfo)
func NewCommandableAzureFunctionWithParams(name string, description string) *CommandableAzureFunction {
	c := CommandableAzureFunction{}
	c.AzureFunction = InheritAzureFunctionWithParams(&c, name, description)
	return &c
}

// Returns body from Azure Function request.
// This method can be overloaded in child classes
// Parameters:
//		- req	Googl Function request
// Returns Parameters from request
func (c *CommandableAzureFunction) GetParameters(req *http.Request) *crun.Parameters {
	return azureutil.AzureFunctionRequestHelper.GetParameters(req)
}

func (c *CommandableAzureFunction) registerCommandSet(commandSet *ccomand.CommandSet) {
	commands := commandSet.Commands()
	for index := 0; index < len(commands); index++ {
		command := commands[index]

		c.RegisterAction(command.Name(), nil, func(w http.ResponseWriter, r *http.Request) {
			correlationId := c.GetCorrelationId(r)
			args := c.GetParameters(r)

			timing := c.Instrument(r.Context(), correlationId, command.Name())
			execRes, execErr := command.Execute(r.Context(), correlationId, args)
			timing.EndTiming(r.Context(), execErr)

			rpcserv.HttpResponseSender.SendResult(w, r, execRes, execErr)
		})
	}
}

// Registers all actions in this Azure Function.
//
// Deprecated: Overloading of this method has been deprecated. Use AzureFunctionService instead.
func (c *CommandableAzureFunction) Register() {
	resCtrl, depErr := c.DependencyResolver.GetOneRequired("controller")
	if depErr != nil {
		panic(depErr)
	}

	controller, ok := resCtrl.(ccomand.ICommandable)
	if !ok {
		c.Logger().Error(context.Background(), "CommandableHttpService", nil, "Can't cast Controller to ICommandable")
		return
	}

	commandSet := controller.GetCommandSet()
	c.registerCommandSet(commandSet)
}
