package main

import (
	"encoding/json"
	"iot-admin/bus"
	"iot-admin/logs"
	"iot-admin/master"
	"iot-admin/source"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	//初始化主站连接
	master.Connects()
	bus.OpenMsgBus()
	web()
}

func web() {
	r := gin.Default()
	r.LoadHTMLGlob("templates/*.html")
	//首页
	r.GET("/", home)
	r.GET("/maters", masters)
	r.POST("/insert_master", insertMaster)
	r.POST("/delete_master", deleteMaster)
	serverConfig := source.SelectServerConfig()
	err := r.Run(":" + strconv.Itoa(serverConfig.WebConfig))
	if err != nil {
		logs.Logger.Errorf("web server start error:%v", err)
	}
}

// 获取首页
func home(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{
		"title": "欢迎来到首页"})
}

// 获取全部主站配置
func masters(c *gin.Context) {
	masterConfigs, err := source.SelectMasters()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "not found master config"})
	}
	c.JSON(http.StatusOK, masterConfigs)
}

// 保存主站配置
func insertMaster(c *gin.Context) {
	data, _ := c.GetRawData()
	var body map[string]string
	_ = json.Unmarshal(data, &body)
	result, err := strconv.ParseInt(body["heart_beat"], 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"success": false})
	} else {
		masterConfig := &source.Master{Id: time.Now().UnixNano(), MasterType: body["type"], MasterName: body["name"], MasterMark: body["description"], Host: body["address"], ClientId: body["client_id"], Username: body["account"], Password: body["password"], LoginTopic: body["login_topic"], LoginResponse: body["reply_topic"], HeartBeat: result}
		source.InsertMaster(masterConfig)
		c.JSON(http.StatusOK, gin.H{"success": true})
	}
}

// 删除主站配置
func deleteMaster(c *gin.Context) {
	data, _ := c.GetRawData()
	var body map[string][]string
	_ = json.Unmarshal(data, &body)
	source.DeleteMasterByIds(body["ids"])
	c.JSON(http.StatusOK, gin.H{"success": true})
}
