package domain

type CommentResult struct {
	Code int // 0 success, 1 fail
	Data interface{}
	Msg  string
}
