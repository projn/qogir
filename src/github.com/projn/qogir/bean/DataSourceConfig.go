package bean

import (
	"errors"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/magiconair/properties"
	"log"
)

type DataSourceConfig struct {
	SourceType   string `properties:"dataSource.sourceType"`
	User         string `properties:"dataSource.user"`
	Password     string `properties:"dataSource.password"`
	DbName       string `properties:"dataSource.dbName"`
	Extend       string `properties:"dataSource.extend"`
	MaxIdleConns int    `properties:"dataSource.maxIdleConns"`
	MaxOpenConns int    `properties:"dataSource.maxOpenConns"`
}

var Db *gorm.DB

func CreateDbBean(configDir string) error {
	properties := properties.MustLoadFile(configDir+"/config/db.properties", properties.UTF8)
	if properties == nil {
		log.Printf("Load db properties error.")
		return errors.New("Load db properties error.")
	}

	var dataSourceConfig DataSourceConfig
	err := properties.Decode(dataSourceConfig)
	if err != nil {
		Logger.Errorf("Load db properties error, error info(%s).", err.Error())
		return err
	}

	Db, err = gorm.Open(dataSourceConfig.SourceType, dataSourceConfig.User+":"+dataSourceConfig.Password+"@/"+dataSourceConfig.DbName+"?"+dataSourceConfig.Extend)
	defer Db.Close()

	if err != nil {
		Logger.Errorf("Create db bean error, error info(%s).", err.Error())
		return err
	}

	Db.DB().SetMaxIdleConns(dataSourceConfig.MaxIdleConns)
	Db.DB().SetMaxOpenConns(dataSourceConfig.MaxOpenConns)

	return nil
}
