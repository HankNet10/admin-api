package main

import (
	"myadmin/model"
	"myadmin/model/user"
)

func main() {
	// model.DataBase.AutoMigrate(ad.AdDetail{})
	// model.DataBase.AutoMigrate(ad.AdPostion{})
	// model.DataBase.AutoMigrate(ad.AdViews{})
	// model.DataBase.AutoMigrate(vod.VodList{})
	// model.DataBase.AutoMigrate(vlog.VlogList{})

	model.DataBase.AutoMigrate(user.UserList{})
}
