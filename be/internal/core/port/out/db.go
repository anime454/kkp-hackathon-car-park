package out

import "context"

type DBPinger interface {
	Ping(ctx context.Context) error
}
