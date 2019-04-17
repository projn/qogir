package bean

import (
	"errors"
	"github.com/go-redis/redis"
	"github.com/magiconair/properties"
	"log"
	"time"
)

type RedisConfig struct {
	Addr         string        `properties:"redis.addr"`
	Password     string        `properties:"redis.password"`
	MaxRetries   int           `properties:"redis.maxRetries"`
	DialTimeout  time.Duration `properties:"redis.dialTimeout"`
	ReadTimeout  time.Duration `properties:"redis.readTimeout"`
	WriteTimeout time.Duration `properties:"redis.writeTimeout"`
	PoolSize     int           `properties:"redis.poolSize"`
	MinIdleConns int           `properties:"redis.minIdleConns"`
	MaxConnAge   time.Duration `properties:"redis.maxConnAge"`
	PoolTimeout  time.Duration `properties:"redis.poolTimeout"`
	IdleTimeout  time.Duration `properties:"redis.idleTimeout"`
}

var RedisClient *redis.Client

func CreateRedisClientBean(configDir string) error {
	properties := properties.MustLoadFile(configDir+"/config/redis.properties", properties.UTF8)
	if properties == nil {
		log.Printf("Load redis properties error")
		return errors.New("Load redis properties error")
	}

	var redisConfig RedisConfig
	err := properties.Decode(redisConfig)
	if err != nil {
		Logger.Errorf("Load redis properties error, error info(%s)", err.Error())
		return err
	}

	RedisClient = redis.NewClient(&redis.Options{
		Addr:         redisConfig.Addr,
		Password:     redisConfig.Password,
		MaxRetries:   redisConfig.MaxRetries,
		DialTimeout:  redisConfig.DialTimeout,
		ReadTimeout:  redisConfig.ReadTimeout,
		WriteTimeout: redisConfig.WriteTimeout,
		PoolSize:     redisConfig.PoolSize,
		MinIdleConns: redisConfig.MinIdleConns,
		MaxConnAge:   redisConfig.MaxConnAge,
		PoolTimeout:  redisConfig.PoolTimeout,
		IdleTimeout:  redisConfig.IdleTimeout,
	})

	return nil
}
