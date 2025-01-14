package api

import (
	"myadmin/model/blog"
	"os"
)

// 用户添加收藏记录
type UserBlogStarAddRequest struct {
	BlogID uint `json:"blog_id"`
	VlogId uint `json:"vlog_id"`
}

type UserBlogCreateUserRequest struct {
	MatchId   uint     `json:"matchid"`
	Topicid   uint     `json:"topicid"`
	Type      uint     `json:"type"`
	VideoPath string   `json:"videopath"`
	Title     string   `json:"title"`
	Essay     string   `json:"essay"`
	Images    []string `json:"images"`
}

// 图文社区图片列表结构体
type apiResultBlogImages struct {
	ID   uint   `json:"id"`
	Path string `json:"path"`
}

// 图文社区 列表数据
type apiResultBlogList struct {
	ID        uint                   `json:"id"`
	User      apiResultUserShow      `json:"user"`
	Detail    string                 `json:"detail"`
	Images    []apiResultBlogImages  `json:"images"`
	Favorites uint                   `json:"favorites"`
	Type      uint                   `json:"type"`
	Star      uint                   `json:"star"`
	Comments  uint                   `json:"comments"`
	CreatedAt int64                  `json:"created_at"`
	HaveEssay uint                   `json:"haveessay"`
	Topic     apiResultBlogTopicShow `json:"topic"`
	Match     apiResultBlogMatchShow `json:"match"`
	Video     apiResultBlogVideoShow `json:"Video"`
	Status    uint                   `json:"status"`
	Refuse    string                 `json:"Reson"`
	Top       uint                   `json:"top"`
}

// 用户收藏列表结构体
type apiResultBlogStar struct {
	Id        uint              `json:"id"`
	UpdatedAt int64             `json:"updated_at"`
	Blog      apiResultBlogList `json:"blog"`
}

func newApiResultBlogList(v blog.BlogList) apiResultBlogList {

	images := make([]apiResultBlogImages, len(v.Images))
	for i, v := range v.Images {
		images[i] = apiResultBlogImages{
			ID:   v.ID,
			Path: os.Getenv("ALI_OSS_DOMAIN") + "/" + v.Path,
		}
	}
	haveEssay := 0
	if len(v.Essay) > 10 {
		haveEssay = 1
	}
	return apiResultBlogList{
		ID:        v.ID,
		Detail:    v.Detail,
		Images:    images,
		User:      newApiResultUserShow(v.User),
		Type:      v.Type,
		Star:      v.Star,
		Comments:  v.Comments,
		Top:       uint(v.Top),
		Favorites: v.Favorites,
		HaveEssay: uint(haveEssay),
		CreatedAt: v.CreatedAt.Unix(),
		Match:     newApiResultBlogMatchShow(v.Match),
		Topic:     newApiResultBlogTopicShow(v.Topic),
		Video:     newApiResultBlogVideoShow(v.Video),
	}
}

func newApiUpResultBlogList(v blog.BlogList) apiResultBlogList {
	images := make([]apiResultBlogImages, len(v.Images))
	if v.Status == 1 { //通过才有
		for i, v := range v.Images {
			images[i] = apiResultBlogImages{
				ID:   v.ID,
				Path: os.Getenv("ALI_OSS_DOMAIN") + "/" + v.Path,
			}
		}
	}
	haveEssay := 0
	if v.Status == 1 {
		return apiResultBlogList{
			ID:        v.ID,
			Detail:    v.Detail,
			Images:    images,
			User:      newApiResultUserShow(v.User),
			Type:      v.Type,
			Star:      v.Star,
			Comments:  v.Comments,
			Favorites: v.Favorites,
			HaveEssay: uint(haveEssay),
			CreatedAt: v.CreatedAt.Unix(),
			Status:    1,
			Match:     newApiResultBlogMatchShow(v.Match),
			Topic:     newApiResultBlogTopicShow(v.Topic),
			Video:     newApiResultBlogVideoShow(v.Video),
		}
	} else { //没通过的不给一些数据
		return apiResultBlogList{
			ID:        v.ID,
			Detail:    v.Detail,
			Images:    nil,
			Refuse:    v.Refuse,
			User:      newApiResultUserShow(v.User),
			Type:      v.Type,
			Star:      0,
			Comments:  0,
			Favorites: 0,
			HaveEssay: 0,
			Status:    uint(v.Status),
			CreatedAt: v.CreatedAt.Unix(),
			Match:     newApiResultBlogMatchShow(v.Match),
		}
	}

}
