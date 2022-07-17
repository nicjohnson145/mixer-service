package main

import (
	"flag"
	"os"
	"github.com/carlmjohnson/requests"
	log "github.com/sirupsen/logrus"
	"github.com/nicjohnson145/mixer-service/pkg/auth"
	"github.com/nicjohnson145/mixer-service/pkg/common"
	"github.com/nicjohnson145/mixer-service/pkg/jwt"
	"context"
)

var (
	username string
	newUser string
	newPassword string
	host string
)

func init() {
	flag.StringVar(&username, "username", "", "User to login as")
	flag.StringVar(&newUser, "new-user", "", "User to register")
	flag.StringVar(&newPassword, "new-password", "mixer", "Password of new user")
	flag.StringVar(&host, "host", "https://mixer.nicjohnson.info", "DNS of mixer server")
}

func main() {
	flag.Parse()

	if username == "" {
		log.Fatal("Must give -username")
	}

	if newUser == "" {
		log.Fatal("Must give -new-user")
	}

	val, ok := os.LookupEnv("MIXER_PASS")
	if !ok {
		log.Fatal("MIXER_PASS must be set")
	}

	loginRequest := auth.LoginRequest{
		Username: username,
		Password: val,
	}

	var loginResp auth.LoginResponse
	err := requests.
		URL(host + common.AuthV1 + "/login").
		BodyJSON(&loginRequest).
		ToJSON(&loginResp).
		Fetch(context.Background())

	if err != nil {
		log.Fatalf("Error logging in: %v", err)
	}

	registerRequest := auth.RegisterNewUserRequest{
		Username: newUser,
		Password: newPassword,
	}

	err = requests.
			URL(host + common.AuthV1 + "/register-user").
			Header(jwt.AuthenticationHeader, loginResp.AccessToken).
			BodyJSON(&registerRequest).
			Fetch(context.Background())
	
	if err != nil {
		log.Fatalf("Error registering user: %v", err)
	}

	log.Infof("Successfully registered %v", newUser)
}
