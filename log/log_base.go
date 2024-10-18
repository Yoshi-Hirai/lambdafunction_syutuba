// log log/slogをベースとしたログシステムパッケージ
package log // パッケージ名はディレクトリ名と同じにする

import (
	"log/slog"
	"os"
	"sync"
)

// ---- Global Variable

// ---- Package Global Variable
var once sync.Once
var jsonLogger *slog.Logger //jsonフォーマットlogger(構造化ログ)

//---- public function ----

// GetInstance (public)logシステムの初期化と(初回のみ)json(構造化)ログパッケージシングルトンの取得
func GetInstance() *slog.Logger {
	once.Do(func() {
		logBaseInit()
	})
	return jsonLogger
}

//---- private function ----

// logBaseInit (private)slogシステムの初期化 & テキスト用ログロガーの設定 & json(構造化)ログロガーSingleton生成関数
func logBaseInit() {
	// テキスト用ログロガー
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	// logを出力するカレントの関数名を出力しようとしたが、関数名の取得がうまくいかなかった。
	logger.Info("TextLogger Create Success.")
	slog.SetDefault(logger)

	// JSON用ログロガー
	jsonLogger = slog.New(slog.NewJSONHandler(os.Stdout, nil))
}
