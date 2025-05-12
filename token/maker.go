package token

import "time"

type Maker interface {
	//Create Token
	CreateToken(username string, duration time.Duration) (string, error)

	//Verify token
	VerifyToken(token string) (*Payload, error)
}
