package main_test

import (
	"fmt"
	"io"
	"myadmin/model"
	"myadmin/model/ad"
	"myadmin/model/admin"
	"myadmin/model/applicationad"
	"myadmin/model/attack"
	"myadmin/model/blog"
	"myadmin/model/config"
	"myadmin/model/seximg"
	"myadmin/model/sexnovel"
	"myadmin/model/suggest"
	"myadmin/model/user"
	"myadmin/model/vip"
	"myadmin/model/vlog"
	"myadmin/model/vod"
	"myadmin/util"
	"myadmin/util/sugar"
	"net/http"
	"os"
	"regexp"
	"testing"

	"github.com/feiin/sensitivewords"
	_ "github.com/joho/godotenv" // 引入.env 变量
)

func TestAutoMigrate1105(t *testing.T) {
	//评论限制长度1000
	model.DataBase.AutoMigrate(user.UserComment{})
	//评论限制长度1000
	model.DataBase.AutoMigrate(blog.BlogComment{})
	//评论限制长度1000
	model.DataBase.AutoMigrate(vlog.VlogComment{})
	//评论限制长度1000
	model.DataBase.AutoMigrate(vod.VodComment{})
}

func TestAutoMigrate1104(t *testing.T) {
	//添加用户表VIP字段
	model.DataBase.AutoMigrate(user.UserList{})
}
func TestAutoMigrate1009(t *testing.T) {
	//添加VIP功能
	model.DataBase.AutoMigrate(user.VipCodeList{})
	model.DataBase.AutoMigrate(user.VipUserList{})
	model.DataBase.AutoMigrate(blog.BlogList{})
}
func TestAutoMigrate1002(t *testing.T) {
	model.DataBase.AutoMigrate(attack.Attack{})
}
func TestAutoMigrate0928(t *testing.T) {
	model.DataBase.AutoMigrate(sexnovel.SexnovelLabel{})
	model.DataBase.AutoMigrate(sexnovel.SexnovelListLabel{})
	model.DataBase.AutoMigrate(sexnovel.SexnovelHistory{})
}

func TestAutoMigrate0907(t *testing.T) {
	//添加色图
	model.DataBase.AutoMigrate(seximg.SeximgType{})
	model.DataBase.AutoMigrate(seximg.SexImage{})
	model.DataBase.AutoMigrate(seximg.Seximg{})
	model.DataBase.AutoMigrate(seximg.SeximgStar{})
	//添加小说
	model.DataBase.AutoMigrate(sexnovel.SexnovelType{})
	model.DataBase.AutoMigrate(sexnovel.Sexnovel{})
	model.DataBase.AutoMigrate(sexnovel.SexnovelChapter{})
	model.DataBase.AutoMigrate(sexnovel.SexnovelContent{})
	model.DataBase.AutoMigrate(sexnovel.SexnovelStar{})
	//添加图片,视频类型
	model.DataBase.AutoMigrate(user.TaskImage{})
}
func TestAutoMigrate0813(t *testing.T) {
	//应用任务表
	model.DataBase.AutoMigrate(user.TaskList{})
	model.DataBase.AutoMigrate(user.TaskUser{})
	model.DataBase.AutoMigrate(user.TaskImage{})
}
func TestAutoMigrate0812(t *testing.T) {
	//应用列表添加AppID字段
	model.DataBase.AutoMigrate(applicationad.ApplicationAd{})
}
func TestAutoMigrate0805(t *testing.T) {
	//添加签到记录表
	model.DataBase.AutoMigrate(user.UserSignList{})
	//添加用户积分变动记录表
	//model.DataBase.AutoMigrate(user.UserIntegralLog{})
	//用户分享IP
	model.DataBase.AutoMigrate(user.UserShare{})
}
func TestAutoMigrate0801(t *testing.T) {
	model.DataBase.AutoMigrate(user.UserInviteList{})
}
func TestAutoMigrate0730(t *testing.T) {
	model.DataBase.AutoMigrate(vod.VodList{})
}
func TestAutoMigrate0724(t *testing.T) {
	model.DataBase.AutoMigrate(user.AppList{})
	model.DataBase.AutoMigrate(user.AppShare{})
}

// 首次同步数据表内容
func TestMigrateAll(t *testing.T) {
	model.DataBase.AutoMigrate(ad.AdDetail{})
	model.DataBase.AutoMigrate(ad.AdPostion{})
	model.DataBase.AutoMigrate(ad.AdViews{})

	/*
	 * 应用广告
	 */
	model.DataBase.AutoMigrate(applicationad.ApplicationAd{})
	model.DataBase.AutoMigrate(applicationad.ApplicationType{})
	model.DataBase.AutoMigrate(applicationad.ApplicationViews{})
	/*
	 * 博客功能
	 */
	model.DataBase.AutoMigrate(blog.AwardMatch{})
	model.DataBase.AutoMigrate(blog.BlogComment{})
	model.DataBase.AutoMigrate(blog.BlogImage{})
	model.DataBase.AutoMigrate(blog.BlogList{})
	model.DataBase.AutoMigrate(blog.BlogMatch{})
	model.DataBase.AutoMigrate(blog.BlogStar{})
	model.DataBase.AutoMigrate(blog.BlogTopic{})
	model.DataBase.AutoMigrate(blog.BlogVideo{})
	model.DataBase.AutoMigrate(blog.MatchRank{})

	/*
	 * 配置列表
	 */
	model.DataBase.AutoMigrate(config.ConfigList{})

	/*
	 * 系统权限表
	 */
	model.DataBase.AutoMigrate(admin.SysAuth{})
	model.DataBase.AutoMigrate(admin.SysRole{})
	model.DataBase.AutoMigrate(admin.SysUser{})

	model.DataBase.AutoMigrate(suggest.SuggestList{})
	/*
	 * 用户功能表
	 */
	model.DataBase.AutoMigrate(user.BlackIp{})
	model.DataBase.AutoMigrate(user.DirtyWord{})
	model.DataBase.AutoMigrate(user.Report{})
	model.DataBase.AutoMigrate(user.UserAvatar{})
	model.DataBase.AutoMigrate(user.UserChatList{})
	model.DataBase.AutoMigrate(user.UserComment{})
	model.DataBase.AutoMigrate(user.UserContentList{})
	model.DataBase.AutoMigrate(user.UserFollow{})
	model.DataBase.AutoMigrate(user.UserInviteList{})
	model.DataBase.AutoMigrate(user.UserList{})
	model.DataBase.AutoMigrate(user.UserName{})
	model.DataBase.AutoMigrate(user.UserNotice{})
	model.DataBase.AutoMigrate(user.UserNoticeRead{})
	model.DataBase.AutoMigrate(user.UserPasswd{})
	model.DataBase.AutoMigrate(user.UserShare{})
	model.DataBase.AutoMigrate(user.UserSms{})
	model.DataBase.AutoMigrate(user.UserUploader{})
	model.DataBase.AutoMigrate(user.AppList{})
	model.DataBase.AutoMigrate(user.AppShare{})

	/*
	 * Vip 长视频功能
	 */
	model.DataBase.AutoMigrate(vip.VipList{})
	model.DataBase.AutoMigrate(vip.VipUnlock{})

	/*
	 * 短视频功能
	 */
	model.DataBase.AutoMigrate(vlog.VlogComment{})
	model.DataBase.AutoMigrate(vlog.VlogList{})
	model.DataBase.AutoMigrate(vlog.VlogStar{})
	/*
	 * 长视频功能
	 */
	model.DataBase.AutoMigrate(vod.VodComment{})
	model.DataBase.AutoMigrate(vod.VodHistory{})
	model.DataBase.AutoMigrate(vod.VodLabel{})
	model.DataBase.AutoMigrate(vod.VodList{})
	model.DataBase.AutoMigrate(vod.VodListLabel{})
	model.DataBase.AutoMigrate(vod.VodListUser{})
	model.DataBase.AutoMigrate(vod.VodStar{})
	model.DataBase.AutoMigrate(vod.VodTopic{})
	model.DataBase.AutoMigrate(vod.VodType{})
}

func TestMtsSubmit(t *testing.T) {

	m := util.MtsSubmitJobsInput{
		Bucket:   os.Getenv("ALI_OSS_ORIGIN"),
		Location: os.Getenv("ALI_OSS_REGION"),
		Object:   "vod/567085ed73301972b57cc7ee4787785a",
	}
	jobid, err := util.MtsSubmitJobs(m, "vodhls/{FileName}", []byte("0123456789ABCDEF"))
	fmt.Println(jobid, err)
}

// 对图片进行 根据不包含后缀的文件名先sha1再md5的32位hex进行aes128cbc文件加密
// 测试生成加密的图片内容
func TestAesEncFile(t *testing.T) {
	url := "https://mogushipin.oss-cn-hangzhou.aliyuncs.com/vlogfile/109/a92d002f28a096760b7acc1909747ab5?x-oss-process=video/snapshot,t_5000,f_jpg,w_400,h_0"
	resp, err := http.Get(url)
	if err != nil {
		t.Log(err)
		t.Fail()
	}
	// resp.Body
	if w, e := os.OpenFile("1enc.jpg", os.O_WRONLY|os.O_CREATE, 0644); err != nil {
		t.Log(e)
		t.Fail()
	} else {
		data, err := io.ReadAll(resp.Body)
		key := util.GenAES128Key("1")
		fmt.Println(key)
		edata := util.CBCEncrypter(key, key, data)
		if err != nil {
			t.Log(err)
			t.Fail()
		}
		w.Write(edata)
	}
}

// 测试获取评论的数量
func TestTotalComment(t *testing.T) {
}

// 生成所有邀请码
func TestUpdateAllInviteCode(t *testing.T) {
	var users []user.UserList
	model.DataBase.Where("id > ?", 600000).Find(&users) //二次更新
	// model.DataBase.Find(&users, []int{1, 2, 3})
	// println(users[len(users)-2].Phone)
	for _, user := range users {
		if user.InviteCode == "" {
			user.InviteCode = sugar.GetInvCodeByUIDUnique(uint64(user.ID))
			model.DataBase.Model(&user).Update("invite_code", user.InviteCode)
		}
	}
}

// 修改评论
func TestUpdateComment(t *testing.T) {
	var vlogcomments []user.UserComment
	model.DataBase.Where("status = ?", 1).Where("type = ?", 2).Find(&vlogcomments) //查到短视频的评论 都改成社区的评论ID
	for _, vlogM := range vlogcomments {
		vlogM.Type = 3
		vlogM.Save()
	}
}

// 测试敏感词过滤
func TestDutyWord(t *testing.T) {
	var list []user.DirtyWord
	query := model.DataBase.Model(user.DirtyWord{})
	query.Find(&list)
	sensitive := sensitivewords.New()
	for _, word := range list {
		sensitive.AddWord(word.Name)
	}
	str := sensitive.Filter("阿海测试啊+q/+v,尼玛,哈哈")
	fmt.Printf("Filter:%v\n", str)
}
func isValidString(s string) bool {
	// 使用正则表达式匹配中文、英文和数字
	reg := regexp.MustCompile("^[a-zA-Z0-9\\p{Han}]+$")
	return reg.MatchString(s)
}

// 测试只有中文字和英文字母
func TestDutyString(t *testing.T) {
	testString := "Hello你𝘁 𝗼 𝗽	小蓝好123"
	if isValidString(testString) {
		fmt.Println("字符串合法")
	} else {
		fmt.Println("字符串不合法")
	}

}

// 临时测试
func TestToU(t *testing.T) {
	model.DataBase.First(&user.UserList{}, 52).Update("deny_comment", 2)
}

// 修改blog评论数 分三个方法避免数据过多
func TestUpdateBlogComment150(t *testing.T) {
	var blogs []blog.BlogList
	model.DataBase.Where("status = ?", 1).Where("comments >?", 150).Find(&blogs) //查询所有博客
	for _, item := range blogs {
		var total int64 = 0
		model.DataBase.Model(&user.UserComment{}).Where("c_id = ? And type = 3", item.ID).Where("status = 1").Count(&total)
		model.DataBase.Model(&item).Update("comments", total) //认证UP
	}
}
func TestUpdateBlogComment100(t *testing.T) {
	var blogs []blog.BlogList
	model.DataBase.Where("status = ?", 1).Where("comments >?", 100).Find(&blogs) //查询所有博客
	for _, item := range blogs {
		var total int64 = 0
		model.DataBase.Model(&user.UserComment{}).Where("c_id = ? And type = 3", item.ID).Where("status = 1").Count(&total)
		model.DataBase.Model(&item).Update("comments", total) //认证UP
	}
}
func TestUpdateBlogComment70(t *testing.T) {
	var blogs []blog.BlogList
	model.DataBase.Where("status = ?", 1).Where("comments >?", 70).Find(&blogs) //查询所有博客
	for _, item := range blogs {
		var total int64 = 0
		model.DataBase.Model(&user.UserComment{}).Where("c_id = ? And type = 3", item.ID).Where("status = 1").Count(&total)
		model.DataBase.Model(&item).Update("comments", total) //认证UP
	}
}

func TestUpdateBlogComment50(t *testing.T) {
	var blogs []blog.BlogList
	model.DataBase.Where("status = ?", 1).Where("comments >?", 50).Find(&blogs) //查询所有博客
	for _, item := range blogs {
		var total int64 = 0
		model.DataBase.Model(&user.UserComment{}).Where("c_id = ? And type = 3", item.ID).Where("status = 1").Count(&total)
		model.DataBase.Model(&item).Update("comments", total) //认证UP
	}
}
func TestUpdateBlogComment30(t *testing.T) {
	var blogs []blog.BlogList
	model.DataBase.Where("status = ?", 1).Where("comments >?", 30).Find(&blogs) //查询所有博客
	for _, item := range blogs {
		var total int64 = 0
		model.DataBase.Model(&user.UserComment{}).Where("c_id = ? And type = 3", item.ID).Where("status = 1").Count(&total)
		model.DataBase.Model(&item).Update("comments", total) //认证UP
	}
}
func TestUpdateBlogComment20(t *testing.T) {
	var blogs []blog.BlogList
	model.DataBase.Where("status = ?", 1).Where("comments >?", 20).Find(&blogs) //查询所有博客
	for _, item := range blogs {
		var total int64 = 0
		model.DataBase.Model(&user.UserComment{}).Where("c_id = ? And type = 3", item.ID).Where("status = 1").Count(&total)
		model.DataBase.Model(&item).Update("comments", total) //认证UP
	}
}
