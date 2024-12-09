package source

import (
	"iot-admin/logs"
	"sync"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var IotDbClient = &DbClient{}

type DbClient struct {
	lock sync.RWMutex
	db   *gorm.DB
}

func init() {
	db, err := gorm.Open(sqlite.Open("./service.db"), &gorm.Config{})
	if err != nil {
		logs.Logger.Errorf("create sqlite link err:%v", err)
		return
	}
	logs.Logger.Debugf("create sqlite link success")
	IotDbClient.db = db
	IotDbClient.validate()
}

func (d *DbClient) validate() {
	d.lock.Lock()
	defer d.lock.Unlock()
	//验证主站表是否存在
	masterFlag := d.db.Migrator().HasTable(&Master{})
	if !masterFlag {
		err := d.db.Migrator().CreateTable(&Master{})
		if err != nil {
			logs.Logger.Errorf("create sqlite table master err:%v", err)
		}
	}
	busConfigFlag := d.db.Migrator().HasTable(&BusConfig{})
	if !busConfigFlag {
		err := d.db.Migrator().CreateTable(&BusConfig{})
		if err != nil {
			logs.Logger.Errorf("create sqlite table BusConfig err:%v", err)
		}
	}
	serverConfigFlag := d.db.Migrator().HasTable(&ServerConfig{})
	if !serverConfigFlag {
		err := d.db.Migrator().CreateTable(&ServerConfig{})
		if err != nil {
			logs.Logger.Errorf("create sqlite table ServerConfig err:%v", err)
		}
	}
}

func SelectMasters() ([]Master, error) {
	IotDbClient.lock.Lock()
	defer IotDbClient.lock.Unlock()
	var masters []Master
	result := IotDbClient.db.Find(&masters)
	if result.Error != nil {
		return nil, result.Error
	}
	return masters, nil
}

func InsertMaster(master *Master) {
	IotDbClient.lock.Lock()
	defer IotDbClient.lock.Unlock()
	IotDbClient.db.Save(master)
}

func DeleteMasterByIds(ids []string) {
	IotDbClient.lock.Lock()
	defer IotDbClient.lock.Unlock()
	IotDbClient.db.Delete(&Master{}, ids)
}

// SelectServerConfig 查询通用配置
func SelectServerConfig() *ServerConfig {
	IotDbClient.lock.Lock()
	defer IotDbClient.lock.Unlock()
	serverConfig := &ServerConfig{}
	IotDbClient.db.Take(serverConfig)
	return serverConfig
}
