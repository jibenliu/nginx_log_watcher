package watcher

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
	"strings"
	"sync"
)

var (
	once sync.Once
	app  *AppConfig
)

func init() {
	once.Do(func() {
		getApp()
		InitCache()
		TickerFlushCache()
	})
}

type AppConfig struct {
	debugMode    bool
	debugLogFile string
	reportRate   reportRate
}

type reportRate struct {
	minutePeriodErrorCount int
	continuousErrorCount   int
}

func getApp() {
	viper.SetConfigName("watcher") // name of config file (without extension)
	viper.SetConfigType("yaml")    // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath(".")       // path to look for the config file in
	viper.AddConfigPath(".")       // optionally look for config in the working directory
	err := viper.ReadInConfig()    // Find and read the config file
	if err != nil {                // Handle errors reading the config file
		if err != os.ErrNotExist {
			Error("read yaml error: " + err.Error())
			panic("配置文件读取异常")
		} else {
			Error("read yaml error: " + err.Error())
			app = &AppConfig{
				debugMode:    false,
				debugLogFile: "/var/log/nginx_log_reporter.log",
				reportRate: reportRate{
					minutePeriodErrorCount: 30,
					continuousErrorCount:   5,
				},
			}
		}
	}
	app = &AppConfig{
		debugMode:    viper.GetBool("DEBUG_MODE"),
		debugLogFile: viper.GetString("DEBUG_MODE_FILE"),
		reportRate: reportRate{
			minutePeriodErrorCount: viper.GetInt("REPORT_RATE.MINUTE_PERIOD_ERROR_COUNT"),
			continuousErrorCount:   viper.GetInt("REPORT_RATE.CONTINUOUS_ERROR_COUNT"),
		},
	}
	if app.debugMode {
		_, err := os.Stat(app.debugLogFile)
		if err == os.ErrNotExist {
			dir := app.debugLogFile[0:strings.LastIndex(app.debugLogFile, "/")]
			err := os.MkdirAll(dir, 0644)
			if err != nil {
				Error("创建日志目录失败，错误原因： " + err.Error())
				panic("创建日志目录失败")
			}
			initLog(app.debugLogFile)
		}
	}

	Debug(fmt.Sprintf("读取到配置或者告警默认值为：MINUTE_PERIOD_ERROR_COUNT %d CONTINUOUS_ERROR_COUNT %d", app.reportRate.minutePeriodErrorCount, app.reportRate.continuousErrorCount))
}
