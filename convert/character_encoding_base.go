// convert コンバートパッケージ
package convert // パッケージ名はディレクトリ名と同じにする

import (
	"io"
	"strings"

	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

// ---- Global Variable

// ---- Package Global Variable

//---- public function ----

// UTF-8 から ShiftJIS
func Utf8ToSjis(str string) (string, error) {
	ret, err := io.ReadAll(transform.NewReader(strings.NewReader(str), japanese.ShiftJIS.NewEncoder()))
	if err != nil {
		return "", err
	}
	return string(ret), err
}

// ShiftJIS から UTF-8
func SjisToUtf8(str string) (string, error) {
	ret, err := io.ReadAll(transform.NewReader(strings.NewReader(str), japanese.ShiftJIS.NewDecoder()))
	if err != nil {
		return "", err
	}
	return string(ret), err
}

// UTF-8 から EUC-JP
func Utf8ToEucjp(str string) (string, error) {
	ret, err := io.ReadAll(transform.NewReader(strings.NewReader(str), japanese.EUCJP.NewEncoder()))
	if err != nil {
		return "", err
	}
	return string(ret), err
}

// EUC-JP から UTF-8
func EucjpToUtf8(str string) (string, error) {
	ret, err := io.ReadAll(transform.NewReader(strings.NewReader(str), japanese.EUCJP.NewDecoder()))
	if err != nil {
		return "", err
	}
	return string(ret), err
}

//---- private function ----
