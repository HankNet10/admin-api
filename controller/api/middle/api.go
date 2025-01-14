package middle

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"myadmin/controller/api/bind"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type ginWrite struct {
	key [16]byte
	iv  [16]byte
	gin.ResponseWriter
}

func (w *ginWrite) Write(body []byte) (int, error) {
	plaintext := body
	key := w.key[:]
	iv := w.iv[:]
	block, err := aes.NewCipher(key)
	if err != nil {
		return 0, err
	}

	var paddedPlaintext []byte
	padding := block.BlockSize() - len(plaintext)%block.BlockSize()
	if padding == 0 {
		paddedPlaintext = []byte(plaintext)
	} else { //PKCS7填充
		padtext := bytes.Repeat([]byte{byte(padding)}, padding)
		paddedPlaintext = append([]byte(plaintext), padtext...)
	}
	blockMode := cipher.NewCBCEncrypter(block, iv)
	aesPlainText := make([]byte, len(paddedPlaintext))
	blockMode.CryptBlocks(aesPlainText, paddedPlaintext)

	m := bind.Message{
		Data: base64.StdEncoding.EncodeToString(aesPlainText),
	}

	jsonMsg, err := json.Marshal(m)
	if err != nil {
		return 0, err
	}
	return w.ResponseWriter.Write(jsonMsg)
}

// 处理用户是否登录验证
func Api() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 验证签名
		queryList := strings.Split(fmt.Sprintln(c.Request.URL), "&sign=")
		if len(queryList) != 2 {
			c.JSON(http.StatusForbidden, bind.Message{Message: "fail(1)"})
			c.Abort()
			return
		}

		// 验证app key
		appKey := c.GetHeader("x-appid")
		if appKey == "" {
			c.JSON(http.StatusForbidden, bind.Message{Message: "fail(3)"})
			c.Abort()
			return
		}

		// 真实验证
		c.Set("appKey", appKey)
		//跳过加密
		istest := c.GetHeader("x-test")
		if istest != "1" {
			queryBase64 := make([]byte, base64.StdEncoding.EncodedLen(len(queryList[0])))
			base64.StdEncoding.Encode(queryBase64, []byte(queryList[0]))
			sign := sha256.Sum256(queryBase64)
			signHex := hex.EncodeToString(sign[4:10])
			if signHex != strings.TrimSpace(queryList[1]) {
				c.JSON(http.StatusForbidden, bind.Message{Message: "fail(2)"})
				c.Abort()
				return
			}

			// 加密
			writer := &ginWrite{
				key:            md5.Sum([]byte(signHex)),
				iv:             md5.Sum([]byte(appKey)),
				ResponseWriter: c.Writer,
			}
			c.Writer = writer
		}
	}
}
