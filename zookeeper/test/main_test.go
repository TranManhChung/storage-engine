package test

import (
	"os"
	"testing"
	"time"

	"github.com/zk"
)

const (
	perms           = 0x1f
	totalPack       = 2000
	totalCCUeachIns = 100
	lockPath        = "/lock"
	groupPath       = "/group"
	totalInstance   = 1
)

var (
	acl   = []zk.ACL{{perms, "world", "anyone"}}
	conns = make([]*zk.Conn, totalInstance) // create a connection for each instance
)

type Group struct {
	Packages []int `json:"package"`
}

func TestMain(m *testing.M) {
	var err error
	// new connection to zk
	for i := 0; i < len(conns); i++ {
		if conns[i], _, err = zk.Connect([]string{"172.17.0.2"}, time.Second); err != nil {
			panic(err)
		}
	}
	os.Exit(m.Run())
}
