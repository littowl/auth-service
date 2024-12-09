package domain

import "auth-service/transport/transport_http/handlers"

type DomainRepository interface {
	Register(string, []byte) error
	GetUser(string) (handlers.User, error)
	ChangePassword(string, string, []byte) error
}
