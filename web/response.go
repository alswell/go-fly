package web

type abortInfo struct {
	Code int
	Msg  string
}

func Abort(code int, msg string) {
	panic(&abortInfo{code, msg})
}

func NotFound(x string) {
	Abort(404, x+" not found")
}

func ServerError(x string) {
	Abort(500, "internal error: "+x)
}
