package admin

import (
	"myadmin/model"
	"myadmin/model/user"
	"net/http"

	"github.com/gin-gonic/gin"
)

var UserChatList user.UserChatList
var UserContentList user.UserContentList

//聊天关系列表
func UserChatListList(c *gin.Context) {
	request := user.UserChatListParam{}
	c.BindQuery(&request)
	if request.Limit > 100 || request.Limit < 1 {
		request.Limit = 20
	}

	list, total := UserChatList.List(request)
	c.JSON(http.StatusOK, gin.H{
		"list":  list,
		"total": total,
	})
}

//具体消息列表
func UserContentListList(c *gin.Context) {
	request := user.UserChatListParam{}
	c.BindQuery(&request)
	if request.Limit > 100 || request.Limit < 1 {
		request.Limit = 20
	}
	list, total := UserContentList.List(request)
	c.JSON(http.StatusOK, gin.H{
		"list":  list,
		"total": total,
	})
}

//回复用户消息
func UserChatAdminSend(c *gin.Context) {

	request := user.UserMessagePostRequest{}
	c.ShouldBindJSON(&request)
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
	// 判断之前是否有聊天关系
	createChat := user.UserChatList{}
	query := model.DataBase.Model(createChat)
	query.Where("user_id = ? And send_id = ?", request.UserId, request.SendUserID)
	query.Find(&createChat)
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
	result := model.DataBase.Save(&createContent)
	if result.Error != nil {
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

//
func UserMssageRead(c *gin.Context) {
	request := user.UserMessagePostRequest{}
	c.ShouldBindJSON(&request)
	code := UserChatList.Read(request)

	c.JSON(http.StatusOK, gin.H{
		"code":    code,
		"message": "完成",
	})
}
