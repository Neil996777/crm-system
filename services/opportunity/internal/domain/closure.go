package domain

import (
	"errors"
	"strings"
	"time"
)

var (
	ErrEarlyWonBlocked        = errors.New("early won blocked")
	ErrLostReasonRequired     = errors.New("lost reason required")
	ErrTerminalRecordReadOnly = errors.New("terminal record read only")
)

func CloseWon(current Opportunity, contractSigned bool, contractID string) (Opportunity, error) {
	if err := EnsureOpenForMutation(current); err != nil {
		return Opportunity{}, err
	}
	contractID = strings.TrimSpace(contractID)
	if current.Stage != StageContractNegotiation || !contractSigned || contractID == "" {
		return Opportunity{}, ErrEarlyWonBlocked
	}
	closed := current
	closed.Stage = StageWon
	closed.WonContractID = contractID
	closed.LostReasonCode = ""
	closed.LostReasonDetail = ""
	return closed, nil
}

func CloseLost(current Opportunity, reasonCode, reasonDetail string) (Opportunity, error) {
	if err := EnsureOpenForMutation(current); err != nil {
		return Opportunity{}, err
	}
	reasonCode = strings.TrimSpace(reasonCode)
	reasonDetail = strings.TrimSpace(reasonDetail)
	if reasonCode == "" {
		return Opportunity{}, ErrLostReasonRequired
	}
	closed := current
	closed.Stage = StageLost
	closed.WonContractID = ""
	closed.LostReasonCode = reasonCode
	closed.LostReasonDetail = reasonDetail
	return closed, nil
}

func ApplyClosureDates(opportunity Opportunity, closeDate time.Time, closedAt time.Time) Opportunity {
	opportunity.CloseDate = closeDate
	opportunity.ClosedAt = closedAt
	return opportunity
}

func EnsureOpenForMutation(opportunity Opportunity) error {
	if IsTerminalStage(opportunity.Stage) {
		return ErrTerminalRecordReadOnly
	}
	return nil
}

func IsTerminalStage(stage string) bool {
	return stage == StageWon || stage == StageLost
}
