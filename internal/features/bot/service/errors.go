package service

import "errors"

var (
	ErrBotAlreadyExists = errors.New("bot already exist")
	ErrInvalidBot       = errors.New("invalid bot")
)
