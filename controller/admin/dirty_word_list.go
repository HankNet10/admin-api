package admin

import (
	"myadmin/model"
	"myadmin/model/user"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

var dirtyWordModel user.DirtyWord

func DirtyWordList(c *gin.Context) {
	request := user.DirtyWordListParam{}
	c.BindQuery(&request)
	if request.Limit > 100 || request.Limit < 1 {
		request.Limit = 20
	}
	request.Sort = "-id"
	list, total := dirtyWordModel.List(request)
	c.JSON(http.StatusOK, gin.H{
		"list":  list,
		"total": total,
	})
}

func DirtyWordCreate(c *gin.Context) {
	request := user.DirtyWord{}
	c.ShouldBindJSON(&request)
	Item := user.DirtyWord{}
	model.DataBase.Where("name = ?", request.Name).First(&Item)
	if Item.ID > 0 {
		c.JSON(http.StatusOK, gin.H{
			"message": "有重复敏感词 " + Item.Name,
		})
		return
	}
	request.Save()

	c.JSON(http.StatusOK, gin.H{
		"message": "创建成功！",
	})
}

func DirtyWordDelete(c *gin.Context) {
	id := c.Query("id")
	rid, err := strconv.Atoi(id)
	if err != nil || rid <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "参数错误",
		})
		return
	}
	model := user.DirtyWord{}.SelectByID(uint(rid))
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
