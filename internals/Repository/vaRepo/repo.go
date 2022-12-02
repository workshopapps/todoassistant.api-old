package vaRepo

import (
	"context"
	"test-va/internals/entity/vaEntity"
)

type VARepo interface {
	Persist(ctx context.Context, req *vaEntity.CreateVAReq) error
	FindByEmail(ctx context.Context, email string) (*vaEntity.FindByEmailRes, error)
	FindById(ctx context.Context, id string) (*vaEntity.FindByIdRes, error)
	DeleteUser(ctx context.Context, id string) error
	UpdateUser(ctx context.Context, req *vaEntity.EditVaReq, id string) error
	ChangePassword(ctx context.Context, password *vaEntity.ChangeVAPassword) error
	GetUserAssignedToVa(ctx context.Context, userId string) ([]*vaEntity.VAStruct, error)
}
