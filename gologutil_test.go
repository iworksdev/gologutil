package gologutil

import (
	"testing"
)

//单元测试
func TestLog(t *testing.T) {
	defer Sync()

	for i := 0; i < 2; i++ {
		Debug("debug log", "level", "DEBUG LOGS")
		DebugWithFields("test fields", map[string]string{"111": "aaa", "222": "bbb"})

		Info("info log", "level", "INFO LOGS")
		InfoWithFields("test fields", map[string]string{"333": "ccc", "444": "ddd"})

		Warn("warn log", "level", "WARN LOGS")
		WarnWithFields("test fields", map[string]string{"555": "eee", "666": "fff"})

		//Fatal("fatal log", "level", "FATAL LOGS")
		//FatalWithFields("test fields", map[string]string{"777": "ggg", "888": "hhh"})

	}

	PrintConf()
}
