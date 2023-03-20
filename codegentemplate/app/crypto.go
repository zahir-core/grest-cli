package app

import "grest.dev/grest"

func Crypto() CryptoInterface {
	if crpto == nil {
		crpto = &cryptoImpl{}
		crpto.configure()
	}
	return crpto
}

type CryptoInterface interface {
	grest.CryptoInterface
}

var crpto *cryptoImpl

// cryptoImpl implement CryptoInterface embed from grest.cryptoImpl for simplicity
type cryptoImpl struct {
	grest.Crypto
}

func NewCrypto(keys ...string) *cryptoImpl {
	c := &cryptoImpl{}
	c.configure()
	if len(keys) > 0 {
		c.Key = keys[0]
	}
	if len(keys) > 1 {
		c.Salt = keys[1]
	}
	if len(keys) > 2 {
		c.Info = keys[2]
	}
	if len(keys) > 3 {
		c.JWTKey = keys[3]
	}
	return c
}

func (c *cryptoImpl) configure() {
	c.Key = CRYPTO_KEY
	c.Salt = CRYPTO_SALT
	c.Info = CRYPTO_INFO
	c.JWTKey = JWT_KEY
}
