package service

import "time"

type Queue interface {
	IsLocked() (bool, time.Duration)
}
