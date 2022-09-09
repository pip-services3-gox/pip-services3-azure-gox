package clients_test

import (
	"context"

	tdata "github.com/pip-services3-gox/pip-services3-azure-gox/test/data"
	cdata "github.com/pip-services3-gox/pip-services3-commons-gox/data"
)

type IDummyClient interface {
	GetDummies(ctx context.Context, correlationId string, filter cdata.FilterParams, paging cdata.PagingParams) (result cdata.DataPage[tdata.Dummy], err error)
	GetDummyById(ctx context.Context, correlationId string, dummyId string) (result tdata.Dummy, err error)
	CreateDummy(ctx context.Context, correlationId string, dummy tdata.Dummy) (result tdata.Dummy, err error)
	UpdateDummy(ctx context.Context, correlationId string, dummy tdata.Dummy) (result tdata.Dummy, err error)
	DeleteDummy(ctx context.Context, correlationId string, dummyId string) (result tdata.Dummy, err error)
}
