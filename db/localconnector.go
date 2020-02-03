package db

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

const LocalFilePath = "./tmp/dat.json"

type LocalData struct {
	Peppers    map[int]string `json:"peppers"`
	HmacSecret string         `json:"hmacSecret"`
	RedisPW    string
}

var localData LocalData

func createHash(key string) string {
	hasher := md5.New()
	hasher.Write([]byte(key))
	return hex.EncodeToString(hasher.Sum(nil))
}

func encrypt(data []byte, passphrase string) []byte {
	block, _ := aes.NewCipher([]byte(createHash(passphrase)))
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		panic(err.Error())
	}
	ciphertext := gcm.Seal(nonce, nonce, data, nil)
	return ciphertext
}

func decrypt(data []byte, passphrase string) []byte {
	key := []byte(createHash(passphrase))
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err.Error())
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}
	nonceSize := gcm.NonceSize()
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		panic(err.Error())
	}
	return plaintext
}

func encryptFile(filename string, data []byte, passphrase string) {
	f, _ := os.Create(filename)
	defer f.Close()
	f.Write(encrypt(data, passphrase))
}

func decryptFile(filename string, passphrase string) []byte {
	data, _ := ioutil.ReadFile(filename)
	return decrypt(data, passphrase)
}

func DecryptLocalSecrets(password string) {
	fmt.Println("Starting the application...")
	// check if file exists
	var localDataFileExists bool = fileExists(LocalFilePath)
	if !localDataFileExists {
		var pepper [32]byte
		if _, err := io.ReadFull(rand.Reader, pepper[:]); err != nil {
			panic(err)
		}
		var hmacSec [32]byte
		if _, err := io.ReadFull(rand.Reader, hmacSec[:]); err != nil {
			panic(err)
		}
		var redisPW [32]byte
		if _, err := io.ReadFull(rand.Reader, hmacSec[:]); err != nil {
			panic(err)
		}
		localData = LocalData{
			Peppers:    map[int]string{1: base64.URLEncoding.EncodeToString(pepper[:])},
			HmacSecret: base64.URLEncoding.EncodeToString(hmacSec[:]),
			RedisPW:    base64.URLEncoding.EncodeToString(redisPW[:]),
		}
		fmt.Println("New local data created!")
	} else {
		plaintext := decryptFile(LocalFilePath, password)
		json.Unmarshal(plaintext, &localData)
	}
}

func EncryptLocalSecrets(password string) {
	file, _ := json.MarshalIndent(localData, "", " ")
	encryptFile(LocalFilePath, file, password)
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func GetPeppers() map[int]string {
	return localData.Peppers
}
func GetHmacSecret() string {
	return localData.HmacSecret
}
func GetRedisPW() string {
	return localData.RedisPW
}
