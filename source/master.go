package source

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
}
