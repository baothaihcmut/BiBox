package services

import (
	"context"

	"golang.org/x/crypto/bcrypt"
)

type PasswordService interface {
	HashPassword(ctx context.Context, password string) (string, error)
	ComparePassword(ctx context.Context, hashedPassword, password string) error
}

type PasswordServiceImpl struct{}

func NewPasswordService() PasswordService {
	return &PasswordServiceImpl{}
}
func (p *PasswordServiceImpl) HashPassword(ctx context.Context, password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}
func (p *PasswordServiceImpl) ComparePassword(ctx context.Context, hashedPassword, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return err
	}
	return nil
}
