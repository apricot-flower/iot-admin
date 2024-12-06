package master

import (
	"iot-admin/logs"
	"iot-admin/source"
	"strings"
)

const (
	TcpType  string = "tcp"
	MqttType string = "mqtt"
)

var Clients = make(map[string]Linker)

type Linker interface {
	connect() //连接
	Login()   //登录
	LogOut()  //退出登录
}

func Connects() {
	logs.Logger.Debug("start initializing the main station connection……")
	var masterType string
	masters, err := source.SelectMasters()
	if err != nil {
		logs.Logger.Errorf("select master config err: %v", err)
	}
	for _, master := range masters {
		masterType = strings.ToLower(master.MasterType)
		var client Linker
		switch masterType {
		case TcpType:
			//create tcp client
			client = &TcpLinker{Master: master}
			client.connect()
		case MqttType:
			//create mqtt client
			client = &MqttLinker{Master: master}
			client.connect()
		default:
			logs.Logger.Errorf("%s unknown master type: %s", master.MasterName, masterType)
			continue
		}
		Clients[master.MasterName] = client
	}

}
