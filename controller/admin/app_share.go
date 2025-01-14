package admin

import (
	"myadmin/model"
	"myadmin/model/user"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func AppShareList(c *gin.Context) {
	request := user.AppShareParam{}
	c.BindQuery(&request)
	if request.Page < 1 {
		request.Page = 1
	}

	if request.Limit > 100 || request.Limit < 1 {
		request.Limit = 20
	}
	request.Sort = "-id"
	list, total := user.Appshare.List(request)
	c.JSON(http.StatusOK, gin.H{
		"list":  list,
		"total": total,
	})
}

func AppShareCreate(c *gin.Context) {
	request := user.AppShare{}
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

func AppShareUpdate(c *gin.Context) {
	request := user.AppShare{}
	c.ShouldBindJSON(&request)
	if result := model.DataBase.Save(request); result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": result.Error.Error(),
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"message": "修改成功！",
		})
	}
}

func AppShareDelete(c *gin.Context) {
	id := c.Query("id")
	rid, err := strconv.Atoi(id)
	if err != nil || rid <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "参数错误",
		})
		return
	}
	model := user.Appshare.SelectByID(uint(rid))
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
