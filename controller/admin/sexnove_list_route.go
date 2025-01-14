package admin

import (
	"encoding/json"
	"log"
	"myadmin/model/sexnovel"
	"myadmin/util/redis"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

var sexnovelModel sexnovel.Sexnovel
var sexnovelChapterModel sexnovel.SexnovelChapter
var sexnovelContentModel sexnovel.SexnovelContent

func SexnovelList(c *gin.Context) {
	request := sexnovel.SexnovelParam{}
	c.BindQuery(&request)
	if request.Limit > 100 || request.Limit < 1 {
		request.Limit = 20
	}
	list, total := sexnovelModel.List(request)
	c.JSON(http.StatusOK, gin.H{
		"list":  list,
		"total": total,
	})
}

func SexnovelCreate(c *gin.Context) {
	request := sexnovel.Sexnovel{}
	c.ShouldBindJSON(&request)
	request.Save()

	c.JSON(http.StatusOK, gin.H{
		"message": "创建成功！",
	})
}

func SexnovelCreate1(c *gin.Context) {
	request := sexnovel.SexnovelUploadParam{}
	c.ShouldBindJSON(&request)
	index := strings.LastIndex(request.Title, ".")
	title := request.Title[:index]
	model := sexnovel.Sexnovel{
		AppID:        uint(request.AppID),
		ChapterCount: 1,
		Title:        title,
		TypeID:       uint(request.TypeID),
		//Path:         request.Path,
		Status: request.Status,
	}
	model.Save()
	if model.ID > 0 {
		content := sexnovel.SexnovelContent{
			ChapterID: int(model.ID),
			IsLong:    false,
			Content:   request.Path,
		}
		content.Save()
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "创建成功！",
	})
}

func SexnovelUpdate(c *gin.Context) {
	var request sexnovel.Sexnovel
	c.ShouldBindJSON(&request)
	model := sexnovelModel.SelectByID(request.ID)
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

func SexnovelDelete(c *gin.Context) {
	id := c.Query("id")
	rid, err := strconv.Atoi(id)
	if err != nil || rid <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "参数错误",
		})
		return
	}
	model := sexnovelModel.SelectByID(uint(rid))
	if model == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "删除失败",
		})
		return
	}
	modelContent := sexnovelContentModel.SelectContentInfo(false, model.ID)
	if modelContent == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "删除失败",
		})
		return
	}
	modelContent.Delete()
	model.Delete()
	c.JSON(http.StatusOK, gin.H{
		"message": "删除成功！",
	})
}

// 章节管理
func SexnovelChapterList(c *gin.Context) {
	request := sexnovel.SexnovelChapterParam{}
	c.BindQuery(&request)
	//if request.Limit > 100 || request.Limit < 1 {
	request.Limit = 99999
	//}
	request.Page = 1
	list, total := sexnovelChapterModel.List(request)
	c.JSON(http.StatusOK, gin.H{
		"list":  list,
		"total": total,
	})
}

func SexnovelChapterCreate(c *gin.Context) {
	request := sexnovel.SexnovelChapter{}
	c.ShouldBindJSON(&request)
	request.Save()

	c.JSON(http.StatusOK, gin.H{
		"message": "创建成功！",
	})
}

func SexnovelChapterUpdate(c *gin.Context) {
	var request sexnovel.SexnovelChapter
	c.ShouldBindJSON(&request)
	model := sexnovelChapterModel.SelectByID(request.ID)
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

func SexnovelChapterDelete(c *gin.Context) {
	id := c.Query("id")
	rid, err := strconv.Atoi(id)
	if err != nil || rid <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "参数错误",
		})
		return
	}
	model := sexnovelChapterModel.SelectByID(uint(rid))
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

// 内容管理
func SexnovelContent(c *gin.Context) {
	id, _ := strconv.Atoi(c.Query("id"))
	model := sexnovelContentModel.SelectContentByID(uint(id))
	c.JSON(http.StatusOK, gin.H{
		"list":  model,
		"total": 1,
	})
}

// 内容
func SexnovelContentInfo(c *gin.Context) {
	id, _ := strconv.Atoi(c.Query("id"))
	islong, _ := strconv.ParseBool(c.Query("islong"))
	model := sexnovelContentModel.SelectContentInfo(islong, uint(id))
	c.JSON(http.StatusOK, gin.H{
		"list":  model,
		"total": 1,
	})
}

func SexnovelContentCreate(c *gin.Context) {
	request := sexnovel.SexnovelContent{}
	c.ShouldBindJSON(&request)
	request.Save()

	c.JSON(http.StatusOK, gin.H{
		"message": "创建成功！",
	})
}

func SexnovelContentUpdate(c *gin.Context) {
	var request sexnovel.SexnovelContent
	c.ShouldBindJSON(&request)
	model := sexnovelContentModel.SelectByID(request.ID)
	if model == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "修改失败",
		})
		return
	}
	request.Save()
	cKey := "api:sexnovel:content:1:" + strconv.FormatBool(model.IsLong) + ":" + strconv.Itoa(model.ChapterID)
	if jsonData, err := redis.Get(cKey); err == nil {
		err = json.Unmarshal([]byte(jsonData), &model)
		if err == nil {
			resultData := gin.H{
				"code":    200,
				"data":    request,
				"message": "",
			}
			jsonData, err := json.Marshal(resultData)
			if err != nil {
				log.Panic("jsonData set err", err)
			}
			redis.Set(cKey, jsonData, 72*time.Hour)
		}
	}
	cKey = "api:sexnovel:content:2:" + strconv.FormatBool(model.IsLong) + ":" + strconv.Itoa(model.ChapterID)
	if jsonData, err := redis.Get(cKey); err == nil {
		err = json.Unmarshal([]byte(jsonData), &model)
		if err == nil {
			resultData := gin.H{
				"code":    200,
				"data":    request,
				"message": "",
			}
			jsonData, err := json.Marshal(resultData)
			if err != nil {
				log.Panic("jsonData set err", err)
			}
			redis.Set(cKey, jsonData, 72*time.Hour)
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "修改成功！",
	})
}

func SexnovelContentDelete(c *gin.Context) {
	id := c.Query("id")
	rid, err := strconv.Atoi(id)
	if err != nil || rid <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "参数错误",
		})
		return
	}
	model := sexnovelContentModel.SelectByID(uint(rid))
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

// 新增标签
func SexnovelListLabelCreate(c *gin.Context) {
	request := sexnovel.SexnovelListLabel{}
	c.ShouldBindJSON(&request)
	request.Save()
	c.JSON(http.StatusOK, gin.H{
		"message": "创建成功！",
		"data":    request,
	})
}

// 删除视频关联作者
func SexnovelListLabelDelete(c *gin.Context) {
	id := c.Query("id")
	rid, err := strconv.Atoi(id)
	if err != nil || rid <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "参数错误",
		})
		return
	}
	model := sexnovel.SexnovelListLabel{}
	model.ID = uint(rid)
	model.Delete()
	c.JSON(http.StatusOK, gin.H{
		"message": "删除成功！",
	})
}
