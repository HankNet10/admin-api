package admin

import (
	"crypto/rand"
	"encoding/base64"
	"myadmin/model/vod"
	"myadmin/util"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

var vodListModel vod.VodList

func VodListList(c *gin.Context) {
	request := vod.VodListParam{}
	c.BindQuery(&request)
	request.Order = "-id"
	list, total := vodListModel.List(request)
	c.JSON(http.StatusOK, gin.H{
		"list":  list,
		"total": total,
	})
}

func VodListCreate(c *gin.Context) {
	request := vod.VodList{}
	c.ShouldBindJSON(&request)
	request.Save(true)
	c.JSON(http.StatusOK, gin.H{
		"message": "创建成功！",
	})
}

// 如果文件上传oss但是没提交到数据库
func VodListExist(c *gin.Context) {
	request := vod.VodList{}
	c.ShouldBindJSON(&request)
	vod := vod.VodListModel.SelectByOssName(request.OssName)
	if vod != nil {
		c.JSON(http.StatusOK, gin.H{
			"message": "文件存在",
		})
		return
	}
	request.Save(true)
	c.JSON(http.StatusOK, gin.H{
		"message": "补充成功！",
	})
}

func VodListUpdate(c *gin.Context) {
	request := vod.VodList{}
	c.ShouldBindJSON(&request)

	request.Update(true)
	c.JSON(http.StatusOK, gin.H{
		"message": "修改成功！",
	})
}

func VodListDelete(c *gin.Context) {
	id := c.Query("id")
	rid, err := strconv.Atoi(id)
	if err != nil || rid <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "参数错误",
		})
		return
	}
	model := vodListModel.SelectByID(uint(rid))
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

// 新增标签
func VodListLabelCreate(c *gin.Context) {
	request := vod.VodListLabel{}
	c.ShouldBindJSON(&request)
	request.Save()
	c.JSON(http.StatusOK, gin.H{
		"message": "创建成功！",
		"data":    request,
	})
}

// 删除视频关联作者
func VodListLabelDelete(c *gin.Context) {
	id := c.Query("id")
	rid, err := strconv.Atoi(id)
	if err != nil || rid <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "参数错误",
		})
		return
	}
	model := vod.VodListLabel{}
	model.ID = uint(rid)
	model.Delete()
	c.JSON(http.StatusOK, gin.H{
		"message": "删除成功！",
	})
}

// 新增作者
func VodListUserCreate(c *gin.Context) {
	request := vod.VodListUser{}
	c.ShouldBindJSON(&request)
	request.Save()
	c.JSON(http.StatusOK, gin.H{
		"message": "创建成功！",
		"data":    request,
	})
}

// 删除视频关联作者
func VodListUserDelete(c *gin.Context) {
	id := c.Query("id")
	rid, err := strconv.Atoi(id)
	if err != nil || rid <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "参数错误",
		})
		return
	}
	model := vod.VodListUser{}
	model.ID = uint(rid)
	model.Delete()
	c.JSON(http.StatusOK, gin.H{
		"message": "删除成功！",
	})
}

// 重新提交转码任务
func VodSubmitJob(c *gin.Context) {
	id := c.Query("id")
	rid, err1 := strconv.Atoi(id)
	if err1 != nil || rid <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "参数错误",
		})
		return
	}
	vod := vodListModel.SelectByID(uint(rid))
	if vod == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "视频数据不存在",
		})
		return
	}
	if len(vod.OssName) < 1 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "原始视频文件不存在",
		})
		return
	}
	if vod.JobStatus == 1 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "转码任务正在进行中",
		})
		return
	}
	// 提交转码任务
	hlspath := "vodhls/" + strconv.Itoa(int(vod.ID)) + "/" + strconv.Itoa(int(time.Now().Unix())) + "/index"
	hlskey := make([]byte, 16)
	rand.Read(hlskey)
	hlsinput := util.MtsSubmitJobsInput{
		Bucket:   os.Getenv("ALI_OSS_ORIGIN"),
		Location: os.Getenv("ALI_OSS_REGION"),
		Object:   vod.OssName,
	}
	jobID, jobErr := util.MtsSubmitJobs(hlsinput, hlspath, hlskey)
	if jobErr != nil {
		vod.JobErr = time.Now().Format("2006-01-02 15:04:05 || 手动提交") + jobErr.Error() + "\r\n"
		vod.JobStatus = 2
	} else {
		vod.JobID = jobID
		vod.JobStatus = 1
	}
	err := vod.Save(true)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    1,
			"message": "保存转码信息失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    1,
		"message": "提交转码任务",
	})
}

// 获取转码结果
func VodJobResult(c *gin.Context) {
	id := c.Query("id")
	rid, err := strconv.Atoi(id)
	if err != nil || rid <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "参数错误",
		})
		return
	}
	vod := vodListModel.SelectByID(uint(rid))
	if vod == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "视频数据不存在",
		})
		return
	}
	if vod.JobStatus != 1 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "转码任务未在进行中",
		})
		return
	}
	if vod.JobID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "jobid不正确。",
		})
		return
	}

	// 查询状态。
	reError := func(msg string) {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    0,
			"message": msg,
		})
	}
	response, err := util.MtsQueryJob(vod.JobID)
	if err != nil {
		reError(err.Error())
		return
	}

	if len(response.JobList.Job) != 1 {
		reError("返回了预料之外的数据")
		return
	}
	job := response.JobList.Job[0]
	// 需要判断转码任务是否在进行中
	if job.State != "TranscodeSuccess" {
		reError("转码状态为" + job.State + " 进度为" + strconv.Itoa(int(job.Percent)))
		return
	}

	// 处理媒体的详细信息
	fps, _ := strconv.ParseUint(job.Output.Properties.Fps, 10, 32)
	vod.HlsFps = uint(fps)
	bitrate, _ := strconv.ParseUint(job.Output.Properties.Bitrate, 10, 32)
	vod.HlsBitrate = uint(bitrate)
	width, _ := strconv.ParseUint(job.Output.Properties.Width, 10, 32)
	vod.HlsWidth = uint(width)
	height, _ := strconv.ParseUint(job.Output.Properties.Height, 10, 32)
	vod.HlsHeight = uint(height)
	filesize, _ := strconv.ParseUint(job.Output.Properties.FileSize, 10, 32)
	vod.HlsSize = uint(filesize)
	duration, _ := strconv.ParseUint(job.Output.Properties.Duration, 10, 32)
	vod.HlsDuration = uint(duration)
	// 播放相关的
	vod.HlsPath = job.Output.OutputFile.Object
	vod.HlsKey, err = base64.StdEncoding.DecodeString(job.Output.Encryption.Key)
	if err != nil {
		reError("解码错误")
		return
	}
	// 这里填写内容
	// vod.FormatM3u8()

	vod.JobStatus = 3
	err = vod.Save(true)
	if err != nil {
		reError(err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    1,
		"message": "校验完成",
	})
}
