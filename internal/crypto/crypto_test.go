package crypto

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateProvateKey(t *testing.T) {
	privateKey := GeneratePrivateKey()

	assert.Equal(t, len(privateKey.Bytes()), privateKeyLen)

	publicKey := privateKey.PublicKey()

	assert.Equal(t, len(publicKey.Bytes()), publicKeyLen)
}
