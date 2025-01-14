package admin

import (
	"myadmin/model/ad"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func AdPostionList(c *gin.Context) {
	request := ad.AdPostionParam{}
	c.BindQuery(&request)
	request.Sort = "-id"
	list, total := ad.AdPostionModel.List(request)
	c.JSON(http.StatusOK, gin.H{
		"list":  list,
		"total": total,
	})
}

func AdPostionCreate(c *gin.Context) {
	request := ad.AdPostion{}
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

func AdPostionUpdate(c *gin.Context) {
	request := ad.AdPostion{}
	c.ShouldBindJSON(&request)
	request.Save()
	c.JSON(http.StatusOK, gin.H{
		"message": "修改成功！",
	})
}

func AdPostionDelete(c *gin.Context) {
	id := c.Query("id")
	rid, err := strconv.Atoi(id)
	if err != nil || rid <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "参数错误",
		})
		return
	}
	model := ad.AdPostionModel.SelectByID(uint(rid))
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
