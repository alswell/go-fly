package api

import (
	"fmt"
	"github.com/alswell/go-fly/examples/web/db"
	"github.com/alswell/go-fly/web"
)

func init() {
	web.RegisterRouter(web.GET, "/hosts",
		func(param struct {
			IDs   []int  `form:"id"`
			Token string `header:"Token"`
		}) interface{} {
			fmt.Println("get hosts", param.IDs, param.Token)
			if len(param.IDs) == 0 {
				return db.GetHosts()
			}

			var r []*db.Host
			for _, id := range param.IDs {
				r = append(r, db.GetHostByID(id))
			}
			return r
		},
	)
	web.RegisterRouter(web.GET, "/hosts/:id",
		func(param struct {
			ID    int    `uri:"id"`
			Token string `header:"Token"`
		}) *db.Host {
			fmt.Println("get host", param.ID, param.Token)
			return db.GetHostByID(param.ID)
		},
	)
	web.RegisterRouter(web.POST, "/hosts",
		func(param struct {
			Name  string `json:"name"`
			Token string `header:"Token"`
		}) bool {
			fmt.Println("create host", param.Name, param.Token)
			return db.CreateHost(param.Name)
		},
	)
}
