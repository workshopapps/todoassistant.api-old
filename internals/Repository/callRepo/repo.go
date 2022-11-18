package callRepo

import (
	"context"
	"test-va/internals/entity/callEntity"
)

type CallRepository interface{
	GetCalls(ctx context.Context) ([]*callEntity.CallRes, error)
}