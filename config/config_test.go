package config

import (
	"testing"
)

func TestCfg(t *testing.T) {

	Init("development", "")
	Init("development", "F:\\go\\mygopkg")
}
