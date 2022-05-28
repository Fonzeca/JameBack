package guard_userhub

import (
	"crypto/rand"
)

type KeyGeneratorUserHub struct {
}

func (t KeyGeneratorUserHub) SecureRandomBytes(length int) ([]byte, error) {
	var randomBytes = make([]byte, length)
	_, err := rand.Read(randomBytes)
	return randomBytes, err
}
