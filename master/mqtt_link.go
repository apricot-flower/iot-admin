package master

import (
	"encoding/json"
	"iot-admin/logs"
	"iot-admin/mqtt_protocol"
	"iot-admin/source"
	"sync"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var _ Linker = (*MqttLinker)(nil)

type MqttLinker struct {
	source.Master
	lock   sync.Mutex
	client mqtt.Client
}

func (m *MqttLinker) LogOut() {
	logout := &mqtt_protocol.LoginOrLogout{FrameType: "logout", Ts: time.Now().UnixNano(), ClientId: m.ClientId, Username: m.Username, Password: m.Password, HeartBeat: m.HeartBeat}
	m.Send(m.LoginTopic, 1, false, logout)
	logs.Logger.Debugf("%s logouting", m.MasterName)
}

func (m *MqttLinker) Login() {
	login := &mqtt_protocol.LoginOrLogout{FrameType: "login", Ts: time.Now().UnixNano(), ClientId: m.ClientId, Username: m.Username, Password: m.Password, HeartBeat: m.HeartBeat}
	m.Send(m.LoginTopic, 1, false, login)
	logs.Logger.Debugf("%s logining", m.MasterName)
}

func (m *MqttLinker) Send(topic string, qos byte, retained bool, payload interface{}) {
	m.lock.Lock()
	defer m.lock.Unlock()
	jsonData, err := json.Marshal(payload)
	if err != nil {
		logs.Logger.Errorf("data to json error: %v", err)
		return
	}
	m.client.Publish(topic, qos, retained, jsonData)
}

func (m *MqttLinker) connect() {
	m.lock.Lock()
	defer m.lock.Unlock()
	opts := mqtt.NewClientOptions()
	opts.AddBroker(m.Host)
	opts.SetClientID(m.ClientId)
	opts.SetUsername(m.Username)
	opts.SetPassword(m.Password)
	opts.SetDefaultPublishHandler(m.messagePubHandler)
	opts.SetAutoReconnect(true)
	opts.SetCleanSession(true)
	opts.SetConnectTimeout(10 * time.Second)
	opts.SetKeepAlive(60 * time.Second)
	opts.OnConnect = m.connectHandler
	opts.OnConnectionLost = m.connectionLostHandler
	m.client = mqtt.NewClient(opts)
	if token := m.client.Connect(); token.Wait() && token.Error() != nil {
		logs.Logger.Errorf("%s connect err: %s", m.MasterName, token.Error())
	}
	token := m.client.Subscribe(m.LoginResponse, 1, nil)
	token.Wait()
	if token.Error() != nil {
		logs.Logger.Debugf("%s subscribe init_topic:<%s> err: %s", m.MasterName, m.LoginResponse, token.Error())
	}
}

// 连接丢失
func (m *MqttLinker) connectionLostHandler(_ mqtt.Client, err error) {
	logs.Logger.Debugf("%s connect lost: %s", m.MasterName, err.Error())
}

// 连接成功
func (m *MqttLinker) connectHandler(_ mqtt.Client) {
	logs.Logger.Debugf("%s connect success", m.MasterName)
	m.Login()
}

// 收到消息
func (m *MqttLinker) messagePubHandler(_ mqtt.Client, msg mqtt.Message) {
	logs.Logger.Debugf("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())
}
