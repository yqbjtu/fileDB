package domain

type Result struct {
	Code int
	Data interface{}
	Msg  string
}

func (result Result) ToJson() (interface{}, error) {

	return "", nil
}

// NewErrorRespWithErr NewErrorResponse 生成 Error Result
func NewErrorRespWithErr(code int, err error) *Result {
	if code >= 0 {
		code = -1
	}
	return &Result{
		Code: code,
		Msg:  err.Error(),
		Data: nil,
	}
}

func NewErrorRespWithMsg(code int, msg string) *Result {
	if code >= 0 {
		code = -1
	}
	return &Result{
		Code: code,
		Msg:  msg,
		Data: nil,
	}
}

// NewSuccessResp NewSuccessResponse 生成 Success Response
func NewSuccessResp(data interface{}) *Result {
	return &Result{
		Code: 0,
		Msg:  "success",
		Data: data,
	}
}
func NewSuccessRespWithMsg(data interface{}, msg string) *Result {
	return &Result{
		Code: 0,
		Msg:  msg,
		Data: data,
	}
}
