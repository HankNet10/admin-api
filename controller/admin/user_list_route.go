package admin

import (
	"myadmin/model"
	"myadmin/model/user"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

var UserListModel user.UserList

func UserListList(c *gin.Context) {
	request := user.UserListParam{}
	c.BindQuery(&request)
	if request.Limit > 100 || request.Limit < 1 {
		request.Limit = 20
	}
	list, total := UserListModel.List(request)
	c.JSON(http.StatusOK, gin.H{
		"list":  list,
		"total": total,
	})
}

func UserListCreate(c *gin.Context) {
	request := user.UserList{}
	c.ShouldBindJSON(&request)
	err := request.Save()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "创建成功！",
	})
}

func UserListUpdate(c *gin.Context) {
	request := user.UserList{}
	c.ShouldBindJSON(&request)
	request.Save()

	c.JSON(http.StatusOK, gin.H{
		"message": "修改成功！",
	})
}

func UserListDelete(c *gin.Context) {
	id := c.Query("id")
	rid, err := strconv.Atoi(id)
	if err != nil || rid <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "参数错误",
		})
		return
	}
	model := UserListModel.SelectByID(uint(rid))
	if model == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "用户不存在",
		})
		return
	}
	model.Delete()
	c.JSON(http.StatusOK, gin.H{
		"message": "删除成功！",
	})
}

func UserDisableComment(c *gin.Context) {
	id := c.Query("id")
	rid, err := strconv.Atoi(id)
	if err != nil || rid <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "参数错误",
		})
		return
	}
	// 对用户禁言。
	userinfo := UserListModel.SelectByID(uint(rid))
	if userinfo == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "用户不存在",
		})
		return
	}
	userinfo.RejuectComment()
	user.DeleteCommentByUserID(id)
	c.JSON(http.StatusOK, gin.H{
		"message": "清理完成",
	})
}

// 邀请列表 后台使用
func UserInviteList(c *gin.Context) {
	request := user.UserInviteParam{}
	c.BindQuery(&request)
	if request.Limit > 100 || request.Limit < 1 {
		request.Limit = 20
	}
	request.Order = "-id"
	model := user.UserInviteList{}
	list, total := model.List(request)
	c.JSON(http.StatusOK, gin.H{
		"list":  list,
		"total": total,
	})
}

// 后台用户头像列表上传
func UserAvatarCreate(c *gin.Context) {
	request := user.UserAvatar{}
	c.ShouldBindJSON(&request)
	request.Save()

	c.JSON(http.StatusOK, gin.H{
		"message": "创建成功！",
	})
}

// 后台用户头像列表
func UserAvatarList(c *gin.Context) {
	request := user.UserListParam{}
	c.BindQuery(&request)
	if request.Limit > 100 || request.Limit < 1 {
		request.Limit = 20
	}
	request.Order = "-id"
	model := user.UserAvatar{}
	list, total := model.AdminList(request)
	c.JSON(http.StatusOK, gin.H{
		"list":  list,
		"total": total,
	})
}

// 后台用户头像列表
func UserAvatarDelete(c *gin.Context) {
	id := c.Query("id")
	rid, err := strconv.Atoi(id)
	if err != nil || rid <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "参数错误",
		})
		return
	}
	mUserAvatar := user.UserAvatar{}
	result := model.DataBase.First(&mUserAvatar, uint(rid))
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "删除失败",
		})
		return
	}
	model.DataBase.Delete(&mUserAvatar)
	c.JSON(http.StatusOK, gin.H{
		"message": "删除成功！",
	})
}

// 获取签到列表
func UserSignList(c *gin.Context) {
	request := user.UserSignListParam{}
	c.BindQuery(&request)
	if request.Page < 1 {
		request.Page = 1
	}

	if request.Limit > 100 || request.Limit < 1 {
		request.Limit = 20
	}
	request.Sort = "-id"
	list, total := user.UserSignListModel.List(request)
	c.JSON(http.StatusOK, gin.H{
		"list":  list,
		"total": total,
	})
}

func TaskListList(c *gin.Context) {
	request := user.TaskListParam{}
	c.BindQuery(&request)
	if request.Limit > 100 || request.Limit < 1 {
		request.Limit = 20
	}
	request.Order = "-id"
	list, total := user.TaskListModel.TaskList(request)
	c.JSON(http.StatusOK, gin.H{
		"list":  list,
		"total": total,
	})
}

func TaskListCreate(c *gin.Context) {
	request := user.TaskList{}
	c.ShouldBindJSON(&request)
	request.TaskListSave()

	c.JSON(http.StatusOK, gin.H{
		"message": "创建成功！",
	})
}

func TaskListUpdate(c *gin.Context) {
	var request user.TaskList
	c.ShouldBindJSON(&request)
	model := user.TaskListModel.TaskListSelectByID(request.ID)
	if model == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "博客不存在",
		})
		return
	}
	request.TaskListSave()
	c.JSON(http.StatusOK, gin.H{
		"message": "修改成功！",
	})
}

func TaskListDelete(c *gin.Context) {
	id := c.Query("id")
	rid, err := strconv.Atoi(id)
	if err != nil || rid <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "参数错误",
		})
		return
	}
	model := user.TaskListModel.TaskListSelectByID(uint(rid))
	if model == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "博客不存在",
		})
		return
	}
	model.TaskListDelete()
	c.JSON(http.StatusOK, gin.H{
		"message": "删除成功！",
	})
}

func TaskUserList(c *gin.Context) {
	request := user.TaskUserParam{}
	c.BindQuery(&request)
	if request.Limit > 100 || request.Limit < 1 {
		request.Limit = 20
	}
	list, total := user.TaskUserModel.TaskUserList(request)
	c.JSON(http.StatusOK, gin.H{
		"list":  list,
		"total": total,
	})
}

func TaskUserCreate(c *gin.Context) {
	request := user.TaskUser{}
	c.ShouldBindJSON(&request)
	request.TaskUserSave()

	c.JSON(http.StatusOK, gin.H{
		"message": "创建成功！",
	})
}

func TaskUserUpdate(c *gin.Context) {
	var request user.TaskUser
	c.ShouldBindJSON(&request)
	model := user.TaskUserModel.TaskUserSelectByID(request.ID)
	if model == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "数据不存在",
		})
		return
	}
	request.TaskUserSave()
	c.JSON(http.StatusOK, gin.H{
		"message": "修改成功！",
	})
}

func TaskUserDelete(c *gin.Context) {
	id := c.Query("id")
	rid, err := strconv.Atoi(id)
	if err != nil || rid <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "参数错误",
		})
		return
	}
	model := user.TaskUserModel.TaskUserSelectByID(uint(rid))
	if model == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "博客不存在",
		})
		return
	}
	model.TaskUserDelete()
	c.JSON(http.StatusOK, gin.H{
		"message": "删除成功！",
	})
}
