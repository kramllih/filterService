package mockdb

import (
	"encoding/json"
	"errors"
	"fmt"
	"sync"

	"github.com/kramllih/filterService/config"
	"github.com/kramllih/filterService/internal/database"
	"github.com/kramllih/filterService/internal/logger"
)

type mockClient struct {
	log *logger.Logger

	mu        sync.RWMutex
	Approvals map[string][]byte
	Rejected  map[string][]byte
	Messages  map[string][]byte
}

func init() {
	database.RegisterType("mockDB", NewDB)
}

func NewDB(cfg *config.ConfigNamespace) (database.Client, error) {

	return &mockClient{
		Approvals: make(map[string][]byte),
		Rejected:  make(map[string][]byte),
		Messages:  make(map[string][]byte),
		log:       logger.NewLogger("mockDB"),
	}, nil
}

func (m *mockClient) StoreApproval(id string, approval []byte) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.Approvals[id]; ok {
		return errors.New("id already exists in database")
	}

	m.Approvals[id] = approval

	return nil

}
func (m *mockClient) GetApproval(id string) (*database.Approval, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if app, ok := m.Approvals[id]; ok {

		approval := database.Approval{}

		if err := json.Unmarshal(app, &approval); err != nil {
			return nil, fmt.Errorf("error decoding data: %w", err)
		}
		return &approval, nil
	}

	return nil, errors.New("id does not exist in database")

}
func (m *mockClient) GetAllApprovals() ([]*database.Approval, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	approvals := []*database.Approval{}

	for _, v := range m.Approvals {
		approval := database.Approval{}

		if err := json.Unmarshal(v, &approval); err != nil {
			return nil, fmt.Errorf("error decoding data: %w", err)
		}

		approvals = append(approvals, &approval)

	}

	return approvals, nil

}
func (m *mockClient) UpdateApprovals(id string, approval []byte) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.Approvals[id]; ok {
		m.Approvals[id] = approval
		return nil
	}

	return errors.New("id does not exist in database")
}
func (m *mockClient) DeleteApprovals(id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.Approvals[id]; ok {
		delete(m.Approvals, id)
		return nil
	}

	return errors.New("id does not exist in database")

}

func (m *mockClient) StoreReject(id string, reject []byte) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.Rejected[id]; ok {
		return errors.New("id already exists in database")
	}

	m.Approvals[id] = reject

	return nil

}
func (m *mockClient) GetAllRejected() ([]*database.Message, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	messages := []*database.Message{}

	for _, v := range m.Rejected {
		message := database.Message{}

		if err := json.Unmarshal(v, &message); err != nil {
			return nil, fmt.Errorf("error decoding data: %w", err)
		}

		messages = append(messages, &message)

	}

	return messages, nil
}

func (m *mockClient) StoreMessage(id string, message []byte) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.Messages[id]; ok {
		return errors.New("id already exists in database")
	}

	m.Messages[id] = message

	return nil

}
func (m *mockClient) GetMessage(id string) (*database.Message, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if app, ok := m.Messages[id]; ok {

		message := database.Message{}

		if err := json.Unmarshal(app, &message); err != nil {
			return nil, fmt.Errorf("error decoding data: %w", err)
		}
		return &message, nil
	}

	return nil, errors.New("id does not exist in database")
}
func (m *mockClient) UpdateMessage(id string, message []byte) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.Messages[id]; ok {
		m.Messages[id] = message
		return nil
	}

	return errors.New("id does not exist in database")

}
func (m *mockClient) GetAllMessages() ([]*database.Message, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	messages := []*database.Message{}

	for _, v := range m.Messages {
		message := database.Message{}

		if err := json.Unmarshal(v, &message); err != nil {
			return nil, fmt.Errorf("error decoding data: %w", err)
		}

		messages = append(messages, &message)

	}

	return messages, nil
}
