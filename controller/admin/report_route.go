package admin

import (
	"myadmin/model/user"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ReportList(c *gin.Context) {
	request := user.ReportListParam{}
	c.BindQuery(&request)
	if request.Page < 1 {
		request.Page = 1
	}

	if request.Limit > 20 || request.Limit < 1 {
		request.Limit = 20
	}
	request.Sort = "-id"
	list, total := user.ReportListModel.List(request) //   .List(request)
	c.JSON(http.StatusOK, gin.H{
		"list":  list,
		"total": total,
	})
}
