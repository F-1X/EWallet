package validator

import (
	"ewallet/pkg/model"
	"regexp"
)

var (
	isValidWalletId = regexp.MustCompile(`^[a-zA-Z0-9]+$`).MatchString
)

func ValidateWallet(value string) error {
	if len(value) != 32 {
		return model.ErrIncorrectWallet.Error()
	}

	if !isValidWalletId(value) {
		return model.ErrIncorrectWallet.Error()
	}

	return nil
}
