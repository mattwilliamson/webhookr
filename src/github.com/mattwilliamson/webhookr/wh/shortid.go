package wh

import (
	"crypto/rand"
	"time"
)

const ID_CHARS = "23456789abcdefghjklmnpqrstuvwxyzABCDEFGHJKLMNPQRSTUVWXYZ"
const ID_MIN_LENGTH = 4

func RandomId() string {
	bytes := make([]byte, ID_MIN_LENGTH + (time.Now().UnixNano() % ID_MIN_LENGTH))
	rand.Read(bytes)

    for i, b := range bytes {
        bytes[i] = ID_CHARS[b % byte(len(ID_CHARS))]
    }

    return string(bytes)
}