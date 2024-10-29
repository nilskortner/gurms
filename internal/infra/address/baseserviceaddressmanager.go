package address

import (
	"gurms/internal/infra/logging/core/factory"
	"gurms/internal/infra/logging/core/logger"
)

var LOGGER logger.Logger = factory.GetLogger("BaseServiceAddressManager")

type BaseServiceAddressManager struct {
}
