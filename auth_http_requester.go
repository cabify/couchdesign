package couchdesign

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/apex/log"
)

type AuthHttpRequester struct {
	baseUrl, username, password string
	httpClient                  *http.Client
}

func NewAuthHttpRequester(user, pass, url string) (*AuthHttpRequester, error) {
	return &AuthHttpRequester{
		username: user,
		password: pass,
		baseUrl:  fmt.Sprintf("http://%s:5984/", url),
		httpClient: &http.Client{
			Timeout: time.Second * 10,
		},
	}, nil
}

func (a *AuthHttpRequester) Get(path string, dest interface{}) error {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s%s", a.baseUrl, path), nil)
	if err != nil {
		return err
	}
	return a.run(req, dest)
}

func (a *AuthHttpRequester) run(req *http.Request, dest interface{}) error {
	if len(a.username) > 0 && len(a.password) > 0 {
		req.SetBasicAuth(a.username, a.password)
	}

	log.WithFields(log.Fields{"URL": req.URL, "method": req.Method}).Debug("Sending request...")

	resp, err := a.httpClient.Do(req)
	if err != nil {
		log.WithError(err).Error("Failed HTTP request!")
		return err
	}

	if resp.StatusCode/100 != 2 {
		return fmt.Errorf("Received response %d for %s", resp.StatusCode, req.URL.String())
	}

	defer resp.Body.Close()

	if dest != nil {
		if err = json.NewDecoder(resp.Body).Decode(dest); err != nil {
			log.WithError(err).Error("Error decoding response!")
			return err
		}
	}

	return nil
}

func (a *AuthHttpRequester) BaseUrl() string {
	return a.baseUrl
}

func (a *AuthHttpRequester) Head(path string) error {
	req, err := http.NewRequest("HEAD", fmt.Sprintf("%s%s", a.baseUrl, path), nil)
	if err != nil {
		return err
	}
	return a.run(req, nil)
}
