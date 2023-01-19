package jsonrpc2_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/0x51-dev/jsonrpc2"
	"github.com/0x51-dev/jsonrpc2/params"
)

func TestHandler(t *testing.T) {
	server := httptest.NewServer(jsonrpc2.NewHandlerFunc(func(request *jsonrpc2.RPCRequest) (any, *jsonrpc2.Error) {
		if strings.HasPrefix(request.Method, "notify_") {
			return nil, nil
		}
		switch request.Method {
		case "sum":
			s, ok := params.GetSlice[json.Number](request.Params)
			if !ok {
				return nil, jsonrpc2.NewInvalidParamsError()
			}
			var sum int64
			for _, v := range s {
				s, _ := v.Int64()
				sum += s
			}
			return sum, nil
		case "subtract":
			var s0, s1 int64
			s, ok := params.GetSlice[json.Number](request.Params)
			if ok && len(s) == 2 {
				s0, _ = s[0].Int64()
				s1, _ = s[1].Int64()
			} else {
				m, ok := params.GetMap[json.Number](request.Params)
				if !ok {
					return nil, jsonrpc2.NewInvalidParamsError()
				}
				s0, _ = m["minuend"].Int64()
				s1, _ = m["subtrahend"].Int64()
			}
			return s0 - s1, nil
		case "get_data":
			return []any{"hello", 5}, nil
		default:
			return nil, jsonrpc2.NewMethodNotFoundError()
		}
	}))

	for _, test := range []struct {
		name     string
		request  string
		response string
	}{
		{
			name:     "positional parameters 1",
			request:  `{"jsonrpc": "2.0", "method": "subtract", "params": [42, 23], "id": 1}`,
			response: `{"jsonrpc": "2.0", "result": 19, "id": 1}`,
		},
		{
			name:     "positional parameters 2",
			request:  `{"jsonrpc": "2.0", "method": "subtract", "params": [23, 42], "id": 1}`,
			response: `{"jsonrpc": "2.0", "result": -19, "id": 1}`,
		},
		{
			name:     "named parameters 1",
			request:  `{"jsonrpc": "2.0", "method": "subtract", "params": {"subtrahend": 23, "minuend": 42}, "id": 3}`,
			response: `{"jsonrpc": "2.0", "result": 19, "id": 3}`,
		},
		{
			name:     "named parameters 2",
			request:  `{"jsonrpc": "2.0", "method": "subtract", "params": {"minuend": 42, "subtrahend": 23}, "id": 4}`,
			response: `{"jsonrpc": "2.0", "result": 19, "id": 4}`,
		},
		{
			name:    "notification 1",
			request: `{"jsonrpc": "2.0", "method": "notify_update", "params": [1,2,3,4,5]}`,
		},
		{
			name:    "notification 2",
			request: `{"jsonrpc": "2.0", "method": "notify_foobar"}`,
		},
		{
			name:     "non-existent method",
			request:  `{"jsonrpc": "2.0", "method": "foobar", "id": "1"}`,
			response: `{"jsonrpc": "2.0", "error": {"code": -32601, "message": "Method not found"}, "id": "1"}`,
		},
		{
			name:     "invalid json",
			request:  `{"jsonrpc": "2.0", "method": "foobar, "params": "bar", "baz]`,
			response: `{"jsonrpc": "2.0", "error": {"code": -32700, "message": "Parse error"}, "id": null}`,
		},
		{
			name:     "invalid request object",
			request:  `{"jsonrpc": "2.0", "method": 1, "params": "bar"}`,
			response: `{"jsonrpc": "2.0", "error": {"code": -32600, "message": "Invalid Request"}, "id": null}`,
		},
		{
			name: "batch invalid json",
			request: `[
				{"jsonrpc": "2.0", "method": "sum", "params": [1,2,4], "id": "1"},
				{"jsonrpc": "2.0", "method"
			]`,
			response: `{"jsonrpc": "2.0", "error": {"code": -32700, "message": "Parse error"}, "id": null}`,
		},
		{
			name:     "batch empty array",
			request:  `[]`,
			response: `{"jsonrpc": "2.0", "error": {"code": -32600, "message": "Invalid Request"}, "id": null}`,
		},
		{
			name:    "batch invalid",
			request: `[1]`,
			response: `[
				{"jsonrpc": "2.0", "error": {"code": -32600, "message": "Invalid Request"}, "id": null}
			]`,
		},
		{
			name:    "batch invalid 2",
			request: `[1, 2, 3]`,
			response: `[
				{"jsonrpc": "2.0", "error": {"code": -32600, "message": "Invalid Request"}, "id": null},
				{"jsonrpc": "2.0", "error": {"code": -32600, "message": "Invalid Request"}, "id": null},
				{"jsonrpc": "2.0", "error": {"code": -32600, "message": "Invalid Request"}, "id": null}
			]`,
		},
		{
			name: "batch",
			request: `[
				{"jsonrpc": "2.0", "method": "sum", "params": [1,2,4], "id": "1"},
				{"jsonrpc": "2.0", "method": "notify_hello", "params": [7]},
				{"jsonrpc": "2.0", "method": "subtract", "params": [42,23], "id": "2"},
				{"foo": "boo"},
				{"jsonrpc": "2.0", "method": "foo.get", "params": {"name": "myself"}, "id": "5"},
				{"jsonrpc": "2.0", "method": "get_data", "id": "9"} 
			]`,
			response: `[
				{"jsonrpc": "2.0", "result": 7, "id": "1"},
				{"jsonrpc": "2.0", "result": 19, "id": "2"},
				{"jsonrpc": "2.0", "error": {"code": -32600, "message": "Invalid Request"}, "id": null},
				{"jsonrpc": "2.0", "error": {"code": -32601, "message": "Method not found"}, "id": "5"},
				{"jsonrpc": "2.0", "result": ["hello", 5], "id": "9"}
			]`,
		},
		{
			name: "batch notify",
			request: `[
				{"jsonrpc": "2.0", "method": "notify_sum", "params": [1,2,4]},
				{"jsonrpc": "2.0", "method": "notify_hello", "params": [7]}
			]`,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			req, _ := http.NewRequest(http.MethodPost, server.URL, bytes.NewReader([]byte(test.request)))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Accept", "application/json")
			if resp, err := http.DefaultClient.Do(req); err != nil {
				t.Error(err)
			} else {
				if data, err := io.ReadAll(resp.Body); err != nil {
					t.Error(err)
				} else {
					if len(test.response) == 0 {
						if resp.StatusCode != 204 {
							t.Error(resp.Status)
						}
					} else {
						var resp any
						if err := json.Unmarshal(data, &resp); err != nil {
							t.Error(err)
						}
						var response any
						if err := json.Unmarshal([]byte(test.response), &response); err != nil {
							t.Error(err)
						}
						if !reflect.DeepEqual(resp, response) {
							t.Error(resp, response)
						}
					}
				}
			}
		})
	}
}
