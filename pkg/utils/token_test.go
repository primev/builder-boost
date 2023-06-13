package utils

import (
	"crypto/ecdsa"
	"testing"

	"github.com/alecthomas/assert"
	"github.com/ethereum/go-ethereum/crypto"
)

func TestGenerateAuthenticationToken(t *testing.T) {
	searcherKey, err := crypto.GenerateKey()
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	type args struct {
		msg string
		key *ecdsa.PrivateKey
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "Normal Test",
			args: args{msg: "hello", key: searcherKey},
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			token, err := GenerateAuthenticationToken(tt.args.msg, tt.args.key)
			if err != nil {
				t.Error(err)
				t.FailNow()
			}
			address, ok := VerifyAuthenticationToken(token, tt.args.msg)

			assert.Equal(t, ok, true)
			assert.Equal(t, address.Hex(), crypto.PubkeyToAddress(tt.args.key.PublicKey).Hex())
		})
	}
}
