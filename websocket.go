package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/lxzan/gws"
	"github.com/tidwall/gjson"
	"net/http"
	"time"
)

func (b *Bot[T]) websocket() {
	header := http.Header{}
	if b.config.AccessToken != "" {
		header.Set("Authorization", "Bearer "+b.config.AccessToken)
	}

	b.config.Logger.Info("连接OneBot中...")
	socket, _, err := gws.NewClient(b, &gws.ClientOption{
		Addr:          fmt.Sprintf("ws://%s", b.config.Addr),
		RequestHeader: header,
	})
	if err != nil {
		b.config.Logger.Error("连接OneBot失败, 5秒后尝试重连...")
		b.config.Logger.Debug(err)
		b.Restart(5)
		return
	}
	b.config.Logger.Info("OneBot连接成功")

	api.logger = b.config.Logger
	api.conn = socket
	go socket.ReadLoop()

	login := WaitCallback(GetLoginInfo())
	b.config.Logger.Infof("登录QQ账号: %s (%d)",
		login.Data["nickname"].String(), login.Data["user_id"].Int())
}

func (b *Bot[T]) OnClose(socket *gws.Conn, err error) {
	b.config.Logger.Error("连接断开, 5秒后尝试重连...")
	b.config.Logger.Debug(err)
	b.Restart(5)
}

func (b *Bot[T]) OnPong(socket *gws.Conn, payload []byte) {
}

func (b *Bot[T]) OnOpen(socket *gws.Conn) {
}

func (b *Bot[T]) OnPing(socket *gws.Conn, payload []byte) {
}

func (b *Bot[T]) OnMessage(socket *gws.Conn, message *gws.Message) {
	defer func(message *gws.Message) {
		_ = message.Close()
	}(message)

	data := gjson.ParseBytes(message.Bytes())
	if data.Get("echo").Exists() {
		action := new(ActionResp)
		_ = json.Unmarshal(message.Bytes(), action)
		action.Data = gjson.GetBytes(message.Bytes(), "data").Map()
		if action.Status != "ok" {
			b.config.Logger.Errorf("API报错: %s", data.String())
		}
		if v, ok := b.api.call.Get(action.Echo); ok {
			v.resp = action
			v.Done()
		}
	} else {
		switch data.Get("post_type").String() {
		case "message": //消息事件
			if data.Get("message_type").String() == "private" { //私聊
				m := new(MessagePrivate)
				_ = json.Unmarshal(message.Bytes(), m)
				for _, handler := range b.event.messagePrivateHandlers {
					go b.cover(func(ctx context.Context) { handler(ctx, m, b.svc) })
				}
			} else if data.Get("message_type").String() == "group" { //群聊
				m := new(MessageGroup)
				_ = json.Unmarshal(message.Bytes(), m)
				for _, handler := range b.event.messageGroupHandlers {
					go b.cover(func(ctx context.Context) { handler(ctx, m, b.svc) })
				}
			}
		case "notice": //通知事件
			switch data.Get("notice_type").String() {
			case "group_upload": //群文件上传
				//GroupFileUploadEvent();
			case "group_admin": //群管理员变动
				//do something
			case "group_decrease": //群成员减少
				//do something
			case "group_increase": //群成员增加
				//do something
			case "group_ban": //群禁言
				//do something
			case "friend_add": //好友添加
				//do something
			case "group_recall": //群消息撤回
				//do something
			case "friend_recall": //好友消息撤回
				//do something
			case "notify": //群通知
				switch data.Get("sub_type").String() {
				case "poke": //戳一戳
					//do something
				case "lucky_king": //红包运气王
					//do something
				case "honor": //群荣誉
					//do something
				}
			}
		case "request": //请求事件
			if data.Get("request_type").String() == "friend" { //加好友请求
				//do something
			} else if data.Get("request_type").String() == "group" { //加群请求or邀请
				//do something
			}
		case "meta_event":
			switch data.Get("meta_event_type").String() {
			case "heartbeat": //ws心跳
				m := new(MetaHeartbeat)
				_ = json.Unmarshal(message.Bytes(), m)
				for _, handler := range b.event.metaHeartbeatHandlers {
					go b.cover(func(ctx context.Context) { handler(ctx, m, b.svc) })
				}
			case "lifecycle": //ws生命周期
				m := new(MetaConnect)
				_ = json.Unmarshal(message.Bytes(), m)
				for _, handler := range b.event.metaConnectHandlers {
					go b.cover(func(ctx context.Context) { handler(ctx, m, b.svc) })
				}
			}
		default:
			b.config.Logger.Infof("无法识别的Event返回: %s", data.String())
		}
	}
}

func (b *Bot[T]) Restart(sleep ...time.Duration) {
	if len(sleep) > 0 {
		time.Sleep(sleep[0] * time.Second)
	}
	b.websocket()
}

func (b *Bot[T]) cover(f func(ctx context.Context)) {
	defer func() {
		if p := recover(); p != nil {
			b.config.Logger.Panic("recovered from panic: ", p)
		}
	}()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	f(ctx)
}
