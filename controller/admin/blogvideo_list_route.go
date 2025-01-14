package admin

import (
	"bytes"
	"crypto/rand"
	"io"
	"myadmin/model/blog"
	"myadmin/util"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var blogVideoModel blog.BlogVideo

// 后台创建与更新视频地址 同一接口处理 一个博客对应一个视频
func BlogVideoCreate(c *gin.Context) {
	var request blog.BlogVideo
	c.ShouldBindJSON(&request)
	if request.BlogID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "必须绑定博客同时上传",
		})
		return
	}
	result := blogVideoModel.SelectByBlogID(request.BlogID)
	if result != nil && result.ID != 0 {
		request.Delete()
	}
	result.Save() //保存新的
	c.JSON(http.StatusOK, gin.H{
		"message": "更新博客视频成功！",
	})
}

//	func BlogVideoDelete(c *gin.Context) {
//		id := c.Query("id")
//		rid, _ := strconv.Atoi(id)
//		var imageModel blog.BlogImage
//		model := imageModel.SelectByID(uint(rid))
//		if model == nil {
//			c.JSON(http.StatusBadRequest, gin.H{
//				"message": "视频不存在",
//			})
//			return
//		}
//		model.Delete()
//		util.OssDeleteObject(model.Path)
//		c.JSON(http.StatusOK, gin.H{
//			"message": "删除成功！",
//		})
//	}
func BlogVideoSubmitJob(c *gin.Context) {
	id := c.Query("id")
	rid, err := strconv.Atoi(id)
	if err != nil || rid <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "参数错误",
		})
		return
	}
	model := blog.BlogVideoModel.SelectByID(uint(rid))
	if model == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "博客不存在",
		})
		return
	}
	blogSubmitJob(model)
	model.Update()
	c.JSON(http.StatusOK, gin.H{
		"message": "提交成功！",
	})
}
func blogSubmitJob(model *blog.BlogVideo) {
	// 提交转码。
	hlspath := "bloghls/" + strconv.Itoa(int(model.ID)) + "/" + strconv.Itoa(int(time.Now().Unix())) + "/index"
	hlskey := make([]byte, 16)
	rand.Read(hlskey)
	hlsinput := util.MtsSubmitJobsInput{
		Bucket:   os.Getenv("ALI_OSS_ORIGIN"),
		Location: os.Getenv("ALI_OSS_REGION"),
		Object:   model.OssName,
	}
	jobID, jobErr := util.MtsSubmitJobs(hlsinput, hlspath, hlskey)
	if jobErr != nil {
		model.JobErr = time.Now().Format("2006-01-02 15:04:05 || 回调提交") + jobErr.Error() + "\r\n"
		model.JobStatus = 2
	} else {
		model.JobID = jobID
		model.JobStatus = 1
	} // 提交成功
}

func BlogCoverSubmitJob(c *gin.Context) {
	id := c.Query("id")
	rid, err := strconv.Atoi(id)
	if err != nil || rid <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "参数错误",
		})
		return
	}
	model := blog.BlogVideoModel.SelectByID(uint(rid))
	if model == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "博客不存在",
		})
		return
	}
	meta, err := util.OssGetObjectMeta(model.OssName)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "保存失败",
		})
		return
	}
	size, _ := strconv.Atoi(meta.Get("Content-Length"))
	model.OssSize = uint(size)
	if model.Cover == "" {
		if url, err := util.OssGetUrl(model.OssName, oss.Process("video/snapshot,t_5000,f_jpg,w_400,h_0")); err == nil {
			if resp, err := http.Get(url); err == nil {
				if resp.StatusCode == 200 {
					if httpData, err := io.ReadAll(resp.Body); err == nil {
						nameCover := uuid.New().String()
						namePath := "blogcover/" + strconv.Itoa(int(model.ID)) + "/" + nameCover + ".jpg"
						key := util.MD5Byte(nameCover)
						encByte := util.CBCEncrypter(key, key, httpData)
						util.OssPutObject(namePath, bytes.NewBuffer(encByte))
						model.Cover = namePath
					}
				}
			}
		}
	}
	model.Update()
	c.JSON(http.StatusOK, gin.H{
		"message": "提交成功！",
		"cover":   model.Cover,
	})
}
