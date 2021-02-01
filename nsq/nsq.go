package nsq

import (
	"sync"

	"github.com/nsqio/nsq/nsqd"
)

type RqProgram struct {
	Once sync.Once
	Nsqd *nsqd.NSQD
}

//f：回调函数，在关闭之前执行函数
func (p *RqProgram) Stop(f func()) error {
	f()
	p.Once.Do(func() {
		p.Nsqd.Exit()
	})
	return nil
}

var Pg RqProgram
