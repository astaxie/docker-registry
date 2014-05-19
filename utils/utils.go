package utils

import (
	"encoding/base64"
	"fmt"
	"github.com/Unknwon/com"
	"github.com/Unknwon/goconfig"
	"log"
	"os"
	"strings"
)

var Cfg *goconfig.ConfigFile

// LoadConfig loads configuration file.
func LoadConfig(cfgPath string) {
	if !com.IsExist(cfgPath) {
		os.Create(cfgPath)
	}
	var err error
	Cfg, err = goconfig.LoadConfigFile(cfgPath)
	if err != nil {
		log.Fatalf("Fail to load configuration file: %v", err)
	}
}

// 下面两个方法用来处理用户名和密码加密:Authorization: Basic ZnNrOmZzaw==
// encode the auth string
func EncodeAuth(userName string, password string) string {
	authStr := userName + ":" + password
	msg := []byte(authStr)
	encoded := make([]byte, base64.StdEncoding.EncodedLen(len(msg)))
	base64.StdEncoding.Encode(encoded, msg)
	return string(encoded)
}

// decode the auth string
func DecodeAuth(authStr string) (string, string, error) {
	decLen := base64.StdEncoding.DecodedLen(len(authStr))
	decoded := make([]byte, decLen)
	authByte := []byte(authStr)
	n, err := base64.StdEncoding.Decode(decoded, authByte)
	if err != nil {
		return "", "", err
	}
	if n > decLen {
		return "", "", fmt.Errorf("Something went wrong decoding auth config")
	}
	arr := strings.SplitN(string(decoded), ":", 2)
	if len(arr) != 2 {
		return "", "", fmt.Errorf("Invalid auth configuration file")
	}
	password := strings.Trim(arr[1], "\x00")
	return arr[0], password, nil
}

func IsDirExists(path string) bool {
	fi, err := os.Stat(path)

	if err != nil {
		return os.IsExist(err)
	} else {
		return fi.IsDir()
	}

	panic("not reached")
}

//func IsFileExists(path string) bool {
//	if _, err := os.Stat(path); err == nil {
//		os.Remove(path)
//	}
//	panic("not reached")
//}

func Substr(str string, start, length int) string {
	rs := []rune(str)
	rl := len(rs)
	end := 0

	if start < 0 {
		start = rl - 1 + start
	}
	end = start + length

	if start > end {
		start, end = end, start
	}

	if start < 0 {
		start = 0
	}
	if start > rl {
		start = rl
	}
	if end < 0 {
		end = 0
	}
	if end > rl {
		end = rl
	}

	return string(rs[start:end])
}
