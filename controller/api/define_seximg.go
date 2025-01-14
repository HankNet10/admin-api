package api

import (
	"math/rand"
	"myadmin/model/seximg"
	"os"
)

// 用户添加收藏记录
type UserSeximgStarAddRequest struct {
	SeximgID uint `json:"seximg_id"`
}

type UserSeximgCreateUserRequest struct {
	Typeid uint     `json:"Typeid"`
	Type   uint     `json:"type"`
	Title  string   `json:"title"`
	Images []string `json:"images"`
}

// 色图图片列表结构体
type apiResultSeximgImages struct {
	ID   uint   `json:"id"`
	Path string `json:"path"`
}

type apiResultSeximgTypeShow struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

// 色图社区 列表数据
type apiResultSeximg struct {
	ID uint `json:"id"`
	//User       apiResultUserShow       `json:"user"`
	Title      string                  `json:"title"`
	Cover      string                  `json:"cover"`
	Images     []apiResultSeximgImages `json:"images"`
	Favorites  uint                    `json:"favorites"`
	Star       uint                    `json:"star"`
	Comments   uint                    `json:"comments"`
	Looks      uint                    `json:"looks"`
	Count      uint                    `json:"count"`
	CreatedAt  int64                   `json:"created_at"`
	SeximgType string                  `json:"seximg_type"`
	Status     uint                    `json:"status"`
	Top        uint                    `json:"top"`
	Sort       uint                    `json:"sort"`
}

// 用户收藏列表结构体
type apiResultSeximgStar struct {
	Id        uint            `json:"id"`
	UpdatedAt int64           `json:"updated_at"`
	Seximg    apiResultSeximg `json:"seximg"`
}

func newApiResultSeximg(v seximg.Seximg) apiResultSeximg {

	images := make([]apiResultSeximgImages, len(v.Images))
	for i, v := range v.Images {
		images[i] = apiResultSeximgImages{
			ID:   v.ID,
			Path: os.Getenv("ALI_OSS_DOMAIN") + "/" + v.Path,
		}
	}
	return apiResultSeximg{
		ID:     v.ID,
		Title:  v.Title,
		Cover:  os.Getenv("ALI_OSS_DOMAIN") + "/" + v.Cover,
		Images: images,
		//User:       newApiResultUserShow(v.User),
		Looks:      v.Looks*5 + uint(rand.Intn(9)),
		Star:       v.Star,
		Comments:   v.Comments,
		Top:        uint(v.Top),
		Favorites:  v.Favorites,
		Count:      v.ImgCount,
		CreatedAt:  v.CreatedAt.Unix(),
		SeximgType: v.SeximgType.Name,
		Sort:       v.Sort,
	}
}

// 生成结构体
func newApiResultSeximgTypeShow(ul seximg.SeximgType) apiResultSeximgTypeShow {
	if ul.ID == 0 {
		return apiResultSeximgTypeShow{}
	}
	return apiResultSeximgTypeShow{
		ID:   ul.ID,
		Name: ul.Name,
	}
}

func newApiUpResultSeximgList(v seximg.Seximg) apiResultSeximg {
	images := make([]apiResultSeximgImages, len(v.Images))
	if v.Status == 1 { //通过才有
		for i, v := range v.Images {
			images[i] = apiResultSeximgImages{
				ID:   v.ID,
				Path: os.Getenv("ALI_OSS_DOMAIN") + "/" + v.Path,
			}
		}
	}
	if v.Status == 1 {
		return apiResultSeximg{
			ID:    v.ID,
			Title: v.Title,
			//Images: images,
			//User:       newApiResultUserShow(v.User),
			Star:       v.Star,
			Comments:   v.Comments,
			Favorites:  v.Favorites,
			Count:      v.ImgCount,
			CreatedAt:  v.CreatedAt.Unix(),
			Status:     1,
			SeximgType: v.SeximgType.Name,
		}
	} else { //没通过的不给一些数据
		return apiResultSeximg{
			ID:    v.ID,
			Title: v.Title,
			//Images: nil,
			//User:       newApiResultUserShow(v.User),
			Star:       0,
			Comments:   0,
			Favorites:  0,
			Count:      0,
			Status:     uint(v.Status),
			CreatedAt:  v.CreatedAt.Unix(),
			SeximgType: v.SeximgType.Name,
		}
	}

}
