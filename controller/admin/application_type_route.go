package admin

import (
	"myadmin/model/applicationad"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func ApplicationTypeList(c *gin.Context) {
	request := applicationad.ApplicationTypeParam{}
	c.BindQuery(&request)
	request.Sort = "-id"
	list, total := applicationad.ApplicationTypeModel.List(request)
	c.JSON(http.StatusOK, gin.H{
		"list":  list,
		"total": total,
	})
}

func ApplicationTypeCreate(c *gin.Context) {
	request := applicationad.ApplicationType{}
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

func ApplicationTypeUpdate(c *gin.Context) {
	request := applicationad.ApplicationType{}
	c.ShouldBindJSON(&request)
	request.Save()
	c.JSON(http.StatusOK, gin.H{
		"message": "修改成功！",
	})
}

func ApplicationTypeDelete(c *gin.Context) {
	id := c.Query("id")
	rid, err := strconv.Atoi(id)
	if err != nil || rid <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "参数错误",
		})
		return
	}
	model := applicationad.ApplicationTypeModel.SelectByID(uint(rid))
	if model == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "数据不存在",
		})
		return
	}
	model.Delete()
	c.JSON(http.StatusOK, gin.H{
		"message": "删除成功！",
	})
}
