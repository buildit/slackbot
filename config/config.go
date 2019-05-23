package config

import (
	"github.com/kelseyhightower/envconfig"
	log "github.com/sirupsen/logrus"
)

// Config contains environment variables used to configure the app
type Config struct {
	ListeningPort     string `default:"4390"`
	VerificationToken string
	OauthToken        string
	PollBucket        string `default:"POLL"`
}

type AppInsightsConfig struct {
	InstrumentationKey string
}

var Env Config
var AppInsights AppInsightsConfig

func init() {
	err := envconfig.Process("appsetting_slackbot", &Env)
	if Env.OauthToken == "" || Env.VerificationToken == "" {
		log.Error("Missing Environment.  APPSETTING_SLACKBOT_OAUTHTOKEN and APPSETTING_SLACKBOT_VERIFICATIONTOKEN are Required.  Exiting...")
		log.Fatal(err.Error())
	}
	if err != nil {
		log.Fatal(err.Error())
	}
}

func init() {
	err := envconfig.Process("appsetting_appinsights", &AppInsights)
	if err != nil {
		log.Fatal(err.Error())
	}
}
