package config

import (
    "testing"
)

func TestConfig(t *testing.T) {
    f := "fc.yaml"
	Cfg = LoadConfig(f)
    if Cfg == nil {
        t.Error("global configuration is nil")
    }
    Cfg.Print()
    t.Log("finish configuration test")
}
