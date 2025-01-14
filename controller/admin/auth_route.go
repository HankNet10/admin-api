package admin

import (
	"myadmin/model/admin"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

var authModel admin.SysAuth

type authInput struct {
	RoleID string
	Path   string
	Method int
	Remark string
}

func AuthCreate(c *gin.Context) {
	var request authInput
	c.ShouldBindJSON(&request)
	newModel := admin.SysAuth{
		RoleID: request.RoleID,
		Path:   request.Path,
		Method: request.Method,
		Remark: request.Remark,
	}
	newModel.Save()
	c.JSON(http.StatusOK, gin.H{
		"message": "创建成功！",
	})
}

func AuthList(c *gin.Context) {
	var param admin.SysAuthParam
	c.BindQuery(&param)
	result, total := authModel.List(param)
	c.JSON(http.StatusOK, gin.H{
		"list":  result,
		"total": total,
		"limit": param.Limit,
	})
}

type authUpdateInput struct {
	ID     uint
	RoleID string
	Path   string //授权路径
	Method int    // 1all, 2get,3,post,4put,5delete
	Remark string
}

func AuthUpdate(c *gin.Context) {
	var request authUpdateInput
	c.ShouldBindJSON(&request)
	model := authModel.SelectByID(request.ID)
	if model == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "角色不存在",
		})
		return
	}
	model.RoleID = request.RoleID
	model.Path = request.Path
	model.Method = request.Method
	model.Remark = request.Remark
	model.Save()
	c.JSON(http.StatusOK, gin.H{
		"message": "修改成功！",
	})
}

func AuthDelete(c *gin.Context) {
	id := c.Query("id")
	rid, err := strconv.Atoi(id)
	if err != nil || rid <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "参数错误",
		})
		return
	}
	auth := authModel.SelectByID(uint(rid))
	if auth == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "授权不存在",
		})
		return
	}
	auth.Delete()
	c.JSON(http.StatusOK, gin.H{
		"message": "删除成功！",
	})
}
