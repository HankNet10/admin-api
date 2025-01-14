package api

import (
	"myadmin/model/ad"
	"myadmin/model/applicationad"
	"myadmin/model/blog"
	"myadmin/model/user"
	"myadmin/model/vip"
	"myadmin/model/vod"
	"os"
	"strconv"
)

// 长视频分类结构
type apiResultTypeListChild struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

// 长视频分类列表数据
type apiResultTypeList struct {
	ID    uint                     `json:"id"`
	Name  string                   `json:"name"`
	Child []apiResultTypeListChild `json:"child"`
}

// 用户公开数据列表
type apiResultUserShow struct {
	ID     uint   `json:"id"`
	Type   uint8  `json:"type"`
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
	Gender uint8  `json:"gender"`
	IsVip  uint8  `json:"isvip"`
}

// 博客话题
type apiResultBlogTopicShow struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

// 博客比赛
type apiResultBlogMatchShow struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

// 博客视频
type apiResultBlogVideoShow struct {
	ID    uint   `json:"id"`
	Path  string `json:"path"`
	Cover string `json:"cover"`
}

type apiPinyinUserShow struct {
	Pinyin string
	User   apiResultUserShow
}

// 快速生成用户公开结构体
func newApiResultUserShow(ul user.UserList) apiResultUserShow {
	if ul.Name == "" {
		ul.Name = "默认用户"
	}
	avatar := ""
	if ul.Avatar != "" {
		avatar = os.Getenv("ALI_OSS_DOMAIN") + "/" + ul.Avatar
	}
	return apiResultUserShow{
		ID:     ul.ID,
		Type:   ul.Type,
		Name:   ul.Name,
		Avatar: avatar,
		IsVip:  ul.IsVip,
	}
}

// 生成博客话题结构体
func newApiResultBlogTopicShow(ul blog.BlogTopic) apiResultBlogTopicShow {
	if ul.ID == 0 {
		return apiResultBlogTopicShow{}
	}
	return apiResultBlogTopicShow{
		ID:   ul.ID,
		Name: ul.Name,
	}
}

// 生成博客视频结构体
func newApiResultBlogVideoShow(ul blog.BlogVideo) apiResultBlogVideoShow {
	if ul.ID == 0 {
		return apiResultBlogVideoShow{}
	}
	return apiResultBlogVideoShow{
		ID:    ul.ID,
		Path:  os.Getenv("PLIST_DOMAIN") + "/play/" + strconv.Itoa(int(ul.ID)) + "/blog.plist",
		Cover: os.Getenv("ALI_OSS_DOMAIN") + "/" + ul.Cover,
	}
}

// 生成博客比赛结构体
func newApiResultBlogMatchShow(ul blog.BlogMatch) apiResultBlogMatchShow {
	if ul.ID == 0 {
		return apiResultBlogMatchShow{}
	}
	return apiResultBlogMatchShow{
		ID:   ul.ID,
		Name: ul.Name,
	}
}

// 长视频数据列表结构体
type apiResultVodList struct {
	ID          uint                `json:"id"`
	Title       string              `json:"title"`
	Cover       string              `json:"cover"`
	Views       uint                `json:"views"`
	HlsDuration uint                `json:"duration"`
	Users       []apiResultUserShow `json:"users"`
	Favorites   uint                `json:"favorites"`
}

// 长视频数据列表结构体
type apiResultTopicVodList struct {
	ID          uint                `json:"id"`
	Title       string              `json:"title"`
	Cover       string              `json:"cover"`
	Views       uint                `json:"views"`
	HlsDuration uint                `json:"duration"`
	Labels      []apiResultVodLabel `json:"labels"`
}

// 格式化用户列表数据
func newapiResultVodList(v vod.VodList) apiResultVodList {
	if v.Cover != "" {
		v.Cover = os.Getenv("ALI_OSS_DOMAIN") + "/" + v.Cover
	}
	users := make([]apiResultUserShow, len(v.UserList))
	for i, v := range v.UserList {
		users[i] = apiResultUserShow{
			ID:     v.User.ID,
			Type:   v.User.Type,
			Name:   v.User.Name,
			Avatar: os.Getenv("ALI_OSS_DOMAIN") + "/" + v.User.Avatar,
		}
	}

	return apiResultVodList{
		ID:          v.ID,
		Title:       v.Title,
		Cover:       v.Cover,
		Views:       v.Views,
		HlsDuration: v.HlsDuration,
		Users:       users,
		Favorites:   v.Favorites,
	}
}

// 格式化视频专题列表
func newapiResultTopicVodList(v vod.VodList) apiResultTopicVodList {
	cdnstr := os.Getenv("ALI_OSS_DOMAIN")
	if v.Cover != "" {
		v.Cover = cdnstr + "/" + v.Cover
	}
	// 组合标签数据
	labels := make([]apiResultVodLabel, len(v.LabelList))
	for i, v := range v.LabelList {
		labels[i] = apiResultVodLabel{
			ID:   v.Label.ID,
			Name: v.Label.Name,
		}
	}
	return apiResultTopicVodList{
		ID:          v.ID,
		Title:       v.Title,
		Cover:       v.Cover,
		Views:       v.Views,
		HlsDuration: v.HlsDuration,
		Labels:      labels,
	}
}

// 长视频标签结构体
type apiResultVodLabel struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

// 长视频播放详情数据
type apiResultVodInfo struct {
	ID        uint                `json:"id"`
	Title     string              `json:"title"`
	Number    string              `json:"nubmer"`
	PlayUrl   string              `json:"play_url"`
	Views     uint                `json:"views"`
	Favorites uint                `json:"favorites"`
	Comments  uint                `json:"comments"`
	Cover     string              `json:"cover"`
	Labels    []apiResultVodLabel `json:"labels"`
	Users     []apiResultUserShow `json:"users"`
	CreatedAt int64               `json:"created_at"`
}

// 用户播放历史结构体
type apiResultVodHistory struct {
	Id        uint             `json:"id"`
	LastSeen  uint             `json:"last_seen"`
	UpdatedAt int64            `json:"updated_at"`
	Vod       apiResultVodList `json:"vod"`
}

// 用户收藏列表结构体
type apiResultVodStar struct {
	Id        uint             `json:"id"`
	UpdatedAt int64            `json:"updated_at"`
	Vod       apiResultVodList `json:"vod"`
}

// 短视频 列表数据
type apiResultVlogList struct {
	ID        uint              `json:"id"`
	Title     string            `json:"title"`
	Cover     string            `json:"cover"`
	Favorites uint              `json:"favorites"`
	Views     uint              `json:"views"`
	Comments  uint              `json:"comments"`
	PlayUrl   string            `json:"play_url"`
	User      apiResultUserShow `json:"user"`
}

func newApiResultVlogList(v blog.BlogList) apiResultVlogList {
	return apiResultVlogList{
		ID:        v.ID,
		Title:     v.Detail,
		Cover:     os.Getenv("ALI_OSS_DOMAIN") + "/" + v.Video.Cover,
		Favorites: v.Favorites,
		Comments:  v.Comments,
		User:      newApiResultUserShow(v.User),
		PlayUrl:   "/play/" + strconv.Itoa(int(v.Video.ID)) + "/vlog.plist.m3u8",
	}
}

// 用户收藏列表结构体
type apiResultVlogStar struct {
	Id        uint              `json:"id"`
	UpdatedAt int64             `json:"updated_at"`
	Vlog      apiResultVlogList `json:"vlog"`
}

// 广告列表
type apiResultAdList struct {
	ID      uint   `json:"id"`
	Name    string `json:"name"`
	Value   string `json:"file"`
	Action  string `json:"action"`
	Weight  uint   `json:"weight"`
	Comment string `json:"comment"`
}

func newapiResultAdList(v *ad.AdDetail) apiResultAdList {
	return apiResultAdList{
		ID:      v.ID,
		Name:    v.Name,
		Value:   os.Getenv("ALI_OSS_DOMAIN") + "/" + v.Value,
		Action:  v.Action,
		Weight:  v.Weight,
		Comment: v.Comment,
	}
}

// 应用类型
// 广告类型列表
type apiResultAdTypeList struct {
	Name    string `json:"name"`
	Weight  uint   `json:"weight"`
	KeyWord string `json:"keyword"`
}

func newapiResultAdTypeList(v *applicationad.ApplicationType) apiResultAdTypeList {
	return apiResultAdTypeList{
		Name:    v.Value,
		Weight:  v.Weight,
		KeyWord: v.Name,
	}
}

// 应用分组列表
type apiResultAppGroupList struct {
	ID     uint   `json:"id"`
	Name   string `json:"name"`
	Value  string `json:"file"`
	Action string `json:"action"`
	Weight uint   `json:"weight"`
	InJump uint8  `json:"injump"`
}

func newapiResultAppGroupList(v *applicationad.ApplicationAd) apiResultAppGroupList {
	photohost := os.Getenv("ALI_OSS_DOMAIN") + "/"
	return apiResultAppGroupList{
		ID:     v.ID,
		Name:   v.Name,
		Value:  photohost + v.Value,
		Action: v.Action,
		Weight: v.Weight,
		InJump: v.InJump,
	}
}

// 应用列表
type apiResultApplicationList struct {
	ID              uint   `json:"id"`
	Name            string `json:"name"`
	Value           string `json:"file"`
	Action          string `json:"action"`
	Weight          uint   `json:"weight"`
	Desc            string `json:"desc"`
	Desc2           string `json:"desc2"`
	IsHot           uint8  `json:"ishot"`
	IsRecommend     uint8  `json:"isrecommend"`
	IsRank          uint8  `json:"isrank"`
	HotWeight       uint8  `json:"hotw"`
	RankWeight      uint8  `json:"rankw"`
	RecommendWeight uint8  `json:"recommendw"`
	Tag             uint8  `json:"tag"`
	InJump          uint8  `json:"injump"`
}

func newapiResultApplicationList(v *applicationad.ApplicationAd) apiResultApplicationList {
	photohost := os.Getenv("ALI_OSS_DOMAIN") + "/"
	return apiResultApplicationList{
		ID:              v.ID,
		Name:            v.Name,
		Value:           photohost + v.Value,
		Action:          v.Action,
		Weight:          v.Weight,
		Desc:            v.Desc,
		IsHot:           v.IsHot,
		IsRecommend:     v.IsRecommend,
		IsRank:          v.IsRank,
		HotWeight:       v.HotValue,
		RankWeight:      v.RankValue,
		RecommendWeight: v.RecommendValue,
		Tag:             v.Tag,
		Desc2:           v.Desc2,
		InJump:          v.InJump,
	}
}

type UserCommentAddRequest struct {
	BlogID    uint   `json:"blog_id"`
	VlogID    uint   `json:"vlog_id"`
	VodID     uint   `json:"vod_id"`
	SeximgID  uint   `json:"seximg_id"`
	Comment   string `json:"comment"`
	ReUserId  uint   `json:"reuser_id"`
	CommentId uint   `json:"commment_id"`
}
type apiResultCommentList struct {
	ID         uint              `json:"id"`
	CId        uint              `json:"c_id"`
	CType      uint              `json:"c_type"`
	Comment    string            `json:"comment"`
	User       apiResultUserShow `json:"user"`
	CreatedAt  int64             `json:"created_at"`
	Count      uint              `json:"count"`
	ParentId   uint              `json:"pid"`
	ReUserName string            `json:"reusername"`
	ReUserID   uint              `json:"reuserid"`
	RID        uint              `json:"rid"`
}

// VIP返回列表

type apiResultVipList struct {
	ID        uint
	TjIndex   int    `json:"tj_index"`
	DzIndex   int    `json:"dz_index"`
	JqIndex   int    `json:"jq_index"`
	YzIndex   int    `json:"yz_index"`
	Introduce string `json:"introduce"`

	Title   string `json:"title"`
	Cover   string `json:"cover"`
	PlayUrl string `json:"play_url"`

	Unlock      bool `json:"unlock"`
	Views       uint `json:"views"`
	HlsDuration uint `json:"duration"`
}

func newApiResultVipList(v vip.VipList) apiResultVipList {
	return apiResultVipList{
		ID:        v.ID,
		TjIndex:   v.TjIndex,
		DzIndex:   v.DzIndex,
		JqIndex:   v.JqIndex,
		YzIndex:   v.YzIndex,
		Introduce: v.Introduce,

		Title:   v.Vod.Title,
		Cover:   os.Getenv("ALI_OSS_DOMAIN") + "/" + v.Vod.Cover,
		PlayUrl: "/play/" + strconv.Itoa(int(v.Vod.ID)) + "/vod.plist",

		Unlock:      false,
		Views:       v.Vod.Views,
		HlsDuration: v.Vod.HlsDuration,
	}
}

// 前端返回体
type apiResultAvatarList struct {
	ID   uint   `json:"id"`
	Path string `json:"path"`
	Sort uint   `json:"sort"`
}
type apiResultChatContentList struct {
	ID        uint   `json:"id"`
	UserId    uint   `json:"user_id"`
	CreatedAt int64  `json:"created_at"`
	Text      string `json:"text"`
	Img       string `json:"img"`
}

// 任务图片列表结构体
type apiResultTaskUserImages struct {
	ID   uint   `json:"id"`
	Type uint   `json:"type"`
	Path string `json:"path"`
}

// 用户任务列表数据
type apiResultTaskUser struct {
	ID        uint                      `json:"id"`
	Images    []apiResultTaskUserImages `json:"images"`
	Videopath string                    `json:"videopath" `
	CreatedAt int64                     `json:"created_at"`
	Content   string                    `json:"content" `
	Status    uint                      `json:"status"`
	Refuse    string                    `json:"reson"`
}
type TaskUserCreateUserRequest struct {
	TaskID    uint     `json:"taskid"`
	Content   string   `json:"content"`
	Images    []string `json:"images"`
	Videopath string   `json:"videopath"`
}
type VipActivationRequest struct {
	Code string `json:"code"`
}
type VipWatchRequest struct {
	Vid int `json:"vid"`
}
