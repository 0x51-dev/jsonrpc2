package jsonrpc2

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

func NewHandlerFunc(h HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			writeError(w, nil, NewMethodNotFoundError())
			return
		}
		if r.Header.Get("Content-Type") != "application/json" ||
			r.Header.Get("Accept") != "application/json" {
			writeError(w, nil, NewInvalidRequestError())
			return
		}

		data, err := io.ReadAll(r.Body)
		if err != nil {
			writeError(w, nil, NewInternalError())
			return
		}

		var s []json.RawMessage
		if err := json.Unmarshal(data, &s); err == nil {
			if len(s) == 0 {
				writeError(w, nil, NewInvalidRequestError())
				return
			}
			var responses []*Response
			for _, v := range s {
				resp := decodeRequest(v, h)
				if resp != nil {
					responses = append(responses, resp)
				}
			}
			if len(responses) == 0 {
				w.WriteHeader(http.StatusNoContent)
				return
			}
			data, _ := json.Marshal(responses)
			w.WriteHeader(http.StatusOK)
			w.Write(data)
			return
		}

		resp := decodeRequest(data, h)
		if resp != nil {
			data, _ := json.Marshal(resp)
			if resp.Error != nil {
				w.WriteHeader(statusCode(resp.Error))
			} else {
				w.WriteHeader(http.StatusOK)
			}
			w.Write(data)
		} else {
			w.WriteHeader(http.StatusNoContent)
		}
	}
}

func writeError(w http.ResponseWriter, id any, err error) {
	switch e := err.(type) {
	case *Error:
		data, _ := json.Marshal(Response{
			JSONRPC: VERSION,
			ID:      id,
			Error:   e,
		})
		w.WriteHeader(statusCode(e))
		w.Write(data)
	default:
		data, _ := json.Marshal(Response{
			JSONRPC: VERSION,
			ID:      id,
			Error:   NewInternalError(),
		})
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(data)
	}
}

type HandlerFunc func(request *RPCRequest) (any, *Error)

type RPCRequest struct {
	Method string `json:"method"`
	Params any    `json:"params,omitempty"`
}

func decodeRequest(data []byte, h HandlerFunc) *Response {
	var id any
	var sID stringIdentifier
	if err := json.Unmarshal(data, &sID); err == nil {
		id = sID.ID
		if sID.JSONRPC != VERSION {
			return &Response{
				JSONRPC: VERSION,
				ID:      id,
				Error:   NewInvalidRequestError(),
			}
		}
	} else {
		var iID intIdentifier
		decoder := json.NewDecoder(bytes.NewReader(data))
		decoder.UseNumber()
		if err := decoder.Decode(&iID); err != nil {
			switch err.(type) {
			case *json.SyntaxError:
				return &Response{
					JSONRPC: VERSION,
					Error:   NewParseError(),
				}
			}
			return &Response{
				JSONRPC: VERSION,
				Error:   NewInvalidRequestError(),
			}
		}

		id = iID.ID
		if iID.JSONRPC != VERSION {
			return &Response{
				JSONRPC: VERSION,
				ID:      id,
				Error:   NewInvalidRequestError(),
			}
		}
	}

	var request RPCRequest
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.UseNumber()
	if err := decoder.Decode(&request); err != nil {
		switch err.(type) {
		case *json.UnmarshalTypeError:
			return &Response{
				JSONRPC: VERSION,
				ID:      id,
				Error:   NewInvalidRequestError(),
			}
		default:
			return &Response{
				JSONRPC: VERSION,
				ID:      id,
				Error:   NewParseError(),
			}
		}
	}
	result, err := h(&request)
	if err != nil {
		return &Response{
			JSONRPC: VERSION,
			ID:      id,
			Error:   err,
		}
	}
	if result != nil {
		return &Response{
			JSONRPC: VERSION,
			ID:      id,
			Result:  result,
		}
	}
	return nil
}

type intIdentifier struct {
	JSONRPC string `json:"jsonrpc"`
	ID      int    `json:"id"`
}

type stringIdentifier struct {
	JSONRPC string  `json:"jsonrpc"`
	ID      *string `json:"id"`
}
