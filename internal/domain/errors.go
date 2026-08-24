package domain

import "errors"

var (
	ErrNotFound   = errors.New("resource not found")
	ErrConflict   = errors.New("resource conflict")
	ErrInvalid    = errors.New("invalid request")
	ErrTransition = errors.New("invalid state transition")
)

func IsNotFound(err error) bool { return err == ErrNotFound }
func IsConflict(err error) bool { return err == ErrConflict }
func IsInvalid(err error) bool  { return err == ErrInvalid }
