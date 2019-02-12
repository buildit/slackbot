package config

import (
	"github.com/kelseyhightower/envconfig"
	"log"
)

// Config contains environment variables used to configure the app
type Config struct {
	VerificationToken string `envconfig:"SLACKBOT_VERIFICATIONTOKEN"`
	OauthToken        string `envconfig:"SLACKBOT_OAUTHTOKEN"`
}

var Env Config

func init() {
	err := envconfig.Process("slackbot", &Env)
	if err != nil {
		log.Fatal(err.Error())
	}
}
