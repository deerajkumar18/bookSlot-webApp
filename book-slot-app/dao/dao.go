package dao

import (
	"context"
)

type EventsInfoDao interface {
	GetSlotsAvailabilityByEventID(ctx context.Context) (SlotsAvailable int, err error)
	UpdateSlotsAvailability(requestType string, ctx context.Context) (err error)
}

type SlotsInfoDao interface {
	InsertSlotsInfo(ctx context.Context) (err error)
	DeleteSlotsInfo(ctx context.Context) (err error)
}
