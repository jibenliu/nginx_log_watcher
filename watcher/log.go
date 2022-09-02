package watcher

import (
	"fmt"
	"github.com/fatih/color"
	_ "github.com/fatih/color"
	"log"
	"os"
)

var outPutFile bool

func initLog(file string) {
	logFile, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0766)
	if err != nil {
		Error("创建日志文件失败，错误原因： " + err.Error())
		panic(err)
	}
	log.SetOutput(logFile) // 将文件设置为log输出的文件
	log.SetPrefix("[nginx_log_watcher]")
	log.SetFlags(log.LstdFlags | log.Lshortfile | log.LUTC)
	outPutFile = true
	return
}
func log2File(output, level string) {
	if outPutFile {
		log.Println(fmt.Sprintf("[%s] %s", level, output))
	}
}
func Debug(output string) {
	color.Cyan(output)
	log2File(output, "DEBUG")
}

func Warn(output string) {
	color.Magenta(output)
	log2File(output, "Warn")
}

func Error(output string) {
	color.Red(output)
	log2File(output, "Error")
	os.Exit(1)
}
