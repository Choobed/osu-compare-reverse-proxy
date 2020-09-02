package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

const osuAPIURL = "https://osu.ppy.sh/api/v2"

type authRes struct {
	TokenType   string  `json:"token_type"`
	ExpiresIn   float64 `json:"expires_in"`
	AccessToken string  `json:"access_token"`
}

func main() {
	http.HandleFunc("/", getUser)

	if err := http.ListenAndServe(":3001", nil); err != nil {
		log.Fatal(err)
	}
}

func getAuthRes() authRes {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal(err)
	}

	osuAPIClientID := os.Getenv("OSU_API_CLIENT_ID")
	osuAPIClientSecret := os.Getenv("OSU_API_CLIENT_SECRET")

	body := strings.NewReader(`
	{
			"client_id": "` + osuAPIClientID + `",
			"client_secret": "` + osuAPIClientSecret + `",
			"grant_type": "client_credentials",
			"scope": "public"
		}
	`)

	res, err := http.Post("https://osu.ppy.sh/oauth/token", "application/json", body)

	if err != nil {
		log.Fatal(err)
	}

	data, err := ioutil.ReadAll(res.Body)

	if err != nil {
		log.Fatal(err)
	}

	res.Body.Close()

	var authRes authRes
	json.Unmarshal(data, &authRes)

	return authRes
}

func getUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	username := r.URL.Query().Get("username")
	userLookupURL := fmt.Sprintf("%s/users/%s/osu", osuAPIURL, username)
	req, err := http.NewRequest("GET", userLookupURL, nil)

	if err != nil {
		log.Fatal(err)
	}

	authRes := getAuthRes()
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", authRes.AccessToken))

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		log.Fatal(err)
	}

	data, err := ioutil.ReadAll(res.Body)

	if err != nil {
		log.Fatal(err)
	}

	res.Body.Close()

	w.Write(data)
}
