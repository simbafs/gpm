package log

import (
	"errors"
	"io/fs"
	"os"
	"path"

	"github.com/gin-gonic/gin"
	"github.com/op/go-logging"
	Config "github.com/simba-fs/gpm/config"
)

var (
	ErrLogIsNotDir = errors.New("log is not a directory")
)
var loggers = map[string]gin.HandlerFunc{}
var formatter = logging.MustStringFormatter(
	`%{color}%{module} %{longfunc} ▶  %{level:.4s} %{color:reset} %{message}`,
)
var formatterWithColor = logging.MustStringFormatter(
	`%{module} %{longfunc} ▶  %{level:.4s} %{message}`,
)

var levels = map[string]logging.Level{
	"critical": logging.CRITICAL,
	"error":    logging.ERROR,
	"warning":  logging.WARNING,
	"notice":   logging.NOTICE,
	"info":     logging.INFO,
	"debug":    logging.DEBUG,
	"0":        logging.CRITICAL,
	"1":        logging.ERROR,
	"2":        logging.WARNING,
	"3":        logging.NOTICE,
	"4":        logging.INFO,
	"5":        logging.DEBUG,
}

func init() {
	logging.SetFormatter(formatter)
}

func Init(c *Config.Config) {
	if c.Log != "" {
		w, err := os.OpenFile(path.Join(c.Log, "gpm.log"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			panic(err)
		}
		b := logging.NewLogBackend(w, "", 0)
		logging.SetBackend(b)
		logging.SetFormatter(formatterWithColor)
	}
}

func SetLevel(level string) {
	logging.SetLevel(levels[level], "")
}

func NewLog(module string) *logging.Logger {
	return logging.MustGetLogger(module)
}

func NewLogMiddleware(config *Config.Config) gin.HandlerFunc {
	f, err := os.Open(config.Log)
	if os.IsNotExist(err) {
		os.Mkdir(config.Log, fs.ModeDir|0744)
	} else if s, _ := f.Stat(); !s.IsDir() {
		panic(ErrLogIsNotDir)
	}

	return func(c *gin.Context) {
		// get proxy route
		name, ok := "", false
		for k, v := range config.Host {
			if v.From == c.Request.Host {
				name, ok = k, true
				break
			}
		}
		if !ok {
			name = "unknown"
			return
		}

		logger, ok := loggers[name]
		if !ok {
			writer, err := os.OpenFile(path.Join(config.Log, "host-" + name + ".log"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				panic(err)
			}

			logger = gin.LoggerWithWriter(writer)
			loggers[name] = logger
		}

		logger(c)
	}
}
