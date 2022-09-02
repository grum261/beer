package models

import "errors"

var ErrInvalidRequestStatus = errors.New("invalid friend request status")

type RequestStatus int8

const (
	StatusDeclined RequestStatus = iota - 1
	StatusSent
	StatusAccepted
)

func (r RequestStatus) Validate() error {
	switch r {
	case StatusDeclined, StatusAccepted, StatusSent:
		return nil
	default:
		return ErrInvalidRequestStatus
	}
}
