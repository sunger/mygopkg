package db

import (
	"testing"
)

func TestBsModel(t *testing.T) {

	InitDb("development")

	// 如果测试结果符合预期的输出信息
	t.Log("数据库初始化正常.")
}
