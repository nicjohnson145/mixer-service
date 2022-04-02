package main

import (
	"encoding/json"
	"flag"
	"github.com/nicjohnson145/mixer-service/pkg/auth"
	"github.com/nicjohnson145/mixer-service/pkg/common"
	"github.com/nicjohnson145/mixer-service/pkg/drink"
	"github.com/nicjohnson145/mixer-service/pkg/jwt"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

type OldDrink struct {
	UUID           string   `json:"uuid"`
	Name           string   `json:"name"`
	PrimaryAlchol  string   `json:"primary_alcohol"`
	PreferredGlass string   `json:"preferred_glass"`
	Ingredients    []string `json:"ingredients"`
	Instructions   string   `json:"instructions"`
}

type Input struct {
	Drinks []OldDrink `json:"drinks"`
}

const (
	PasswordKey = "MIXER_PASSWORD"
)

var (
	url      string
	username string
	input    string
)

func init() {
	flag.StringVar(&url, "url", "http://localhost:30000", "URL of mixer service")
	flag.StringVar(&username, "username", "", "Username to authenticate to mixer service with")
	flag.StringVar(&input, "input", "", "Path to input file")
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func validateCommandLine() {
	if username == "" {
		log.Fatal("-username required")
	}
	if input == "" {
		log.Fatal("-input required")
	}
	if _, ok := os.LookupEnv(PasswordKey); !ok {
		log.Fatal("MIXER_PASSWORD must be set")
	}
}

func parseAndConvertOldJson() []drink.CreateDrinkRequest {
	bytes, err := ioutil.ReadFile(input)
	check(err)
	var inp Input
	err = json.Unmarshal(bytes, &inp)
	check(err)

	requests := make([]drink.CreateDrinkRequest, 0, len(inp.Drinks))
	for _, d := range inp.Drinks {
		requests = append(
			requests,
			drink.CreateDrinkRequest{
				DrinkData: drink.DrinkData{
					Name:           d.Name,
					PrimaryAlcohol: d.PrimaryAlchol,
					PreferredGlass: d.PreferredGlass,
					Ingredients:    d.Ingredients,
					Instructions:   d.Instructions,
					Publicity:      drink.DrinkPublicityPublic,
				},
			},
		)
	}

	return requests
}

func getAuthToken() string {
	body := auth.LoginRequest{
		Username: username,
		Password: os.Getenv(PasswordKey),
	}
	bodyBytes, err := json.Marshal(body)
	check(err)
	c := &http.Client{}
	req, err := http.NewRequest("POST", url+common.AuthV1+"/login", strings.NewReader(string(bodyBytes)))
	check(err)
	req.Header.Add("Content-type", "application/json")
	resp, err := c.Do(req)
	check(err)

	if resp.StatusCode != 200 {
		log.Fatal("got non 200 status code logging in: ", resp.StatusCode)
	}

	defer resp.Body.Close()
	var logResp auth.LoginResponse
	err = json.NewDecoder(resp.Body).Decode(&logResp)
	check(err)

	return logResp.AccessToken
}

func createDrink(c *http.Client, token string, r drink.CreateDrinkRequest) {
	bodyBytes, err := json.Marshal(r)
	check(err)
	req, err := http.NewRequest("POST", url+common.DrinksV1+"/create", strings.NewReader(string(bodyBytes)))
	check(err)
	req.Header.Add(jwt.AuthenticationHeader, token)
	req.Header.Add("Content-type", "application/json")
	resp, err := c.Do(req)
	check(err)

	if resp.StatusCode != 200 {
		log.Fatal("failed to create drink ", r.DrinkData.Name)
	}
}

func main() {
	flag.Parse()

	validateCommandLine()

	// Make the conversions
	requests := parseAndConvertOldJson()

	// Authenticate
	token := getAuthToken()

	// Make some requests
	c := &http.Client{}
	for _, req := range requests {
		createDrink(c, token, req)
	}
}
