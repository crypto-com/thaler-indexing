package factory

import (
	"time"

	random "github.com/brianvoe/gofakeit/v5"
)

func RandomUTCTime() time.Time {
	return random.DateRange(time.Unix(0, 0).UTC(), time.Now().UTC())
}
