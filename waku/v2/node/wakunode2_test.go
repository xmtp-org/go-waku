package node

import (
	"context"
	"net"
	"testing"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/status-im/go-waku/tests"
	"github.com/stretchr/testify/require"
)

func TestWakuNode2(t *testing.T) {
	hostAddr, _ := net.ResolveTCPAddr("tcp", "0.0.0.0:0")

	key, err := tests.RandomHex(32)
	require.NoError(t, err)
	prvKey, err := crypto.HexToECDSA(key)
	require.NoError(t, err)

	ctx := context.Background()

	wakuNode, err := New(ctx,
		WithPrivateKey(prvKey),
		WithHostAddress(hostAddr),
		WithWakuRelay(),
	)
	require.NoError(t, err)

	err = wakuNode.Start()
	defer wakuNode.Stop()

	require.NoError(t, err)
}

/*

The test is flaky. Sometimes it gets through all 1100 messages in 2-3 seconds, sometimes it gets stuck anything from 65 to 900 after 20s. All messages publish without error, but not all messages are read. I don't quite know why, but it's obviously something serious we are going to want to get to the bottom of. Also seems to work more reliably on my super fast laptop than in CI, which is maybe a hint.

I can bump the number of messages to 2k, and it will get stuck at a number > 1057, which is why I think it's a separate issue from the one I identified (which always got stuck at the exact same point).

*/
// func Test1100(t *testing.T) {
// 	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
// 	defer cancel()

// 	hostAddr1, _ := net.ResolveTCPAddr("tcp", "0.0.0.0:0")
// 	key1, _ := tests.RandomHex(32)
// 	prvKey1, _ := crypto.HexToECDSA(key1)

// 	hostAddr2, _ := net.ResolveTCPAddr("tcp", "0.0.0.0:0")
// 	key2, _ := tests.RandomHex(32)
// 	prvKey2, _ := crypto.HexToECDSA(key2)

// 	wakuNode1, err := New(ctx,
// 		WithPrivateKey(prvKey1),
// 		WithHostAddress(hostAddr1),
// 		WithWakuRelay(),
// 	)
// 	require.NoError(t, err)
// 	err = wakuNode1.Start()
// 	require.NoError(t, err)
// 	defer wakuNode1.Stop()

// 	wakuNode2, err := New(ctx,
// 		WithPrivateKey(prvKey2),
// 		WithHostAddress(hostAddr2),
// 		WithWakuRelay(),
// 	)
// 	require.NoError(t, err)
// 	err = wakuNode2.Start()
// 	require.NoError(t, err)
// 	defer wakuNode2.Stop()

// 	err = wakuNode2.DialPeer(ctx, wakuNode1.ListenAddresses()[0].String())
// 	require.NoError(t, err)

// 	time.Sleep(2 * time.Second)

// 	sub1, err := wakuNode1.Relay().Subscribe(ctx)
// 	require.NoError(t, err)
// 	sub2, err := wakuNode1.Relay().Subscribe(ctx)
// 	require.NoError(t, err)

// 	wg := sync.WaitGroup{}

// 	wg.Add(3)
// 	go func() {
// 		defer wg.Done()

// 		ticker := time.NewTimer(20 * time.Second)
// 		defer ticker.Stop()

// 		msgCnt := 0
// 		for {
// 			select {
// 			case <-ticker.C:
// 				if msgCnt != 1100 {
// 					require.Fail(t, "Timeout Sub1", msgCnt)
// 				}
// 			case <-sub1.C:
// 				msgCnt++
// 				if msgCnt == 1100 {
// 					return
// 				}
// 			}
// 		}
// 	}()

// 	go func() {
// 		defer wg.Done()

// 		ticker := time.NewTimer(20 * time.Second)
// 		defer ticker.Stop()

// 		msgCnt := 0
// 		for {
// 			select {
// 			case <-ticker.C:
// 				if msgCnt != 1100 {
// 					require.Fail(t, "Timeout Sub2", msgCnt)
// 				}
// 			case <-sub2.C:
// 				msgCnt++
// 				if msgCnt == 1100 {
// 					return
// 				}
// 			}
// 		}
// 	}()

// 	go func() {
// 		defer wg.Done()
// 		for i := 1; i <= 1100; i++ {
// 			msg := createTestMsg(0)
// 			msg.Payload = []byte(fmt.Sprint(i))
// 			msg.Timestamp = int64(i)
// 			if err := wakuNode2.Publish(ctx, msg); err != nil {
// 				require.Fail(t, "Could not publish all messages")
// 			}
// 		}
// 	}()

// 	wg.Wait()

// }
