package bean

import (
	"errors"
	"github.com/magiconair/properties"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
	"log"
	"time"
)

type MongoDbConfig struct {
	Hosts                  []string        `properties:"mongodb.hosts"`
	Username               string        `properties:"mongodb.username"`
	Password               string        `properties:"mongodb.password"`
	ConnectTimeout         time.Duration `properties:"mongodb.connectTimeout"`
	HeartbeatInterval      time.Duration `properties:"mongodb.heartbeatInterval"`
	LocalThreshold         time.Duration `properties:"mongodb.localThreshold"`
	MaxConnIdleTime        time.Duration `properties:"mongodb.maxConnIdleTime"`
	MaxPoolSize            uint16        `properties:"mongodb.maxPoolSize"`
	ReadConcern            string        `properties:"mongodb.readConcern"`
	ServerSelectionTimeout time.Duration `properties:"mongodb.serverSelectionTimeout"`
	SocketTimeout          time.Duration `properties:"mongodb.socketTimeout"`
	WriteConcern           string        `properties:"mongodb.writeConcern,default=majority"`
}

var mongoDbClient *mongo.Client

func CreateMongoDbClientBean(configDir string) error {
	properties := properties.MustLoadFile(configDir+"/config/mongodb.properties", properties.UTF8)
	if properties == nil {
		log.Printf("Load mongodb properties error.")
		return errors.New("Load mongodb properties error.")
	}

	var mongoDbConfig MongoDbConfig
	err := properties.Decode(mongoDbConfig)
	if err != nil {
		Logger.Errorf("Load mongodb properties error, error info(%s)", err.Error())
		return err
	}

	if len(mongoDbConfig.Hosts) == 0 {
		log.Printf("Invaild mongodb address info.")
		return errors.New("Invaild mongodb address info.")
	}

	mongoDbClient, err = mongo.NewClient(&options.ClientOptions{
		Hosts: mongoDbConfig.Hosts,
		Auth: &options.Credential{
			Username: mongoDbConfig.Username,
			Password: mongoDbConfig.Password,
		},
		ConnectTimeout:         &mongoDbConfig.ConnectTimeout,
		HeartbeatInterval:      &mongoDbConfig.HeartbeatInterval,
		LocalThreshold:         &mongoDbConfig.LocalThreshold,
		MaxConnIdleTime:        &mongoDbConfig.MaxConnIdleTime,
		MaxPoolSize:            &mongoDbConfig.MaxPoolSize,
		ReadConcern:            readconcern.New(readconcern.Level(mongoDbConfig.ReadConcern)),
		ServerSelectionTimeout: &mongoDbConfig.ServerSelectionTimeout,
		SocketTimeout:          &mongoDbConfig.SocketTimeout,
		WriteConcern:           writeconcern.New(writeconcern.WTagSet(mongoDbConfig.WriteConcern)),
	})

	if err != nil {
		Logger.Errorf("Create mongodb client bean error, error info(%s)", err.Error())
		return err
	}

	return nil
}
