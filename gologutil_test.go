package gologutil

import "testing"

//单元测试
func TestLog(t *testing.T) {
	defer Sync()

	for i := 0; i < 10; i++ {
		Debug("debug log", "level", "DEBUG LOGS")
		DebugWithFields("test fields", map[string]string{"111": "aaa", "222": "bbb"})

		Info("info log", "level", "INFO LOGS")
		InfoWithFields("test fields", map[string]string{"333": "ccc", "444": "ddd"})
	}

	PrintConf()
}
