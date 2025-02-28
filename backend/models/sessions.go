package models

import "errors"

type Session struct {
	SessionId  string `json:"sessionId"`
	UserId     int    `json:"userId"`
	ExpireTime int64  `json:"expireTime"`
}

var ErrInvalidSession = errors.New("invalid session")
var ErrSessionExpired = errors.New("session expired")
var ErrInvalidUser = errors.New("invalid user")
