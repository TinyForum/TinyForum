package mask

import "errors"

var (
	// ErrNilPointer 传入 nil 指针
	ErrNilPointer = errors.New("mask: input is nil")
	// ErrNonPointer 传入非指针类型
	ErrNonPointer = errors.New("mask: input must be a pointer to struct")
)
