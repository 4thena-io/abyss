package database

import (
	"errors"
	"fmt"
	"net"
	"testing"
)

type stateError string

func (e stateError) Error() string    { return string(e) }
func (e stateError) SQLState() string { return string(e) }
func TestIsTransient(t *testing.T) {
	for _, tc := range []struct {
		err  error
		want bool
	}{
		{&net.OpError{Op: "dial", Err: errors.New("refused")}, true},
		{fmt.Errorf("connect: %w", stateError("57P03")), true},
		{stateError("28P01"), false},
		{stateError("3D000"), false},
		{&MigrationError{Err: stateError("57P03")}, false},
		{errors.New("invalid configuration"), false},
	} {
		if got := IsTransient(tc.err); got != tc.want {
			t.Errorf("%v: got %v", tc.err, got)
		}
	}
}
