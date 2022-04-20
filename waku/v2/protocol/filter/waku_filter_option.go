package filter

import (
	"context"
	"time"

	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/status-im/go-waku/waku/v2/utils"
	"go.uber.org/zap"
)

type (
	FilterSubscribeParameters struct {
		Host         host.Host
		SelectedPeer peer.ID
		Log          *zap.SugaredLogger
	}

	FilterSubscribeOption func(*FilterSubscribeParameters)

	FilterParameters struct {
		Timeout time.Duration
	}

	Option func(*FilterParameters)
)

func WithTimeout(timeout time.Duration) Option {
	return func(params *FilterParameters) {
		params.Timeout = timeout
	}
}

func WithPeer(p peer.ID) FilterSubscribeOption {
	return func(params *FilterSubscribeParameters) {
		params.SelectedPeer = p
	}
}

func WithAutomaticPeerSelection() FilterSubscribeOption {
	return func(params *FilterSubscribeParameters) {
		p, err := utils.SelectPeer(params.Host, string(FilterID_v20beta1), params.Log)
		if err == nil {
			params.SelectedPeer = *p
		} else {
			params.Log.Info("Error selecting peer: ", err)
		}
	}
}

func WithFastestPeerSelection(ctx context.Context) FilterSubscribeOption {
	return func(params *FilterSubscribeParameters) {
		p, err := utils.SelectPeerWithLowestRTT(ctx, params.Host, string(FilterID_v20beta1), params.Log)
		if err == nil {
			params.SelectedPeer = *p
		} else {
			params.Log.Info("Error selecting peer: ", err)
		}
	}
}

func DefaultOptions() []Option {
	return []Option{
		WithTimeout(24 * time.Hour),
	}
}

func DefaultSubscribtionOptions() []FilterSubscribeOption {
	return []FilterSubscribeOption{
		WithAutomaticPeerSelection(),
	}
}
