package comm

import (
	"testing"
)

func TestBsMi(t *testing.T) {

	jmUtil := JiaJieMi{}

	//加密密码
	jmpwd := jmUtil.Jia("123")
	t.Log(jmpwd)
}
