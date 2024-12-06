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
	//验证主站表是否存在
	masterFlag := d.db.Migrator().HasTable(&Master{})
	if !masterFlag {
		err := d.db.Migrator().CreateTable(&Master{})
		if err != nil {
			logs.Logger.Errorf("create sqlite table master err:%v", err)
		}
	}
}

func SelectMasters() ([]Master, error) {
	var masters []Master
	result := IotDbClient.db.Find(&masters)
	if result.Error != nil {
		return nil, result.Error
	}
	return masters, nil
}

func InsertMaster(master *Master) {
	IotDbClient.db.Save(master)
}

func DeleteMasterByIds(ids []string) {
	IotDbClient.db.Delete(&Master{}, ids)
}
