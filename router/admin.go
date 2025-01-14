package router

import (
	"myadmin/controller/admin"

	"github.com/gin-gonic/gin"
)

func initAdmin(r *gin.Engine, s string) {
	adminRoute := r.Group(s) // 全局通用工具类路由。
	{
		adminRoute.POST("/login", admin.Login) // 用户登录接口
		adminRoute.Use(JWTAdminMiddleware())   // 进行账号jwt解码

		adminRoute.GET("/account", admin.Account)                        // 获取自己的信息
		adminRoute.POST("/account/changepassword", admin.ChangePassword) // 修改用户密码
		adminRoute.Use(AuthorizeAdminMiddleware())                       // 进行账号鉴权

		/*
		 *
		 * - - - - - - - - - -|
		 * 管理员后台处理接口
		 * - - - - - - - - - -|
		 *
		 */
		adminRoute.POST("/admin/create", admin.AdminCreate)          // 创建管理员
		adminRoute.GET("/admin/list", admin.AdminList)               // 管理员列表
		adminRoute.PUT("/admin/update", admin.AdminUpdate)           // 修改管理员
		adminRoute.DELETE("/admin/delete", admin.AdminDelete)        // 删除管理员
		adminRoute.POST("/role/create", admin.RoleCreate)            // 创建角色
		adminRoute.GET("/role/list", admin.RoleList)                 // 创建角色
		adminRoute.PUT("/role/update", admin.RoleUpdate)             // 修改角色
		adminRoute.DELETE("/role/delete", admin.RoleDelete)          // 删除角色
		adminRoute.POST("/auth/create", admin.AuthCreate)            // 创建授权
		adminRoute.GET("/auth/list", admin.AuthList)                 // 授权列表
		adminRoute.PUT("/auth/update", admin.AuthUpdate)             // 修改授权
		adminRoute.DELETE("/auth/delete", admin.AuthDelete)          // 删除授权
		adminRoute.POST("/admin/uploadfile", admin.AndminFileUoload) // 文件上传
		/*
		 *
		 * - - - - - - - - - -|
		 * 长视频后台处理接口
		 * - - - - - - - - - -|
		 *
		 */
		adminRoute.GET("/vod/type/list", admin.VodTypeList)          // 类型列表
		adminRoute.POST("/vod/type/create", admin.VodTypeCreate)     // 新增类型
		adminRoute.PUT("/vod/type/update", admin.VodTypeUpdate)      // 更新类型
		adminRoute.DELETE("/vod/type/delete", admin.VodTypeDelete)   // 删除类型
		adminRoute.GET("/vod/label/list", admin.VodLabelList)        // 标签列表
		adminRoute.POST("/vod/label/create", admin.VodLabelCreate)   // 新建标签
		adminRoute.PUT("/vod/label/update", admin.VodLabelUpdate)    // 更新标签
		adminRoute.DELETE("/vod/label/delete", admin.VodLabelDelete) // 删除标签
		// 文件上传流程 | 计算文件md5 -> 验证存在 -> 数据库新增 -> 阿里云回调自动提交转码
		adminRoute.POST("/vod/list/exist", admin.VodListExist)     // 查验列表数据是否已存在数据
		adminRoute.GET("/vod/list/list", admin.VodListList)        // 视频列表
		adminRoute.POST("/vod/list/create", admin.VodListCreate)   // 新增视频
		adminRoute.PUT("/vod/list/update", admin.VodListUpdate)    // 更新视频
		adminRoute.DELETE("/vod/list/delete", admin.VodListDelete) // 删除视频
		// !!! 异常情况才会手动点击调用
		adminRoute.GET("/vod/media/submit", admin.VodSubmitJob) // 手动提交转码任务
		adminRoute.GET("/vod/media/result", admin.VodJobResult) // 手动获取转码结果
		// !!! 视频关联更多数据
		adminRoute.POST("/vod/list/label/create", admin.VodListLabelCreate)   // 视频关联标签
		adminRoute.DELETE("/vod/list/label/delete", admin.VodListLabelDelete) // 视频删除标签
		adminRoute.POST("/vod/list/user/create", admin.VodListUserCreate)     // 视频关联用户
		adminRoute.DELETE("/vod/list/user/delete", admin.VodListUserDelete)   // 视频删除用户
		// !!!
		//专题管理
		adminRoute.GET("/vod/topic/list", admin.VodTopicList)
		adminRoute.POST("/vod/topic/create", admin.VodTopicCreate)
		adminRoute.PUT("/vod/topic/update", admin.VodTopicUpdate)
		adminRoute.DELETE("/vod/topic/delete", admin.VodTopicDelete)
		// 博客管理
		adminRoute.GET("/blog/list/list", admin.BlogListList)
		adminRoute.POST("/blog/list/create", admin.BlogListCreate)
		adminRoute.PUT("/blog/list/update", admin.BlogListUpdate)
		adminRoute.PUT("/blog/list/updatevideo", admin.BlogVideoCreate) //后台上传博客视频
		adminRoute.GET("/blog/media/submit", admin.BlogVideoSubmitJob)  // 提交博客视频
		adminRoute.GET("/blog/media/cover", admin.BlogCoverSubmitJob)   // 提交博客封面生成

		adminRoute.DELETE("/blog/list/delete", admin.BlogListDelete)
		adminRoute.POST("/blog/image/create", admin.BlogImageCreate) // 博客图片
		adminRoute.DELETE("/blog/image/delete", admin.BlogImageDelete)
		//热搜 纯redis存储
		adminRoute.GET("/search/list/list", admin.HotListList)
		adminRoute.DELETE("/search/list/refreshall", admin.HotReplaceAll)
		adminRoute.DELETE("/search/list/delete", admin.HotDelete)
		adminRoute.PUT("/search/list/update", admin.HotUpdate)
		//话题管理
		adminRoute.GET("/blog/topic/list", admin.BlogTopicList)
		adminRoute.POST("/blog/topic/create", admin.BlogTopicCreate)
		adminRoute.PUT("/blog/topic/update", admin.BlogTopicUpdate)
		adminRoute.DELETE("/blog/topic/delete", admin.BlogTopicDelete)
		//比赛管理
		adminRoute.GET("/blog/match/list", admin.BlogMatchList)
		adminRoute.POST("/blog/match/create", admin.BlogMatchCreate)
		adminRoute.POST("/blog/match/statusUpdate", admin.BlogMatchStatusUpdate)
		adminRoute.PUT("/blog/match/update", admin.BlogMatchUpdate)
		adminRoute.DELETE("/blog/match/delete", admin.BlogMatchDelete)
		adminRoute.GET("/blog/match/ranklist", admin.MatchRankList)
		//比赛奖品
		adminRoute.GET("/blog/award/list", admin.AwardMatchList)
		adminRoute.POST("/blog/award/create", admin.AwardMatchCreate)
		adminRoute.PUT("/blog/award/update", admin.AwardMatchUpdate)
		adminRoute.DELETE("/blog/award/delete", admin.AwardMatchDelete)
		//评论
		adminRoute.GET("/user/comment/list", admin.UserCommentList)
		adminRoute.DELETE("/user/comment/delete", admin.UserCommentDelete) //删除单条评论
		adminRoute.PUT("/user/comment/update", admin.UserCommentUpdate)    // 更新单条数据
		adminRoute.GET("/user/silience", admin.UserDisableComment)         // 删除用户的所有评论并禁言
		adminRoute.POST("/user/comment/allow", admin.AllowLen)             // 批量通过此页

		// Vlog
		adminRoute.POST("/vlog/list/create", admin.VlogListCreate)
		adminRoute.GET("/vlog/list/list", admin.VlogListList)
		adminRoute.PUT("/vlog/list/update", admin.VlogListUpdate)
		adminRoute.DELETE("/vlog/list/delete", admin.VlogListDelete)
		adminRoute.GET("/vlog/media/submit", admin.VlogSubmitJob) // 手动提交转码任务
		adminRoute.GET("/vlog/media/result", admin.VlogJobResult) // 手动获取转码结果
		//审核up管理
		adminRoute.GET("/user/up/list", admin.UserUploaderList)
		adminRoute.PUT("/user/up/update", admin.UserUploaderUpdate)
		adminRoute.DELETE("/user/up/delete", admin.UserUploaderDelete)
		adminRoute.POST("/user/up/allow", admin.AllowUp) // 批量通过此页

		//审核敏感词
		adminRoute.GET("/user/dirtyword/list", admin.DirtyWordList)
		adminRoute.POST("/user/dirtyword/create", admin.DirtyWordCreate)
		adminRoute.DELETE("/user/dirtyword/delete", admin.DirtyWordDelete)
		//黑名单IP
		adminRoute.GET("/user/blackip/list", admin.BlackIpList)
		adminRoute.POST("/user/blackip/create", admin.BlackIpCreate)
		adminRoute.DELETE("/user/blackip/delete", admin.BlackIpDelete)
		//通知管理
		adminRoute.GET("/user/notice/list", admin.UserNoticeList)
		adminRoute.DELETE("/user/notice/delete", admin.UserNoticeDelete)
		adminRoute.POST("/user/notice/sendadmin", admin.UserAdminNoticeCreate)
		adminRoute.PUT("/user/notice/update", admin.UserNoticeUpdate)
		// 用户管理
		adminRoute.GET("/user/list/list", admin.UserListList)
		adminRoute.POST("/user/list/create", admin.UserListCreate)
		adminRoute.PUT("/user/list/update", admin.UserListUpdate)
		adminRoute.DELETE("/user/list/delete", admin.UserListDelete)
		//头像管理
		adminRoute.POST("/user/avatar/upload", admin.UserAvatarCreate) // 用户头像列表上传
		adminRoute.GET("/user/avatar/list", admin.UserAvatarList)      // 用户头像列表
		adminRoute.DELETE("/user/avatar/delete", admin.UserAvatarDelete)
		//站内信
		adminRoute.GET("/user/chat/content", admin.UserContentListList)
		adminRoute.GET("/user/chat/userlist", admin.UserChatListList)
		adminRoute.POST("/user/chat/post", admin.UserChatAdminSend)
		adminRoute.POST("/user/chat/read", admin.UserMssageRead)
		//邀请关系
		adminRoute.GET("/user/invite/list", admin.UserInviteList)
		// 配置
		adminRoute.GET("/config/list/list", admin.ConfigListList)
		adminRoute.POST("/config/list/create", admin.ConfigListCreate)
		adminRoute.PUT("/config/list/update", admin.ConfigListUpdate)
		adminRoute.DELETE("/config/list/delete", admin.ConfigListDelete)
		adminRoute.GET("/config/sts", admin.ConfigSts) // 获取oss key

		// 广告
		adminRoute.GET("/ad/postion/list", admin.AdPostionList)
		adminRoute.POST("/ad/postion/create", admin.AdPostionCreate)
		adminRoute.PUT("/ad/postion/update", admin.AdPostionUpdate)
		adminRoute.DELETE("/ad/postion/delete", admin.AdPostionDelete)

		// - 广告
		adminRoute.GET("/ad/detail/list", admin.AdDetailList)
		adminRoute.POST("/ad/detail/create", admin.AdDetailCreate)
		adminRoute.PUT("/ad/detail/update", admin.AdDetailUpdate)
		adminRoute.DELETE("/ad/detail/delete", admin.AdDetailDelete)
		adminRoute.GET("/ad/views/list", admin.AdViewsList)
		adminRoute.PUT("/ad/detail/replaceaction", admin.AdDetailReplaceAction)

		// - 应用中心
		adminRoute.GET("/ad/appad/list", admin.ApplicationAdList)
		adminRoute.POST("/ad/appad/create", admin.ApplicationAdCreate)
		adminRoute.PUT("/ad/appad/update", admin.ApplicationAdUpdate)
		adminRoute.DELETE("/ad/appad/delete", admin.ApplicationAdDelete)
		adminRoute.GET("/ad/appadviews/list", admin.ApplicationAdViewsList)

		adminRoute.GET("/ad/apptype/list", admin.ApplicationTypeList)
		adminRoute.POST("/ad/apptype/create", admin.ApplicationTypeCreate)
		adminRoute.PUT("/ad/apptype/update", admin.ApplicationTypeUpdate)
		adminRoute.DELETE("/ad/apptype/delete", admin.ApplicationTypeDelete)

		// - 反馈列表
		adminRoute.GET("/suggest/list", admin.SuggestList)
		// - 举报列表
		adminRoute.GET("/report/list", admin.ReportList)
		// - 用户昵称
		adminRoute.GET("/username/list", admin.UserNameList)
		adminRoute.DELETE("/username/delete", admin.UserNameDelete)
		adminRoute.PUT("/username/update", admin.UserNameUpdate)
		adminRoute.POST("/username/allow", admin.UserNameAllowPage) // 批量通过此页
		// - VIP列表编辑
		adminRoute.POST("/vip/list", admin.VipPost)
		adminRoute.GET("/vip/list", admin.VipList)
		adminRoute.PUT("/vip/list", admin.VipPut)
		adminRoute.DELETE("/vip/list", admin.VipDelete)
		// - AppList列表编辑
		adminRoute.GET("/app/list", admin.AppListList)
		adminRoute.POST("/app/list", admin.AppListCreate)
		adminRoute.PUT("/app/list", admin.AppListUpdate)
		adminRoute.DELETE("/app/list", admin.AppListDelete)
		// - AppShare列表编辑
		adminRoute.GET("/app_share/list", admin.AppShareList)
		adminRoute.POST("/app_share/list", admin.AppShareCreate)
		adminRoute.PUT("/app_share/list", admin.AppShareUpdate)
		adminRoute.DELETE("/app_share/list", admin.AppShareDelete)

		adminRoute.GET("/user/signlist", admin.UserSignList)

		// - TaskList任务管理
		adminRoute.GET("/task/list", admin.TaskListList)
		adminRoute.POST("/task/list", admin.TaskListCreate)
		adminRoute.PUT("/task/list", admin.TaskListUpdate)
		adminRoute.DELETE("/task/list", admin.TaskListDelete)

		// - TaskList任务管理
		adminRoute.GET("/taskuser/list", admin.TaskUserList)
		adminRoute.POST("/taskuser/list", admin.TaskUserCreate)
		adminRoute.PUT("/taskuser/list", admin.TaskUserUpdate)
		adminRoute.DELETE("/taskuser/list", admin.TaskUserDelete)

		// - Seximg色图管理
		adminRoute.GET("/seximg/type", admin.SeximgTypeList)
		adminRoute.POST("/seximg/type", admin.SeximgTypeCreate)
		adminRoute.PUT("/seximg/type", admin.SeximgTypeUpdate)
		adminRoute.DELETE("/seximg/type", admin.SeximgTypeDelete)
		adminRoute.GET("/seximg/list", admin.SeximgList)
		adminRoute.POST("/seximg/list", admin.SeximgCreate)
		adminRoute.PUT("/seximg/list", admin.SeximgUpdate)
		adminRoute.DELETE("/seximg/list", admin.SeximgDelete)
		adminRoute.POST("/seximg/image", admin.SexImageCreate)
		adminRoute.DELETE("/seximg/image", admin.SexImageDelete)

		// - Sexnovel小说管理
		adminRoute.GET("/sexnovel/type", admin.SexnovelTypeList)
		adminRoute.POST("/sexnovel/type", admin.SexnovelTypeCreate)
		adminRoute.PUT("/sexnovel/type", admin.SexnovelTypeUpdate)
		adminRoute.DELETE("/sexnovel/type", admin.SexnovelTypeDelete)
		adminRoute.GET("/sexnovel/list", admin.SexnovelList)
		adminRoute.POST("/sexnovel/list", admin.SexnovelCreate)
		adminRoute.POST("/sexnovel/upload", admin.SexnovelCreate1)
		adminRoute.PUT("/sexnovel/list", admin.SexnovelUpdate)
		adminRoute.DELETE("/sexnovel/list", admin.SexnovelDelete)
		adminRoute.GET("/sexnovelchapter/list", admin.SexnovelChapterList)
		adminRoute.POST("/sexnovelchapter/list", admin.SexnovelChapterCreate)
		adminRoute.PUT("/sexnovelchapter/list", admin.SexnovelChapterUpdate)
		adminRoute.DELETE("/sexnovelchapter/list", admin.SexnovelChapterDelete)
		adminRoute.GET("/sexnovelcontent/list", admin.SexnovelContent)
		adminRoute.GET("/sexnovelcontent/content", admin.SexnovelContentInfo)
		adminRoute.POST("/sexnovelcontent/list", admin.SexnovelContentCreate)
		adminRoute.PUT("/sexnovelcontent/list", admin.SexnovelContentUpdate)
		adminRoute.DELETE("/sexnovelcontent/list", admin.SexnovelContentDelete)
		adminRoute.GET("/sexnovel/label/list", admin.SexnovelLabelList)                 // 标签列表
		adminRoute.POST("/sexnovel/label/create", admin.SexnovelLabelCreate)            // 新建标签
		adminRoute.PUT("/sexnovel/label/update", admin.SexnovelLabelUpdate)             // 更新标签
		adminRoute.DELETE("/sexnovel/label/delete", admin.SexnovelLabelDelete)          // 删除标签
		adminRoute.POST("/sexnovel/list/label/create", admin.SexnovelListLabelCreate)   // 视频关联标签
		adminRoute.DELETE("/sexnovel/list/label/delete", admin.SexnovelListLabelDelete) // 视频删除标签

		// - VIP管理
		adminRoute.GET("/vipcode/list", admin.VipCodeList)
		adminRoute.POST("/vipcode/list", admin.VipCodeListCreate)
		adminRoute.POST("/vipcode/createmore", admin.VipCodeListCreateMore)
		// - VIP用户管理
		adminRoute.GET("/vipuser/list", admin.VipUserList)
		adminRoute.POST("/vipuser/list", admin.VipUserListCreate)
		adminRoute.PUT("/vipuser/list", admin.VipUserListUpdate)
		adminRoute.DELETE("/vipuser/list", admin.VipUserListDelete)

		//攻击管理
		adminRoute.GET("/attack/list", admin.AttackList)
		adminRoute.POST("/attack/list", admin.AttackCreate)
		adminRoute.PUT("/attack/list", admin.AttackUpdate)
		adminRoute.DELETE("/attack/list", admin.AttackDelete)
	}
}
