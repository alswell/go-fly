package db

import (
	"github.com/alswell/go-fly/examples/web/errors"
)

type Disk struct {
	ID     int    `json:"id"`
	HostID int    `json:"host_id"`
	Dev    string `json:"dev"`
}

func GetDisksByHostID(id int) []*Disk {
	if id != 1 {
		errors.HostNotFound(id)
	}
	return []*Disk{
		{2, 1, "/dev/sda"},
		{3, 1, "/dev/sdb"},
	}
}
