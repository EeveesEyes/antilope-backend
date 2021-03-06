package util

import (
	"encoding/json"
	"github.com/EeveesEyes/owasp-password-strength-test"
	"math/rand"
)

func GetRandString(n int) string {
	p := make([]byte, n)
	rand.Read(p)
	return string(p)
}

func TestPasswordStrength(password string) (owasp.TestResult, error) {
	passwordConfig := owasp.DefaultPasswordConfig()
	jsonResult, _ := passwordConfig.TestPassword(password)
	var result owasp.TestResult
	err := json.Unmarshal(jsonResult, &result)
	return result, err
}
