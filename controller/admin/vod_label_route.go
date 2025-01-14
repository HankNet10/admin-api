package admin

import (
	"myadmin/model/vod"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

var vodLabelModel vod.VodLabel

func VodLabelList(c *gin.Context) {
	request := vod.VodLabelParam{}
	c.BindQuery(&request)
	list, total := vodLabelModel.List(request)
	c.JSON(http.StatusOK, gin.H{
		"list":  list,
		"total": total,
	})
}

func VodLabelCreate(c *gin.Context) {
	request := vod.VodLabel{}
	c.ShouldBindJSON(&request)
	request.Save()

	c.JSON(http.StatusOK, gin.H{
		"message": "创建成功！",
	})
}

func VodLabelUpdate(c *gin.Context) {
	var request vod.VodLabel
	c.ShouldBindJSON(&request)
	model := vodLabelModel.SelectByID(request.ID)
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

func VodLabelDelete(c *gin.Context) {
	id := c.Query("id")
	rid, err := strconv.Atoi(id)
	if err != nil || rid <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "参数错误",
		})
		return
	}
	model := vodLabelModel.SelectByID(uint(rid))
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
