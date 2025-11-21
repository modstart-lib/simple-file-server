package errors

type BusinessError struct {
	Msg    string
	Detail interface{}
	Map    map[string]interface{}
	Err    error
}
