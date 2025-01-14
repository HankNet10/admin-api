package router

import (
	"myadmin/controller/api"
	"myadmin/controller/api/middle"
	"myadmin/controller/util"
	"time"

	"github.com/gin-gonic/gin"
)

func initApi(r *gin.Engine, s string) {
	apiRoute := r.Group(s) // 全局通用工具类路由。
	{

		// vip 列表
		apiRoute.GET("/vod/vip/list", JWTIsUserMiddleware(), CacheGetResult(20*time.Minute), api.VipList) // 获取VIP播放列表
		apiRoute.GET("/vod/vip/unlock", JWTUserMiddleware(), api.VipUnlock)                               // 获取VIP播放列表
		// 用户分享 老版本用户分享 可不处理保持原样
		apiRoute.POST("/user/share", api.UserShareAdd)
		apiRoute.GET("/user/share", api.UserShareList)
		// 新用户邀请
		apiRoute.GET("/user/sharelist", JWTUserMiddleware(), api.UserShareListNew) //用户列表

		// 长视频接口
		apiRoute.GET("/vodlabel/all", api.VodLabelAll)                                   // 所有标签
		apiRoute.GET("/vod/list", api.VodList)                                           // 视频列表
		apiRoute.GET("/vod/topiclist", CacheGetResult(10*time.Minute), api.VodTopicList) // 专题视频列表
		apiRoute.GET("/vod/info", api.VodInfo)                                           // 播放详情

		//话题
		apiRoute.GET("/vod/topic/list", api.VodTopicListGet)
		apiRoute.GET("/vod/topic/detail", api.VodTopicDetailsGet)

		apiRoute.GET("/vod/hotspot", CacheGetResult(20*time.Minute), api.Hotspot) // 热门视频排行
		apiRoute.GET("/vod/index", CacheGetResult(20*time.Minute), api.VodIndex)
		apiRoute.GET("/vod/type", CacheGetResult(20*time.Minute), api.TypeList)

		apiRoute.GET("/vod/24upload", CacheGetResult(50*time.Minute), api.Vod24Upload)           // 24小时更新
		apiRoute.GET("/vod/recommend", CacheGetResult(30*time.Minute), api.VodRecommend)         // 猜你喜欢
		apiRoute.GET("/vod/clever", api.VodClever)                                               // 智能搜索
		apiRoute.GET("/vod/list/follow", JWTUserMiddleware(), api.VodFollowList)                 // 我关注的视频
		apiRoute.GET("/vod/hotsearchword", CacheGetResult(60*time.Minute), api.HotSearchVodList) // 热搜
		// 用户视频播放记录与收藏
		apiRoute.PUT("/user/vod/history/add", JWTUserMiddleware(), api.VodHistoryAdd)
		apiRoute.GET("/user/vod/history/list", JWTUserMiddleware(), api.VodHistoryList)
		apiRoute.DELETE("/user/vod/history/delete", JWTUserMiddleware(), api.VodHistoryDelete)
		apiRoute.PUT("/user/vod/star/add", JWTUserMiddleware(), api.VodStarAdd) //用户收藏
		apiRoute.GET("/user/vod/star/list", JWTUserMiddleware(), api.VodStarList)
		apiRoute.DELETE("/user/vod/star/delete", JWTUserMiddleware(), api.VodStarDelete)
		// 用户社区播放记录与收藏
		apiRoute.PUT("/user/blog/star/add", JWTUserMiddleware(), api.BlogStarAdd) //用户收藏
		apiRoute.GET("/user/blog/star/list", JWTUserMiddleware(), api.BlogStarList)
		apiRoute.GET("/user/blog/star/isLike", JWTUserMiddleware(), api.BlogIsLike)
		apiRoute.DELETE("/user/blog/star/delete", JWTUserMiddleware(), api.BlogStarDelete)

		// 用户短视频播放记录与收藏
		apiRoute.GET("/vlog/list", CacheGetResult(15*time.Minute), api.VlogList)  // 获取vlog接口
		apiRoute.PUT("/user/vlog/star/add", JWTUserMiddleware(), api.BlogStarAdd) //用户收藏
		apiRoute.GET("/user/vlog/star/list", JWTUserMiddleware(), api.VlogStarList)
		apiRoute.DELETE("/user/vlog/star/delete", JWTUserMiddleware(), api.BlogStarDelete)

		//话题
		apiRoute.GET("/blog/topic/list", CacheGetResult(10*time.Minute), api.TopicListGet)
		apiRoute.GET("/blog/topic/detail", CacheGetResult(10*time.Minute), api.TopicDetailsGet)
		//比赛
		apiRoute.GET("/blog/match/reward", CacheGetResult(10*time.Minute), api.AwardDetailsGet)
		apiRoute.GET("/blog/match/list", CacheGetResult(10*time.Minute), api.MatchListGet)
		apiRoute.GET("/blog/match/detail", CacheGetResult(10*time.Minute), api.MatchDetailsGet)
		apiRoute.GET("/blog/match/ranklist", CacheGetResult(10*time.Minute), api.MatchEndRankList)
		// 社区图文
		apiRoute.GET("/blog/upuserlist", JWTUserMiddleware(), api.UserUpBlogList)
		apiRoute.POST("/blog/postbyuser", JWTUserMiddleware(), api.BlogCreateByUser)
		apiRoute.GET("/blog/list", api.BlogList)
		apiRoute.GET("/blog/search", api.BlogSearch)
		apiRoute.GET("/blog/info", CacheGetResult(20*time.Minute), api.BlogInfo)         // 获取blog 详情 接口
		apiRoute.DELETE("/blog/delete", JWTUserMiddleware(), api.BlogUserDelete)         //删除帖子
		apiRoute.GET("/blog/essay", CacheGetResult(30*time.Minute), api.BlogArticleInfo) // 获取blog 详情长文
		//用户评论
		apiRoute.GET("/blog/comment/list", api.CommentList) //里面的缓存
		apiRoute.POST("/blog/comment/add", JWTUserMiddleware(), api.CommentAdd)
		apiRoute.GET("/vlog/comment/list", api.CommentList)
		apiRoute.POST("/vlog/comment/add", JWTUserMiddleware(), api.CommentAdd)
		apiRoute.GET("/vod/comment/list", api.CommentList)
		apiRoute.POST("/vod/comment/add", JWTUserMiddleware(), api.CommentAdd)
		apiRoute.GET("/user/comment/list", JWTUserMiddleware(), api.CommentUserList)
		apiRoute.DELETE("/user/comment/delete", JWTUserMiddleware(), api.CommentDlete) //删除帖子

		apiRoute.GET("/user/comment/jobcroncomment", api.JobDleteComment)   //5分钟定时任务
		apiRoute.GET("/user/comment/jobcroncomment1", api.JobDleteComment1) //30分钟定时任务
		apiRoute.GET("/user/comment/jobcroncomment2", api.JobDleteComment2) //60分钟定时任务
		// 用户相关操作
		apiRoute.POST("/user/register", api.UserRegister)                          //用户注册
		apiRoute.POST("/user/resetpasswd", api.UserResetPasswd)                    //重设密码
		apiRoute.POST("/user/sms", api.SmsSubmit)                                  //用户注册短信验证码
		apiRoute.POST("/user/login", api.UserLogin)                                //用户登录
		apiRoute.POST("/user/unregister", JWTUserMiddleware(), api.UserUnRegister) //用户注销
		apiRoute.GET("/user/verify", JWTUserMiddleware(), api.UserVerify)          //用户验证接口 在我的页面 每日登录验证一次

		//UP主操作
		apiRoute.POST("/user/uploaderpost", JWTUserMiddleware(), api.UserUploaderPost) //申请up主
		apiRoute.GET("/user/uploaderstatu", JWTUserMiddleware(), api.UserGetUploader)  //获取用户状态
		apiRoute.GET("/user/getuptoken", JWTUserMiddleware(), api.UserSTSTOken)        //得到oss 签名STS

		//私信操作
		apiRoute.POST("/user/chat/post", JWTUserMiddleware(), api.UserChatSend)              //发送私信
		apiRoute.GET("/user/chat/getgfmessage", JWTUserMiddleware(), api.UserChatListGet)    //获得具体私信内容
		apiRoute.GET("/user/chat/getunread", JWTUserMiddleware(), api.UserUnredGet)          //获取未读数量和最后一条
		apiRoute.GET("/user/chat/getuserunread", JWTUserMiddleware(), api.UserUnredGet1)     //获取未读数量和最后一条
		apiRoute.POST("/user/chat/read", JWTUserMiddleware(), api.UserReadPost)              //已读
		apiRoute.GET("/user/chat/list", JWTUserMiddleware(), api.UserChatListList)           //获取聊天列表
		apiRoute.GET("/user/chat/getusermessage", JWTUserMiddleware(), api.UserChatListGet1) //获得具体私信内容

		apiRoute.GET("/user/actor", CacheGetResult(120*time.Minute), api.UserActor) //女优列表
		apiRoute.GET("/user/list", CacheGetResult(60*time.Minute), api.UserList)    //用户列表

		apiRoute.GET("/user/follow", api.UserFollow)   //我关注的人，或者关注我的人
		apiRoute.GET("/user/getgift", api.UserGetGift) //抽奖

		apiRoute.GET("/user/other", JWTIsUserMiddleware(), api.UserOther) //查看其他人的用户详情
		apiRoute.GET("/user/info", JWTUserMiddleware(), api.UserInfo)
		apiRoute.PUT("/user/follow/add", JWTUserMiddleware(), api.UserAddFollow)          // 关注其他用户
		apiRoute.DELETE("/user/follow/delete", JWTUserMiddleware(), api.UserDeleteFollow) // 关注其他用户
		apiRoute.PUT("/user/update", JWTUserMiddleware(), api.UserEdit)                   // 用户更新信息
		apiRoute.POST("/suggest/add", api.SuggestAdd)                                     // 用户提交建议
		apiRoute.POST("/report/add", api.ReportAdd)                                       // 用户提交举报

		//用户姓名

		//用户头像
		apiRoute.POST("/user/updateavatar", JWTIsUserMiddleware(), api.UserUpdateAvatar)
		apiRoute.GET("/user/avatarlist", api.UserGetAvatar)

		// 广告相关
		apiRoute.GET("/ads", api.AdDetailList) // 获取广告列表
		apiRoute.GET("/adclick", api.AdClick)  // 提交广告点击
		// 应用接口
		apiRoute.GET("/app/groupitems", api.AppItemDetailList)                                       // 获取应用分组列表
		apiRoute.GET("/app/allitems", api.AppItemAllDetailList)                                      // 获取应用所有列表
		apiRoute.GET("/app/click", api.APPItemClick)                                                 // 提交应用点击
		apiRoute.GET("/app/detail", CacheGetResult(30*time.Minute), api.AppItemDetail)               // 详情
		apiRoute.GET("/app/licationtypes", CacheGetResult(30*time.Minute), api.ApplicationTypeLists) // 获取类型列表
		apiRoute.GET("/app/belonginfo", api.AppItemBelong)                                           // 广告主详情
		apiRoute.GET("/app/belongclick", api.AppItemView)                                            // 广告主点击数据
		// apiRoute.GET("/cdnhost", api.ConfigCdnHost)    // 获取加速域名 稍等删除
		apiRoute.GET("/config", api.ConfigList)           // 获取自定义配置列表
		apiRoute.GET("/configgroup", api.ConfigGroupList) // 获取自定义配置组列表
		apiRoute.GET("/captcha", util.Captcha)            // 图文验证码
		apiRoute.POST("/cache", util.CacheSet)            // 存储数据 字符串类型
		apiRoute.GET("/cache", util.CacheGet)             // 获取存储数据 字符串类型

		apiRoute.POST("/download/shareName", middle.Api(), api.AppShareAdd)
		apiRoute.PUT("/download/shareRegister", middle.Api(), api.AppUpdate)

		apiRoute.GET("/user/signday", JWTUserMiddleware(), api.SignDay)      //用户签到记录
		apiRoute.GET("/user/signin", JWTUserMiddleware(), api.SignIn)        //用户签到
		apiRoute.GET("/newusershare", JWTUserMiddleware(), api.NewUserShare) //用户点击分享记录
		//apiRoute.GET("/adwatchs", JWTUserMiddleware(), api.AdWatchs)         //用户观看广告得积分

		apiRoute.GET("/task/list", api.TaskListList)                             //任务列表
		apiRoute.GET("/task/info", JWTUserMiddleware(), api.TaskListInfo)        //任务详情
		apiRoute.GET("/taskuser/info", JWTUserMiddleware(), api.TaskUserInfo)    //任务详情
		apiRoute.POST("/task/taskuseradd", JWTUserMiddleware(), api.TaskUserAdd) //用户提交任务

		//色图管理
		apiRoute.GET("/seximgtype/list", api.SeximgTypeList)
		apiRoute.GET("/seximg/list", api.SeximgList)
		apiRoute.GET("/seximg/info", api.SeximgInfo)
		apiRoute.GET("/user/seximg/comment/list", api.CommentList)
		apiRoute.POST("/user/seximg/comment/add", JWTUserMiddleware(), api.CommentAdd)
		apiRoute.GET("/user/seximg/star/isLike", JWTUserMiddleware(), api.SeximgIsLike)
		apiRoute.PUT("/user/seximg/star/add", JWTUserMiddleware(), api.SeximgStarAdd)
		apiRoute.GET("/user/seximg/star/list", JWTUserMiddleware(), api.SeximgStarList)
		apiRoute.DELETE("/user/seximg/star/delete", JWTUserMiddleware(), api.SeximgStarDelete)

		//小说管理
		apiRoute.GET("/sexnoveltype/list", api.SexnovelTypeList)
		apiRoute.GET("/sexnovel/label", CacheGetResult(20*time.Minute), api.SexnovelLabel)
		apiRoute.GET("/sexnovel/list", api.SexnovelList)
		apiRoute.GET("/sexnovel/chapte/list", api.SexnovelChapterList)
		apiRoute.GET("/sexnovel/info", api.SexnovelInfo)
		apiRoute.GET("/sexnovel/content/info", api.SexnovelContent)
		apiRoute.GET("/user/sexnovel/star/isLike", JWTUserMiddleware(), api.SexnovelIsLike)
		apiRoute.PUT("/user/sexnovel/star/add", JWTUserMiddleware(), api.SexnovelStarAdd)
		apiRoute.GET("/user/sexnovel/star/list", JWTUserMiddleware(), api.SexnovelStarList)
		apiRoute.DELETE("/user/sexnovel/star/delete", JWTUserMiddleware(), api.SexnovelStarDelete)
		// 用户小说观看记录与收藏
		apiRoute.PUT("/user/sexnovel/history/add", JWTUserMiddleware(), api.SexnovelHistoryAdd)
		apiRoute.GET("/user/sexnovel/history/list", JWTUserMiddleware(), api.SexnovelHistoryList)
		apiRoute.DELETE("/user/sexnovel/history/delete", JWTUserMiddleware(), api.SexnovelHistoryDelete)

		//VIP码
		apiRoute.POST("/user/vip/list", JWTUserMiddleware(), api.VipActivation)
		apiRoute.GET("/user/vip/list", JWTUserMiddleware(), api.VipUserInfo)
		apiRoute.POST("/user/vip/watch", JWTUserMiddleware(), api.VipWatch)

		//获取攻击列表
		apiRoute.GET("/user/attack/list", api.AttackList)
	}
}
