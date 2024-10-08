package main

import (
	"context"
)

type eventHandlers[T any] struct {
	pluginEnableHandler []func(svc T)

	metaConnectHandlers   []func(ctx context.Context, meta *MetaConnect, svc T)
	metaHeartbeatHandlers []func(ctx context.Context, meta *MetaHeartbeat, svc T)

	messagePrivateHandlers []func(ctx context.Context, message *MessagePrivate, svc T)
	messageGroupHandlers   []func(ctx context.Context, message *MessageGroup, svc T)

	noticeFriendAddHandlers           []func(ctx context.Context, notice *NoticeFriendAdd, svc T)
	noticeFriendMessageRecallHandlers []func(ctx context.Context, notice *NoticeFriendMessageRecall, svc T)
	// TODO
}

func (b *Bot[T]) OnPluginEnable(f func(svc T)) {
	b.event.pluginEnableHandler = append(b.event.pluginEnableHandler, f)
}

func (b *Bot[T]) OnMetaConnect(f func(ctx context.Context, meta *MetaConnect, svc T)) {
	b.event.metaConnectHandlers = append(b.event.metaConnectHandlers, f)
}

func (b *Bot[T]) OnMetaHeartbeat(f func(ctx context.Context, meta *MetaHeartbeat, svc T)) {
	b.event.metaHeartbeatHandlers = append(b.event.metaHeartbeatHandlers, f)
}

func (b *Bot[T]) OnMessagePrivate(f func(ctx context.Context, message *MessagePrivate, svc T)) {
	b.event.messagePrivateHandlers = append(b.event.messagePrivateHandlers, f)
}

func (b *Bot[T]) OnMessageGroup(f func(ctx context.Context, message *MessageGroup, svc T)) {
	b.event.messageGroupHandlers = append(b.event.messageGroupHandlers, f)
}

func (b *Bot[T]) OnNoticeFriendAdd(f func(ctx context.Context, notice *NoticeFriendAdd, svc T)) {
	b.event.noticeFriendAddHandlers = append(b.event.noticeFriendAddHandlers, f)
}

func (b *Bot[T]) OnNoticeFriendMessageRecall(f func(ctx context.Context, notice *NoticeFriendMessageRecall, svc T)) {
	b.event.noticeFriendMessageRecallHandlers = append(b.event.noticeFriendMessageRecallHandlers, f)
}
