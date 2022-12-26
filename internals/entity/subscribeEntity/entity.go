package subscribeEntity

type SubscribeReq struct {
	Email string `json:"email" validate:"email"`
}

type SubscribeRes struct {
	Email string `json:"email" validate:"email"`
}

type ContactUsReq struct {
	Name    string `json:"name"`
	Email   string `json:"email" validate:"email"`
	Message string `json:"message"`
}

type ContactUsRes struct {
	Email   string `json:"email" validate:"email"`
	Name    string `json:"name"`
	Message string `json:"message"`
}
