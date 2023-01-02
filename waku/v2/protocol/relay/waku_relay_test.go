package relay

import (
	"context"
	"crypto/rand"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/waku-org/go-waku/tests"
	"github.com/waku-org/go-waku/waku/v2/protocol/pb"
	"github.com/waku-org/go-waku/waku/v2/timesource"
	"github.com/waku-org/go-waku/waku/v2/utils"
)

func TestWakuRelay(t *testing.T) {
	testTopic := "/waku/2/go/relay/test"

	port, err := tests.FindFreePort(t, "", 5)
	require.NoError(t, err)

	host, err := tests.MakeHost(context.Background(), port, rand.Reader)
	require.NoError(t, err)

	relay, err := NewWakuRelay(context.Background(), host, nil, 0, timesource.NewDefaultClock(), utils.Logger())
	defer relay.Stop()
	require.NoError(t, err)

	sub, err := relay.subscribe(testTopic)
	defer sub.Cancel()
	require.NoError(t, err)

	topics := relay.Topics()
	require.Equal(t, 1, len(topics))
	require.Equal(t, testTopic, topics[0])

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		defer cancel()

		_, err := sub.Next(ctx)
		require.NoError(t, err)
	}()

	msg := &pb.WakuMessage{
		Payload:      []byte{1},
		Version:      0,
		ContentTopic: "test",
		Timestamp:    0,
	}
	_, err = relay.PublishToTopic(context.Background(), msg, testTopic)
	require.NoError(t, err)

	<-ctx.Done()
}

func TestWakuRelay_Unsubscribe(t *testing.T) {
	testTopic := "/waku/2/go/relay/test"

	port, err := tests.FindFreePort(t, "", 5)
	require.NoError(t, err)

	host, err := tests.MakeHost(context.Background(), port, rand.Reader)
	require.NoError(t, err)

	relay, err := NewWakuRelay(context.Background(), host, nil, 0, timesource.NewDefaultClock(), utils.Logger())
	defer relay.Stop()
	require.NoError(t, err)

	sub, err := relay.SubscribeToTopic(context.Background(), testTopic)
	require.NoError(t, err)

	sub.Unsubscribe()
}
