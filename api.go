package main

import (
	"encoding/json"
	"fmt"
	cmap "github.com/endymx/concurrent-map"
	"github.com/lxzan/gws"
	"sync"
	"time"
)

type API struct {
	logger Logger

	conn     *gws.Conn
	actionId int64
	call     cmap.ConcurrentMap[string, *Call]
}

var api = &API{
	call: cmap.ConcurrentMap[string, *Call]{},
}

func send(action string, params map[string]any) string {
	if api.conn == nil {
		api.logger.Errorf("QQ机器人未成功登录，无法发送")
		return ""
	}

	api.actionId++
	echo := fmt.Sprintf("%d%d", time.Now().Unix(), api.actionId)
	data, _ := json.Marshal(ActionReq{
		Action: action,
		Params: params,
		Echo:   echo,
	})

	api.logger.Debugf("发送api：%s", string(data))
	err := api.conn.WriteMessage(gws.OpcodeText, data)
	if err != nil {
		api.logger.Errorf("发送给QQ机器人消息失败：%s", err)
		return ""
	}
	return echo
}

type Call struct {
	resp *ActionResp
	*sync.WaitGroup
}

// WaitCallback 等待回调, 5秒超时后自动退出等待
func WaitCallback(echo string) (resp *ActionResp) {
	if echo == "" {
		return nil
	}
	timeout := func() {
		time.AfterFunc(time.Second*5, func() {
			if v, ok := api.call.Get(echo); ok {
				v.Done()
			}
		})
	}

	if v, ok := api.call.Get(echo); !ok {
		wg := &Call{
			WaitGroup: new(sync.WaitGroup),
		}
		wg.Add(1)
		api.call.Set(echo, wg)

		timeout()
		wg.Wait()
		resp = wg.resp
	} else {
		timeout()
		v.Wait()
		resp = v.resp
	}

	api.call.Remove(echo)
	return
}

//##################################信息类###################################

// SendPrivateMsg 发送私聊消息
func SendPrivateMsg[MT MessageType](userId int64, msg MT, autoEscape ...bool) string {
	m := map[string]any{}
	m["user_id"] = userId
	m["message"] = msg
	if len(autoEscape) == 1 {
		m["auto_escape"] = autoEscape[0]
	}
	return send("send_private_msg", m)
}

// SendGroupMsg 发送群消息
func SendGroupMsg[MT MessageType](groupId int64, msg MT, autoEscape ...bool) string {
	m := map[string]any{}
	m["group_id"] = groupId
	m["message"] = msg
	if len(autoEscape) == 1 {
		m["auto_escape"] = autoEscape[0]
	}
	return send("send_group_msg", m)
}

// SendGroupForwardMsg 发送合并转发消息
func SendGroupForwardMsg[MT MessageType](groupId int64, msg MT) string {
	m := map[string]any{}
	m["group_id"] = groupId
	m["messages"] = msg
	return send("send_group_forward_msg", m)
}

// SendMsg 发送消息
func SendMsg[MT MessageType](messageType any, userId int64, groupId int64, msg MT, autoEscape ...bool) string {
	m := map[string]any{}
	m["message_type"] = messageType
	m["user_id"] = userId
	m["group_id"] = groupId
	m["message"] = msg
	if len(autoEscape) == 1 {
		m["auto_escape"] = autoEscape[0]
	}
	return send("send_msg", m)
}

// DeleteMsg 撤回消息
func DeleteMsg(messageId int32) string {
	m := map[string]any{}
	m["message_id"] = messageId
	return send("delete_msg", m)
}

func getMsg(messageId int32) string {
	m := map[string]any{}
	m["message_id"] = messageId
	return send("get_msg", m)
}

func getForwardMsg(id int32) string {
	m := map[string]any{}
	m["id"] = id
	return send("get_forward_msg", m)
}

//##################################群组管理类###################################

func setGroupKick(groupId int64, userId int64, request bool) string {
	m := map[string]any{}
	m["group_id"] = groupId
	m["user_id"] = userId
	m["reject_add_request"] = request
	return send("set_group_kick", m)
}

// SetGroupBan 群组单人禁言
func SetGroupBan(groupId int64, userId int64, duration int) string {
	m := map[string]any{}
	m["group_id"] = groupId
	m["user_id"] = userId
	m["duration"] = duration
	return send("set_group_ban", m)
}

func setGroupABan(groupId int64, anonymous any, flag string, duration int) string {
	m := map[string]any{}
	m["group_id"] = groupId
	m["anonymous"] = anonymous
	m["flag"] = flag
	m["duration"] = duration
	return send("set_group_anonymous_ban", m)
}

func setGroupAllBan(groupId int64, enable bool) string {
	m := map[string]any{}
	m["group_id"] = groupId
	m["enable"] = enable
	return send("set_group_whole_ban", m)
}

func setGroupAdmin(groupId int64, userId int64, enable bool) string {
	m := map[string]any{}
	m["group_id"] = groupId
	m["user_id"] = userId
	m["enable"] = enable
	return send("set_group_admin", m)
}

func setGroupAnonymous(groupId int64, enable bool) string {
	m := map[string]any{}
	m["group_id"] = groupId
	m["enable"] = enable
	return send("set_group_anonymous", m)
}

func setGroupCard(groupId int64, userId int64, card string) string {
	m := map[string]any{}
	m["group_id"] = groupId
	m["user_id"] = userId
	m["card"] = card
	return send("set_group_card", m)
}

func SetGroupName(groupId int64, groupName string) string {
	m := map[string]any{}
	m["group_id"] = groupId
	m["group_name"] = groupName
	return send("set_group_name", m)
}

func setGroupLeave(groupId int64, isDismiss bool) string {
	m := map[string]any{}
	m["group_id"] = groupId
	m["is_dismiss"] = isDismiss
	return send("set_group_leave", m)
}

func setGroupTitle(groupId int64, userId int32, specialTitle string, duration int) string {
	m := map[string]any{}
	m["group_id"] = groupId
	m["user_id"] = userId
	m["special_title"] = specialTitle
	m["duration"] = duration
	return send("set_group_special_title", m)
}

//##################################杂项类###################################
/**
 * 发送好友赞
 *
 * @param int user_id 要设置的 QQ 号
 * @param string times 赞的次数，每个好友每天最多 10 次
 * @return mixed (true|null|false)只要不返回false并且网络没有断开，而且服务端接收正常，数据基本上可以看做100%能发过去
 */
func sendLike(userId int64, times int) string {
	m := map[string]any{}
	m["times"] = times
	m["user_id"] = userId
	return send("send_like", m)
}

/**
 * 处理加好友请求
 *
 * @param string flag 加好友请求的 flag（需从上报的数据中获得）
 * @param bool approve 是否同意请求
 * @param string remark 添加后的好友备注（仅在同意时有效）
 * @return mixed (true|null|false)只要不返回false并且网络没有断开，而且服务端接收正常，数据基本上可以看做100%能发过去
 */
func setFriendAdd(flag string, approve bool, remark string) string {
	m := map[string]any{}
	m["flag"] = flag
	m["approve"] = approve
	m["remark"] = remark
	return send("set_friend_add_request", m)
}

/**
 * 处理加群请求／邀请
 *
 * @param string flag 加好友请求的 flag（需从上报的数据中获得）
 * @param string type add 或 invite，请求类型（需要和上报消息中的 sub_type 字段相符）
 * @param bool approve 是否同意请求
 * @param string reason 拒绝理由（仅在拒绝时有效）
 * @return mixed (true|null|false)只要不返回false并且网络没有断开，而且服务端接收正常，数据基本上可以看做100%能发过去
 */
func setGroupAdd(flag string, typea string, approve bool, reason string) string {
	m := map[string]any{}
	m["flag"] = flag
	m["type"] = typea
	m["approve"] = approve
	m["reason"] = reason
	return send("set_group_add_request", m)
}

func GetLoginInfo() string {
	return send("get_login_info", nil)
}

func getStrangerInfo(userId int64, noCache bool) string {
	m := map[string]any{}
	m["user_id"] = userId
	m["no_cache"] = noCache
	return send("get_stranger_info", m)
}

func getFriendList() string {
	return send("get_friend_list", nil)
}

func getFriendInfo(groupId int64, noCache bool) string {
	m := map[string]any{}
	m["group_id"] = groupId
	m["no_cache"] = noCache
	return send("get_group_info", m)
}

func getGroupList() string {
	return send("get_group_list", nil)
}

func getGroupMemberInfo(groupId int64, userId int64, noCache bool) string {
	m := map[string]any{}
	m["group_id"] = groupId
	m["user_id"] = userId
	m["no_cache"] = noCache
	return send("get_group_member_info", m)
}

func getGroupMemberList(groupId int64) string {
	m := map[string]any{}
	m["group_id"] = groupId
	return send("get_group_member_list", m)
}

func getGroupHonorInfo(groupId int64, typea string) string {
	m := map[string]any{}
	m["group_id"] = groupId
	m["type"] = typea
	return send("get_group_honor_info", m)
}
