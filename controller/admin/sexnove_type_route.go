package admin

import (
	"myadmin/model/sexnovel"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

var sexnovelTypeModel sexnovel.SexnovelType

func SexnovelTypeList(c *gin.Context) {
	request := sexnovel.SexnovelTypeParam{}
	c.BindQuery(&request)
	list, total := sexnovelTypeModel.List(request)
	c.JSON(http.StatusOK, gin.H{
		"list":  list,
		"total": total,
	})
}

func SexnovelTypeCreate(c *gin.Context) {
	request := sexnovel.SexnovelType{}
	c.ShouldBindJSON(&request)
	request.Save()

	c.JSON(http.StatusOK, gin.H{
		"message": "创建成功！",
	})
}

type sexnovelTypeUpdateInput struct {
	ID     uint
	Name   string
	Status uint8
	Sort   uint8
	Icon   string
}

func SexnovelTypeUpdate(c *gin.Context) {
	var request sexnovelTypeUpdateInput
	c.ShouldBindJSON(&request)
	model := sexnovelTypeModel.SelectByID(request.ID)
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

func SexnovelTypeDelete(c *gin.Context) {
	id := c.Query("id")
	rid, err := strconv.Atoi(id)
	if err != nil || rid <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "参数错误",
		})
		return
	}
	model := sexnovelTypeModel.SelectByID(uint(rid))
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
