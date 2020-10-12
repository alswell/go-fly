package web

import (
	"fmt"
	"os"
	"reflect"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	GET    = "GET"
	POST   = "POST"
	PUT    = "PUT"
	DELETE = "DELETE"
)

type UrlRouter struct {
	Method  string
	Url     string
	Handler interface{}
}

var routers []*UrlRouter

var PreRoute func(ctx *gin.Context)

func RegisterRouter(method, url string, handler interface{}) {
	routers = append(routers, &UrlRouter{Method: method, Url: url, Handler: handler})
}

func (r *UrlRouter) preCheck() {
	f := reflect.TypeOf(r.Handler)
	if f.NumIn() > 1 {
		fmt.Println("handlers accept only 1 parameter")
		os.Exit(255)
	}
	if f.NumIn() == 1 && f.In(0).Kind() != reflect.Struct {
		fmt.Println("parameter should be struct")
		os.Exit(255)
	}
}

func (r *UrlRouter) frame(ctx *gin.Context) {
	code := 200
	var resp interface{}
	defer func() {
		ctx.JSON(code, resp)
	}()

	if PreRoute != nil {
		PreRoute(ctx)
	}

	f := reflect.TypeOf(r.Handler)
	arg := reflect.New(f.In(0))
	args := []reflect.Value{arg.Elem()}
	ctx.ShouldBindUri(arg.Interface())
	ctx.ShouldBindHeader(arg.Interface())
	if ctx.Request.Method == GET || ctx.Request.Method == DELETE {
		ctx.ShouldBindQuery(arg.Interface())
	} else {
		ctx.ShouldBindJSON(arg.Interface())
	}
	var result []reflect.Value
	if TryCatch(func() {
		result = reflect.ValueOf(r.Handler).Call(args)
	}, func(i interface{}) {
		info := i.(*abortInfo)
		code = info.Code
		resp = info.Msg
	}) {
		return
	}
	if len(result) == 0 {
		resp = "finish"
		return
	}
	switch result[0].Kind() {
	case reflect.Interface:
		resp = result[0].Interface()
	case reflect.Bool:
		if result[0].Bool() {
			resp = "success"
		} else {
			code = 500
			resp = "fail"
		}
	case reflect.Int:
		code = int(result[0].Int())
		if len(result) == 2 {
			if code == 200 {
				resp = result[1].Interface()
			} else {
				switch result[1].Kind() {
				case reflect.String:
					resp = "err msg: " + result[1].String()
				default:
					resp = result[1].Interface()
				}
			}
		} else {
			resp = "finish"
		}
	case reflect.String:
		resp = result[0].String()
	default:
		resp = result[0].Interface()
	}
}

func StartWebServer(port int) {
	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()
	engine.Use(gin.Logger())

	for _, r := range routers {
		r.preCheck()
		engine.Handle(r.Method, r.Url, r.frame)
	}

	for {
		err := engine.Run(fmt.Sprintf("0.0.0.0:%d", port))
		time.Sleep(time.Second)
		fmt.Printf("web server fail: %s, restart\n", err)
	}
}
