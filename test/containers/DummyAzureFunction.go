package containers_test

import (
	"context"
	"encoding/json"
	"net/http"

	azurecont "github.com/pip-services3-gox/pip-services3-azure-gox/containers"
	tbuild "github.com/pip-services3-gox/pip-services3-azure-gox/test/build"
	tdata "github.com/pip-services3-gox/pip-services3-azure-gox/test/data"
	tlogic "github.com/pip-services3-gox/pip-services3-azure-gox/test/logic"
	azureutil "github.com/pip-services3-gox/pip-services3-azure-gox/utils"
	cconv "github.com/pip-services3-gox/pip-services3-commons-gox/convert"
	cdata "github.com/pip-services3-gox/pip-services3-commons-gox/data"
	cerr "github.com/pip-services3-gox/pip-services3-commons-gox/errors"
	crefer "github.com/pip-services3-gox/pip-services3-commons-gox/refer"
	cvalid "github.com/pip-services3-gox/pip-services3-commons-gox/validate"
	rpcserv "github.com/pip-services3-gox/pip-services3-rpc-gox/services"
)

type DummyAzureFunction struct {
	controller tlogic.IDummyController
	*azurecont.AzureFunction
}

func NewDummyAzureFunction() *DummyAzureFunction {
	c := DummyAzureFunction{}
	c.AzureFunction = azurecont.InheritAzureFunctionWithParams(&c, "dummy", "Dummy Azure function")
	c.DependencyResolver.Put(context.Background(), "controller", crefer.NewDescriptor("pip-services-dummies", "controller", "default", "*", "*"))
	c.AddFactory(tbuild.NewDummyFactory())
	return &c
}

func (c *DummyAzureFunction) SetReferences(ctx context.Context, references crefer.IReferences) {
	c.AzureFunction.SetReferences(ctx, references)
	resCtrl, depErr := c.DependencyResolver.GetOneRequired("controller")
	if depErr != nil {
		panic(depErr)
	}

	c.controller = resCtrl.(tlogic.IDummyController)
}

func (c *DummyAzureFunction) getPageByFilter(res http.ResponseWriter, req *http.Request) {
	params := req.URL.Query()
	paginParams := make(map[string]string, 0)

	paginParams["skip"] = params.Get("skip")
	paginParams["take"] = params.Get("take")
	paginParams["total"] = params.Get("total")

	delete(params, "skip")
	delete(params, "take")
	delete(params, "total")

	result, err := c.controller.GetPageByFilter(
		req.Context(),
		c.GetCorrelationId(req),
		cdata.NewFilterParamsFromValue(params),
		cdata.NewPagingParamsFromTuples(paginParams),
	)

	rpcserv.HttpResponseSender.SendResult(res, req, result, err)
}

func (c *DummyAzureFunction) getOneById(res http.ResponseWriter, req *http.Request) {
	correlationId := c.GetCorrelationId(req)
	var body map[string]string

	err := azureutil.AzureFunctionRequestHelper.DecodeBody(req, &body)

	if err != nil {
		err := cerr.NewInternalError(correlationId, "JSON_CNV_ERR", "Cant convert from JSON to Dummy").WithCause(err)
		rpcserv.HttpResponseSender.SendError(res, req, err)
		return
	}

	defer req.Body.Close()

	result, err := c.controller.GetOneById(
		req.Context(),
		correlationId,
		body["dummy_id"])

	rpcserv.HttpResponseSender.SendResult(res, req, result, err)
}

func (c *DummyAzureFunction) create(res http.ResponseWriter, req *http.Request) {
	correlationId := c.GetCorrelationId(req)

	dummy, err := c.getDummy(correlationId, req)

	if err != nil {
		err := cerr.NewInternalError(correlationId, "JSON_CNV_ERR", "Cant convert from JSON to Dummy").WithCause(err)
		rpcserv.HttpResponseSender.SendError(res, req, err)
		return
	}

	result, err := c.controller.Create(
		req.Context(),
		correlationId,
		dummy,
	)

	rpcserv.HttpResponseSender.SendCreatedResult(res, req, result, err)
}

func (c *DummyAzureFunction) update(res http.ResponseWriter, req *http.Request) {
	correlationId := c.GetCorrelationId(req)

	dummy, err := c.getDummy(correlationId, req)

	if err != nil {
		err := cerr.NewInternalError(correlationId, "JSON_CNV_ERR", "Cant convert from JSON to Dummy").WithCause(err)
		rpcserv.HttpResponseSender.SendError(res, req, err)
		return
	}

	result, err := c.controller.Update(
		req.Context(),
		correlationId,
		dummy,
	)
	rpcserv.HttpResponseSender.SendResult(res, req, result, err)
}

func (c *DummyAzureFunction) deleteById(res http.ResponseWriter, req *http.Request) {
	correlationId := c.GetCorrelationId(req)

	var body map[string]string

	err := azureutil.AzureFunctionRequestHelper.DecodeBody(req, &body)
	defer req.Body.Close()

	if err != nil {
		err := cerr.NewInternalError(correlationId, "JSON_CNV_ERR", "Cant convert from JSON to Dummy").WithCause(err)
		rpcserv.HttpResponseSender.SendError(res, req, err)
		return
	}

	dummyId := body["dummy_id"]

	result, err := c.controller.DeleteById(
		req.Context(),
		c.GetCorrelationId(req),
		dummyId,
	)
	rpcserv.HttpResponseSender.SendDeletedResult(res, req, result, err)
}

func (c *DummyAzureFunction) getDummy(correlationId string, req *http.Request) (tdata.Dummy, error) {
	var body map[string]any
	var dummy tdata.Dummy

	err := azureutil.AzureFunctionRequestHelper.DecodeBody(req, &body)
	defer req.Body.Close()

	if err != nil {
		return tdata.Dummy{}, err
	}

	dummyBytes, err := json.Marshal(body["dummy"])

	if err != nil {
		return tdata.Dummy{}, err
	}

	err = json.Unmarshal(dummyBytes, &dummy)

	if err != nil {
		return tdata.Dummy{}, err
	}

	return dummy, nil
}

func (c *DummyAzureFunction) Register() {

	c.RegisterAction(
		"get_dummies",
		cvalid.NewObjectSchema().WithOptionalProperty(
			"body", cvalid.NewObjectSchema().WithOptionalProperty(
				"filter", cvalid.NewFilterParamsSchema())).WithOptionalProperty(
			"paging", cvalid.NewPagingParamsSchema()).Schema,
		c.getPageByFilter,
	)

	c.RegisterAction(
		"get_dummy_by_id",
		cvalid.NewObjectSchema().WithRequiredProperty("body", cvalid.NewObjectSchema().WithRequiredProperty("dummy_id", cconv.String)).Schema,
		c.getOneById,
	)

	c.RegisterAction(
		"create_dummy",
		cvalid.NewObjectSchema().WithRequiredProperty("body", cvalid.NewObjectSchema().WithRequiredProperty("dummy", tdata.NewDummySchema())).Schema,
		c.create,
	)

	c.RegisterAction(
		"update_dummy",
		cvalid.NewObjectSchema().WithRequiredProperty("body", cvalid.NewObjectSchema().WithRequiredProperty("dummy", tdata.NewDummySchema())).Schema,
		c.update,
	)

	c.RegisterAction(
		"delete_dummy",
		cvalid.NewObjectSchema().WithRequiredProperty("body", cvalid.NewObjectSchema().WithRequiredProperty("dummy_id", cconv.String)).Schema,
		c.deleteById,
	)
}
