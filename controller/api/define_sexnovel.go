package api

import (
	"math/rand"
	"myadmin/model/sexnovel"
	"os"
)

// 用户添加收藏记录
type UserSexnovelStarAddRequest struct {
	SexnovelID uint `json:"sexnovel_id"`
}

type UserSexnovelCreateUserRequest struct {
	Typeid uint     `json:"Typeid"`
	Type   uint     `json:"type"`
	Title  string   `json:"title"`
	Images []string `json:"images"`
}

// 色图图片列表结构体
type apiResultSexnovelImages struct {
	ID   uint   `json:"id"`
	Path string `json:"path"`
}

type apiResultSexnovelTypeShow struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type apiResultSexnovelChapterShow struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

// 列表数据
type apiResultSexnovel struct {
	ID           uint   `json:"id"`
	Title        string `json:"title"`
	Cover        string `json:"cover"`
	Favorites    uint   `json:"favorites"`
	Star         uint   `json:"star"`
	Looks        uint   `json:"looks"`
	CreatedAt    int64  `json:"created_at"`
	SexnovelType string `json:"sexnovel_type"`
	Status       uint   `json:"status"`
	Top          uint   `json:"top"`
	Sort         uint   `json:"sort"`
	ChapterCount uint   `json:"chapter_count"`
	NewChapter   string `json:"new_chapter"`
	IsLong       bool   `json:"is_long"`
}

// 用户收藏列表结构体
type apiResultSexnovelStar struct {
	Id        uint              `json:"id"`
	UpdatedAt int64             `json:"updated_at"`
	Sexnovel  apiResultSexnovel `json:"sexnovel"`
}

func newApiResultSexnovel(v sexnovel.Sexnovel) apiResultSexnovel {
	return apiResultSexnovel{
		ID:           v.ID,
		Title:        v.Title,
		NewChapter:   v.NewChapter,
		ChapterCount: v.ChapterCount,
		Cover:        os.Getenv("ALI_OSS_DOMAIN") + "/" + v.Cover,
		Looks:        v.Looks*5 + uint(rand.Intn(9)),
		Star:         v.Star,
		Top:          uint(v.Top),
		Favorites:    v.Favorites,
		CreatedAt:    v.CreatedAt.Unix(),
		SexnovelType: v.SexnovelType.Name,
		Sort:         v.Sort,
		IsLong:       v.IsLong,
	}
}

// 生成结构体
func newApiResultSexnovelTypeShow(ul sexnovel.SexnovelType) apiResultSexnovelTypeShow {
	if ul.ID == 0 {
		return apiResultSexnovelTypeShow{}
	}
	return apiResultSexnovelTypeShow{
		ID:   ul.ID,
		Name: ul.Name,
	}
}

// 生成章节结构体
func newApiResultSexnovelChapterShow(ul sexnovel.SexnovelChapter) apiResultSexnovelChapterShow {
	if ul.ID == 0 {
		return apiResultSexnovelChapterShow{}
	}
	return apiResultSexnovelChapterShow{
		ID:   ul.ID,
		Name: ul.Title,
	}
}

func newApiUpResultSexnovelList(v sexnovel.Sexnovel) apiResultSexnovel {
	if v.Status == 1 {
		return apiResultSexnovel{
			ID:           v.ID,
			Title:        v.Title,
			Star:         v.Star,
			Favorites:    v.Favorites,
			CreatedAt:    v.CreatedAt.Unix(),
			Status:       1,
			SexnovelType: v.SexnovelType.Name,
		}
	} else { //没通过的不给一些数据
		return apiResultSexnovel{
			ID:           v.ID,
			Title:        v.Title,
			Star:         0,
			Favorites:    0,
			Status:       uint(v.Status),
			CreatedAt:    v.CreatedAt.Unix(),
			SexnovelType: v.SexnovelType.Name,
		}
	}
}

// 标签结构体
type apiResultSexnovelLabel struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

// 详情数据
type apiResultSexnovelInfo struct {
	ID         uint                     `json:"id"`
	Title      string                   `json:"title"`
	Views      uint                     `json:"views"`
	Favorites  uint                     `json:"favorites"`
	IsLong     bool                     `json:"is_long"`
	Cover      string                   `json:"cover"`
	Labels     []apiResultSexnovelLabel `json:"labels"`
	NewChapter string                   `json:"new_chapter"`
	CreatedAt  int64                    `json:"created_at"`
}

// 用户小说观看历史结构体
type apiResultSexnovelHistory struct {
	Id        uint                  `json:"id"`
	UpdatedAt int64                 `json:"updated_at"`
	Sexnovel  apiResultSexnovelInfo `json:"sexnovel"`
}

// 格式化用户列表数据
func newapiResultSexnovelList(v sexnovel.Sexnovel) apiResultSexnovelInfo {
	if v.Cover != "" {
		v.Cover = os.Getenv("ALI_OSS_DOMAIN") + "/" + v.Cover
	}
	return apiResultSexnovelInfo{
		ID:         v.ID,
		Title:      v.Title,
		Cover:      v.Cover,
		Views:      v.Looks,
		Favorites:  v.Favorites,
		IsLong:     v.IsLong,
		NewChapter: v.NewChapter,
	}
}
