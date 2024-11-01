// convert コンバートパッケージ
package convert // パッケージ名はディレクトリ名と同じにする

import (
	"log/slog"
	"regexp"
	"strconv"
)

// ---- Global Variable

// ---- Package Global Variable

//---- public function ----

// 文字列から数字(int)を抜き出す
func ExtractInt64(str string) int64 {
	rex := regexp.MustCompile("[0-9]+")
	str = rex.FindString(str)
	intVal, err := strconv.ParseInt(str, 10, 32)
	if err != nil {
		slog.String("error", err.Error())
		return 0
	}
	return intVal
}

//---- private function ----
