package pins

import (
	"fmt"
	"math/rand"
	"time"
)

func GetRandomPin() string {
	rand.Seed(time.Now().UnixNano())
	randomNumber := rand.Intn(900000) + 100000
	return fmt.Sprintf("%d", randomNumber)
}
