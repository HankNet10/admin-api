package util

import (
	"encoding/base64"
	"image/color"
	"math/rand"
	"strconv"

	// "myadmin/util/captcha"
	"myadmin/util/redis"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/penndev/gopkg/captcha"
)

type ModelCaptcha struct {
	CaptchaId  string
	CaptchaUrl string
}

// @Summary 生成验证码
// @Schemes 生成用户验证码
// @Description 返回base64验证码,与验证码id。
// @Tags 工具
// @Accept json
// @Produce json
// @Success 200 {object} ModelCaptcha
// @Router /util/captcha [get]
func Captcha(c *gin.Context) {
	rand.Seed(time.Now().UnixMicro())
	code := rand.Intn(9000) + 1000

	option := captcha.Option{
		Width:     100,
		Height:    40,
		DPI:       90,
		Text:      strconv.Itoa(code),
		FontSize:  20,
		TextColor: color.RGBA{0, 0, 0, 255},
	}
	buf, err := captcha.NewPngImg(option)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}
	data := base64.StdEncoding.EncodeToString(buf.Bytes())
	genId := uuid.New().String()
	redis.Set("captcha-"+genId, code, time.Minute*5)
	c.JSON(http.StatusOK, ModelCaptcha{
		CaptchaUrl: "data:image/png;base64," + data,
		CaptchaId:  genId,
	})

}
