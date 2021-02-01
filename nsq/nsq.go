package nsq

import (
	"context"
	"os/exec"
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

// 运行着的exe文件
type ExeFile struct {
	CancelFn      context.CancelFunc //取消函数
	BackGroundCmd *exec.Cmd
	Name          string //模块名称，不包括版本
	Port          string //启动端口，方便调试查看
}

// 运行着的exe文件集合，key：名称+版本
var Maper = make(map[string]*ExeFile, 0)

// 运行着的exe占用端口集合，k:端口，v：文件名称
var MaperPort = make(map[int]string, 0)
