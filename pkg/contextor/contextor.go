package contextor

import "context"

// New takes a context and checks if its nil, if the provided
// context is not nil it will return otherwise it will create
// a new context with context.Background() and return.
func New(ctx context.Context) context.Context {
	if ctx == nil {
		return context.Background()
	}

	return ctx
}
