package admin

import (
	"myadmin/model/seximg"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

var seximgTypeModel seximg.SeximgType

func SeximgTypeList(c *gin.Context) {
	request := seximg.SeximgTypeParam{}
	c.BindQuery(&request)
	list, total := seximgTypeModel.List(request)
	c.JSON(http.StatusOK, gin.H{
		"list":  list,
		"total": total,
	})
}

func SeximgTypeCreate(c *gin.Context) {
	request := seximg.SeximgType{}
	c.ShouldBindJSON(&request)
	request.Save()

	c.JSON(http.StatusOK, gin.H{
		"message": "创建成功！",
	})
}

type seximgTypeUpdateInput struct {
	ID     uint
	Name   string
	Status uint8
	Sort   uint8
	Icon   string
}

func SeximgTypeUpdate(c *gin.Context) {
	var request seximgTypeUpdateInput
	c.ShouldBindJSON(&request)
	model := seximgTypeModel.SelectByID(request.ID)
	if model == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "分类不存在",
		})
		return
	}
	model.Name = request.Name
	model.Status = request.Status
	model.Sort = request.Sort
	model.Save()
	c.JSON(http.StatusOK, gin.H{
		"message": "修改成功！",
	})
}

func SeximgTypeDelete(c *gin.Context) {
	id := c.Query("id")
	rid, err := strconv.Atoi(id)
	if err != nil || rid <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "参数错误",
		})
		return
	}
	model := seximgTypeModel.SelectByID(uint(rid))
	if model == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "分类不存在",
		})
		return
	}
	model.Delete()
	c.JSON(http.StatusOK, gin.H{
		"message": "删除成功！",
	})
}
