package jwxtClient

import (
	"server/backend/jwxtClient/request"
)

type jwxtClient struct {
	*request.HttpClient
	username string
	isLogin  bool
}

func NewClient() *jwxtClient {
	return &jwxtClient{
		HttpClient: request.NewClient(),
		username:   "NULL",
		isLogin:    false,
	}
}
