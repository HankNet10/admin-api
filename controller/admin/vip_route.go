package admin

import (
	"myadmin/model/vip"
	"myadmin/util/sugar"
	"net/http"

	"github.com/gin-gonic/gin"
)

func VipPost(c *gin.Context) {
	m := vip.VipList{}
	c.ShouldBindJSON(&m)
	// 检查是否存在id的vod
	vip.VipListSave(m)
	c.JSON(http.StatusOK, gin.H{
		"message": "创建成功！",
	})
}

func VipPut(c *gin.Context) {
	m := vip.VipList{}
	c.ShouldBindJSON(&m)
	// 检查是否存在id的vod
	vip.VipListSave(m)
	c.JSON(http.StatusOK, gin.H{
		"message": "创建成功！",
	})
}

func VipList(c *gin.Context) {
	var request vip.VipListParam
	c.BindQuery(&request)
	list, total := vip.VipListData(request)
	c.JSON(http.StatusOK, gin.H{
		"list":  list,
		"total": total,
	})
}

func VipDelete(c *gin.Context) {
	id := sugar.StringToUint(c.Query("id"))
	vip.VipListDelete(id)
	c.JSON(http.StatusOK, gin.H{
		"message": "删除成功",
	})
}
