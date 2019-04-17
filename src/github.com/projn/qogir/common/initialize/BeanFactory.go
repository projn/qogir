package initialize

import (
	"../../bean"
	"github.com/kataras/iris/core/errors"
	"github.com/magiconair/properties"
	"log"
)

type RunTimeProperties struct {
	AppName                string `properties:"system.application.name"`
	BeanSwitchRedisSingle  bool   `properties:"system.bean.switch.redis.single"`
	BeanSwitchRedisCluster bool   `properties:"system.bean.switch.redis.cluster"`
	BeanSwitchMqConsumer   bool   `properties:"system.bean.switch.mq.consumer"`
	BeanSwitchMqProducer   bool   `properties:"system.bean.switch.mq.producer"`
	BeanSwitchWebsocket    bool   `properties:"system.bean.switch.websocket"`
	BeanSwitchDataSource   bool   `properties:"system.bean.switch.dataSource"`
	BeanSwitchGrpc         bool   `properties:"system.bean.switch.grpc"`
	BeanSwitchMongodb      bool   `properties:"system.bean.switch.mongodb"`
	I18nDir                string `properties:"system.i18n.dir"`
	TokenSecretKey         string `properties:"system.tokenSecretKey"`
	ApiAccessRoleSendMsg   string `properties:"system.api.access.role.sendMsg"`
	ApiAccessRoleActuator  string `properties:"system.api.access.role.actuator"`
	WsMsgIds               string `properties:"system.ws.msg.ids"`
}

func (runTimeProperties *RunTimeProperties) getProperties(configDir string) (*RunTimeProperties, error) {
	properties := properties.MustLoadFile(configDir+"/application.properties", properties.UTF8)
	if properties == nil {
		log.Printf("Load application properties error.")
		return nil, errors.New("Load application properties error.")
	}

	err := properties.Decode(runTimeProperties)
	if err != nil {
		log.Printf("Load application properties error, error info(%s).", err.Error())
		return nil, err
	}

	return runTimeProperties, nil
}

var GlobalProperties RunTimeProperties

func InitServiceBean(configDir string) {
	if configDir == "" {
		log.Printf("Invaild config dir.")
		return
	}

	var runTimeProperties RunTimeProperties
	GlobalProperties, err := runTimeProperties.getProperties(configDir)
	if err != nil {
		return
	}

	err = bean.CreateLoggerBean(configDir)
	if err != nil {
		return
	}

	if GlobalProperties.BeanSwitchDataSource {
		err = bean.CreateDbBean(configDir)
		if err != nil {
			return
		}
	}

	if GlobalProperties.BeanSwitchDataSource {
		err = bean.CreateDbBean(configDir)
		if err != nil {
			return
		}
	}

	if GlobalProperties.BeanSwitchDataSource {
		err = bean.CreateDbBean(configDir)
		if err != nil {
			return
		}
	}

	if GlobalProperties.BeanSwitchDataSource {
		err = bean.CreateDbBean(configDir)
		if err != nil {
			return
		}
	}

}
