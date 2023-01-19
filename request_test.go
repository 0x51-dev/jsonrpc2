package jsonrpc2_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/0x51-dev/jsonrpc2"
)

func TestNewRequest(t *testing.T) {
	server := httptest.NewServer(jsonrpc2.NewHandlerFunc(func(request *jsonrpc2.RPCRequest) (any, *jsonrpc2.Error) {
		return nil, nil
	}))
	for _, req := range []*jsonrpc2.Request{
		jsonrpc2.NewRequest("notify_null"),
		jsonrpc2.NewRequestWithID(0, "notify_id_int"),
		jsonrpc2.NewRequestWithIDString("some-id", "notify_string"),
	} {
		data, err := json.Marshal(req)
		if err != nil {
			t.Fatal(err)
		}
		req, _ := http.NewRequest(http.MethodPost, server.URL, bytes.NewReader(data))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Accept", "application/json")
		if resp, err := http.DefaultClient.Do(req); err != nil {
			t.Error(err)
		} else {
			if resp.StatusCode != http.StatusNoContent {
				data, _ := io.ReadAll(resp.Body)
				t.Error(resp.StatusCode, string(data))
			}
		}
	}
	for _, req := range []*jsonrpc2.Request{
		{ID: 0.0},
		{ID: true},
		// etc.
	} {
		data, err := json.Marshal(req)
		if err != nil {
			t.Fatal(err)
		}
		req, _ := http.NewRequest(http.MethodPost, server.URL, bytes.NewReader(data))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Accept", "application/json")
		if resp, err := http.DefaultClient.Do(req); err != nil {
			t.Error(err)
		} else {
			if resp.StatusCode != http.StatusBadRequest {
				data, _ := io.ReadAll(resp.Body)
				t.Error(resp.StatusCode, string(data))
			}
		}
	}
}
