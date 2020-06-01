package factoryinit

import (
	"time"

	random "github.com/brianvoe/gofakeit/v5"
)

func init() {
	random.Seed(time.Now().UnixNano())
}
