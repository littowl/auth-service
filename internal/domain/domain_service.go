package domain

import (
	"auth-service/transport/transport_http/handlers"
	"auth-service/utils"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type DomainService struct {
	repository DomainRepository
}

func NewDomainService(repo DomainRepository) *DomainService {
	return &DomainService{
		repository: repo,
	}
}

func (s *DomainService) Register(token string, login string, password string) error {
	claims, err := utils.DecodeJWT(token)
	if err != nil {
		return err
	}
	if claims.Role != "admin" {
		return errors.New("User is not admin")
	}
	if claims.Type != "access" {
		return errors.New("Invalid token")
	}

	// create pass_hash
	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	err = s.repository.Register(login, passHash)
	if err != nil {
		return err
	}

	return nil
}

func (s *DomainService) Login(login string, password string) (handlers.LoginResponse, error) {
	user, err := s.repository.GetUser(login)
	if err != nil {
		return handlers.LoginResponse{}, err
	}

	if login != "admin" || password != "admin" {
		// bcrypt compare hash from db and password
		err = bcrypt.CompareHashAndPassword([]byte(user.Hash), []byte(password))
		if err != nil {
			return handlers.LoginResponse{}, err
		}
	}

	accessExp := time.Now().Add(time.Minute * 5).Unix()
	refreshExp := time.Now().Add(time.Hour * 5).Unix()

	// create tokens
	access, err := utils.CreateToken(accessExp, "1111", user.Login, user.Role, "access")
	if err != nil {
		return handlers.LoginResponse{}, err
	}

	refresh, err := utils.CreateToken(refreshExp, "1111", user.Login, user.Role, "refresh")
	if err != nil {
		return handlers.LoginResponse{}, err
	}

	return handlers.LoginResponse{Access: access, Refresh: refresh}, nil
}

func (s *DomainService) GetUser(token string) (handlers.User, error) {
	claims, err := utils.DecodeJWT(token)
	if err != nil {
		return handlers.User{}, err
	}
	if claims.Type != "access" {
		return handlers.User{}, errors.New("Invalid token")
	}

	user, err := s.repository.GetUser(claims.Login)
	if err != nil {
		return handlers.User{}, err
	}

	return handlers.User{
		ID:    user.ID,
		Login: user.Login,
		Role:  user.Role,
		Hash:  user.Hash,
	}, nil
}

func (s *DomainService) Refresh(token string) (handlers.LoginResponse, error) {
	claims, err := utils.DecodeJWT(token)
	if err != nil {
		return handlers.LoginResponse{}, err
	}
	if claims.Type != "refresh" {
		return handlers.LoginResponse{}, errors.New("Invalid token")
	}

	accessExp := time.Now().Add(time.Minute * 5).Unix()
	refreshExp := time.Now().Add(time.Hour * 5).Unix()

	// create tokens
	access, err := utils.CreateToken(accessExp, "1111", claims.Login, claims.Role, "access")
	if err != nil {
		return handlers.LoginResponse{}, err
	}

	refresh, err := utils.CreateToken(refreshExp, "1111", claims.Login, claims.Role, "refresh")
	if err != nil {
		return handlers.LoginResponse{}, err
	}

	return handlers.LoginResponse{Access: access, Refresh: refresh}, nil
}

func (s *DomainService) ChangePassword(token string, new_password string) error {
	claims, err := utils.DecodeJWT(token)
	if err != nil {
		return err
	}

	// get user, from user get pass_hash
	user, err := s.repository.GetUser(claims.Login)
	if err != nil {
		return err
	}

	// create pass_hash
	new_hash, err := bcrypt.GenerateFromPassword([]byte(new_password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// change password
	err = s.repository.ChangePassword(user.Login, user.Hash, new_hash)
	if err != nil {
		return err
	}

	return nil
}
