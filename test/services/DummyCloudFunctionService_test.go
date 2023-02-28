package services_test

import (
	"context"
	"testing"

	cconf "github.com/pip-services3-gox/pip-services3-commons-gox/config"
	"github.com/stretchr/testify/assert"
)

type DummyAzureFunctionServiceTest struct {
	fixture       *DummyAzureFunctionFixture
	funcContainer *DummyAzureFunction
}

func newDummyAzureFunctionServiceTest() *DummyAzureFunctionServiceTest {
	return &DummyAzureFunctionServiceTest{}
}

func (c *DummyAzureFunctionServiceTest) setup(t *testing.T) {
	config := cconf.NewConfigParamsFromTuples(
		"logger.descriptor", "pip-services:logger:console:default:1.0",
		"service.descriptor", "pip-services-dummies:service:azurefunc:default:1.0",
	)

	ctx := context.Background()

	c.funcContainer = NewDummyAzureFunction()
	c.funcContainer.Configure(ctx, config)
	err := c.funcContainer.Open(ctx, "")
	assert.Nil(t, err)

	c.fixture = NewDummyAzureFunctionFixture(c.funcContainer.GetHandler())
}

func (c *DummyAzureFunctionServiceTest) teardown(t *testing.T) {
	err := c.funcContainer.Close(context.Background(), "")
	assert.Nil(t, err)
}

func TestCrudOperationsAzureService(t *testing.T) {
	c := newDummyAzureFunctionServiceTest()
	if c == nil {
		return
	}

	c.setup(t)
	t.Run("CRUD Operations", c.fixture.TestCrudOperations)
	c.teardown(t)
}
