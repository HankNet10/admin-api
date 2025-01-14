package api

//
//
//  用户上传参数
//
//

// 注册用户需要的参数
type UserRegisterRequest struct {
	Phone        string `json:"phone"`
	Password     string `json:"password"`
	Verification string `json:"verification"`
	InviteCode   string `json:"invitecode"`
	CaptchaId    string `json:"captchaId"`
	Captcha      string `json:"captcha"`
}

// 用户登录需要的参数
type UserLoginRequest struct {
	Phone     string `json:"phone"`
	Password  string `json:"password"`
	CaptchaId string `json:"captchaId"`
	Captcha   string `json:"captcha"`
}

// 用户添加长视频播放记录
type UserVodHistroyAddRequest struct {
	VodID   uint `json:"vod_id"`
	VodTime uint `json:"vod_time"`
}

// 用户添加长视频收藏记录
type UserVodStarAddRequest struct {
	VodID uint `json:"vod_id"`
}

type UserFollowAddRequest struct {
	FollowID uint `json:"follow_id"`
}

type UserVodCommentAddRequest struct {
	VodID     uint   `json:"vod_id"`
	UserID    uint   `json:"user_id"`
	Comment   string `json:"comment"`
	ReUserId  uint   `json:"reuser_id"`
	CommentId uint   `json:"commment_id"`
}

type UserReportRequest struct {
	ReportId   string `json:"reportId"`
	CId        string `json:"cId"`
	CType      string `json:"cType"`
	Text       string `json:"text"`
	ReportedId string `json:"reportedId"`
	RType      string `json:"rType"`
}
type UserVlogCommentAddRequest struct {
	VlogID    uint   `json:"vlog_id"`
	Comment   string `json:"comment"`
	ReUserId  uint   `json:"reuser_id"`
	CommentId uint   `json:"commment_id"`
}

// 用户添加长视频收藏记录
type UserVlogStarAddRequest struct {
	VlogID uint `json:"vlog_id"`
}

// 用户登录需要的参数
type UserSmsRequest struct {
	Phone     string `json:"phone"`
	CaptchaId string `json:"captchaId"`
	Captcha   string `json:"captcha"`
}

// 用户登录需要的参数
type UserEditRequest struct {
	Name         string `json:"name"`
	Gender       string `json:"gender"`
	Birthday     string `json:"birthday"`
	Introduction string `json:"introduction"`
}

// 用户提交邀请参数
type UserSharePostRequest struct {
	ShareUserID uint   `json:"shareUserID"`
	DeviceId    string `json:"deviceID"`
	AccessKey   string `json:"accessKey"`
}

// 发送私信参数
type UserMessagePostRequest struct {
	SendUserID uint   `json:"sendUserID"`
	ReUserId   uint   `json:"reUserId"`
	Text       string `json:"text"`
	ImgPath    string `json:"imgPath"`
}

// 申请up主
type UserUploaderPostRequest struct {
	ImagePath string `json:"imagepath"`
	Introduce string `json:"introduce"`
	Verify    string `json:"verify"`
}

type UserAvatarRquestParam struct {
	Id uint `json:"id"`
}

// 用户添加小说观看记录
type UserSexnovelHistroyAddRequest struct {
	SexnovelID uint `json:"sexnovel_id"`
}
