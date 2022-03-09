package naming

import (
	"context"
)

type Resolver struct {
	id  string
	dis *Discovery
}

func (r *Resolver) Fetch(ctx context.Context) (map[string][]*Instance, bool) {

}
