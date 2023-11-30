package pins

import (
	"fmt"
	"math/rand"
)

func GetRandomPin() string {
	randomNumber := rand.Intn(900000) + 100000
	return fmt.Sprintf("%d", randomNumber)
}
