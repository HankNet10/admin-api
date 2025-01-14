package admin

import (
	"myadmin/model/admin"
	"myadmin/util/bcrypt"
	"myadmin/util/jwt"
	"myadmin/util/redis"
	"net/http"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// User login structure
type loginInput struct {
	Username  string // 用户名
	Password  string // 密码
	Captcha   string // 验证码
	CaptchaId string // 验证码ID
}

func Login(c *gin.Context) {
	var request loginInput
	c.ShouldBindJSON(&request)
	// 创建验证码
	code, err := redis.Pull("captcha-" + request.CaptchaId)
	if err != nil || code != request.Captcha {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "验证码错误",
		})
		return
	}
	//用户验证失败返回
	resultErr := func() {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "用户名或密码错误",
		})
	}

	// 查询账号是否存在
	var au admin.SysUser
	res := au.SelectUserByName(request.Username)
	//初始化admin的初始管理员。
	if res == nil && request.Username == "admin" {
		res = &admin.SysUser{
			Name:     request.Username,
			NickName: request.Username,
			Password: bcrypt.GeneratePassword(request.Password),
			Status:   1,
		}
		res.Save()
	}
	//如果用户还是不存在
	if res == nil {
		resultErr()
		return
	}
	//比较验证密码
	if !bcrypt.ComparePassword(request.Password, res.Password) {
		resultErr()
		return
	}
	// 验证通过生成JWT
	exp := time.Now().AddDate(0, 0, 7).Unix()
	token, err := jwt.EncodeToken(res.Name, exp)
	if err != nil {
		resultErr()
		return
	}

	c.JSON(http.StatusOK, jwt.AccessToken{
		AccessToken: token,
		TokenType:   jwt.TokenType,
		ExpiresIn:   exp,
	})
}

func Account(c *gin.Context) {
	name := c.MustGet("AuthAdmin").(string)
	var au admin.SysUser
	result := au.SelectUserByName(name)
	if result == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "账号异常",
		})
		return
	}
	c.JSON(http.StatusOK, result)
}

// User login structure
type changePasswordInput struct {
	Password    string // 密码
	NewPassword string // 验证码
}

func ChangePassword(c *gin.Context) {
	// 获取登录用户。
	name := c.MustGet("AuthAdmin").(string)
	var au admin.SysUser
	result := au.SelectUserByName(name)
	if result == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "账号异常",
		})
		return
	}
	// 对比密码，判断老密码是否正常。
	var request changePasswordInput
	c.ShouldBindJSON(&request)
	if !bcrypt.ComparePassword(request.Password, result.Password) {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "验证原密码错误，请检查。",
		})
		return
	}
	//
	hashed := bcrypt.GeneratePassword(request.NewPassword)
	result.Password = hashed
	result.Save()
	c.JSON(http.StatusOK, gin.H{
		"message": "修改完成",
	})
}

func AdminList(c *gin.Context) {
	// name := c.MustGet("AuthAdmin").(string)
	asu := admin.SysUserParan{}
	c.BindQuery(&asu)

	var au admin.SysUser
	result := au.List(asu)
	c.JSON(http.StatusOK, gin.H{
		"list":  result,
		"total": len(result),
	})
}

var adminModel admin.SysUser

type adminInput struct {
	ID       uint
	Name     string
	NickName string
	Password string
	RoleID   uint
	Email    string
	Status   uint8
	Remark   string
}

func AdminCreate(c *gin.Context) {
	var request adminInput
	c.ShouldBindJSON(&request)

	au := adminModel.SelectUserByName(request.Name)
	if au != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "用户已存在",
		})
		return
	}

	newAu := admin.SysUser{
		Name:     request.Name,
		NickName: request.NickName,
		Password: bcrypt.GeneratePassword(request.Password),
		RoleID:   request.RoleID,
		Email:    request.Email,
		Remark:   request.Remark,
		Status:   request.Status,
	}
	newAu.Save()
	c.JSON(http.StatusOK, gin.H{
		"message": "创建成功！",
	})
}

func AdminUpdate(c *gin.Context) {
	var request adminInput
	c.ShouldBindJSON(&request)

	au := adminModel.SelectUserByID(request.ID)
	if au == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "用户信息不存在",
		})
		return
	}

	if request.NickName != "" {
		au.NickName = request.NickName
	}

	if request.RoleID > 0 {
		au.RoleID = request.RoleID
	}

	if request.Email != "" {
		au.Email = request.Email
	}

	if request.Status > 0 {
		au.Status = request.Status
	}

	if request.Remark != "" {
		au.Remark = request.Remark
	}

	if request.Password != "" { //后台直接修改密码
		hashed := bcrypt.GeneratePassword(request.Password)
		au.Password = hashed
	}

	au.Save()
	c.JSON(http.StatusOK, gin.H{
		"message": "保存成功",
	})
}

func AdminDelete(c *gin.Context) {
	id := c.Query("id")
	rid, err := strconv.Atoi(id)
	if err != nil || rid <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "参数错误",
		})
		return
	}
	au := adminModel.SelectUserByID(uint(rid))
	if au == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "用户不存在",
		})
		return
	}
	au.Delete()

	c.JSON(http.StatusOK, gin.H{
		"message": "删除成功！",
	})
}

// 临时文件上传不重要文件 聊天图片
func AndminFileUoload(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}
	filepth := "../../adminupload/" + filepath.Base(file.Filename)
	if err := c.SaveUploadedFile(file, filepth); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"data": file.Filename,
	})
}
