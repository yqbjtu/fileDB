package domain

type Result struct {
	Code int
	Data interface{}
	Msg  string
}

func (result Result) ToJson() (interface{}, error) {

	return "", nil
}
