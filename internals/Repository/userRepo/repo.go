package userRepo

import "test-va/internals/entity/userEntity"

type UserRepository interface {
	Persist(req *userEntity.CreateUserReq) error
}
