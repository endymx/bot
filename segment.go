package main

import (
	"fmt"
)

// NewTextSegment 文本
func NewTextSegment(text string) Message {
	return Message{
		Type: "text",
		Data: map[string]string{
			"text": text,
		},
	}
}

// NewFaceSegment 表情
func NewFaceSegment(id int) Message {
	return Message{
		Type: "face",
		Data: map[string]string{
			"id": fmt.Sprintf("%d", id),
		},
	}
}

// NewImageSegment 图片
func NewImageSegment(file string, opts ...string) Message {
	msg := Message{
		Type: "image",
		Data: map[string]string{
			"file": file,
		},
	}
	for i, opt := range opts {
		switch i {
		case 0:
			msg.Data.(map[string]string)["type"] = opt
		case 1:
			msg.Data.(map[string]string)["cache"] = opt
		case 2:
			msg.Data.(map[string]string)["proxy"] = opt
		case 3:
			msg.Data.(map[string]string)["timeout"] = opt
		}
	}

	return msg
}

// NewRecordSegment 语音
func NewRecordSegment(file string) Message {
	return Message{
		Type: "record",
		Data: map[string]string{
			"file": file,
		},
	}
}

// NewVideoSegment 视频
func NewVideoSegment(file string) Message {
	return Message{
		Type: "video",
		Data: map[string]string{
			"file": file,
		},
	}
}

// NewAtSegment @某人
func NewAtSegment(qq int64) Message {
	return Message{
		Type: "at",
		Data: map[string]string{
			"qq": fmt.Sprintf("%d", qq),
		},
	}
}

// NewRpsSegment 猜拳魔法表情
func NewRpsSegment() Message {
	return Message{
		Type: "rps",
	}
}

// NewDiceSegment 掷骰子魔法表情
func NewDiceSegment() Message {
	return Message{
		Type: "dice",
	}
}

// NewShakeSegment 窗口抖动（戳一戳）
func NewShakeSegment() Message {
	return Message{
		Type: "shake",
	}
}

// NewPokeSegment 戳一戳
func NewPokeSegment(_type, id int) Message {
	return Message{
		Type: "poke",
		Data: map[string]string{
			"type": fmt.Sprintf("%d", _type),
			"id":   fmt.Sprintf("%d", id),
		},
	}
}

// NewReplySegment 回复
func NewReplySegment(id int32) Message {
	return Message{
		Type: "reply",
		Data: map[string]string{
			"id": fmt.Sprintf("%d", id),
		},
	}
}

// NewNodeSegment 消息节点
func NewNodeSegment[MT MessageType](uin, name string, message MT) Message {
	return Message{
		Type: "node",
		Data: Node[MT]{
			Uin:     uin,
			Name:    name,
			Content: message,
		},
	}
}
