package xstring

import (
	"math/rand"
	"strconv"
	"time"
)

// 卡号加 *
func CardNumberToStar(cardNumber string) string {
	l := len(cardNumber)
	if l <= 6 {
		return cardNumber
	}
	return cardNumber[:4] + star[l] + cardNumber[len(cardNumber)-4:]
}

// 生成有效的卡号
func GenerateCardNo() string {
	rand.New(rand.NewSource(time.Now().UnixNano()))
	prefix := rand.Intn(5) + 51

	// 剩下的 9 位随机生成（总共16位，去掉前2位prefix和1位校验位）
	cardNumber := strconv.Itoa(prefix)
	for i := 0; i < 14; i++ {
		cardNumber += strconv.Itoa(rand.Intn(10))
	}

	// 生成的前15位+校验位需要通过Luhn校验
	for i := 0; i <= 9; i++ {
		completeCard := cardNumber + strconv.Itoa(i)
		if luhnCheck(completeCard) {
			return completeCard
		}
	}

	return ""
}

// Luhn算法：用于验证卡号是否有效
func luhnCheck(num string) bool {
	sum := 0
	alt := false
	for i := len(num) - 1; i >= 0; i-- {
		n, _ := strconv.Atoi(string(num[i]))
		if alt {
			n *= 2
			if n > 9 {
				n -= 9
			}
		}
		sum += n
		alt = !alt
	}
	return sum%10 == 0
}
