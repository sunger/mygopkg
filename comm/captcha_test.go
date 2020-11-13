package comm

import (
	"testing"
)

func TestCp(t *testing.T) {

	id, b64s, err := DriverDigitFunc()
	if err != nil {
		t.Log(err)
		return
	}
	t.Log(id, b64s)
}
