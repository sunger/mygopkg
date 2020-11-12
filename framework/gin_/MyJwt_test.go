package gin_

import (
	"testing"
)

func TestBsModel(t *testing.T) {

	user:=&UserInfo{ID: "12315748"}
	token,_:=CreateToken(user)
	t.Log(token)
	// 如果测试结果符合预期的输出信息
	t.Log("token生成正常.")

	user2,err := ParseToken(token)
	if err!=nil{
		t.Log("token解析错误.",err)
	}
	t.Log("token解析正常.",user2.ID)
}
