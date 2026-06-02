package domain

import (
	"errors"
	"math"
	"strconv"
	"strings"
	"time"
)

const (
	PaymentStatusUnpaid        = "Unpaid"
	PaymentStatusPartiallyPaid = "PartiallyPaid"
	PaymentStatusPaid          = "Paid"
	PaymentCurrencyCNY         = "CNY"
)

var (
	ErrInvalidAmount          = errors.New("invalid amount")
	ErrOverpaymentBlocked     = errors.New("overpayment blocked")
	ErrSingleCurrencyRequired = errors.New("single currency required")
)

type PaymentPlan struct {
	ID            string
	ContractID    string
	DueAmount     string
	DueDate       time.Time
	Currency      string
	Status        string
	ArchivedAt    time.Time
	ArchivedBy    string
	ArchiveReason string
	Version       int
	UpdatedAt     time.Time
}

type ActualPayment struct {
	ID              string
	ContractID      string
	IdempotencyKey  string
	Amount          string
	PaymentDate     time.Time
	Note            string
	Currency        string
	PaymentStatus   string
	RemainingAmount string
	Version         int
	UpdatedAt       time.Time
}

func NewPaymentPlan(contractID, dueAmount, dueDate, currency string) (PaymentPlan, error) {
	contractID = strings.TrimSpace(contractID)
	amount := normalizeAmount(dueAmount)
	if contractID == "" || amount == "" || dueDate == "" {
		return PaymentPlan{}, ErrInvalidAmount
	}
	if err := ensureSingleCurrency(currency); err != nil {
		return PaymentPlan{}, err
	}
	date, err := ParseDate(dueDate)
	if err != nil {
		return PaymentPlan{}, ErrValidation
	}
	return PaymentPlan{ContractID: contractID, DueAmount: amount, DueDate: date, Currency: PaymentCurrencyCNY, Status: PaymentStatusUnpaid, Version: 1}, nil
}

func NewActualPayment(contract Contract, idempotencyKey, amount, paymentDate, note, currency string, paidBefore string) (ActualPayment, error) {
	idempotencyKey = strings.TrimSpace(idempotencyKey)
	normalizedAmount := normalizeAmount(amount)
	if idempotencyKey == "" || normalizedAmount == "" || strings.TrimSpace(paymentDate) == "" {
		return ActualPayment{}, ErrInvalidAmount
	}
	if err := ensureSingleCurrency(currency); err != nil {
		return ActualPayment{}, err
	}
	date, err := ParseDate(paymentDate)
	if err != nil {
		return ActualPayment{}, ErrValidation
	}
	contractCents := amountCents(contract.Amount)
	paidBeforeCents := amountCents(paidBefore)
	paymentCents := amountCents(normalizedAmount)
	if paymentCents <= 0 {
		return ActualPayment{}, ErrInvalidAmount
	}
	paidAfter := paidBeforeCents + paymentCents
	if paidAfter > contractCents {
		return ActualPayment{}, ErrOverpaymentBlocked
	}
	remaining := contractCents - paidAfter
	status := PaymentStatusPartiallyPaid
	if remaining == 0 {
		status = PaymentStatusPaid
	}
	return ActualPayment{
		ContractID:      contract.ID,
		IdempotencyKey:  idempotencyKey,
		Amount:          normalizedAmount,
		PaymentDate:     date,
		Note:            strings.TrimSpace(note),
		Currency:        PaymentCurrencyCNY,
		PaymentStatus:   status,
		RemainingAmount: centsAmount(remaining),
		Version:         1,
	}, nil
}

func ensureSingleCurrency(currency string) error {
	currency = strings.TrimSpace(currency)
	if currency == "" || currency == PaymentCurrencyCNY {
		return nil
	}
	return ErrSingleCurrencyRequired
}

func amountCents(value string) int64 {
	amount, _ := strconv.ParseFloat(strings.TrimSpace(value), 64)
	return int64(math.Round(amount * 100))
}

func centsAmount(value int64) string {
	return strconv.FormatFloat(float64(value)/100, 'f', 2, 64)
}
