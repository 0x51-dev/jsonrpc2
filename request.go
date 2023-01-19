package jsonrpc2

const VERSION = "2.0"

// A rpc call is represented by sending a Request object to a Server.
type Request struct {
	// A String specifying the version of the JSON-RPC protocol. MUST be exactly "2.0".
	JSONRPC string `json:"jsonrpc"`
	// A String containing the name of the method to be invoked. Method names that begin with the word rpc followed by
	// a period character (U+002E or ASCII 46) are reserved for rpc-internal methods and extensions and MUST NOT be used
	// for anything else.
	Method string `json:"method"`
	// A Structured value that holds the parameter values to be used during the invocation of the method. This member
	// MAY be omitted.
	Params any `json:"params,omitempty"`
	// An identifier established by the Client that MUST contain a String, Number, or NULL value if included. If it is
	// not included it is assumed to be a notification. The value SHOULD normally not be Null and Numbers SHOULD NOT
	// contain fractional parts.
	ID any `json:"id"`
}

func NewRequest(method string, params ...any) *Request {
	return &Request{
		ID:      nil,
		Method:  method,
		Params:  params,
		JSONRPC: VERSION,
	}
}

func NewRequestWithID(id int, method string, params ...any) *Request {
	return &Request{
		ID:      id,
		Method:  method,
		Params:  params,
		JSONRPC: VERSION,
	}
}

func NewRequestWithIDString(id string, method string, params ...any) *Request {
	return &Request{
		ID:      id,
		Method:  method,
		Params:  params,
		JSONRPC: VERSION,
	}
}
