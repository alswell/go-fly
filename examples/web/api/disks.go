package api

import (
	"github.com/alswell/go-fly/examples/web/db"
	"github.com/alswell/go-fly/web"
)

func init() {
	web.RegisterRouter(web.GET, "/hosts/:id/disks",
		func(param struct {
			ID    int    `uri:"id"`
			Token string `header:"Token"`
		}) []*db.Disk {
			return db.GetDisksByHostID(param.ID)
		},
	)
	web.RegisterRouter(web.POST, "/hosts/:id/disks",
		func() string {
			web.Abort(500, "not implemented")
			return ""
		},
	)

}
