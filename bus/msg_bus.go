package bus

import (
	"iot-admin/source"
	"net"
)

var AppMessageBus *MsgBus

// MsgBus 消息总线
type MsgBus struct {
	Port   int //端口
	listen net.Listener
}

func (m *MsgBus) open() {
}

func OpenMsgBus() {
	AppMessageBus = new(MsgBus)
	serverConfig := source.SelectServerConfig()
	AppMessageBus.Port = serverConfig.BusPort
	AppMessageBus.open()
}
