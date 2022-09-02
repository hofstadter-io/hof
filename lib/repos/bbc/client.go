package bbc

import (
	"os"

	"github.com/ktrysmt/go-bitbucket"
)

const TokenEnv = "BITBUCKET_TOKEN"
const UsernameEnv = "BITBUCKET_USERNAME"
const PasswordEnv = "BITBUCKET_PASSWORD"

func NewClient() (client *bitbucket.Client, err error) {

	// if token is empty, it means the same things as unauthenticated setup
	if token := os.Getenv(TokenEnv); token != "" {
		token := os.Getenv(TokenEnv)
		client = bitbucket.NewOAuthbearerToken(token)
	} else if password := os.Getenv(PasswordEnv); password != "" {
		username := os.Getenv(UsernameEnv)
		client = bitbucket.NewBasicAuth(username, password)
	} else {
		// if token is empty, it means the same things as unauthenticated setup
		client = bitbucket.NewOAuthbearerToken("")
	}
	return client, err
}
