package admin

import (
	"myadmin/model/suggest"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SuggestList(c *gin.Context) {
	request := suggest.SuggestListParam{}
	c.BindQuery(&request)
	if request.Page < 1 {
		request.Page = 1
	}

	if request.Limit > 20 || request.Limit < 1 {
		request.Limit = 20
	}
	request.Sort = "-id"
	list, total := suggest.SuggestListModel.List(request) //   .List(request)
	c.JSON(http.StatusOK, gin.H{
		"list":  list,
		"total": total,
	})
}
