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

//Encode the authorization string
func EncodeBasicAuth(username string, password string) string {
	auth := username + ":" + password
	msg := []byte(auth)
	authorization := make([]byte, base64.StdEncoding.EncodedLen(len(msg)))
	base64.StdEncoding.Encode(authorization, msg)
	return string(authorization)
}

// decode the authorization string
func DecodeBasicAuth(authorization string) (username string, password string, err error) {
	basic := strings.Split(strings.TrimSpace(authorization), " ")
	if len(basic) <= 1 {
		return "", "", err
	}

	decLen := base64.StdEncoding.DecodedLen(len(basic[1]))
	decoded := make([]byte, decLen)
	authByte := []byte(basic[1])
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

	username = arr[0]
	password = strings.Trim(arr[1], "\x00")

	return username, password, nil
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

func RemoveDuplicateString(s *[]string) {
	found := make(map[string]bool)
	j := 0

	for i, val := range *s {
		if _, ok := found[val]; !ok {
			found[val] = true
			(*s)[j] = (*s)[i]
			j++
		}
	}

	*s = (*s)[:j]
}
