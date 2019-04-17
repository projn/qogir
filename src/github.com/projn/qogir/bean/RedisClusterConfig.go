package bean

import (
	"errors"
	"github.com/go-redis/redis"
	"github.com/magiconair/properties"
	"log"
	"time"
)

type RedisClusterConfig struct {
	Addrs        []string        `properties:"redis.addrs"`
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

var RedisClusterClient *redis.ClusterClient

func CreateRedisClusterClientBean(configDir string) error {
	properties := properties.MustLoadFile(configDir+"/config/redis-cluster.properties", properties.UTF8)
	if properties == nil {
		log.Printf("Load redis properties error.")
		return errors.New("Load redis properties error.")
	}

	var redisClusterConfig RedisClusterConfig
	err := properties.Decode(redisClusterConfig)
	if err != nil {
		Logger.Errorf("Load redis properties error, error info(%s)", err.Error())
		return err
	}

	if len(redisClusterConfig.Addrs)==0 {
		log.Printf("Invaild redis address info.")
		return errors.New("Invaild redis address info.")
	}

	RedisClusterClient = redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:        redisClusterConfig.Addrs,
		Password:     redisClusterConfig.Password,
		MaxRetries:   redisClusterConfig.MaxRetries,
		DialTimeout:  redisClusterConfig.DialTimeout,
		ReadTimeout:  redisClusterConfig.ReadTimeout,
		WriteTimeout: redisClusterConfig.WriteTimeout,
		PoolSize:     redisClusterConfig.PoolSize,
		MinIdleConns: redisClusterConfig.MinIdleConns,
		MaxConnAge:   redisClusterConfig.MaxConnAge,
		PoolTimeout:  redisClusterConfig.PoolTimeout,
		IdleTimeout:  redisClusterConfig.IdleTimeout,
	})

	return nil
}
