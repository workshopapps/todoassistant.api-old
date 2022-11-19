package cryptoService

import "golang.org/x/crypto/bcrypt"

type CryptoSrv interface {
	HashPassword(password string) (string, error)
	ComparePassword(hashed, plain string) error
}

type cryptoSrv struct {
}

func (c cryptoSrv) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

func (c cryptoSrv) ComparePassword(hashed, plain string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(plain))
}

func NewCryptoSrv() CryptoSrv {
	return &cryptoSrv{}
}
