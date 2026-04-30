package security

import (
	"crypto/rand"
	"math/big"
	"strconv"
	"strings"
)

func GenerateOTP() string {
	var otp strings.Builder
	for range 6 {
		num, _ := rand.Int(rand.Reader, big.NewInt(10))
		otp.WriteString(strconv.Itoa(int(num.Int64())))
	}
	return otp.String()
}
