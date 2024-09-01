package domain

type CommonResult struct {
	Code int
	Data interface{}
	Msg  string
}

func (result CommonResult) ToJson() (interface{}, error) {

	return "", nil
}

// NewErrorRespWithErr NewErrorResponse 生成 Error CommonResult
func NewErrorRespWithErr(code int, err error) *CommonResult {
	if code >= 0 {
		code = -1
	}
	return &CommonResult{
		Code: code,
		Msg:  err.Error(),
		Data: nil,
	}
}

func NewErrorRespWithMsg(code int, msg string) *CommonResult {
	if code >= 0 {
		code = -1
	}
	return &CommonResult{
		Code: code,
		Msg:  msg,
		Data: nil,
	}
}

// NewSuccessResp NewSuccessResponse 生成 Success Response
func NewSuccessResp(data interface{}) *CommonResult {
	return &CommonResult{
		Code: 0,
		Msg:  "success",
		Data: data,
	}
}
func NewSuccessRespWithMsg(data interface{}, msg string) *CommonResult {
	return &CommonResult{
		Code: 0,
		Msg:  msg,
		Data: data,
	}
}
