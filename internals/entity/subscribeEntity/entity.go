package subscribeEntity

type SubscribeReq struct {
	Email    string `json:"email" validate:"email"`
}

type SubscribeRes struct {
	Email    string `json:"email" validate:"email"`
}
