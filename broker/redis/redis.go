package redis

import (
	"github.com/kalandramo/lulu-ext/broker"

	"github.com/kalandramo/lulu-ext/broker/redis/option"
	"github.com/kalandramo/lulu-ext/broker/redis/pubsub"
	"github.com/kalandramo/lulu-ext/broker/redis/stream"
)

func NewBroker(driverType option.DriverType, opts ...broker.Option) broker.Broker {
	switch driverType {
	case option.DriverTypeStream:
		return stream.NewBroker(opts...)
	case option.DriverTypePubSub:
		return pubsub.NewBroker(opts...)
	default:
		return pubsub.NewBroker(opts...)
	}
}
