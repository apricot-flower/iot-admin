package source

// Master 主站配置
type Master struct {
	Id            int64  `json:"id" gorm:"size:64;PRIMARY_KEY;NOT NULL"`
	MasterType    string `json:"master_type" gorm:"size:10;NOT NULL"`
	MasterName    string `json:"master_name" gorm:"size:255;NOT NULL"`
	MasterMark    string `json:"master_mark" gorm:"size:255"`
	Host          string `json:"host" gorm:"size:255;NOT NULL"`
	ClientId      string `json:"client_id" gorm:"size:100;NOT NULL"`
	Username      string `json:"username" gorm:"size:255;NOT NULL"`
	Password      string `json:"password" gorm:"size:255;NOT NULL"`
	LoginTopic    string `json:"login_topic" gorm:"size:255"`
	LoginResponse string `json:"login_response" gorm:"size:255"`
	HeartBeat     int64  `json:"heart_beat" gorm:"size:255"` //心跳周期，单位：秒
}

// BusConfig 子APP配置
type BusConfig struct {
	Id            int64  `json:"id" gorm:"size:64;PRIMARY_KEY;NOT NULL"`  //id
	AppName       string `json:"app_name" gorm:"size:255;NOT NULL"`       //app名称
	AppKey        string `json:"app_key" gorm:"size:255;NOT NULL"`        //app编码
	AppMark       string `json:"app_mark" gorm:"size:255;NOT NULL"`       //app说明
	PhysicalModel string `json:"physical_model" gorm:"size:255;NOT NULL"` //对应的物模型
}

// ServerConfig 通用配置
type ServerConfig struct {
	WebConfig   int    `json:"web_config" gorm:"size:255;default:8080"` //web端口
	BusPort     int    `json:"bus_port" gorm:"size:255;default:8314"`   //消息总线端口
	IotFlag     string `json:"iot_flag" gorm:"size:255;NOT NULL"`       //边缘计算框架的唯一编码
	NativePlace string `json:"native_place" gorm:"size:1000"`           //安装位置

}
