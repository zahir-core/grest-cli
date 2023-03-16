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
	return c
}

func (c *cryptoImpl) configure() {
	c.Key = CRYPTO_KEY
	c.Salt = CRYPTO_SALT
	c.Info = CRYPTO_INFO
}

func (c *cryptoImpl) Encrypt(text string) (string, error) {
	encrypted, err := grest.NewCrypto(c.Key, c.Salt, c.Info).Encrypt(text)
	if err != nil {
		return encrypted, err
	}
	return CRYPTO_PREFIX + encrypted, nil
}

func (c *cryptoImpl) Decrypt(text string) (string, error) {
	prefixLength := len([]rune(CRYPTO_PREFIX))
	if CRYPTO_PREFIX != "" && CRYPTO_PREFIX == text[:prefixLength] {
		return grest.NewCrypto(c.Key, c.Salt, c.Info).Decrypt(text[prefixLength:])
	}
	return grest.NewCrypto(c.Key, c.Salt, c.Info).Decrypt(text)
}
