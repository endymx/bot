package main

import "github.com/tidwall/gjson"

// Event 事件
type Event struct {
	Time     float64 `json:"time"`
	SelfId   int64   `json:"self_id"`
	PostType string  `json:"post_type"`
}

// Status 状态
type Status struct {
	Online bool `json:"online"`
	Good   bool `json:"good"`
}

// ActionReq 动作请求
type ActionReq struct {
	Action string         `json:"action"`
	Params map[string]any `json:"params,omitempty"`
	Echo   string         `json:"echo,omitempty"`
}

// ActionResp 动作响应
type ActionResp struct {
	Status  string                  `json:"status"`
	Retcode int                     `json:"retcode"`
	Data    map[string]gjson.Result `json:"-"`
	Echo    string                  `json:"echo,omitempty"`
}

// MetaConnect 连接
type MetaConnect struct {
	Event
	MetaEventType string `json:"meta_event_type"`
	SubType       string `json:"sub_type"`
}

// MetaHeartbeat 心跳
type MetaHeartbeat struct {
	Event
	MetaEventType string `json:"meta_event_type"`
	Status        Status `json:"status"`
	Interval      int    `json:"interval"`
}

// MessageType 支持的消息类型
type MessageType interface {
	Message | []Message | string
}

// Message 消息
type Message struct {
	Type string `json:"type"`
	Data any    `json:"data"`
}

type Sender struct {
	UserId   int64  `json:"user_id"`
	NickName string `json:"nick_name"`
	Card     string `json:"card"`
	Sex      string `json:"sex"`
	Age      int32  `json:"age"`
	Area     string `json:"area"`
	Level    string `json:"level"`
	Role     string `json:"role"`
	Title    string `json:"title"`
}

type Anonymous struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
	Flag string `json:"flag"`
}

// MessagePrivate 私聊消息
type MessagePrivate struct {
	Event
	MessageType string    `json:"message_type"`
	SubType     string    `json:"sub_type"`
	MessageId   int32     `json:"message_id"`
	UserId      int64     `json:"user_id"`
	Message     []Message `json:"message"`
	RawMessage  string    `json:"raw_message"`
	Font        int32     `json:"font"`
	Sender      Sender    `json:"sender"`
}

// MessageGroup 群消息
type MessageGroup struct {
	Event
	MessageType string    `json:"message_type"`
	SubType     string    `json:"sub_type"`
	MessageId   int32     `json:"message_id"`
	UserId      int64     `json:"user_id"`
	GroupId     int64     `json:"group_id"`
	Anonymous   Anonymous `json:"anonymous"`
	Message     []Message `json:"message"`
	RawMessage  string    `json:"raw_message"`
	Font        int32     `json:"font"`
	Sender      Sender    `json:"sender"`
}

type NoticeEvent struct {
	Event
	NoticeType string `json:"notice_type"`
}

// NoticeFriendAdd 好友添加
type NoticeFriendAdd struct {
	NoticeEvent
	UserId int64 `json:"user_id"`
}

// NoticeFriendMessageRecall 好友消息撤回
type NoticeFriendMessageRecall struct {
	NoticeEvent
	GroupId    int64  `json:"group_id"`
	MessageId  string `json:"message_id"`
	UserId     int64  `json:"user_id"`
	OperatorId int64  `json:"operator_id"`
}

// NoticePrivateMessageDelete 私聊消息删除
type NoticePrivateMessageDelete struct {
	NoticeEvent
	MessageId string `json:"message_id"`
	UserId    int64  `json:"user_id"`
}

// NoticeGroupMember 群成员增加/减少
type NoticeGroupMember struct {
	NoticeEvent
	SubType    string `json:"sub_type"`
	UserId     int64  `json:"user_id"`
	GroupId    int64  `json:"group_id"`
	OperatorId int64  `json:"operator_id"`
}

// NoticeGroupMessageRecall 群消息撤回
type NoticeGroupMessageRecall struct {
	NoticeEvent
	GroupId    int64  `json:"group_id"`
	MessageId  string `json:"message_id"`
	UserId     int64  `json:"user_id"`
	OperatorId int64  `json:"operator_id"`
}

type File struct {
	Id    string `json:"id"`
	Name  string `json:"name"`
	Size  int64  `json:"size"`
	Busid int64  `json:"busid"`
}

// NoticeGroupFileUpload 群文件上传
type NoticeGroupFileUpload struct {
	NoticeEvent
	GroupId int64 `json:"group_id"`
	UserId  int64 `json:"user_id"`
	File    File  `json:"file"`
}

// NoticeGroupAdminChange 群管理员变动
type NoticeGroupAdminChange struct {
	NoticeEvent
	SubType    string `json:"sub_type"`
	GroupId    int64  `json:"group_id"`
	UserId     int64  `json:"user_id"`
	OperatorId int64  `json:"operator_id"`
}

// NoticeGroupMute 群禁言
type NoticeGroupMute struct {
	NoticeEvent
	SubType    string `json:"sub_type"`
	GroupId    int64  `json:"group_id"`
	UserId     int64  `json:"user_id"`
	OperatorId int64  `json:"operator_id"`
	Duration   int64  `json:"duration"`
}

// NoticeGroupPoke 群内戳一戳
type NoticeGroupPoke struct {
	NoticeEvent
	SubType  string `json:"sub_type"`
	GroupId  int64  `json:"group_id"`
	UserId   int64  `json:"user_id"`
	TargetId int64  `json:"target_id"`
}

// NoticeGroupRedEnvelopeLuckKing 群红包运气王
type NoticeGroupRedEnvelopeLuckKing struct {
	NoticeEvent
	SubType  string `json:"sub_type"`
	GroupId  int64  `json:"group_id"`
	UserId   int64  `json:"user_id"`
	TargetId int64  `json:"target_id"`
}

// NoticeGroupMemberHonorChanged 群成员荣誉变更
type NoticeGroupMemberHonorChanged struct {
	NoticeEvent
	SubType   string `json:"sub_type"`
	GroupId   int64  `json:"group_id"`
	UserId    int64  `json:"user_id"`
	HonorType string `json:"honor_type"`
}

// RequestFriendAdd 加好友请求
type RequestFriendAdd struct {
	Event
	RequestType string `json:"request_type"`
	UserId      int64  `json:"user_id"`
	Comment     string `json:"comment"`
	Flag        string `json:"flag"`
}

// RequestGroupAdd 加群请求／邀请
type RequestGroupAdd struct {
	Event
	RequestType string `json:"request_type"`
	SubType     string `json:"sub_type"`
	GroupId     int64  `json:"group_id"`
	UserId      int64  `json:"user_id"`
	Comment     string `json:"comment"`
	Flag        string `json:"flag"`
}

// GetSelfInfo 获取机器人自身信息
type GetSelfInfo struct {
	ActionResp
	Data struct {
		UserId          int64  `json:"user_id"`
		UserName        string `json:"user_name"`
		UserDisplayname string `json:"user_displayname"`
	} `json:"data"`
}

// GetUserInfo 获取用户信息
type GetUserInfo struct {
	ActionResp
	Data struct {
		UserId          int64  `json:"user_id"`
		UserName        string `json:"user_name"`
		UserDisplayname string `json:"user_displayname"`
		UserRemark      string `json:"user_remark"`
	} `json:"data"`
}

// GetFriendList 获取好友列表
type GetFriendList struct {
	ActionResp
	Data []struct {
		UserId          int64  `json:"user_id"`
		UserName        string `json:"user_name"`
		UserDisplayname string `json:"user_displayname"`
		UserRemark      string `json:"user_remark"`
	} `json:"data"`
}

// Node 合并消息节点
type Node[MT MessageType] struct {
	Uin     string `json:"uin"`
	Name    string `json:"name"`
	Content MT     `json:"content"`
}
