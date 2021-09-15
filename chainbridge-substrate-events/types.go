package events

import (
	"fmt"

	"github.com/selendra/selendra-bridge/go-substrate-rpc-client/v3/types"
	"github.com/vedhavyas/go-subkey/scale"
)

type BlockRewardInfo struct {
	Seed          types.U256 `json:"seed"`
	OnlineTarget  types.U256 `json:"onlineTarget"`
	ComputeTarget types.U256 `json:"computeTarget"`
}

type PayoutReason byte

const (
	OnlineReward  PayoutReason = 0
	ComputeReward PayoutReason = 1
)

func (v *PayoutReason) Decode(decoder scale.Decoder) error {
	b, err := decoder.ReadOneByte()
	vb := PayoutReason(b)
	switch vb {
	case OnlineReward, ComputeReward:
		*v = vb
	default:
		return fmt.Errorf("unknown VoteThreshold enum: %v", vb)
	}
	return err
}

func (v PayoutReason) Encode(encoder scale.Encoder) error {
	return encoder.PushByte(byte(v))
}
