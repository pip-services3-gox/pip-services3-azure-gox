package clients_test

import (
	"context"

	azureclient "github.com/pip-services3-gox/pip-services3-azure-gox/clients"
	tdata "github.com/pip-services3-gox/pip-services3-azure-gox/test/data"
	cdata "github.com/pip-services3-gox/pip-services3-commons-gox/data"
	rpcclients "github.com/pip-services3-gox/pip-services3-rpc-gox/clients"
)

type DummyAzureFunctionClient struct {
	*azureclient.AzureFunctionClient
}

func NewDummyAzureFunctionClient() *DummyAzureFunctionClient {
	return &DummyAzureFunctionClient{
		AzureFunctionClient: azureclient.NewAzureFunctionClient(),
	}
}

func (c *DummyAzureFunctionClient) GetDummies(ctx context.Context, correlationId string, filter cdata.FilterParams, paging cdata.PagingParams) (result cdata.DataPage[tdata.Dummy], err error) {
	timing := c.Instrument(ctx, correlationId, "dummies.get_dummies")

	response, err := c.Call(ctx, "dummies.get_dummies", correlationId, nil)
	if err != nil {
		return cdata.DataPage[tdata.Dummy]{}, err
	}

	defer timing.EndTiming(ctx, err)
	return rpcclients.HandleHttpResponse[cdata.DataPage[tdata.Dummy]](response, correlationId)
}

func (c *DummyAzureFunctionClient) GetDummyById(ctx context.Context, correlationId string, dummyId string) (result tdata.Dummy, err error) {
	timing := c.Instrument(ctx, correlationId, "dummies.get_dummy_by_id")

	response, err := c.Call(ctx, "dummies.get_dummy_by_id", correlationId, cdata.NewAnyValueMapFromTuples("dummy_id", dummyId))
	if err != nil {
		return tdata.Dummy{}, err
	}

	defer timing.EndTiming(ctx, err)
	return rpcclients.HandleHttpResponse[tdata.Dummy](response, correlationId)
}

func (c *DummyAzureFunctionClient) CreateDummy(ctx context.Context, correlationId string, dummy tdata.Dummy) (result tdata.Dummy, err error) {
	timing := c.Instrument(ctx, correlationId, "dummies.create_dummy")

	response, err := c.Call(ctx, "dummies.create_dummy", correlationId, cdata.NewAnyValueMapFromTuples("dummy", dummy))
	if err != nil {
		return tdata.Dummy{}, err
	}

	defer timing.EndTiming(ctx, err)
	return rpcclients.HandleHttpResponse[tdata.Dummy](response, correlationId)
}

func (c *DummyAzureFunctionClient) UpdateDummy(ctx context.Context, correlationId string, dummy tdata.Dummy) (result tdata.Dummy, err error) {
	timing := c.Instrument(ctx, correlationId, "dummies.update_dummy")

	response, err := c.Call(ctx, "dummies.update_dummy", correlationId, cdata.NewAnyValueMapFromTuples("dummy", dummy))
	if err != nil {
		return tdata.Dummy{}, err
	}

	defer timing.EndTiming(ctx, err)
	return rpcclients.HandleHttpResponse[tdata.Dummy](response, correlationId)
}

func (c *DummyAzureFunctionClient) DeleteDummy(ctx context.Context, correlationId string, dummyId string) (result tdata.Dummy, err error) {
	timing := c.Instrument(ctx, correlationId, "dummies.delete_dummy")

	response, err := c.Call(ctx, "dummies.delete_dummy", correlationId, cdata.NewAnyValueMapFromTuples("dummy_id", dummyId))
	if err != nil {
		return tdata.Dummy{}, err
	}

	defer timing.EndTiming(ctx, err)
	return rpcclients.HandleHttpResponse[tdata.Dummy](response, correlationId)
}
