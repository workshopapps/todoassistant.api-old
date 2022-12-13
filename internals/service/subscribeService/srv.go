package subscribeService

import (
	"context"
	"log"
	"test-va/internals/Repository/subscribeRepo"
	"test-va/internals/entity/ResponseEntity"
	"test-va/internals/entity/emailEntity"
	"test-va/internals/entity/subscribeEntity"
	"test-va/internals/service/emailService"
	"time"
)

type SubscribeService interface {
	PersistEmail(req *subscribeEntity.SubscribeReq) (*subscribeEntity.SubscribeRes, *ResponseEntity.ServiceError)
}

type subscribeSrv struct {
	repo     subscribeRepo.SubscribeRepository
	emailSrv emailService.EmailService
}

func NewSubscribeSrv(repo subscribeRepo.SubscribeRepository, emailSrv emailService.EmailService) SubscribeService {
	return &subscribeSrv{repo: repo, emailSrv: emailSrv}
}

func (t *subscribeSrv) PersistEmail(req *subscribeEntity.SubscribeReq) (*subscribeEntity.SubscribeRes, *ResponseEntity.ServiceError) {
	var message emailEntity.SendEmailReq

	// create context of 1 minute
	ctx, cancelFunc := context.WithTimeout(context.TODO(), time.Minute*1)
	defer cancelFunc()

	result, err1 := t.repo.CheckEmail(ctx, req)
	if result != nil {
		return nil, ResponseEntity.NewCustomServiceError("Already subscribed", err1)
	}

	message.EmailAddress = req.Email
	message.EmailSubject = "Subject: Subscription To Ticked Newsletter\n"
	message.EmailBody = CreateMessageBody()

	err := t.emailSrv.SendMail(message)
	if err != nil {
		return nil, ResponseEntity.NewInternalServiceError(err)
	}

	err = t.repo.PersistEmail(ctx, req)
	if err != nil {
		log.Println("From subcribe ", err)
		return nil, ResponseEntity.NewInternalServiceError(err)
	}
	data := subscribeEntity.SubscribeRes{
		Email: req.Email,
	}

	return &data, nil
}

// Auxillary function
func CreateMessageBody() string {
	subject := "Subscription to Ticked!\n\n"
	mainBody := "Thank you for subscribing to our newsletter!\n\nGet ready for an awesome ride"

	message := subject + mainBody
	return string(message)
}
