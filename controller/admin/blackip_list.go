package admin

import (
	"myadmin/model/user"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

var BlackIpModel user.BlackIp

func BlackIpList(c *gin.Context) {
	request := user.BlackIpListParam{}
	c.BindQuery(&request)
	if request.Limit > 100 || request.Limit < 1 {
		request.Limit = 20
	}
	request.Sort = "-id"
	list, total := BlackIpModel.List(request)
	c.JSON(http.StatusOK, gin.H{
		"list":  list,
		"total": total,
	})
}

func BlackIpCreate(c *gin.Context) {
	request := user.BlackIp{}
	c.ShouldBindJSON(&request)
	request.Save()

	c.JSON(http.StatusOK, gin.H{
		"message": "创建成功！",
	})
}

func BlackIpDelete(c *gin.Context) {
	id := c.Query("id")
	rid, err := strconv.Atoi(id)
	if err != nil || rid <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "参数错误",
		})
		return
	}
	model := user.BlackIp{}.SelectByID(uint(rid))
	if model == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "不存在",
		})
		return
	}
	model.Delete()
	c.JSON(http.StatusOK, gin.H{
		"message": "删除成功！",
	})
}
