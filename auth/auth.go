package auth

import (
	"fmt"
	"strings"

	"github.com/dockboard/docker-registry/models"
	"github.com/dockboard/docker-registry/utils"
)

type AuthError string

func (e AuthError) Error() string {
	return string(e)
}

func BaseAuth(authBasic string) (authUsername string, authPasswd string, err error) {
	var authorizationBasic string
	authBasic = strings.TrimSpace(authBasic)
	authorizationBasic = utils.Substr(authBasic, 6, len(authBasic))
	fmt.Println("authorizationBasic", authorizationBasic, ";")
	authUsername, authPasswd, authErr := utils.DecodeAuth(authorizationBasic)
	fmt.Println("authUsername", authUsername, ";", authPasswd)
	if authErr != nil {
		return "", "", AuthError("Auth Error")
	}

	user := &models.User{Username: authUsername, Password: authPasswd}
	_, err = models.Engine.Get(user)
	if err != nil {
		return "", "", AuthError("Auth Error")
	}

	return authUsername, authPasswd, nil
}
