package infrastructure

import "golang.org/x/crypto/bcrypt"

type PasswordService interface {
	Hash(password string) (string, error)
	Compare(hashedPassword, password string) bool
}

type PasswordServiceImpl struct{}

func NewPasswordService() PasswordService {
	return &PasswordServiceImpl{}
}

// Compare implements PasswordService.
func (p *PasswordServiceImpl) Compare(hashedPassword string, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

// Hash implements PasswordService.
func (p *PasswordServiceImpl) Hash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}
