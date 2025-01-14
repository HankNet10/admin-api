package admin

import (
	"myadmin/model"
	"myadmin/model/user"
	"myadmin/util/redis"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func UserCommentList(c *gin.Context) {
	request := user.CommentParam{}
	c.BindQuery(&request)
	if request.Limit > 100 || request.Limit < 1 {
		request.Limit = 20
	}
	list, total := user.CommentModel.List(request)
	c.JSON(http.StatusOK, gin.H{
		"list":  list,
		"total": total,
	})
}
func UserCommentDelete(c *gin.Context) {
	id := c.Query("id")
	rid, err := strconv.Atoi(id)
	if err != nil || rid <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "参数错误",
		})
		return
	}
	model := user.CommentModel.SelectByID(uint(rid))
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

func UserCommentUpdate(c *gin.Context) {
	request := user.UserComment{}
	c.ShouldBindJSON(&request)
	if request.ID < 1 || request.UserID < 1 || request.Comment == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "选择修改数据",
		})
		return
	}
	request.Save()
	c.JSON(http.StatusOK, gin.H{
		"message": "修改成功！",
	})
	cacheKey := "api:commentlist:" + strconv.FormatUint(uint64(request.CID), 10) + ":" + strconv.FormatUint(uint64(request.Type), 10) + ":"
	redis.PullPrefix(cacheKey) //删除缓存
}

// 批量通过此页
func AllowLen(c *gin.Context) {
	request := user.RequestUserAllowComment{}
	c.ShouldBindJSON(&request)
	model.DataBase.Model(user.UserComment{}).Where("id in ?", request.Ids).Where("status = ?", 0).Updates(map[string]interface{}{"status": 1})
	cacheKey := "api:commentlist:"
	redis.PullPrefix(cacheKey) //删除所有评论缓存
}
