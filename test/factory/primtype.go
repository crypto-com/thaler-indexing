package factory

import (
	"time"

	random "github.com/brianvoe/gofakeit/v5"
)

func RandomInt64Ptr() *int64 {
	value := random.Int64()
	return &value
}

func RandomPveInt64Ptr() *int64 {
	value := int64(random.Uint64())
	return &value
}

func RandomUint64Ptr() *uint64 {
	value := random.Uint64()
	return &value
}

func RandomUint32Ptr() *uint32 {
	value := random.Uint32()
	return &value
}

func RandomTimePtr() *time.Time {
	value := random.DateRange(time.Unix(0, 0).UTC(), time.Now().UTC())
	return &value
}

func RandomEmailPtr() *string {
	value := random.Email()
	return &value
}
