package database

import (
	"fmt"

	"github.com/kramllih/filterService/config"
)

type Client interface {
	StoreApproval(string, []byte) error
	GetApproval(string) (*Approval, error)
	GetAllApprovals() ([]*Approval, error)
	UpdateApprovals(string, []byte) error
	DeleteApprovals(string) error

	StoreReject(string, []byte) error
	GetAllRejected() ([]*Message, error)

	StoreMessage(string, []byte) error
	GetMessage(string) (*Message, error)
	UpdateMessage(string, []byte) error
	GetAllMessages() ([]*Message, error)
}

type Factory func(config *config.ConfigNamespace) (Client, error)

var cache = map[string]Factory{}

func RegisterType(name string, f Factory) {
	if cache[name] != nil {
		panic(fmt.Errorf("queue type  '%v' exists already", name))
	}
	cache[name] = f
}

func FindFactory(name string) Factory {
	return cache[name]
}

func Load(config *config.ConfigNamespace) (Client, error) {

	factory := FindFactory(config.Name())
	if factory == nil {
		return nil, fmt.Errorf("database type %v undefined", config.Name())
	}

	return factory(config)
}
