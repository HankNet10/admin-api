package admin

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"io"
	"myadmin/model/vlog"
	"myadmin/util"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func VlogListCreate(c *gin.Context) {
	request := vlog.VlogList{}
	c.ShouldBindJSON(&request)
	request.Save()

	c.JSON(http.StatusOK, gin.H{
		"message": "创建成功！",
	})
}

func VlogListList(c *gin.Context) {
	request := vlog.VlogListParam{}
	c.BindQuery(&request)
	if request.Limit > 100 || request.Limit < 1 {
		request.Limit = 20
	}
	list, total := vlog.VlogListModel.List(request)
	c.JSON(http.StatusOK, gin.H{
		"list":  list,
		"total": total,
	})
}

func VlogListUpdate(c *gin.Context) {
	var request vlog.VlogList
	c.ShouldBindJSON(&request)
	model := vlog.VlogListModel.SelectByID(request.ID)
	if model == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "数据不存在",
		})
		return
	}
	model.Status = request.Status
	model.Title = request.Title
	model.UserID = request.UserID
	// 判断是否是新的ossname 如果是新的要重新转码。
	if request.OssName != model.OssName {
		// 验证文件是否存在
		meta, err := util.OssGetObjectMeta(request.OssName)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "保存失败",
			})
			return
		}
		model.OssName = request.OssName
		size, _ := strconv.Atoi(meta.Get("Content-Length"))
		model.OssSize = uint(size)
		vlogSubmitJob(model)
	}
	if request.Cover == "" {
		if url, err := util.OssGetUrl(model.OssName, oss.Process("video/snapshot,t_5000,f_jpg,w_400,h_0")); err == nil {
			if resp, err := http.Get(url); err == nil {
				if resp.StatusCode == 200 {
					if httpData, err := io.ReadAll(resp.Body); err == nil {
						nameCover := uuid.New().String()
						namePath := "vlogcover/" + strconv.Itoa(int(model.ID)) + "/" + nameCover + ".jpg"
						key := util.GenAES128Key(nameCover)
						encByte := util.CBCEncrypter(key, key, httpData)
						util.OssPutObject(namePath, bytes.NewBuffer(encByte))
						model.Cover = namePath
					}
				}
			}
		}
	}
	model.Update()
	if model.Status == 0 {
		model.StatusDown()
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "修改成功！",
	})
}

func VlogListDelete(c *gin.Context) {
	id := c.Query("id")
	rid, err := strconv.Atoi(id)
	if err != nil || rid <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "参数错误",
		})
		return
	}
	model := vlog.VlogListModel.SelectByID(uint(rid))
	if model == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "博客不存在",
		})
		return
	}
	model.Delete()
	c.JSON(http.StatusOK, gin.H{
		"message": "删除成功！",
	})
}

func vlogSubmitJob(model *vlog.VlogList) {
	// 提交转码。
	hlspath := "vloghls/" + strconv.Itoa(int(model.ID)) + "/" + strconv.Itoa(int(time.Now().Unix())) + "/index"
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

func VlogSubmitJob(c *gin.Context) {
	id := c.Query("id")
	rid, err := strconv.Atoi(id)
	if err != nil || rid <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "参数错误",
		})
		return
	}
	model := vlog.VlogListModel.SelectByID(uint(rid))
	if model == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "博客不存在",
		})
		return
	}
	vlogSubmitJob(model)
	model.Update()
	c.JSON(http.StatusOK, gin.H{
		"message": "提交成功！",
	})
}

func VlogJobResult(c *gin.Context) {
	id := c.Query("id")
	rid, err := strconv.Atoi(id)
	if err != nil || rid <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "参数错误",
		})
		return
	}
	model := vlog.VlogListModel.SelectByID(uint(rid))
	if model == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "博客不存在",
		})
		return
	}
	// 检查结果是否正确。
	if model.JobStatus != 1 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "转码任务未在进行中",
		})
		return
	}
	if model.JobID == "" {
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
	response, err := util.MtsQueryJob(model.JobID)
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
	model.HlsFps = uint(fps)
	bitrate, _ := strconv.ParseUint(job.Output.Properties.Bitrate, 10, 32)
	model.HlsBitrate = uint(bitrate)
	width, _ := strconv.ParseUint(job.Output.Properties.Width, 10, 32)
	model.HlsWidth = uint(width)
	height, _ := strconv.ParseUint(job.Output.Properties.Height, 10, 32)
	model.HlsHeight = uint(height)
	filesize, _ := strconv.ParseUint(job.Output.Properties.FileSize, 10, 32)
	model.HlsSize = uint(filesize)
	duration, _ := strconv.ParseUint(job.Output.Properties.Duration, 10, 32)
	model.HlsDuration = uint(duration)
	// 播放相关的
	model.HlsPath = job.Output.OutputFile.Object
	model.HlsKey, err = base64.StdEncoding.DecodeString(job.Output.Encryption.Key)
	if err != nil {
		reError("解码错误")
		return
	}
	// 这里填写内容
	// model.FormatM3u8()
	model.JobStatus = 3

	model.Update()
	c.JSON(http.StatusOK, gin.H{
		"message": "更新成功！",
	})
}
