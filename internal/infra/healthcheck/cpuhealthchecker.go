package healthcheck

import (
	"gurms/internal/infra/logging/core/factory"
	"gurms/internal/infra/logging/core/logger"
)

var CHCLOGGER logger.Logger = factory.GetLogger("CpuHealthManager")

type CpuHealthChecker struct {
}
