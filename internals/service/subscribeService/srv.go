package subscribeService

import (
	"context"
	"log"
	"test-va/internals/Repository/subscribeRepo"
	"test-va/internals/entity/ResponseEntity"
	"test-va/internals/entity/subscribeEntity"
	"time"
)

type SubscribeService interface{
	PersistEmail(req *subscribeEntity.SubscribeReq)(*subscribeEntity.SubscribeRes, *ResponseEntity.ServiceError)
}

type subscribeSrv struct{
	repo 	subscribeRepo.SubscribeRepository
}
func NewSubscribeSrv(repo subscribeRepo.SubscribeRepository) SubscribeService{
	return &subscribeSrv{repo: repo}
}

func (t *subscribeSrv) PersistEmail(req *subscribeEntity.SubscribeReq) (*subscribeEntity.SubscribeRes, *ResponseEntity.ServiceError){
	// create context of 1 minute
	ctx, cancelFunc := context.WithTimeout(context.TODO(), time.Minute*1)
	defer cancelFunc()

	err := t.repo.PersistEmail(ctx, req)
	if err != nil {
		log.Println("From subcribe ",err)
		return nil, ResponseEntity.NewInternalServiceError(err)
	}
	data:= subscribeEntity.SubscribeRes{
		Email: req.Email,
	}

	return &data, nil
}
