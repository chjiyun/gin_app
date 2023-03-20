package myError

import "gin_app/app/common"

func New(text string) error {
	return &errString{s: text}
}

func NewET(errType common.ErrType) error {
	return &errString{s: errType.String()}
}

type errString struct {
	s string
}

func (e errString) Error() string {
	return e.s
}
