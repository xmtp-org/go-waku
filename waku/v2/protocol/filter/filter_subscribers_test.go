package filter

import (
	"testing"
	"time"

	"github.com/libp2p/go-libp2p-core/test"
	"github.com/status-im/go-waku/waku/v2/protocol/pb"
	"github.com/stretchr/testify/assert"
)

const TOPIC = "/test/topic"

func TestAppend(t *testing.T) {
	subs := NewSubscribers(10 * time.Second)
	peerId, err := test.RandPeerID()
	assert.NoError(t, err)
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
		break
	}
	assert.True(t, hasMatch)
}

func TestRemove(t *testing.T) {
	subs := NewSubscribers(10 * time.Second)
	peerId, err := test.RandPeerID()
	assert.NoError(t, err)
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
		break
	}
	assert.False(t, hasMatch)
}

func TestRemoveBogus(t *testing.T) {
	subs := NewSubscribers(10 * time.Second)
	peerId, err := test.RandPeerID()
	assert.NoError(t, err)
	contentTopic := "topic1"
	request := pb.FilterRequest{
		Subscribe:      true,
		Topic:          TOPIC,
		ContentFilters: []*pb.FilterRequest_ContentFilter{{ContentTopic: contentTopic}},
	}
	subs.Append(Subscriber{peerId, "request_1", request})
	// This will panic with this error:
	// panic: runtime error: index out of range [1] with length 1
	subs.RemoveContentFilters(peerId, []*pb.FilterRequest_ContentFilter{{ContentTopic: "does not exist"}, {ContentTopic: contentTopic}})
}
