package util

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"

	"github.com/gin-gonic/gin"
)

//
func GenAES128Key(id string) []byte {
	// 计算文件的aes解密的 key 和 iv
	sha1Enc := sha1.New()
	sha1Enc.Write([]byte(id))
	sha1_hash := hex.EncodeToString(sha1Enc.Sum(nil))
	md5Enc := md5.New()
	md5Enc.Write([]byte(sha1_hash))
	md5Key := md5Enc.Sum(nil)
	// 对文本内容进行加密
	// contentEncry := CBCEncrypter(md5Key, md5Key, []byte(content))
	return md5Key
}
func MD5(v string) string {
	d := []byte(v)
	m := md5.New()
	m.Write(d)
	return hex.EncodeToString(m.Sum(nil))
}
func MD5Byte(v string) []byte {
	d := []byte(v)
	m := md5.New()
	m.Write(d)
	return m.Sum(nil)
}

// param key 必须为 16 24 32 对应 128 192 256位加密
// param iv 必须为 16 byte大小
// text 加密的字节内容
// return 返回加密后的数据
func CBCEncrypter(key, iv, text []byte) []byte {
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}
	// 填充文本内容
	paddText := PKCS7Padding(text, block.BlockSize())
	blockMode := cipher.NewCBCEncrypter(block, iv)

	// 加密
	result := make([]byte, len(paddText))
	blockMode.CryptBlocks(result, paddText)
	// 返回密文
	return result
}

// PKCS7Padding 填充模式
func PKCS7Padding(text []byte, blockSize int) []byte {
	// 计算待填充的长度
	padding := blockSize - len(text)%blockSize
	var paddingText []byte
	if padding == 0 {
		// 已对齐，填充一整块数据，每个数据为 blockSize
		paddingText = bytes.Repeat([]byte{byte(blockSize)}, blockSize)
	} else {
		// 未对齐 填充 padding 个数据，每个数据为 padding
		paddingText = bytes.Repeat([]byte{byte(padding)}, padding)
	}
	return append(text, paddingText...)
}

// 返回gin请求的完整域名部分 例 https://baidu.com
func GetGinRequestHostUrl(c *gin.Context) string {
	schemes := "http"
	if c.Request.Header.Get("X-Forwarded-Proto") != "" {
		schemes = c.Request.Header.Get("X-Forwarded-Proto")
	}
	return schemes + "://" + c.Request.Host
}
