package common

import (
	"gurms/internal/infra/property/constants"
)

type AddressProperties struct {
	advertiseStrategy int
	advertiseHost     string
	attachPortToHost  bool
}

func NewAddressProperties() *AddressProperties {
	return &AddressProperties{
		advertiseStrategy: constants.PRIVATE_ADDRESS,
		advertiseHost:     "",
		attachPortToHost:  true,
	}
}
