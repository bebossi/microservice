package event

import "time"

type TransactionCreated struct {
	Name      string
	Payload   interface{}
	CreatedAt time.Time
}

func NewTransactionCreated() *TransactionCreated {
	return &TransactionCreated{
		Name: "TransactionCreated",
	}
}

func (t *TransactionCreated) GetName() string {
	return t.Name
}

func (t *TransactionCreated) GetPayload() interface{} {
	return t.Payload
}


func (t *TransactionCreated) SetPayload(payload interface{}) {
	t.Payload = payload
}

func (t *TransactionCreated) GetDateTime() time.Time {
	return time.Now()
}
