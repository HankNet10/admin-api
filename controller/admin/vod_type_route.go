package admin

import (
	"myadmin/model/vod"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

var vodTypeModel vod.VodType

func VodTypeList(c *gin.Context) {
	request := vod.VodTypeParam{}
	c.BindQuery(&request)
	list, total := vodTypeModel.List(request)
	c.JSON(http.StatusOK, gin.H{
		"list":  list,
		"total": total,
	})
}

func VodTypeCreate(c *gin.Context) {
	request := vod.VodType{}
	c.ShouldBindJSON(&request)
	request.Save()

	c.JSON(http.StatusOK, gin.H{
		"message": "创建成功！",
	})
}

type vodTypeUpdateInput struct {
	ID     uint
	Name   string // 密码
	Parent uint
	Status uint8
	Sort   uint8
	Icon   string
}

func VodTypeUpdate(c *gin.Context) {
	var request vodTypeUpdateInput
	c.ShouldBindJSON(&request)
	model := vodTypeModel.SelectByID(request.ID)
	if model == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "分类不存在",
		})
		return
	}
	model.Name = request.Name
	model.Parent = request.Parent
	model.Status = request.Status
	model.Sort = request.Sort
	model.Save()
	c.JSON(http.StatusOK, gin.H{
		"message": "修改成功！",
	})
}

func VodTypeDelete(c *gin.Context) {
	id := c.Query("id")
	rid, err := strconv.Atoi(id)
	if err != nil || rid <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "参数错误",
		})
		return
	}
	model := vodTypeModel.SelectByID(uint(rid))
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
