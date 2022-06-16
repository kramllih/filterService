package bbolt

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/kramllih/filterService/config"
	"github.com/kramllih/filterService/internal/database"
	"github.com/kramllih/filterService/internal/logger"
	"go.etcd.io/bbolt"
)

const (
	Approvals string = "approvals"
	Rejected  string = "rejected"
	Messages  string = "messages"
)

type bolt struct {
	DB  *bbolt.DB
	log *logger.Logger
}

func init() {
	database.RegisterType("boltdb", NewDB)
}

func NewDB(cfg *config.ConfigNamespace) (database.Client, error) {

	boltConfig := boltCfg{
		Name: "database.db",
	}

	if err := cfg.Config().UnpackRaw(&boltConfig); err != nil {
		return nil, err
	}

	db, err := bbolt.Open(boltConfig.Name, 0600, nil)
	if err != nil {
		return nil, err
	}

	err = db.Update(func(tx *bbolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(Approvals))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}

		_, err = tx.CreateBucketIfNotExists([]byte(Rejected))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}

		_, err = tx.CreateBucketIfNotExists([]byte(Messages))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	client := &bolt{
		DB:  db,
		log: logger.NewLogger("database"),
	}

	return client, nil
}

func (b *bolt) StoreApproval(id string, approval []byte) error {

	err := b.DB.Update(func(tx *bbolt.Tx) error {
		bu := tx.Bucket([]byte(Approvals))
		err := bu.Put([]byte(id), approval)
		return err

	})
	if err != nil {
		b.log.Errorf("Unable to insert approval to database: %s", err)
		return err
	}

	return nil

}
func (b *bolt) GetApproval(id string) (*database.Approval, error) {

	var approval *database.Approval

	err := b.DB.View(func(tx *bbolt.Tx) error {
		bu := tx.Bucket([]byte(Approvals))
		if bu == nil {
			return errors.New("invalid bucket")
		}

		approvalBytes := bu.Get([]byte(id))

		if approvalBytes == nil {
			return nil
		}

		err := json.Unmarshal(approvalBytes, &approval)
		if err != nil {
			return fmt.Errorf("json unmarshal error: %w", err)
		}

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("unable to get approval from database: %s", err)
	}

	return approval, err
}

func (b *bolt) GetAllApprovals() ([]*database.Approval, error) {

	approvals := []*database.Approval{}

	err := b.DB.View(func(tx *bbolt.Tx) error {
		bu := tx.Bucket([]byte(Approvals))
		if bu == nil {
			return errors.New("invalid bucket")
		}

		cursor := bu.Cursor()

		for k, v := cursor.First(); k != nil; k, v = cursor.Next() {
			d := database.Approval{}

			err := json.Unmarshal(v, &d)
			if err != nil {
				return fmt.Errorf("json unmarshal error: %w", err)
			}

			approvals = append(approvals, &d)

		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return approvals, nil

}

func (b *bolt) UpdateApprovals(id string, approval []byte) error {

	err := b.DB.Update(func(tx *bbolt.Tx) error {
		bu := tx.Bucket([]byte(Approvals))
		err := bu.Put([]byte(id), approval)
		return err

	})
	if err != nil {
		b.log.Errorf("Unable to update approval in database: %s", err)
		return err
	}

	return nil
}

func (b *bolt) DeleteApprovals(id string) error {

	err := b.DB.Update(func(tx *bbolt.Tx) error {
		bu := tx.Bucket([]byte(Approvals))
		err := bu.Delete([]byte(id))
		return err

	})
	if err != nil {
		b.log.Errorf("Unable to update approval in database: %s", err)
		return err
	}

	return nil
}

func (b *bolt) StoreReject(id string, message []byte) error {

	err := b.DB.Update(func(tx *bbolt.Tx) error {
		bu := tx.Bucket([]byte(Rejected))
		err := bu.Put([]byte(id), message)
		return err

	})
	if err != nil {
		b.log.Errorf("Unable to insert reject to database: %s", err)
		return err
	}

	return nil
}

func (b *bolt) GetAllRejected() ([]*database.Message, error) {
	rejected := []*database.Message{}

	err := b.DB.View(func(tx *bbolt.Tx) error {
		bu := tx.Bucket([]byte(Rejected))
		if bu == nil {
			return errors.New("invalid bucket")
		}

		cursor := bu.Cursor()

		for k, v := cursor.First(); k != nil; k, v = cursor.Next() {
			d := database.Message{}

			err := json.Unmarshal(v, &d)
			if err != nil {
				return fmt.Errorf("json unmarshal error: %w", err)
			}

			rejected = append(rejected, &d)

		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return rejected, nil
}

func (b *bolt) StoreMessage(id string, message []byte) error {

	err := b.DB.Update(func(tx *bbolt.Tx) error {
		bu := tx.Bucket([]byte(Messages))
		err := bu.Put([]byte(id), message)
		return err

	})
	if err != nil {
		b.log.Errorf("Unable to insert message to database: %s", err)
		return err
	}

	return nil
}

func (b *bolt) GetMessage(id string) (*database.Message, error) {

	var message *database.Message

	err := b.DB.View(func(tx *bbolt.Tx) error {
		bu := tx.Bucket([]byte(Messages))
		if bu == nil {
			return errors.New("invalid bucket")
		}

		approvalBytes := bu.Get([]byte(id))

		if approvalBytes == nil {
			return nil
		}

		err := json.Unmarshal(approvalBytes, &message)
		if err != nil {
			return fmt.Errorf("json unmarshal error: %w", err)
		}

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("unable to get approval from database: %s", err)
	}

	return message, err

}
func (b *bolt) UpdateMessage(id string, message []byte) error {

	err := b.DB.Update(func(tx *bbolt.Tx) error {
		bu := tx.Bucket([]byte(Messages))
		err := bu.Put([]byte(id), message)
		return err

	})
	if err != nil {
		b.log.Errorf("Unable to update message in database: %s", err)
		return err
	}

	return nil
}

func (b *bolt) GetAllMessages() ([]*database.Message, error) {
	messages := []*database.Message{}

	err := b.DB.View(func(tx *bbolt.Tx) error {
		bu := tx.Bucket([]byte(Messages))
		if bu == nil {
			return errors.New("invalid bucket")
		}

		cursor := bu.Cursor()

		for k, v := cursor.First(); k != nil; k, v = cursor.Next() {
			d := database.Message{}

			err := json.Unmarshal(v, &d)
			if err != nil {
				return fmt.Errorf("json unmarshal error: %w", err)
			}

			messages = append(messages, &d)

		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return messages, nil
}
