package domain

import "errors"

var ErrInvalidTransition = errors.New("invalid transition")

func AdvanceStage(current string) (string, error) {
	switch current {
	case StageNewOpportunity:
		return StageNeedsConfirmed, nil
	case StageNeedsConfirmed:
		return StageQuote, nil
	case StageQuote:
		return StageContractNegotiation, nil
	default:
		return "", ErrInvalidTransition
	}
}

func ValidateStageTransition(fromStage, toStage string) error {
	next, err := AdvanceStage(fromStage)
	if err != nil {
		return err
	}
	if next != toStage {
		return ErrInvalidTransition
	}
	return nil
}
