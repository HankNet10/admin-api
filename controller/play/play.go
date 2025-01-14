package play

import (
	"encoding/hex"
	"errors"
	"io/ioutil"
	"myadmin/model"
	"myadmin/model/blog"
	"myadmin/model/vod"
	"myadmin/util"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func fileContent(oPath string) (string, error) {
	body, err := util.OssGetObject(oPath)
	if err != nil {
		return "", err
	}
	bodyByte, err := ioutil.ReadAll(body)
	if err != nil {
		return "", err
	}

	tsPath := os.Getenv("ALI_OSS_V_DOMAIN") + "/" + strings.Replace(oPath, util.PATHNAME, "", 1)
	newBody := strings.Replace(string(bodyByte), "index-", tsPath+"index-", -1)
	return newBody, nil
}

//
// ========================================================================================
//

func vodHlsM3u8(c *gin.Context) (string, error) {
	id := c.Param("id")
	vodDate := vod.VodList{}

	idd, err := strconv.Atoi(id)
	if err != nil {
		return "", errors.New("查询数据错误")
	}
	dresult := model.DataBase.Model(&vodDate).First(&vodDate, idd)
	if dresult.Error != nil {
		return "", errors.New("查询数据错误")
	}
	if vodDate.JobStatus != 3 {
		return "", errors.New("转码状态错误")
	}
	m3u8Content, err := fileContent(vodDate.HlsPath)
	if err != nil {
		return "", errors.New("获取源文件失败")
	}
	keypath := os.Getenv("PLIST_DOMAIN") + path.Dir(c.Request.RequestURI) + "/vod.enc" + "?accesskey=" + hex.EncodeToString(util.GenAES128Key(id))
	return strings.Replace(m3u8Content, "{enc.key}", keypath, 1), nil
}

// vod hls key 内容
func VodHlsKey(c *gin.Context) {
	id := c.Param("id")
	vodDate := vod.VodList{}
	idd, err := strconv.Atoi(id)
	if err != nil {
		c.String(http.StatusNotFound, "404")
		return
	}
	dresult := model.DataBase.Model(&vodDate).Select("HlsKey").First(&vodDate, idd)
	if dresult.Error != nil {
		c.String(http.StatusNotFound, "404")
		return
	}
	c.Writer.Write(vodDate.HlsKey)
}

// 处理 hls m3u8 原本数据
func VodHlsM3u8(c *gin.Context) {
	if m3u8Content, err := vodHlsM3u8(c); err != nil {
		c.String(http.StatusNotFound, "404")
	} else {
		c.String(http.StatusOK, m3u8Content)
	}
}

// 加密后的 hls m3u8 数据
func VodHlsM3u8Enc(c *gin.Context) {
	if m3u8Content, err := vodHlsM3u8(c); err != nil {
		c.String(http.StatusNotFound, "404")
	} else {
		hashKey := util.GenAES128Key(c.Param("id"))
		context := util.CBCEncrypter(hashKey, hashKey, []byte(m3u8Content))
		c.Header("Content-Type", "application/octet-stream")
		c.Writer.Write(context)
	}
}

//
// ========================================================================================
//

// blog 播放方法
func blogHlsM3u8(c *gin.Context) (string, error) {
	id := c.Param("id")

	blogData := blog.BlogVideo{}
	idd, err := strconv.Atoi(id)
	if err != nil {
		return "", errors.New("数据错误")
	}
	dresult := model.DataBase.Model(&blogData).First(&blogData, idd)
	if dresult.Error != nil {
		return "", errors.New("查询数据错误")
	}
	if blogData.JobStatus != 3 {
		return "", errors.New("转码状态错误")
	}
	m3u8Content, err := fileContent(blogData.HlsPath)
	if err != nil {
		return "", errors.New("获取源文件失败")
	}
	keypath := os.Getenv("PLIST_DOMAIN") + path.Dir(c.Request.RequestURI) + "/blog.enc" + "?accesskey=" + hex.EncodeToString(util.GenAES128Key(id))
	return strings.Replace(m3u8Content, "{enc.key}", keypath, 1), nil
}

// blod hls key 内容
func BlogHlsKey(c *gin.Context) {
	id := c.Param("id")
	blogData := blog.BlogVideo{}
	idd, err := strconv.Atoi(id)
	if err != nil {
		c.String(http.StatusNotFound, "404")
		return
	}
	dresult := model.DataBase.Model(&blogData).Select("HlsKey").First(&blogData, idd)
	if dresult.Error != nil {
		c.String(http.StatusNotFound, "404")
		return
	}
	c.Writer.Write(blogData.HlsKey)
}

// 处理hls播放 苹果设备不能加密的兼容问题。
func BlogHlsM3u8(c *gin.Context) {
	if m3u8Content, err := blogHlsM3u8(c); err != nil {
		c.String(http.StatusNotFound, "404")
	} else {
		c.String(http.StatusOK, m3u8Content)
	}
}

// 机密后的m3u8
func BlogHlsM3u8Enc(c *gin.Context) {
	if m3u8Content, err := blogHlsM3u8(c); err != nil {
		c.String(http.StatusNotFound, "404")
	} else {
		hashKey := util.GenAES128Key(c.Param("id"))
		context := util.CBCEncrypter(hashKey, hashKey, []byte(m3u8Content))
		c.Header("Content-Type", "application/octet-stream")
		c.Writer.Write(context)
	}
}
