package preload

import (
	"github.com/unix755/xtools/xCrypto"
	"golang.org/x/crypto/chacha20poly1305"
)

func Encrypt(plaintext []byte, key []byte) (ciphertext []byte, err error) {
	// 比特流加密
	if len(key) != 0 {
		key = xCrypto.ZeroPadding(key, chacha20poly1305.KeySize)
		key = key[0:chacha20poly1305.KeySize]
		return xCrypto.NewChaCha20Poly1305(key, []byte{}).Encrypt(plaintext)
	}
	// 密钥为空, 无需加密
	return plaintext, nil
}

func Decrypt(ciphertext []byte, key []byte) (plaintext []byte, err error) {
	// 比特流解密
	if len(key) != 0 {
		key = xCrypto.ZeroPadding(key, chacha20poly1305.KeySize)
		key = key[0:chacha20poly1305.KeySize]
		ciphertext, err = xCrypto.NewChaCha20Poly1305(key, []byte{}).Decrypt(ciphertext)
		if err != nil {
			return nil, err
		}
	}
	// 密钥为空, 无需解密
	return ciphertext, nil
}
