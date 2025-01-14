package api

import (
	"myadmin/model"
	"myadmin/model/user"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
)

var UserChatList user.UserChatList

// 发送消息
func UserChatSend(c *gin.Context) {
	request := user.UserMessagePostRequest{}
	c.ShouldBindJSON(&request)
	appID := c.Request.Header.Get("x-appid")
	userID := c.MustGet("UserID").(string)
	uID, _ := strconv.Atoi(userID)
	request.SendUserID = uint(uID)
	if appID == "" {
		appID = "1"
	}

	if request.UserId == 0 {
		request.UserId = 7310 //目前就蘑菇能接收 蓝莓还不支持没有官方
	}
	if request.Text == "" && request.ImgPath == "" {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"data":    nil,
			"message": "必须要输入文字或者图片",
		})
		return
	}
	if request.SendUserID == 0 || request.UserId == 0 {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"data":    nil,
			"message": "必须指定用户id",
		})
		return
	}
	// 判断之前是否有聊天关系
	createChat := user.UserChatList{}
	query := model.DataBase.Model(createChat)
	query.Where("user_id = ? And send_id = ?", request.UserId, request.SendUserID)
	query.First(&createChat)
	if createChat.ID == 0 { //创建新的聊天关系
		createChat := user.UserChatList{
			UserID: request.UserId,
			SendID: request.SendUserID,
			Number: 1,
			AppID:  request.AppID,
		}
		result := model.DataBase.Save(&createChat)
		createSendChat := user.UserChatList{ //发送者记录
			UserID: request.SendUserID,
			SendID: request.UserId,
			Number: 0,
			AppID:  request.AppID,
		}
		result2 := model.DataBase.Save(&createSendChat)
		if result.Error != nil || result2.Error != nil {
			c.JSON(http.StatusOK, gin.H{
				"code":    400,
				"data":    nil,
				"message": "创建聊天关系失败，请提交反馈",
			})
			return
		}
	}
	if createChat.ID != 0 { //接受者加1个未读
		model.DataBase.Model(&createChat).Update("number", createChat.Number+1)
	}

	//创建消息
	createContent := user.UserContentList{
		SendID:  request.SendUserID,
		UserID:  request.UserId,
		Text:    request.Text,
		ImgPath: request.ImgPath,
	}
	result2 := model.DataBase.Save(&createContent)
	if result2.Error != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"data":    nil,
			"message": "消息发送失败",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"data":    nil,
		"message": "消息发送成功",
	})
}

// 获取未读消息 和最后一条消息
func UserUnredGet(c *gin.Context) {
	userID := c.MustGet("UserID").(string)
	chatUnRed := user.UserChatList{}

	model.DataBase.Where("user_id = ?", userID).Order("number desc").First(&chatUnRed)
	if chatUnRed.ID != 0 {
		chatLast := user.UserContentList{}
		model.DataBase.Where("user_id = ?", userID).Last(&chatLast)
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"data": gin.H{
				"unReadcount": chatUnRed.Number,
				"lasttext":    chatLast.Text,
				"lasttime":    chatLast.CreatedAt.Unix(),
			},
			"message": "成功",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"unReadcount": 0,
			"lasttext":    "",
			"lasttime":    0,
		},
		"message": "成功",
	})
}

// 获取未读消息 和最后一条消息
func UserUnredGet1(c *gin.Context) {
	userID := c.MustGet("UserID").(string)
	sendID := c.Query("id")
	chatUnRed := user.UserChatList{}

	model.DataBase.Where("user_id = ? and send_id = ?", userID, sendID).First(&chatUnRed)
	if chatUnRed.ID != 0 {
		chatLast := user.UserContentList{}
		model.DataBase.Where("user_id = ? and send_id = ?", userID, sendID).Last(&chatLast)
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"data": gin.H{
				"unReadcount": chatUnRed.Number,
				"lasttext":    chatLast.Text,
				"lasttime":    chatLast.CreatedAt.Unix(),
			},
			"message": "成功",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"unReadcount": 0,
			"lasttext":    "",
			"lasttime":    0,
		},
		"message": "成功",
	})
}

// 获取消息聊天详情列表
func UserChatListList(c *gin.Context) {
	userID := c.MustGet("UserID").(string)
	list, total := UserChatList.UserList(userID)
	for i := 0; i <= len(list)-1; i++ {
		if len(list[i].SendUser.Avatar) > 5 {
			list[i].SendUser.Avatar = os.Getenv("ALI_OSS_DOMAIN") + "/" + list[i].SendUser.Avatar
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"list":  list,
			"total": total,
		},
		"message": "成功",
	})

}

// 获取消息聊天详情列表
func UserChatListGet(c *gin.Context) {
	userID := c.MustGet("UserID").(string)
	var list []user.UserContentList
	query := model.DataBase.Model(user.UserContentList{})
	query.Where("user_id = ? OR send_id = ?", userID, userID)
	//从后往前查50条数据
	query.Order("id desc").Limit(50).Find(&list)
	//查出来的数据chatList再倒序
	data := make([]apiResultChatContentList, len(list))
	//list
	for i := 0; i <= len(list)-1; i++ {
		imgPath := list[i].ImgPath
		if imgPath == "" {
			imgPath = ""
		} else {
			imgPath = os.Getenv("ALI_OSS_DOMAIN") + "/" + list[i].ImgPath
		}
		data[i] = apiResultChatContentList{
			ID:        list[i].ID,
			UserId:    list[i].UserID,
			CreatedAt: list[i].CreatedAt.Unix(),
			Text:      list[i].Text,
			Img:       imgPath,
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"list": data,
		},
		"message": "成功",
	})

}

// 获取消息聊天详情列表
func UserChatListGet1(c *gin.Context) {
	userID := c.MustGet("UserID").(string)
	sendID := c.Query("id")
	var list []user.UserContentList
	query := model.DataBase.Model(user.UserContentList{})
	query.Where("(user_id = ? and send_id = ?) or (user_id = ? and send_id = ?)", userID, sendID, sendID, userID)
	//从后往前查50条数据
	query.Order("id desc").Limit(50).Find(&list)
	//查出来的数据chatList再倒序
	data := make([]apiResultChatContentList, len(list))
	//list
	for i := 0; i <= len(list)-1; i++ {
		imgPath := list[i].ImgPath
		if imgPath == "" {
			imgPath = ""
		} else {
			imgPath = os.Getenv("ALI_OSS_DOMAIN") + "/" + list[i].ImgPath
		}
		data[i] = apiResultChatContentList{
			ID:        list[i].ID,
			UserId:    list[i].UserID,
			CreatedAt: list[i].CreatedAt.Unix(),
			Text:      list[i].Text,
			Img:       imgPath,
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"list": data,
		},
		"message": "成功",
	})

}

// 已读
func UserReadPost(c *gin.Context) {
	userID := c.MustGet("UserID").(string)
	request := user.UserMessagePostRequest{}
	c.ShouldBindJSON(&request)
	chatUnRed := user.UserChatList{}
	if request.SendUserID == 0 {
		request.SendUserID = 7310
	}
	model.DataBase.Where("user_id = ? And send_id = ?", userID, request.SendUserID).First(&chatUnRed)
	if chatUnRed.ID != 0 {
		if chatUnRed.Number == 0 {
			c.JSON(http.StatusOK, gin.H{
				"code":    200,
				"data":    nil,
				"message": "成功",
			})
			return
		}
		chatUnRed.Number = 0
		model.DataBase.Save(&chatUnRed)
		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"data":    nil,
			"message": "成功",
		})
		return
	}
	if chatUnRed.ID == 0 {
		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"data":    nil,
			"message": "成功",
		})
		return
	}

}
