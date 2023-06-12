package model

import (
	"context"

	"github.com/kardianos/service"
)

type Service interface {
	Start(ss service.Service) error // async
	StartSync(ctx context.Context) error
	Stop(ss service.Service) error
}
