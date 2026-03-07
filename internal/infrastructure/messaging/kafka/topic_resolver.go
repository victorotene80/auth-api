package kafka

import "strings"

type TopicResolver struct {
	Prefix string
}

func (r TopicResolver) TopicForMessage(name string) string {
	if r.Prefix == "" {
		return name
	}
	return r.Prefix + name
}

func (r TopicResolver) NormalizeSubscription(name string) string {
	if r.Prefix != "" && strings.HasPrefix(name, r.Prefix) {
		return name
	}
	return r.TopicForMessage(name)
}