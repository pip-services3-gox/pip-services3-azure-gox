package services

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"regexp"

	"github.com/gorilla/mux"
	azureutil "github.com/pip-services3-gox/pip-services3-azure-gox/utils"
	cconf "github.com/pip-services3-gox/pip-services3-commons-gox/config"
	cconv "github.com/pip-services3-gox/pip-services3-commons-gox/convert"
	crefer "github.com/pip-services3-gox/pip-services3-commons-gox/refer"
	cvalid "github.com/pip-services3-gox/pip-services3-commons-gox/validate"
	ccount "github.com/pip-services3-gox/pip-services3-components-gox/count"
	clog "github.com/pip-services3-gox/pip-services3-components-gox/log"
	ctrace "github.com/pip-services3-gox/pip-services3-components-gox/trace"
	rpcserv "github.com/pip-services3-gox/pip-services3-rpc-gox/services"
)

type IAzureFunctionServiceOverrides interface {
	Register()
}

// Abstract service that receives remove calls via Azure Function protocol.
//
// This service is intended to work inside AzureFunction container that
// exposes registered actions externally.
//
// 	Configuration parameters
// 		- dependencies:
//			- controller:	override for Controller dependency
//
// 	References
//		- *:logger:*:*:1.0			(optional) ILogger components to pass log messages
//		- *:counters:*:*:1.0		(optional) ICounters components to pass collected measurements
//
// 	Example:
//		type MyAzureFunctionService struct {
//			*services.AzureFunctionService
//			controller IMyController
//		}
//
//		func NewMyAzureFunctionService() *MyAzureFunctionService {
//			c := MyAzureFunctionService{}
//
//			c.AzureFunctionService = services.InheritAzureFunctionService(&c, "v1.myservice")
//			c.DependencyResolver.Put(context.Background(), "controller", refer.NewDescriptor("mygroup", "controller", "default", "*", "1.0"))
//
//			return &c
//		}
//
//		func (c *MyAzureFunctionService) SetReferences(ctx context.Context, references refer.IReferences) {
//			c.AzureFunctionService.SetReferences(ctx, references)
//			depRes, depErr := c.DependencyResolver.GetOneRequired("controller")
//
//			if depErr == nil && depRes != nil {
//				c.controller = depRes.(IMyController)
//			}
//		}
//
//		func (c *MyAzureFunctionService) Register() {
//			c.RegisterAction(
//				"get_mydata",
//				nil,
//				func(w http.ResponseWriter, r *http.Request) {
//					var body map[string]any
//
//					err := AzureFunctionRequestHelper.DecodeBody(r, &body)
//					defer r.Body.Close()
//
//					result, err := c.controller.DeleteById(
//						r.Context(),
//						c.GetCorrelationId(r),
//						body,
//					)
//					HttpResponseSender.SendDeletedResult(w, r, result, err)
//				},
//			)
//		}
//
//		...
//
//		service := NewMyAzureFunctionService()
//		service.Configure(ctx, config.NewConfigParamsFromTuples(
//			"connection.protocol", "http",
//			"connection.host", "localhost",
//			"connection.port", 8080,
//		))
//
//		service.SetReferences(ctx, refer.NewReferencesFromTuples(
//			refer.NewDescriptor("mygroup", "controller", "default", "default", "1.0"), controller,
//		))
//		service.Open(ctx, "123")
//		fmt.Println("The Azure Function service is running")
//
type AzureFunctionService struct {
	name         string
	actions      []*AzureFunctionAction
	interceptors []func(http.ResponseWriter, *http.Request, http.HandlerFunc)
	opened       bool

	Overrides IAzureFunctionServiceOverrides
	// The dependency resolver.
	DependencyResolver *crefer.DependencyResolver
	// The logger.
	Logger *clog.CompositeLogger
	// The performance counters.
	Counters *ccount.CompositeCounters
	// The tracer.
	Tracer *ctrace.CompositeTracer
}

// Creates an instance of this service.
// Parameters:
//		- name	a service name to generate action cmd.
func NewAzureFunctionService(name string) *AzureFunctionService {
	c := AzureFunctionService{
		name:               name,
		actions:            make([]*AzureFunctionAction, 0),
		interceptors:       make([]func(http.ResponseWriter, *http.Request, http.HandlerFunc), 0),
		opened:             false,
		DependencyResolver: crefer.NewDependencyResolver(),
		Logger:             clog.NewCompositeLogger(),
		Counters:           ccount.NewCompositeCounters(),
		Tracer:             ctrace.NewCompositeTracer(),
	}

	c.Overrides = &c
	return &c
}

// InheritAzureFunctionService creates new instance of AzureFunctionService
func InheritAzureFunctionService(overrides IAzureFunctionServiceOverrides, name string) *AzureFunctionService {
	return &AzureFunctionService{
		name:               name,
		actions:            make([]*AzureFunctionAction, 0),
		interceptors:       make([]func(http.ResponseWriter, *http.Request, http.HandlerFunc), 0),
		opened:             false,
		Overrides:          overrides,
		DependencyResolver: crefer.NewDependencyResolver(),
		Logger:             clog.NewCompositeLogger(),
		Counters:           ccount.NewCompositeCounters(),
		Tracer:             ctrace.NewCompositeTracer(),
	}
}

// Registers all service routes in HTTP endpoint.
// This method is called by the service and must be overridden
// in child structs.
func (c *AzureFunctionService) Register() {}

// Configure the component with specified parameters.
//	see ConfigParams
//	Parameters:
//		- ctx context.Context
//		- config *conf.ConfigParams configuration parameters to set.
func (c *AzureFunctionService) Configure(ctx context.Context, config *cconf.ConfigParams) {
	c.DependencyResolver.Configure(ctx, config)
}

// SetReferences sets references to dependent components.
//	see IReferences
//	Parameters:
//		- ctx context.Context
//		- references IReferences references to locate the component dependencies.
func (c *AzureFunctionService) SetReferences(ctx context.Context, references crefer.IReferences) {
	c.Logger.SetReferences(ctx, references)
	c.Counters.SetReferences(ctx, references)
	c.Tracer.SetReferences(ctx, references)
	c.DependencyResolver.SetReferences(ctx, references)
}

// Instrument method are adds instrumentation to log calls and measure call time.
// It returns a Timing object that is used to end the time measurement.
//	Parameters:
//		- ctx context.Context
//		- correlationId     (optional) transaction id to trace execution through call chain.
//		- name              a method name.
//	Returns: Timing object to end the time measurement.
func (c *AzureFunctionService) Instrument(ctx context.Context, correlationId string, name string) *rpcserv.InstrumentTiming {
	c.Logger.Trace(ctx, correlationId, "Executing %s method", name)
	c.Counters.IncrementOne(ctx, name+".exec_count")

	counterTiming := c.Counters.BeginTiming(ctx, name+".exec_time")
	traceTiming := c.Tracer.BeginTrace(ctx, correlationId, name, "")

	return rpcserv.NewInstrumentTiming(correlationId, name, "exec",
		c.Logger, c.Counters, counterTiming, traceTiming)
}

// IsOpen Checks if the component is opened.
//	Returns: bool true if the component has been opened and false otherwise.
func (c *AzureFunctionService) IsOpen() bool {
	return c.opened
}

// Open method are opens the component.
//	Parameters:
//		- ctx context.Context
//		- correlationId  string (optional) transaction id to trace execution through call chain.
//	Returns: error or nil no errors occured.
func (c *AzureFunctionService) Open(ctx context.Context, correlationId string) error {
	if c.opened {
		return nil
	}

	c.Overrides.Register()

	return nil
}

// Close method are closes component and frees used resources.
//	Parameters:
//		- ctx context.Context
//		- correlationId (optional) transaction id to trace execution through call chain.
//	Returns: error or nil no errors occurred.
func (c *AzureFunctionService) Close(ctx context.Context, correlationId string) error {
	if c.opened {
		return nil
	}

	c.opened = false
	c.actions = nil
	c.interceptors = nil

	return nil
}

// Get all actions supported by the service.
// Returns an array with supported actions.
func (c *AzureFunctionService) GetActions() []*AzureFunctionAction {
	return c.actions
}

func (c *AzureFunctionService) ApplyValidation(schema *cvalid.Schema, action http.HandlerFunc) http.HandlerFunc {
	// Create an action function
	actionWrapper := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				err, ok := rec.(error)
				if !ok {
					msg := cconv.StringConverter.ToString(r)
					err = errors.New(msg)
				}
				c.Logger.Error(r.Context(), c.GetCorrelationId(r), err, "http handler panics with error")
			}
		}()

		// Validate object
		if schema != nil {
			var params = make(map[string]any, 0)
			for k, v := range r.URL.Query() {
				params[k] = v[0]
			}

			for k, v := range mux.Vars(r) {
				params[k] = v
			}

			// Make copy of request
			bodyBuf, bodyErr := ioutil.ReadAll(r.Body)
			if bodyErr != nil {
				rpcserv.HttpResponseSender.SendError(w, r, bodyErr)
				return
			}
			_ = r.Body.Close()
			r.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBuf))
			//-------------------------
			var body any
			_ = json.Unmarshal(bodyBuf, &body)
			params["body"] = body

			correlationId := c.GetCorrelationId(r)
			err := schema.ValidateAndReturnError(correlationId, params, false)
			if err != nil {
				rpcserv.HttpResponseSender.SendError(w, r, err)
				return
			}
		}

		action(w, r)
	}

	return actionWrapper
}

func (c *AzureFunctionService) ApplyInterceptors(action http.HandlerFunc) http.HandlerFunc {
	actionWrapper := action

	for index := len(c.interceptors) - 1; index >= 0; index-- {
		interceptor := c.interceptors[index]
		actionWrapper = func(action http.HandlerFunc) http.HandlerFunc {
			return func(w http.ResponseWriter, r *http.Request) {
				interceptor(w, r, action)
			}
		}(actionWrapper)
	}

	return actionWrapper
}

func (c *AzureFunctionService) GenerateActionCmd(name string) string {
	cmd := name
	if c.name != "" {
		cmd = c.name + "." + cmd
	}

	return cmd
}

// Registers a action in Google Function function.
// Parameters:
//		- name		an action name
//		- schema		a validation schema to validate received parameters.
//		- action		an action function that is called when operation is invoked.
func (c *AzureFunctionService) RegisterAction(name string, schema *cvalid.Schema, action http.HandlerFunc) {
	actionWrapper := c.ApplyValidation(schema, action)
	actionWrapper = c.ApplyInterceptors(actionWrapper)

	registeredAction := &AzureFunctionAction{
		Cmd:    c.GenerateActionCmd(name),
		Schema: schema,
		Action: actionWrapper,
	}

	c.actions = append(c.actions, registeredAction)
}

// Registers an action with authorization.
// Parameters:
//		- name		an action name
//		- schema	a validation schema to validate received parameters.
//		- authorize		an authorization interceptor
//		- action		an action function that is called when operation is invoked.
func (c *AzureFunctionService) RegisterActionWithAuth(name string, schema *cvalid.Schema, authorize func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc),
	action http.HandlerFunc) {
	actionWrapper := c.ApplyValidation(schema, action)

	if authorize != nil {
		nextAction := action
		action = func(w http.ResponseWriter, r *http.Request) {
			authorize(w, r, nextAction)
		}
	}

	actionWrapper = c.ApplyInterceptors(actionWrapper)

	registeredAction := &AzureFunctionAction{
		Cmd:    c.GenerateActionCmd(name),
		Schema: schema,
		Action: actionWrapper,
	}

	c.actions = append(c.actions, registeredAction)
}

// Registers a middleware for actions in Google Function service.
// Parameters:
//		- action	an action function that is called when middleware is invoked.
func (c *AzureFunctionService) RegisterInterceptor(cmd string, action func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc)) {
	interceptorWrapper := func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		currCmd, _ := c.GetCommand(r)
		matched, _ := regexp.MatchString(cmd, currCmd)
		if cmd != "" && !matched {
			next.ServeHTTP(w, r)
		} else {
			action(w, r, next.ServeHTTP)
		}
	}
	c.interceptors = append(c.interceptors, interceptorWrapper)
}

// Returns correlationId from Google Function request.
// This method can be overloaded in child structs
func (c *AzureFunctionService) GetCorrelationId(r *http.Request) string {
	return azureutil.AzureFunctionRequestHelper.GetCorrelationId(r)
}

// Returns command from Google Function request.
// This method can be overloaded in child structs.
// Parameters:
//		- req	the function request
// Returns command from request
func (c *AzureFunctionService) GetCommand(r *http.Request) (string, error) {
	return azureutil.AzureFunctionRequestHelper.GetCommand(r)
}
