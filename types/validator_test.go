package types

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidatorProtoBuf(t *testing.T) {
	val, _ := RandValidator(true, 100)
	testCases := []struct {
		msg      string
		v1       *Validator
		expPass1 bool
		expPass2 bool
	}{
		{"success validator", val, true, true},
		{"failure empty", &Validator{}, false, false},
		{"failure nil", nil, false, false},
	}
	for _, tc := range testCases {
		protoVal, err := tc.v1.ToProto()

		if tc.expPass1 {
			require.NoError(t, err, tc.msg)
		} else {
			require.Error(t, err, tc.msg)
		}

		val, err := ValidatorFromProto(protoVal)
		if tc.expPass2 {
			require.NoError(t, err, tc.msg)
			require.Equal(t, tc.v1, val, tc.msg)
		} else {
			require.Error(t, err, tc.msg)
		}
	}
}

func TestValidatorValidateBasic(t *testing.T) {
	priv := NewMockPV()
	pubKey, _ := priv.GetPubKey()

	blsPubKey := make([]byte, BlsPubKeySize)
	rand.Read(blsPubKey)
	relayer := make([]byte, AddressSize)
	rand.Read(relayer)
	challenger := make([]byte, AddressSize)
	rand.Read(challenger)

	val := NewValidator(pubKey, 1)
	val.SetBlsKey(blsPubKey)
	val.SetRelayerAddress(relayer)
	val.SetChallengerAddress(challenger)

	wrongBlsPubKey := val.Copy()
	wrongBlsPubKey.SetBlsKey([]byte{'a'})

	wrongRelayer := val.Copy()
	wrongRelayer.SetRelayerAddress([]byte{'a'})

	wrongChallenger := val.Copy()
	wrongChallenger.SetChallengerAddress([]byte{'b'})

	testCases := []struct {
		val *Validator
		err bool
		msg string
	}{
		{
			val: NewValidator(pubKey, 1),
			err: false,
			msg: "",
		},
		{
			val: nil,
			err: true,
			msg: "nil validator",
		},
		{
			val: &Validator{
				PubKey: nil,
			},
			err: true,
			msg: "validator does not have a public key",
		},
		{
			val: NewValidator(pubKey, -1),
			err: true,
			msg: "validator has negative voting power",
		},
		{
			val: &Validator{
				PubKey:  pubKey,
				Address: nil,
			},
			err: true,
			msg: "validator address is the wrong size: ",
		},
		{
			val: &Validator{
				PubKey:  pubKey,
				Address: []byte{'a'},
			},
			err: true,
			msg: "validator address is the wrong size: 61",
		},
		{
			val: val,
			err: false,
			msg: "",
		},
		{
			val: wrongBlsPubKey,
			err: true,
			msg: "validator relayer bls key is the wrong size: [97]",
		},
		{
			val: wrongRelayer,
			err: true,
			msg: "validator relayer address is the wrong size: [97]",
		},
		{
			val: wrongChallenger,
			err: true,
			msg: "validator challenger address is the wrong size: [98]",
		},
	}

	for _, tc := range testCases {
		err := tc.val.ValidateBasic()
		if tc.err {
			if assert.Error(t, err) {
				assert.Equal(t, tc.msg, err.Error())
			}
		} else {
			assert.NoError(t, err)
		}
	}
}
