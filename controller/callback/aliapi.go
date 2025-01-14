package callback

import (
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"log"
	"math/rand"
	"myadmin/model/blog"
	"myadmin/model/vlog"
	"myadmin/model/vod"
	"myadmin/util"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

var vodModel vod.VodList

type CompleteMultipartObject struct {
	Key  string `json:"key"`
	Size uint   `json:"size"`
}

type CompleteMultipartOss struct {
	Object  CompleteMultipartObject `json:"object"`
	Version string                  `json:"ossSchemaVersion"`
}

type CompleteMultipartEvents struct {
	Oss       CompleteMultipartOss `json:"oss"`
	EventName string               `json:"eventName"`
}

type CompleteMultipartUpload struct {
	Events []CompleteMultipartEvents `json:"events"`
}

// 这里处理文件大文件上传的回调
func OssUpload(c *gin.Context) {
	reError := func(msg string) {
		log.Println(msg)
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    0,
			"message": msg,
		})
	}

	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		reError("oss回调读取数据出错")
		return
	}

	enc := base64.StdEncoding
	dbuf := make([]byte, enc.DecodedLen(len(body)))
	n, err := enc.Decode(dbuf, body)
	if err != nil {
		reError("oss回调base64解析失败")
		return
	}

	var callback CompleteMultipartUpload
	json.Unmarshal(dbuf[:n], &callback)
	var events CompleteMultipartEvents
	if len(callback.Events) == 1 {
		events = callback.Events[0]
	}
	if len(events.Oss.Object.Key) < 1 {
		reError("oss回调数据错误")
		return
	}
	vod := vodModel.SelectByOssName(events.Oss.Object.Key)
	if vod == nil {
		reError("oss回调数据没找到" + events.Oss.Object.Key)
		return
	}
	if vod.OssSize != events.Oss.Object.Size {
		reError("oss回调文件size对比错误")
		return
	}
	if vod.JobStatus != 0 {
		reError("vod job status 状态错误")
		return
	}

	// 提交转码任务 ------------------
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
		vod.JobErr = time.Now().Format("2006-01-02 15:04:05 || 回调提交") + jobErr.Error() + "\r\n"
		vod.JobStatus = 2
	} else {
		vod.JobID = jobID
		vod.JobStatus = 1
	}
	// ========
	err = vod.Save(true)
	if err != nil {
		log.Println("保存转码信息失败", err)
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    1,
		"message": "文件上传完成",
	})
}

type MtsMessage struct {
	JobID string
	Type  string
}

type MtsSubscription struct {
	Message    string
	MessageMD5 string
}

func MtsNotify(c *gin.Context) {
	reError := func(msg string) {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    0,
			"message": msg,
		})
	}
	var notify MtsSubscription
	err := c.ShouldBindJSON(&notify)
	if err != nil {
		reError(err.Error())
		return
	}
	var message MtsMessage
	err = json.Unmarshal([]byte(notify.Message), &message)
	if err != nil {
		reError(err.Error())
		return
	}
	if message.JobID == "" {
		reError("jobid不正确。")
		return
	}
	// 按照顺序查找任务所属的类型。
	if vod := vod.VodListModel.SelectJobID(message.JobID); vod != nil {
		if vod.JobStatus != 1 {
			reError("工作状态不对")
		}
		// 查询状态。
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

		// 这里填写内容、
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

	} else if vlog := vlog.VlogListModel.SelectJobID(message.JobID); vlog != nil {

		if vlog.JobStatus != 1 {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "转码任务未在进行中",
			})
			return
		}
		if vlog.JobID == "" {
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
		response, err := util.MtsQueryJob(vlog.JobID)
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
		vlog.HlsFps = uint(fps)
		bitrate, _ := strconv.ParseUint(job.Output.Properties.Bitrate, 10, 32)
		vlog.HlsBitrate = uint(bitrate)
		width, _ := strconv.ParseUint(job.Output.Properties.Width, 10, 32)
		vlog.HlsWidth = uint(width)
		height, _ := strconv.ParseUint(job.Output.Properties.Height, 10, 32)
		vlog.HlsHeight = uint(height)
		filesize, _ := strconv.ParseUint(job.Output.Properties.FileSize, 10, 32)
		vlog.HlsSize = uint(filesize)
		duration, _ := strconv.ParseUint(job.Output.Properties.Duration, 10, 32)
		vlog.HlsDuration = uint(duration)
		// 播放相关的
		vlog.HlsPath = job.Output.OutputFile.Object
		vlog.HlsKey, err = base64.StdEncoding.DecodeString(job.Output.Encryption.Key)
		if err != nil {
			reError("解码错误")
			return
		}
		// 这里填写内容
		// vlog.FormatM3u8()
		vlog.JobStatus = 3
		vlog.Update()
		c.JSON(http.StatusOK, gin.H{
			"message": "更新成功！",
		})
	} else if blog := blog.BlogVideoModel.SelectJobID(message.JobID); blog != nil { //博客视频
		if blog.JobStatus != 1 {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "转码任务未在进行中",
			})
			return
		}
		if blog.JobID == "" {
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
		response, err := util.MtsQueryJob(blog.JobID)
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
		blog.HlsFps = uint(fps)
		bitrate, _ := strconv.ParseUint(job.Output.Properties.Bitrate, 10, 32)
		blog.HlsBitrate = uint(bitrate)
		width, _ := strconv.ParseUint(job.Output.Properties.Width, 10, 32)
		blog.HlsWidth = uint(width)
		height, _ := strconv.ParseUint(job.Output.Properties.Height, 10, 32)
		blog.HlsHeight = uint(height)
		filesize, _ := strconv.ParseUint(job.Output.Properties.FileSize, 10, 32)
		blog.HlsSize = uint(filesize)
		duration, _ := strconv.ParseUint(job.Output.Properties.Duration, 10, 32)
		blog.HlsDuration = uint(duration)
		// 播放相关的
		blog.HlsPath = job.Output.OutputFile.Object
		blog.HlsKey, err = base64.StdEncoding.DecodeString(job.Output.Encryption.Key)
		if err != nil {
			reError("解码错误")
			return
		}
		// 这里填写内容
		blog.JobStatus = 3
		blog.Update()
		c.JSON(http.StatusOK, gin.H{
			"message": "更新成功！",
		})
	} else {
		log.Println("未定义的jobid" + message.JobID)
	}
}
