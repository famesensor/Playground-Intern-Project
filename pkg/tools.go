package pkg

import (
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"time"
)

func RandomNumber(max, min int) string {
	rand.Seed(time.Now().UnixNano())
	randomNumber := rand.Intn(max-min+1) + min
	return strconv.Itoa(randomNumber)
}

func RandUUID(max, min, n int) string {
	b := make([]byte, n)
	randN := int64(rand.Intn(max-min+1) + min)
	rand.Seed(time.Now().UnixNano() + randN)
	_, err := rand.Read(b)
	if err != nil {
		log.Fatal(err)
	}
	return fmt.Sprintf("%x-%x-%x-%x-%x-%x-%x-%x-%x-%x",
		b[0:4], b[4:6], b[6:12], b[12:14], b[14:17], b[17:20], b[20:24], b[24:28], b[28:30], b[30:32])
}
