package application

import "context"

type Application interface {
	Start()
	Stop(ctx context.Context) error
	Reload(ctx context.Context)
}
