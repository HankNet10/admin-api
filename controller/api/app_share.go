package api

import (
	"myadmin/controller/api/bind"
	"myadmin/model/user"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type AppShareCreateRequest struct {
	ShareName        string `json:"shareName"`
	TerminalType     string `json:"terminalType"` // 设备型号
	BrowserUUID      string `json:"uuid"`         // 浏览器UUID
	BrowserUserAgent string `json:"userAgent"`    // 浏览器UA
}

type AppShareCreatePutRequest struct {
	ID                uint   `json:"id"`
	DeviceUUID        string `json:"uuid"`      // 设备UUID
	DeviceUserAgent   string `json:"userAgent"` // 设备型号
	DeviceInstallTime string `json:"time"`
}

// 落地页点击下载
func AppShareAdd(c *gin.Context) {
	var request AppShareCreateRequest
	c.ShouldBindJSON(&request)
	appID := c.Request.Header.Get("x-appid")
	if appID == "" {
		appID = "1"
	}
	appid, _ := strconv.Atoi(appID)
	if request.ShareName == "" || request.BrowserUUID == "" || request.TerminalType == "" {
		c.JSON(http.StatusBadRequest, bind.ErrorMessage{Message: "参数错误"})
		return
	}
	//默认官网为 home
	shareName := request.ShareName
	appShare := &user.AppShare{}
	appShare.AppListID = uint(appid)
	appShare.ShareName = shareName
	appShare.BrowserInstallIP = c.ClientIP()
	appShare.BrowserUserAgent = c.Request.UserAgent()
	appShare.Status = 0
	appShare.BrowserUUID = request.BrowserUUID
	appList := user.App.SelectByID(appShare.AppListID)
	if appList == nil || appList.ID == 0 {
		c.JSON(http.StatusBadRequest, bind.ErrorMessage{Message: "未知错误"})
		return
	}

	channel := bindChannel{}
	if request.TerminalType == "android" {
		appShare.Save()
		channel.ID = appShare.ID
		channel.URL = appList.AndroidDownloadUrl
	} else if request.TerminalType == "ios" {
		channel.URL = appList.IosDownloadUrl
	} else if request.TerminalType == "wap" {
		channel.URL = appList.WapUrl
	} else {
		c.JSON(http.StatusBadRequest, bind.ErrorMessage{Message: "参数错误"})
		return
	}
	c.JSON(http.StatusOK, channel)
}

// apk启动调用一次
func AppUpdate(c *gin.Context) {
	var request AppShareCreatePutRequest
	c.ShouldBindJSON(&request)
	appID := c.Request.Header.Get("x-appid")
	if appID == "" {
		appID = "1"
	}
	appid, _ := strconv.Atoi(appID)
	if request.DeviceUUID == "" || request.DeviceInstallTime == "" || request.DeviceUserAgent == "" {
		c.JSON(http.StatusBadRequest, bind.ErrorMessage{Message: "参数错误"})
		return
	}
	appShare := &user.AppShare{}
	if request.ID > 0 {
		appShare = user.Appshare.SelectByID(request.ID)
	}
	//2 ID匹配 3 IP和设备型号匹配 1 apk包安装
	if appShare == nil || appShare.ID == 0 {
		appShare = user.Appshare.SelectByIP(appid, c.ClientIP())
		if appShare == nil || appShare.ID == 0 {
			appShare.Status = 1
			appShare.AppListID = uint(appid)
		} else {
			appShare.Status = 3
		}
	} else {
		appShare.Status = 2
	}
	appShare.DeviceUUID = request.DeviceUUID
	times, err := time.Parse("2006-01-02 15:04:05", request.DeviceInstallTime)
	if err != nil {
		times = time.Now()
	}
	appShare.DeviceInstallTime = times.Format("2006-01-02 15:04:05")
	appShare.DeviceUserAgent = request.DeviceUserAgent
	appShare.DeviceInstallIP = c.ClientIP()
	appShare.Save()
	c.JSON(http.StatusOK, bind.ErrorMessage{Message: "完成"})
}
