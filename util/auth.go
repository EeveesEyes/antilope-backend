package util

import (
	"fmt"
	"github.com/EeveesEyes/antilope-backend/db"
	"github.com/EeveesEyes/antilope-backend/models"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"log"
	"time"
)

var Peppers = map[int]string{
	4:  "8D9Aec5e9B080Bde71cfc80355a65CdC",
	5:  "A6E2d5d4F2fA4a5a76F354C17f3AcDDe",
	6:  "aB76BfAfB5ced4Db4eB2C1006BbE0afF",
	7:  "9922b1B5AeDEDaDDfbCebf210B2F1FAe",
	8:  "f19E234c39cC2DFFe35498f3ecaAC1Df",
	9:  "dB2A65dFf55C8bb12A2623bafecbF6AE",
	10: "0325eCC27AD3dDf3bb947a11acbE657f",
} //make(map[int]string)

const hmacSampleSecret = "A6E2d5d4F2fA4a5a76F354C17f3AcDDeAeDEDaD"

//const ExpirationTime = time.Hour * 1
const ExpirationTime = time.Minute * 20

var LastClear time.Time

func flavourPassword(password string, pepperID int) string {
	return Peppers[pepperID] + password
}

func getLatestPepperID() int {
	keys := make([]int, len(Peppers))
	i := 0
	for k := range Peppers {
		keys[i] = k
		i++
	}
	if len(Peppers) == 0 {
		panic("Missing pepper!")
	}
	maxValue := keys[0]
	for _, v := range keys[1:] {
		if v > maxValue {
			maxValue = v
		}
	}
	return maxValue
}

func HashPassword(password string) (hash string, pepperID int, err error) {
	pepperID = getLatestPepperID()
	flavouredPW := flavourPassword(password, pepperID)
	hashBytes, _ := bcrypt.GenerateFromPassword([]byte(flavouredPW), 14)
	hash = string(hashBytes)
	return
}

func ValidatePassword(password string, user models.User) bool {
	flavouredPW := flavourPassword(password, user.PepperID)
	match := bcrypt.CompareHashAndPassword([]byte(user.Hash), []byte(flavouredPW))
	return match == nil
}

func GenerateJWT(user *models.User) (tokenString string, err error) {
	// Create a new token object, specifying signing method and the claims
	// you would like it to contain.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss": user.Email,
		"aud": "antilope",
		"exp": time.Now().Add(ExpirationTime).Unix(),
		"nbf": time.Now(),
	})
	// Sign and get the complete encoded token as a string using the secret
	tokenString, err = token.SignedString([]byte(hmacSampleSecret))
	return
}

// If userEmail is empty, Issuer will not be checked
func ValidateJWT(tokenString, userEmail string) bool {
	token, err := GetToken(tokenString)
	if err != nil {
		return false
	}

	if !CheckJWTinBlacklist(tokenString) {
		return false
	}
	if claims, ok := GetTokenClaims(token); ok {
		if claims.VerifyIssuer(userEmail, userEmail != "") &&
			claims.VerifyAudience("antilope", true) &&
			claims.VerifyExpiresAt(time.Now().Unix(), true) &&
			claims.VerifyNotBefore(time.Now().Unix(), true) {
			token.Valid = true
			return true
		}
		return false
	} else {
		return false
	}
}

func InvalidateToken(tokenString string) (err error) {
	token, err := GetToken(tokenString)
	if err != nil {
		return err
	}
	claims, _ := GetTokenClaims(token)
	expiration := claims["exp"].(float64)
	println(expiration)
	ex_tmp := int64(expiration)
	db.AddJWTToBlacklist(tokenString, ex_tmp)
	return nil
}

// Utils
func CheckJWTinBlacklist(tokenString string) (ok bool) {
	if LastClear.Add(ExpirationTime).Unix() < time.Now().Unix() {
		deletedCount := db.RemoveExpiredJWTsFromBlacklist()
		LastClear = time.Now()
		fmt.Printf("removed %v from JWT Blacklist\n", deletedCount)
	}
	return db.CheckJWTinBlacklist(tokenString)
}

func GetToken(tokenString string) (token *jwt.Token, err error) {
	token, err = jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(hmacSampleSecret), nil
	})
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	return token, nil
}

func GetTokenClaims(token *jwt.Token) (claims jwt.MapClaims, ok bool) {
	claims, ok = token.Claims.(jwt.MapClaims)
	return claims, ok
}
