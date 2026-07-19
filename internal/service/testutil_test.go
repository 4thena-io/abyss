package service

import "errors"

var errBoom = errors.New("boom")

func ptr[T any](v T) *T { return &v }
