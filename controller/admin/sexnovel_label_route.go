package admin

import (
	"myadmin/model/sexnovel"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

var sexnovelLabelModel sexnovel.SexnovelLabel

func SexnovelLabelList(c *gin.Context) {
	request := sexnovel.SexnovelLabelParam{}
	c.BindQuery(&request)
	list, total := sexnovelLabelModel.List(request)
	c.JSON(http.StatusOK, gin.H{
		"list":  list,
		"total": total,
	})
}

func SexnovelLabelCreate(c *gin.Context) {
	request := sexnovel.SexnovelLabel{}
	c.ShouldBindJSON(&request)
	request.Save()

	c.JSON(http.StatusOK, gin.H{
		"message": "创建成功！",
	})
}

func SexnovelLabelUpdate(c *gin.Context) {
	var request sexnovel.SexnovelLabel
	c.ShouldBindJSON(&request)
	model := sexnovelLabelModel.SelectByID(request.ID)
	if model == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "分类不存在",
		})
		return
	}
	model.Name = request.Name
	model.Index = request.Index
	model.Status = request.Status
	model.Sort = request.Sort
	model.Save()
	c.JSON(http.StatusOK, gin.H{
		"message": "修改成功！",
	})
}

func SexnovelLabelDelete(c *gin.Context) {
	id := c.Query("id")
	rid, err := strconv.Atoi(id)
	if err != nil || rid <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "参数错误",
		})
		return
	}
	model := sexnovelLabelModel.SelectByID(uint(rid))
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
