package txmanager

import (
	"sync"
	"time"

	"github.com/google/uuid"
)

type Transaction struct {
	SentAt    time.Time `json:"sent_at"`
	UpdatedAt time.Time `json:"updated_at"`

	UUID        string `json:"uuid"`
	TxHash      string `json:"txhash"`
	Height      string `json:"height"`
	Status      string `json:"status"`
	RawLog      string `json:"raw_log"`
	Logs        string `json:"logs"`
	GasWanted   string `json:"gas_wanted"`
	GasUsed     string `json:"gas_used"`
	Timestamp   string `json:"timestamp"`
	RawResponse string `json:"raw_response"`

	Code int `json:"code"`
}

type TransactionExecutor interface {
	Execute(tp *TransactionPointer) error
	Fetch(tp *TransactionPointer) error
}

type TransactionPointer struct {
	Tx *Transaction
	mu sync.RWMutex
}

func NewTransaction() *Transaction {
	return &Transaction{
		UUID: uuid.New().String(),
		// Initialize other fields as necessary
	}
}

func (tp *TransactionPointer) BackgroundUpdate() {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		tp.updateTransaction()
	}
}

func ApplyUpdate(tp *TransactionPointer, update *Transaction) {
	tp.mu.Lock()
	defer tp.mu.Unlock()
	tp.Tx = update
}

func (tp *TransactionPointer) ClearTransaction(timeout time.Duration) {
	time.Sleep(timeout)
	tp.mu.Lock()
	defer tp.mu.Unlock()
	tp.Tx = nil
}

func (tp *TransactionPointer) Copy() *Transaction {
	tp.mu.RLock()
	defer tp.mu.RUnlock()
	return &Transaction{
		SentAt:      tp.Tx.SentAt,
		UpdatedAt:   tp.Tx.UpdatedAt,
		UUID:        tp.Tx.UUID,
		TxHash:      tp.Tx.TxHash,
		Height:      tp.Tx.Height,
		Status:      tp.Tx.Status,
		RawLog:      tp.Tx.RawLog,
		Logs:        tp.Tx.Logs,
		GasWanted:   tp.Tx.GasWanted,
		GasUsed:     tp.Tx.GasUsed,
		Timestamp:   tp.Tx.Timestamp,
		RawResponse: tp.Tx.RawResponse,
		Code:        tp.Tx.Code,
	}
}

func (tp *TransactionPointer) updateTransaction() {
	return
}
