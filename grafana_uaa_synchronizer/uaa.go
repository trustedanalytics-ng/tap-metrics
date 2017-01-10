package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/cloudfoundry-incubator/uaago"
)

type uaaToken struct {
	content      string
	refreshAfter time.Time
}

type UaaUser struct {
	Active               bool          `json:"active"`
	Approvals            []interface{} `json:"approvals"`
	Emails               []struct {
		Primary bool   `json:"primary"`
		Value   string `json:"value"`
	} `json:"emails"`
	Groups               []struct {
		Display string `json:"display"`
		Type    string `json:"type"`
		Value   string `json:"value"`
	} `json:"groups"`
	ID                   string `json:"id"`
	Meta                 struct {
				     Created      string `json:"created"`
				     LastModified string `json:"lastModified"`
				     Version      int    `json:"version"`
			     } `json:"meta"`
	Name                 struct{} `json:"name"`
	Origin               string   `json:"origin"`
	PasswordLastModified string   `json:"passwordLastModified"`
	Schemas              []string `json:"schemas"`
	UserName             string   `json:"userName"`
	Verified             bool     `json:"verified"`
	ZoneID               string   `json:"zoneId"`
}

func (user *UaaUser) isAdmin() bool {
	for i := range user.Groups {
		if user.Groups[i].Display == "tap.admin" {
			return true
		}
	}
	return false
}

type UaaOrganization struct{}

type UaaClient interface {
	EnsureCredentials() error
	GetUsers() ([]UaaUser, error)
	//GetOrganizations() ([]*UaaOrganization, error) - waits for multi-org
}

type HttpGetter struct {
	url     string
	hClient *http.Client
}

func (hg *HttpGetter) run(path, token string, destination interface{}) error {
	req, err := http.NewRequest("GET", hg.url + path, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", token)
	resp, err := hg.hClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return json.NewDecoder(resp.Body).Decode(destination)
}

type UaaHttpClient struct {
	url          string
	clientId     string
	clientSecret string

	tokenClient  *uaago.Client
	currentToken *uaaToken

	httpGetter   *HttpGetter
}

func (uhc *UaaHttpClient) GetUsers() ([]UaaUser, error) {
	usersResponse := &UaaUsersResponse{}
	// Probably more than one page users won't require any fixes

	err := uhc.httpGetter.run("/Users?" +
		"sortOrder=descending" +
		"&sortBy=lastModified" +
		"&count=1000",
		uhc.currentToken.content,
		usersResponse)
	if err != nil {
		return nil, err
	}

	return usersResponse.Resources, nil
}

func (uhc *UaaHttpClient) EnsureCredentials() error {
	// Probably no need for protecting it with locking
	if uhc.currentToken == nil || uhc.currentToken.refreshAfter.Before(time.Now()) {
		return uhc.refreshToken()
	}
	return nil
}

func (uhc *UaaHttpClient) refreshToken() error {
	token, expiresIn, err := uhc.tokenClient.GetAuthTokenWithExpiresIn(uhc.clientId, uhc.clientSecret, true)
	if err != nil {
		return err
	}
	timeShift := time.Duration(expiresIn - 10) * time.Second
	uhc.currentToken = &uaaToken{
		content: token,
		refreshAfter: time.Now().Add(timeShift),
	}
	return nil
}

func NewUaaClientFromEnv() (UaaClient, error) {
	uaaUrl := os.Getenv("UAA_URL")
	if uaaUrl == "" {
		return nil, errors.New("No UAA url given")
	}
	uaaUrl = strings.TrimSuffix(uaaUrl, "/")
	clientId := os.Getenv("CLIENT_ID")
	if clientId == "" {
		return nil, errors.New("No UAA Client ID given")
	}
	clientSecret := os.Getenv("CLIENT_SECRET")
	if clientSecret == "" {
		return nil, errors.New("No UAA Client secret given")
	}

	tokenClient, err := uaago.NewClient(uaaUrl)
	handlePotentialSetupError(err, "Unable to setup UAA client")

	hGetter := HttpGetter{
		url: uaaUrl,
		hClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}

	return &UaaHttpClient{
		url: uaaUrl,
		clientId: clientId,
		clientSecret: clientSecret,

		tokenClient: tokenClient,

		httpGetter: &hGetter,
	}, nil
}



// UAA structs

type UaaUsersResponse struct {
	ItemsPerPage int `json:"itemsPerPage"`
	Resources    []UaaUser `json:"resources"`
	Schemas      []string `json:"schemas"`
	StartIndex   int      `json:"startIndex"`
	TotalResults int      `json:"totalResults"`
}

