// Code generated by hertz generator.

package core

import (
	"context"
	"time"

	core "github.com/ClubWeGo/douyin/biz/model/core"
	"github.com/ClubWeGo/douyin/kitex_server"
	"github.com/ClubWeGo/douyin/tools"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

// FeedMethod .
// @router /douyin/feed/ [GET]
func FeedMethod(ctx context.Context, c *app.RequestContext) {
	var err error
	var req core.FeedReq
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	// TODO : 记录ip地址和注册api调用次数，限制统一设备短时间太多的请求，预防爬虫。 redis

	resp := new(core.FeedResp)

	// 字段处理
	var currentUserId int64
	if req.Token != nil { // 可选字段，需要验证是否存在，判断对应指针是否存在
		_, currentUserId, err = tools.ValidateToken(*req.Token) //
		if err != nil {
			currentUserId = 0
			// "无效Token或Token已失效, 此处不做约束，继续执行代码"
		}
	}

	var latestTime = time.Now().UnixNano()
	if req.LatestTime != nil {
		latestTime = (*req.LatestTime) * 1e6 // app传入的是13位毫秒级时间戳，usermicro需传入纳秒级时间戳
	}

	// TODO : 缓存命中逻辑

	// 缓存未命中，去后端调api
	resultList, nextTime, err := kitex_server.GetFeed(latestTime, currentUserId, 30)
	if err != nil {
		msgFailed := "获取视频流失败" + err.Error()
		resp.StatusCode = 1
		resp.StatusMsg = &msgFailed
		c.JSON(consts.StatusOK, resp)
		return
	}

	resp.VideoList = resultList

	msgsucceed := "获取视频流成功"
	resp.StatusMsg = &msgsucceed
	resp.NextTime = nextTime / 1e6

	c.JSON(consts.StatusOK, resp)
}
