package mqtt_protocol

// LoginOrLogout 登录-登出结构体
type LoginOrLogout struct {
	FrameType string `json:"type"`
	Ts        int64  `json:"ts"`
	ClientId  string `json:"clientId"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	HeartBeat int64  `json:"heartBeat"` //心跳周期
}
