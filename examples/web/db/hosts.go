package db

import (
	"github.com/alswell/go-fly/examples/web/errors"
	"github.com/alswell/go-fly/web"
)

type Host struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func GetHostByID(id int) *Host {
	if id != 0 && id != 1 {
		errors.HostNotFound(id)
	}
	return &Host{id, "localhost"}
}

func GetHosts() []*Host {
	return []*Host{{0, "localhost"}, {1, "localhost"}}
}

func CreateHost(name string) bool {
	web.ServerError("create host: " + name)
	return true
}
