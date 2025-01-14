package api

import (
	"myadmin/model/attack"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

var attackModel attack.Attack

// 获取
func AttackList(c *gin.Context) {
	request := attack.AttackParam{}
	c.BindQuery(&request)
	appID := c.Request.Header.Get("x-appid")
	if appID == "" {
		appID = "1"
	}
	request.AppId, _ = strconv.Atoi(appID)
	if request.Page < 1 {
		request.Page = 1
	}
	if request.Limit > 100 || request.Limit < 1 {
		request.Limit = 20
	}
	request.Status = 1
	list, total := attackModel.List(request)

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"total": total,
			"list":  list,
		},
		"message": "",
	})
}
