package filter

import (
	"testing"
	"time"

	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/test"
	"github.com/status-im/go-waku/waku/v2/protocol/pb"
	"github.com/stretchr/testify/assert"
)

const TOPIC = "/test/topic"

func createPeerId(t *testing.T) peer.ID {
	peerId, err := test.RandPeerID()
	assert.NoError(t, err)
	return peerId
}

func TestAppend(t *testing.T) {
	subs := NewSubscribers(10 * time.Second)
	peerId := createPeerId(t)
	contentTopic := "topic1"
	request := pb.FilterRequest{
		Subscribe:      true,
		Topic:          TOPIC,
		ContentFilters: []*pb.FilterRequest_ContentFilter{{ContentTopic: contentTopic}},
	}
	subs.Append(Subscriber{peerId, "request_1", request})

	var hasMatch bool
	for range subs.Items(&contentTopic) {
		hasMatch = true
	}
	assert.True(t, hasMatch)
}

func TestRemove(t *testing.T) {
	subs := NewSubscribers(10 * time.Second)
	peerId := createPeerId(t)
	contentTopic := "topic1"
	request := pb.FilterRequest{
		Subscribe:      true,
		Topic:          TOPIC,
		ContentFilters: []*pb.FilterRequest_ContentFilter{{ContentTopic: contentTopic}},
	}
	subs.Append(Subscriber{peerId, "request_1", request})
	subs.RemoveContentFilters(peerId, request.ContentFilters)

	var hasMatch bool
	for range subs.Items(&contentTopic) {
		hasMatch = true
	}
	assert.False(t, hasMatch)
}

func TestRemovePartial(t *testing.T) {
	subs := NewSubscribers(10 * time.Second)
	peerId := createPeerId(t)
	topic1 := "topic1"
	topic2 := "topic2"
	request := pb.FilterRequest{
		Subscribe:      true,
		Topic:          TOPIC,
		ContentFilters: []*pb.FilterRequest_ContentFilter{{ContentTopic: topic1}, {ContentTopic: topic2}},
	}
	subs.Append(Subscriber{peerId, "request_1", request})
	subs.RemoveContentFilters(peerId, []*pb.FilterRequest_ContentFilter{{ContentTopic: topic1}})

	var hasMatch bool
	for sub := range subs.Items(&topic2) {
		hasMatch = true
		assert.Len(t, sub.filter.ContentFilters, 1)
	}
	assert.True(t, hasMatch)
}

func TestRemoveBogus(t *testing.T) {
	subs := NewSubscribers(10 * time.Second)
	peerId := createPeerId(t)
	contentTopic := "topic1"
	request := pb.FilterRequest{
		Subscribe:      true,
		Topic:          TOPIC,
		ContentFilters: []*pb.FilterRequest_ContentFilter{{ContentTopic: contentTopic}},
	}
	subs.Append(Subscriber{peerId, "request_1", request})
	subs.RemoveContentFilters(peerId, []*pb.FilterRequest_ContentFilter{{ContentTopic: "does not exist"}, {ContentTopic: contentTopic}})

	var hasMatch bool
	for range subs.Items(&contentTopic) {
		hasMatch = true
	}
	assert.False(t, hasMatch)
}
