package http_test

import (
	"context"
	"testing"

	"net/http"

	"github.com/afex/hystrix-go/hystrix"
	"github.com/go-kit/kit/endpoint"
	"github.com/kolbis/corego/errors"
	tlehttp "github.com/kolbis/corego/transport/http"
)

func TestSuccessExecute(t *testing.T) {
	ctx := context.Background()
	req := 1
	res := tlehttp.Execute(ctx, req, getSuccessfullExecutionEndpoint())

	if res.IsCircuitOpen == true {
		t.Error("circuit should not be open")
	}
	if res.IsReachedRateLimit == true {
		t.Error("limit should not have reached")
	}
	if res.Error != nil {
		t.Error("didnt expect an error")
	}
	if res.StatusCode != 200 {
		t.Error("wrong status code")
	}
	if res.Data.(string) != "ok" {
		t.Error("wrong data")
	}
}

func TestCircuitOpenedExecute(t *testing.T) {
	ctx := context.Background()
	req := 1
	res := tlehttp.Execute(ctx, req, getCircuitOpenedExecutionEndpoint())

	if res.IsCircuitOpen != true {
		t.Error("circuit should be open")
	}
	if res.IsReachedRateLimit == true {
		t.Error("limit should not have reached")
	}
	if res.Error == nil {
		t.Error("expected an error")
	}
	if res.Data.(string) != "nok" {
		t.Error("wrong data")
	}
}

func TestServerSideErrorExecute(t *testing.T) {
	ctx := context.Background()
	req := 1
	res := tlehttp.Execute(ctx, req, getServerErrorExecutionEndpoint())

	if res.StatusCode != 500 {
		t.Error("status code should have been 500")
	}
	if res.Error == nil && errors.IsApplicationError(res.Error) {
		t.Error("expected an error")
	}
}

func getSuccessfullExecutionEndpoint() endpoint.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		return "ok", nil
	}
}

func getCircuitOpenedExecutionEndpoint() endpoint.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		return "nok", hystrix.ErrCircuitOpen
	}
}

func getServerErrorExecutionEndpoint() endpoint.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		response := &http.Response{
			StatusCode: 500,
		}
		return response, nil
	}
}
