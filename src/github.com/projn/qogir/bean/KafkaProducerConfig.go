package bean

import (
	"errors"
	"github.com/Shopify/sarama"
	"github.com/magiconair/properties"
	"log"
	"strings"
	"time"
)

type KafkaProducerConfig struct {
	RequiredAcks     sarama.RequiredAcks `properties:"kafka.producer.requiredAcks,default=1"`
	BootstrapServers []string            `properties:"kafka.producer.bootstrapServers"`
	CompressionCodec int                 `properties:"kafka.producer.compressionCodec"`
	RetryMax         int                 `properties:"kafka.producer.retryMax,default=3"`
	ClientId         string              `properties:"kafka.producer.clientId"`
	DialTimeout      time.Duration       `properties:"kafka.producer.dialTimeout,default=30000"`
	WriteTimeout     time.Duration       `properties:"kafka.producer.writeTimeout,default=30000"`
	ReadTimeout      time.Duration       `properties:"kafka.producer.readTimeout,default=30000"`
	MaxMessageBytes  int                 `properties:"kafka.producer.maxMessageBytes,default=1000000"`
}

var syncProducer sarama.SyncProducer
var asyncProducer sarama.AsyncProducer

func CreateKafkaProducerBean(configDir string) error {
	properties := properties.MustLoadFile(configDir+"/config/kafka-produer.properties", properties.UTF8)
	if properties == nil {
		log.Printf("Load kafka properties error.")
		return errors.New("Load kafka properties error.")
	}

	var kafkaProducerConfig KafkaProducerConfig
	err := properties.Decode(kafkaProducerConfig)
	if err != nil {
		Logger.Errorf("Load kafka properties error, error info(%s).", err.Error())
		return err
	}

	if len(kafkaProducerConfig.BootstrapServers) == 0 {
		log.Printf("Invaild kafka address info.")
		return errors.New("Invaild kafka address info.")
	}

	c := &sarama.Config{}

	c.Net.MaxOpenRequests = 5
	c.Net.DialTimeout = kafkaProducerConfig.DialTimeout
	c.Net.ReadTimeout = kafkaProducerConfig.ReadTimeout
	c.Net.WriteTimeout = kafkaProducerConfig.WriteTimeout
	c.Net.SASL.Handshake = true

	c.Metadata.Retry.Max = kafkaProducerConfig.RetryMax
	c.Metadata.Retry.Backoff = 250 * time.Millisecond
	c.Metadata.RefreshFrequency = 10 * time.Minute
	c.Metadata.Full = true

	c.Producer.MaxMessageBytes = kafkaProducerConfig.MaxMessageBytes
	c.Producer.RequiredAcks = kafkaProducerConfig.RequiredAcks
	c.Producer.Timeout = 10 * time.Second
	c.Producer.Partitioner = sarama.NewHashPartitioner
	c.Producer.Retry.Max = kafkaProducerConfig.RetryMax
	c.Producer.Retry.Backoff = 100 * time.Millisecond
	c.Producer.Return.Errors = true
	c.Producer.CompressionLevel = kafkaProducerConfig.CompressionCodec

	c.ClientID = kafkaProducerConfig.ClientId
	c.ChannelBufferSize = 256

	syncProducer, err = sarama.NewSyncProducer(kafkaProducerConfig.BootstrapServers, c)
	if err != nil {
		Logger.Errorf("Create kafka sync producer bean error, error info(%s).", err.Error())
		return err
	}

	asyncProducer, err = sarama.NewAsyncProducer(addrs, c)
	if err != nil {
		Logger.Errorf("Create kafka async producer bean error, error info(%s).", err.Error())
		return err
	}

	return nil
}
