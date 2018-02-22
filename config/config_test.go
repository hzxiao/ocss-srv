package config

import (
	"testing"
)

func TestConfig(t *testing.T) {
	InitConfig("config", ".")
	PrintAll()

	port := GetString("server.port")
	if port == "" {
		t.Error("get string error")
		return
	}
	t.Log(port)
}
