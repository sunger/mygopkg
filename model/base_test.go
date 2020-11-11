package model

import (
	"testing"
)

func TestBsModel(t *testing.T) {

	r := &BsModel{}
	r.CreateId()
	// 如果测试结果符合预期的输出信息
	t.Log("初始化正常.", r.Id)
}
