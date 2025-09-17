package internal

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/franzego/user-service/authentication"
	"github.com/jackc/pgx/v5/pgtype"
	"golang.org/x/crypto/bcrypt"
)

type Servicer interface {
	SignUp(ctx context.Context, username string, email string, passwordhash string, firstname pgtype.Text, lastname pgtype.Text) (string, error)
	Login(ctx context.Context, email string, passwordhash string) (string, error)
}

type Service struct {
	RepoServicer
	authentication.Auth
}

func NewService(svc RepoServicer) *Service {
	return &Service{RepoServicer: svc}
}
func (s *Service) SignUp(ctx context.Context, username string, email string,
	passwordhash string, firstname pgtype.Text, lastname pgtype.Text) (string, error) {
	hashpwd, err := bcrypt.GenerateFromPassword([]byte(passwordhash), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
	}
	user, err := s.RepoServicer.CreateUser(ctx, username, email, string(hashpwd), firstname, lastname)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate keys") {
			return "", errors.New("user already exists")
		}
		return "", err
	}

	tok, err := s.Auth.CreateToken(user.Username)
	if err != nil {
		log.Fatal(err)
	}
	return tok, nil

}

func (s *Service) Login(ctx context.Context, email string, passwordhash string) (string, error) {
	//get user from db first
	user, err := s.RepoServicer.GetUserByEmail(context.Background(), email, passwordhash)
	if err != nil {
		return "", fmt.Errorf("invalid credentials")
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(passwordhash))
	if err != nil {
		fmt.Println("invalid passwords")
	}
	tok, err := s.Auth.CreateToken(user.Username)
	if err != nil {
		log.Fatal(err)
	}
	return tok, nil
}
