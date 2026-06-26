package transport

import "github.com/kalandramo/lulu-ext/broker"

// SubscribeOption holds deferred subscription parameters.
type SubscribeOption struct {
	Handler          broker.Handler
	Binder           broker.Binder
	SubscribeOptions []broker.SubscribeOption
}

// SubscribeOptionMap maps topic names to their SubscribeOption.
type SubscribeOptionMap map[string]*SubscribeOption
