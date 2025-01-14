package admin

import (
	"myadmin/util/redis"
	"net/http"

	"github.com/gin-gonic/gin"
)

//
type HotSearch struct {
	Wd    string  `json:"wd"`
	Score float64 `json:"score"`
}
type HotSeachParam struct {
	Page  int     `form:"page"`
	Limit int     `form:"limit"`
	Wd    string  `json:"wd"`
	Score float64 `json:"score"`
}

func HotListList(c *gin.Context) {
	request := HotSeachParam{}
	c.BindQuery(&request)
	if request.Limit > 100 || request.Limit < 1 {
		request.Limit = 20
	}
	if request.Page < 1 {
		request.Page = 1
	}
	lists := redis.HotSeachList(int64(request.Page), int64(request.Limit))
	results := make([]HotSearch, len(lists))
	for i, v := range lists {
		results[i] = HotSearch{
			Wd:    interfaceToString(v.Member),
			Score: v.Score,
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"list":  results,
		"total": redis.HotSeachCount(),
	})
}
func interfaceToString(i interface{}) string {
	// 使用类型断言将接口值转换为字符串
	if str, ok := i.(string); ok {
		return str
	}
	return ""
}

func HotUpdate(c *gin.Context) {
	var request HotSeachParam
	c.ShouldBindJSON(&request)
	redis.UpdateZsetValue("search", request.Wd, request.Score)
	c.JSON(http.StatusOK, gin.H{
		"message": "更新成功",
	})
}
func HotDelete(c *gin.Context) {
	wd := c.Query("wd")
	redis.DelZsetValue("search", wd)
	c.JSON(http.StatusOK, gin.H{
		"message": "删除成功！",
	})
}

func HotReplaceAll(c *gin.Context) {
	redis.Pull("search")
	c.JSON(http.StatusOK, gin.H{
		"message": "删除成功！",
	})
}
