package admin

import (
	"myadmin/model/attack"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

var attackModel attack.Attack

//var attackParamsModel attack.AttackParams

func AttackList(c *gin.Context) {
	request := attack.AttackParam{}
	c.BindQuery(&request)
	if request.Limit > 100 || request.Limit < 1 {
		request.Limit = 20
	}
	list, total := attackModel.List(request)
	c.JSON(http.StatusOK, gin.H{
		"list":  list,
		"total": total,
	})
}

func AttackCreate(c *gin.Context) {
	request := attack.Attack{}
	c.ShouldBindJSON(&request)
	request.Save()

	c.JSON(http.StatusOK, gin.H{
		"message": "创建成功！",
	})
}

func AttackUpdate(c *gin.Context) {
	var request attack.Attack
	c.ShouldBindJSON(&request)
	model := attackModel.SelectByID(request.ID)
	if model == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "修改失败",
		})
		return
	}
	request.Save()
	c.JSON(http.StatusOK, gin.H{
		"message": "修改成功！",
	})
}

func AttackDelete(c *gin.Context) {
	id := c.Query("id")
	rid, err := strconv.Atoi(id)
	if err != nil || rid <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "参数错误",
		})
		return
	}
	model := attackModel.SelectByID(uint(rid))
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

// 参数管理
// func AttackParamsList(c *gin.Context) {
// 	request := attack.AttackParamRequest{}
// 	c.BindQuery(&request)
// 	list, total := attackParamsModel.List(request)
// 	c.JSON(http.StatusOK, gin.H{
// 		"list":  list,
// 		"total": total,
// 	})
// }

// func AttackParamsCreate(c *gin.Context) {
// 	request := attack.AttackParams{}
// 	c.ShouldBindJSON(&request)
// 	request.Save()

// 	c.JSON(http.StatusOK, gin.H{
// 		"message": "创建成功！",
// 	})
// }

// func AttackParamsUpdate(c *gin.Context) {
// 	var request attack.AttackParams
// 	c.ShouldBindJSON(&request)
// 	model := attackParamsModel.SelectByID(request.ID)
// 	if model == nil {
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"message": "修改失败",
// 		})
// 		return
// 	}
// 	request.Save()
// 	c.JSON(http.StatusOK, gin.H{
// 		"message": "修改成功！",
// 	})
// }

// func AttackParamsDelete(c *gin.Context) {
// 	id := c.Query("id")
// 	rid, err := strconv.Atoi(id)
// 	if err != nil || rid <= 0 {
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"message": "参数错误",
// 		})
// 		return
// 	}
// 	model := attackParamsModel.SelectByID(uint(rid))
// 	if model == nil {
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"message": "删除失败",
// 		})
// 		return
// 	}
// 	model.Delete()
// 	c.JSON(http.StatusOK, gin.H{
// 		"message": "删除成功！",
// 	})
// }
