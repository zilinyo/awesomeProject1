package log

/*
@Time : 2021/5/11 3:03 PM
@Author : apple
@File : log.go
@Software: GoLand
*/
import (
	"fmt"
	dialog "github.com/go-kit/kit/log"
	slog "log"
	"os"
	"time"
)

var Logs dialog.Logger

func init() {
	t := time.Now().Format("2006-01")
	t += "wechat-"
	f := fmt.Sprintf("%s%s.log", "./", t)
	src, err := os.OpenFile(f, os.O_APPEND|os.O_WRONLY|os.O_CREATE, os.ModeAppend|os.ModePerm)
	if err != nil {
		fmt.Println("err", err)
	}
	slog.SetOutput(src)
	Logs = dialog.NewJSONLogger(dialog.NewSyncWriter(src))
	slog.SetOutput(dialog.NewStdlibAdapter(Logs))
	Logs = dialog.With(Logs, "[go-kit]ts", dialog.DefaultTimestamp, "caller", dialog.DefaultCaller)
}
func GetLogger() dialog.Logger {
	return Logs
}
