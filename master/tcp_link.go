package master

import (
	"bufio"
	"dlt698"
	"iot-admin/logs"
	"iot-admin/source"
	"net"
	"sync"
	"time"
)

var _ Linker = (*TcpLinker)(nil)

type TcpLinker struct {
	source.Master
	lock     sync.Mutex
	conn     net.Conn
	reader   *bufio.Reader
	pidIndex uint8
}

func (t *TcpLinker) pid() uint8 {
	defer func() {
		t.pidIndex++
		if t.pidIndex >= 63 {
			t.pidIndex = 0
		}
	}()
	return t.pidIndex
}

func (t *TcpLinker) LogOut() {
	//TODO implement me
	panic("implement me")
}

func (t *TcpLinker) Login() {
	login, err := dlt698.CreateLogin(t.ClientId, 60, t.pid(), 0x00)
	if err != nil {
		logs.Logger.Errorf("%s create login failed: %v", t.MasterName, err)
		return
	}
	t.Send(login)
}

func (t *TcpLinker) Send(data []byte) {
	t.lock.Lock()
	defer t.lock.Unlock()
	_, err := t.conn.Write(data)
	if err != nil {
		logs.Logger.Errorf("%s write frame failed: %v", t.MasterName, err)
		return
	}

}

func (t *TcpLinker) connect() {
	go func() {
		for {
			err := t.link()
			if err != nil {
				time.Sleep(10 * time.Second)
				continue
			} else {
				t.Login()
			}
			break
		}
	}()
}

func (t *TcpLinker) link() error {
	t.lock.Lock()
	defer t.lock.Unlock()
	conn, err := net.Dial("tcp", t.Host)
	if err != nil {
		logs.Logger.Errorf("create tcp_link:<%s> error:%v", t.MasterName, err)
		return err
	}
	t.conn = conn
	t.reader = bufio.NewReader(t.conn)
	logs.Logger.Debugf("%s connect success", t.MasterName)
	t.pidIndex = 0
	return nil
}

func (t *TcpLinker) reConnect() error {
	if t.conn != nil {
		_ = t.conn.Close()
	}
	return t.link()
}
