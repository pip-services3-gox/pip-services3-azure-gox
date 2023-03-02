package services_test

import (
	"context"
	"testing"

	cconf "github.com/pip-services3-gox/pip-services3-commons-gox/config"
	"github.com/stretchr/testify/assert"
)

type DummyCommandableAzureFunctionServiceTest struct {
	fixture       *DummyAzureFunctionFixture
	funcContainer *DummyAzureFunction
}

func newDummyCommandableAzureFunctionServiceTest() *DummyCommandableAzureFunctionServiceTest {
	return &DummyCommandableAzureFunctionServiceTest{}
}

func (c *DummyCommandableAzureFunctionServiceTest) setup(t *testing.T) {
	config := cconf.NewConfigParamsFromTuples(
		"logger.descriptor", "pip-services:logger:console:default:1.0",
		"service.descriptor", "pip-services-dummies:service:commandable-azurefunc:default:1.0",
	)

	ctx := context.Background()

	c.funcContainer = NewDummyAzureFunction()
	c.funcContainer.Configure(ctx, config)
	err := c.funcContainer.Open(ctx, "")
	assert.Nil(t, err)

	c.fixture = NewDummyAzureFunctionFixture(c.funcContainer.GetHandler())
}

func (c *DummyCommandableAzureFunctionServiceTest) teardown(t *testing.T) {
	err := c.funcContainer.Close(context.Background(), "")
	assert.Nil(t, err)
}

func TestCrudOperationsCommandableService(t *testing.T) {
	c := newDummyCommandableAzureFunctionServiceTest()
	if c == nil {
		return
	}

	c.setup(t)
	t.Run("CRUD Operations", c.fixture.TestCrudOperations)
	c.teardown(t)
}
