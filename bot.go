package main

import (
	"context"
	"fmt"
)

type Bot[T any] struct {
	svc    T
	config Config

	plugins []Plugin
	event   *eventHandlers[T]

	api *API
}

type Plugin struct {
	Name string
}

func Create[T any](svc IService[T], config Config) *Bot[T] {
	bot := &Bot[T]{
		svc:    svc.GetService(),
		config: config,

		event: new(eventHandlers[T]),

		api: api,
	}

	if config.Logger == nil {
		config.Logger = SimpleLogger()
	}

	if config.Addr == "" {
		config.Logger.Fatal("无法获得QQ的ws地址")
		return nil
	}
	bot.websocket()

	for _, plugin := range bot.plugins {
		config.Logger.Infof("插件启用：%s", plugin.Name)
	}
	for _, handler := range bot.event.pluginEnableHandler {
		go bot.cover(func(ctx context.Context) { handler(bot.svc) })
	}

	return bot
}

func (b *Bot[T]) GetActionId() int64 {
	return b.api.actionId
}

func (b *Bot[T]) Register(name string) error {
	for _, plugin := range b.plugins {
		if plugin.Name == name {
			return fmt.Errorf("插件重复：%s", name)
		}
	}
	return nil
}
