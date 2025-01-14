package admin

import (
	"myadmin/model/blog"
	"myadmin/util"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

var blogListModel blog.BlogList

func BlogListList(c *gin.Context) {
	request := blog.BlogListParam{}
	c.BindQuery(&request)
	if request.Limit > 100 || request.Limit < 1 {
		request.Limit = 20
	}
	list, total := blogListModel.List(request)
	c.JSON(http.StatusOK, gin.H{
		"list":  list,
		"total": total,
	})
}

func BlogListCreate(c *gin.Context) {
	request := blog.BlogList{}
	c.ShouldBindJSON(&request)
	request.Save()

	c.JSON(http.StatusOK, gin.H{
		"message": "创建成功！",
	})
}

func BlogListUpdate(c *gin.Context) {
	var request blog.BlogList
	c.ShouldBindJSON(&request)
	model := blogListModel.SelectByID(request.ID)
	if model == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "博客不存在",
		})
		return
	}
	request.Save()
	c.JSON(http.StatusOK, gin.H{
		"message": "修改成功！",
	})
}

func BlogListDelete(c *gin.Context) {
	id := c.Query("id")
	rid, err := strconv.Atoi(id)
	if err != nil || rid <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "参数错误",
		})
		return
	}
	model := blogListModel.SelectByID(uint(rid))
	if model == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "博客不存在",
		})
		return
	}
	model.Delete()
	c.JSON(http.StatusOK, gin.H{
		"message": "删除成功！",
	})
}

func BlogImageCreate(c *gin.Context) {
	request := blog.BlogImage{}
	c.ShouldBindJSON(&request)
	request.Save()

	c.JSON(http.StatusOK, gin.H{
		"message": "创建成功！",
	})
}

func BlogImageDelete(c *gin.Context) {
	id := c.Query("id")
	rid, err := strconv.Atoi(id)
	if err != nil || rid <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "参数错误",
		})
		return
	}
	var imageModel blog.BlogImage
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
