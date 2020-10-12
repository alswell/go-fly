package errors

import (
	"fmt"
	"github.com/alswell/go-fly/web"
)

func HostNotFound(id int) {
	web.Abort(404, fmt.Sprintf("host[%d] not found", id))
}
