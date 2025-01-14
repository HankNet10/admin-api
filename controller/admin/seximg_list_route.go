package admin

import (
	"myadmin/model/seximg"
	"myadmin/util"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

var seximgModel seximg.Seximg

func SeximgList(c *gin.Context) {
	request := seximg.SeximgParam{}
	c.BindQuery(&request)
	if request.Limit > 100 || request.Limit < 1 {
		request.Limit = 20
	}
	list, total := seximgModel.List(request)
	c.JSON(http.StatusOK, gin.H{
		"list":  list,
		"total": total,
	})
}

func SeximgCreate(c *gin.Context) {
	request := seximg.Seximg{}
	c.ShouldBindJSON(&request)
	request.Save()

	c.JSON(http.StatusOK, gin.H{
		"message": "创建成功！",
	})
}

func SeximgUpdate(c *gin.Context) {
	var request seximg.Seximg
	c.ShouldBindJSON(&request)
	model := seximgModel.SelectByID(request.ID)
	if model == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "修改失败",
		})
		return
	}
	count := seximg.SexImageModel.SelectCountBySeximgID(model.ID)
	request.ImgCount = uint(count)
	request.Save()
	c.JSON(http.StatusOK, gin.H{
		"message": "修改成功！",
	})
}

func SeximgDelete(c *gin.Context) {
	id := c.Query("id")
	rid, err := strconv.Atoi(id)
	if err != nil || rid <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "参数错误",
		})
		return
	}
	model := seximgModel.SelectByID(uint(rid))
	if model == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "删除失败",
		})
		return
	}
	model.Delete()
	c.JSON(http.StatusOK, gin.H{
		"message": "删除成功！",
	})
}

func SexImageCreate(c *gin.Context) {
	request := seximg.SexImage{}
	c.ShouldBindJSON(&request)
	request.Save()

	c.JSON(http.StatusOK, gin.H{
		"message": "创建成功！",
	})
}

func SexImageDelete(c *gin.Context) {
	id := c.Query("id")
	rid, err := strconv.Atoi(id)
	if err != nil || rid <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "参数错误",
		})
		return
	}
	var imageModel seximg.SexImage
	model := imageModel.SelectByID(uint(rid))
	if model == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "图片不存在",
		})
		return
	}
	model.Delete()
	util.OssDeleteObject(model.Path)
	c.JSON(http.StatusOK, gin.H{
		"message": "删除成功！",
	})
}
