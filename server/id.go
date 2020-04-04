package server

import (
	"fmt"
	"time"

	"github.com/nybuxtsui/log"
	"github.com/sony/sonyflake"
)

var (
	idCh = make(chan uint64)
)

func init() {
	go idWorker()
}

func idWorker() {
	flake := sonyflake.NewSonyflake(sonyflake.Settings{})
	for {
		id, err := flake.NextID()
		if err != nil {
			log.Error("genId failed: %s", err.Error())
			time.Sleep(time.Second)
		} else {
			idCh <- id
		}

	}
}

func genId() string {
	id := <-idCh
	return fmt.Sprintf("%d", id)
}
