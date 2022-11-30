package subscribeService

import (
	"test-va/internals/entity/ResponseEntity"
	"test-va/internals/entity/subscribeEntity"
)

type SubscribeService interface{
	PersistEmail(req *subscribeEntity.SubscribeReq)(*subscribeEntity.SubscribeRes, *ResponseEntity.ServiceError)
}

type subscribeSrv struct{

}
func NewSubscribeSrv() SubscribeService{
	return &subscribeSrv{}
}

func (t *subscribeSrv) PersistEmail(req *subscribeEntity.SubscribeReq) (*subscribeEntity.SubscribeRes, *ResponseEntity.ServiceError){
	return nil, nil
}
