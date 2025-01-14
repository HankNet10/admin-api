package router

import (
	"bytes"
	"encoding/gob"
	"log"
	"myadmin/model/admin"
	"myadmin/util/jwt"
	"myadmin/util/redis"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// 处理允许跨域。
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", os.Getenv("GIN_CORS"))
		c.Header("Access-Control-Allow-Methods", "GET, HEAD, POST, PUT, DELETE, CONNECT, OPTIONS, TRACE, PATCH")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "Access-Control-Allow-Headers, Origin, Accept, X-Requested-With, Content-Type, Access-Control-Request-Method, Access-Control-Request-Headers, x-token, x-appid")
		if c.Request.Method == "OPTIONS" {
			c.Header("Access-Control-Max-Age", "259200")
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}

// Admin
// 检验admin用户后台登录权限
func JWTAdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		abort := func(s string) {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": s,
			})
			c.Abort()
		}
		token := c.Request.Header.Get("x-token")
		if token == "" {
			abort("需要进行身份验证")
			return
		}
		// 进行身份验证
		claims, err := jwt.DecondeToken(token)
		if err != nil {
			abort(err.Error())
			return
		}
		sub, ok := claims["sub"]
		if !ok {
			abort("错误的token [not sub]")
			return
		}
		username, ok := sub.(string)
		if !ok {
			abort("错误的token [not sub]")
			return
		}
		c.Set("AuthAdmin", username)
		c.Next()
	}
}

// Admin
// 处理账户鉴权,获取用户数据库信息
func AuthorizeAdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		abort := func(s string) {
			c.JSON(http.StatusForbidden, gin.H{
				"message": s,
			})
			c.Abort()
		}
		name := c.MustGet("AuthAdmin").(string)
		var au admin.SysUser
		sysUser := au.SelectUserByName(name)
		if sysUser == nil {
			abort("用户已经被清理，请联系管理员！")
			return
		}
		if sysUser.Status != 1 {
			abort("用户已经锁定，请联系管理员！")
			return
		}

		if sysUser.RoleID > 0 {
			//这里开始进行权限验证。
			result := sysUser.Authorize(c.Request.Method, c.Request.URL.Path)
			if !result {
				abort("请联系管理员授权！")
				return
			}
		}
		c.Set("AdminModel", sysUser)
		c.Next()
	}
}

// 处理用户登录。
func JWTUserMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		abort := func(s string) {
			c.JSON(http.StatusOK, gin.H{
				"code":    403,
				"data":    nil,
				"message": s,
			})
			c.Abort()
		}
		token := c.Request.Header.Get("x-token")
		if token == "" {
			abort("需要进行身份验证")
			return
		}
		// 进行身份验证
		claims, err := jwt.DecondeToken(token)
		if err != nil {
			abort(err.Error())
			return
		}
		sub, ok := claims["sub"]
		if !ok {
			abort("错误的token [not sub]")
			return
		}
		userid, ok := sub.(string)
		if !ok {
			abort("错误的token [not sub]")
			return
		}
		c.Set("UserID", userid)
		c.Next()
	}
}

// 处理用户是否登录验证
func JWTIsUserMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("UserID", "0")
		token := c.Request.Header.Get("x-token")
		if claims, err := jwt.DecondeToken(token); err == nil {
			if sub, ok := claims["sub"]; ok {
				if userid, ok := sub.(string); ok {
					c.Set("UserID", userid)
				}
			}
		}
		c.Next()
	}
}

type cacheHttpResponse struct {
	Header http.Header
	Status int
	Body   []byte
}

type cacheGinWrite struct {
	gin.ResponseWriter
	body bytes.Buffer
}

func (w *cacheGinWrite) Write(data []byte) (int, error) {
	w.body.Write(data)
	return w.ResponseWriter.Write(data)
}

// Redis 缓存中间件
func CacheGetResult(redisT time.Duration) gin.HandlerFunc {

	return func(c *gin.Context) {
		if _, ok := c.Get("UserID"); ok {
			c.Next()
			return
		}

		appID := c.Request.Header.Get("x-appid")
		if appID == "" {
			appID = "1"
		}

		uKey := strings.ReplaceAll(c.Request.URL.String(), "/", ":")
		cacheKey := "middle-cgr-" + appID + uKey
		var cResp cacheHttpResponse
		if err := redis.Deserialize(cacheKey, &cResp); err == nil {
			c.Status(cResp.Status)
			for key, values := range cResp.Header {
				for _, val := range values {
					c.Writer.Header().Set(key, val)
				}
			}
			c.Writer.Write(cResp.Body)
			c.Abort()
			return
		}
		ginwrite := cacheGinWrite{ResponseWriter: c.Writer}
		c.Writer = &ginwrite
		c.Next()

		resp := cacheHttpResponse{
			Header: c.Writer.Header().Clone(),
			Status: c.Writer.Status(),
			Body:   ginwrite.body.Bytes(),
		}
		if err := redis.Serialize(cacheKey, resp, redisT); err != nil {
			log.Panic("cache middle", err)
		}
	}
}

func Serialize(value interface{}) ([]byte, error) {
	var b bytes.Buffer
	encoder := gob.NewEncoder(&b)
	if err := encoder.Encode(value); err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

func Deserialize(payload []byte, ptr interface{}) (err error) {
	return gob.NewDecoder(bytes.NewBuffer(payload)).Decode(ptr)
}
