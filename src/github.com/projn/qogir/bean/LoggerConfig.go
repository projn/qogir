package bean

import (
	"errors"
	"github.com/magiconair/properties"
	"github.com/phachon/go-logger"
	"log"
)

type LoggerConfig struct {
	LoggerFileDir string `properties:"logger.loggerFileDir"`
	MaxSize       int64  `properties:"logger.maxSize"`
	MaxLine       int64  `properties:"logger.maxLine"`
	DateSlice     string `properties:"logger.dateSlice"`
	Format        string `properties:"logger.format"`
}

var Logger *go_logger.Logger

func CreateLoggerBean(configDir string) error {
	properties := properties.MustLoadFile(configDir+"/config/logger.properties", properties.UTF8)
	if properties == nil {
		log.Printf("Load db properties error.")
		return errors.New("Load db properties error.")
	}

	var loggerConfig LoggerConfig
	err := properties.Decode(loggerConfig)
	if err != nil {
		Logger.Errorf("Load logger properties error, error info(%s).", err.Error())
		return err
	}

	Logger = go_logger.NewLogger()
	fileConfig := &go_logger.FileConfig{

		LevelFileName: map[int]string{
			Logger.LoggerLevel("error"): loggerConfig.LoggerFileDir + "/run-error.log",
			Logger.LoggerLevel("info"):  loggerConfig.LoggerFileDir + "/run.log",
		},
		MaxSize:    loggerConfig.MaxSize,
		MaxLine:    loggerConfig.MaxLine,
		DateSlice:  loggerConfig.DateSlice,
		JsonFormat: false,
		Format:     loggerConfig.Format,
	}

	Logger.Attach("file", go_logger.LOGGER_LEVEL_DEBUG, fileConfig)
	Logger.SetAsync()

	return nil
}

