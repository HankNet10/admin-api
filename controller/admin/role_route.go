package admin

import (
	"myadmin/model/admin"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

var roleModel admin.SysRole

type roleInput struct {
	Name   string // 密码
	Remark string // 验证码
}

func RoleCreate(c *gin.Context) {
	var request roleInput
	c.ShouldBindJSON(&request)
	au := roleModel.SelectRoleByName(request.Name)
	if au != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "角色已存在",
		})
		return
	}

	newModel := admin.SysRole{
		Name:   request.Name,
		Remark: request.Remark,
	}
	newModel.Save()
	c.JSON(http.StatusOK, gin.H{
		"message": "创建成功！",
	})
}

func RoleList(c *gin.Context) {
	request := admin.SysRoleParam{}
	c.BindQuery(&request)
	list, total := roleModel.List(request)
	c.JSON(http.StatusOK, gin.H{
		"list":  list,
		"total": total,
	})
}

type roleUpdateInput struct {
	ID     uint
	Name   string // 密码
	Remark string // 验证码
}

func RoleUpdate(c *gin.Context) {
	var request roleUpdateInput
	c.ShouldBindJSON(&request)
	re := roleModel.SelectByID(request.ID)
	if re == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "角色不存在",
		})
		return
	}
	re.Name = request.Name
	re.Remark = request.Remark
	re.Save()
	c.JSON(http.StatusOK, gin.H{
		"message": "修改成功！",
	})
}

func RoleDelete(c *gin.Context) {
	id := c.Query("id")
	rid, err := strconv.Atoi(id)
	if err != nil || rid <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "参数错误",
		})
		return
	}
	role := roleModel.SelectByID(uint(rid))
	if role == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "角色不存在",
		})
		return
	}
	role.Delete()
	c.JSON(http.StatusOK, gin.H{
		"message": "删除成功！",
	})
}
