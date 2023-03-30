package app

import (
	"strings"

	"github.com/google/uuid"
	"grest.dev/grest"
)

func Crypto() CryptoInterface {
	if crpto == nil {
		crpto = &cryptoUtil{}
		crpto.configure()
	}
	return crpto
}

type CryptoInterface interface {
	NewToken() string
	NewHash(text string, cost ...int) (string, error)
	CompareHash(hashed, text string) error
	NewJWT(claims any) (string, error)
	ParseAndVerifyJWT(token string, claims any) error
	Encrypt(text string) (string, error)
	Decrypt(text string) (string, error)
	GenerateKey() ([]byte, error)
	PKCS5Padding(ciphertext []byte, blockSize int) []byte
	PKCS5Unpadding(encrypt []byte) ([]byte, error)
}

var crpto *cryptoUtil

// cryptoUtil implement CryptoInterface embed from grest.cryptoUtil for simplicity
type cryptoUtil struct {
	grest.Crypto
}

func NewCrypto(keys ...string) *cryptoUtil {
	c := &cryptoUtil{}
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

func (c *cryptoUtil) configure() {
	c.Key = CRYPTO_KEY
	c.Salt = CRYPTO_SALT
	c.Info = CRYPTO_INFO
	c.JWTKey = JWT_KEY
}

func (c *cryptoUtil) NewToken() string {
	return strings.ReplaceAll(uuid.NewString(), "-", "")
}
