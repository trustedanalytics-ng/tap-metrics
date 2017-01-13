package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/url"
	"net/http"
	"strconv"
	"os"

	"github.com/grafana/grafana/pkg/api/dtos"
)

type GrafanaSyncOperator struct {
	url      string
	authUser *url.Userinfo
	client   *GrafanaClient
}

func (gso *GrafanaSyncOperator) SyncUsers(currentUsers []UaaUser, uaaClient UaaClient) error {
	log.Println("Starting users sync with Grafana")
	users, err := gso.client.GetUsers()
	if err != nil {
		log.Println("Error when recieving Grafana users", err)
		return err
	}

	users = filterOutGrafanaLocalAdmin(users)
	log.Printf("Grafana have %d users (excluding local admin)", len(users))

	err = gso.FixMissingUsersAndFixAdminRoles(currentUsers, users)
	if err != nil {
		log.Println("Error when creating missing users", err)
		return err
	}

	log.Println("Removing zombie users from Grafana")
	err = gso.RemoveZombieUsersFromGrafana(currentUsers, users)
	if err != nil {
		log.Println("Error when removing zombie users from Grafana", err)
		return err
	}
	return nil
}

func (gso *GrafanaSyncOperator) FixMissingUsersAndFixAdminRoles(uaaUsers []UaaUser, grafanaUsers []GrafanaUser) error {
	guMap := make(map[string]*GrafanaUser)
	for i := range grafanaUsers {
		gu := &grafanaUsers[i]
		guMap[gu.Email] = gu
	}
	for _, uu := range uaaUsers {
		gu, present := guMap[uu.email()]
		if !present {
			err := gso.CreateUser(&uu)
			if err != nil {
				return err
			}
		} else {
			if uu.isAdmin() != gu.IsAdmin {
				err := gso.client.SetUserAdmin(gu.Id, uu.isAdmin())
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (gso *GrafanaSyncOperator) RemoveZombieUsersFromGrafana(uaaUsers []UaaUser, grafanaUsers []GrafanaUser) error {
	usersCache := make(map[string]bool)
	for i := range uaaUsers {
		usersCache[uaaUsers[i].email()] = true
	}
	for _, gu := range grafanaUsers {
		_, present := usersCache[gu.Email]
		if !present {
			err := gso.client.DeleteUser(&gu)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (gso *GrafanaSyncOperator) CreateUser(user *UaaUser) error {
	id, err := gso.client.CreateUser(user)
	if err != nil {
		return err
	}
	if user.isAdmin() {
		return gso.client.SetUserAdmin(id, true)
	}
	return nil
}

func (gso *GrafanaSyncOperator) SyncOrganizations(currentOrganizations []UaaOrganization, uaaClient UaaClient) error {
	// TODO - when multi-org will be in place
	return nil
}

func filterOutGrafanaLocalAdmin(users []GrafanaUser) []GrafanaUser {
	for i := range users {
		if users[i].Email == "admin@localhost" {
			return append(users[:i], users[i + 1:]...)
		}
	}
	return users
}

func NewGrafanaSyncOperatorFromEnv() (*GrafanaSyncOperator, error) {
	grafanaUrl := os.Getenv("GRAFANA_URL")
	if grafanaUrl == "" {
		return nil, errors.New("No grafana URL provided")
	}
	user := os.Getenv("GRAFANA_USER")
	pass := os.Getenv("GRAFANA_PASSWORD")
	client := GrafanaClient{
		url: grafanaUrl,
		user: user,
		pass: pass,
		client: &http.Client{},
	}
	return &GrafanaSyncOperator{
		client: &client,
	}, nil
}

type GrafanaClient struct {
	url        string
	user, pass string
	client     *http.Client
}

type GrafanaUser struct {
	Id      int `json:"id"`
	Email   string `json:"email"`
	Name    string `json:"name"`
	Login   string `json:"login"`
	IsAdmin bool `json:"isAdmin"`
}

func (gc *GrafanaClient) makeRequest(method, urlSuffix string, in, out interface{}) error {
	buff := new(bytes.Buffer)
	json.NewEncoder(buff).Encode(in)
	req, err := http.NewRequest(method, gc.url + urlSuffix, buff)
	if err != nil {
		return err
	}
	if gc.user != "" && gc.pass != "" {
		req.SetBasicAuth(gc.user, gc.pass)
	}
	req.Header.Add("Content-Type", "application/json")

	resp, err := gc.client.Do(req)
	if err != nil {
		return nil
	}
	defer resp.Body.Close()
	return json.NewDecoder(resp.Body).Decode(out)
}

func (gc *GrafanaClient) GetUsers() ([]GrafanaUser, error) {
	var users []GrafanaUser
	gc.makeRequest("GET", "/api/users", nil, &users)
	return users, nil
}

func (gc *GrafanaClient) CreateUser(uUser *UaaUser) (int, error) {
	log.Println("In Grafana creating user: " + uUser.UserName)
	form := dtos.AdminCreateUserForm{
		Email: uUser.email(),
		Login: uUser.email(),
		Name: uUser.UserName,
		Password: RandomPassword(16),
	}
	out := grafanaUserCreateResponse{}
	err := gc.makeRequest("POST", "/api/admin/users", &form, &out)
	if err != nil {
		return 0, err
	}
	return out.id, nil
}

func (gc *GrafanaClient) DeleteUser(user *GrafanaUser) error {
	log.Println("In Grafana deleting user: " + user.Login)
	out := grafanaUserCreateResponse{}
	return gc.makeRequest("DELETE", "/api/admin/users/" + strconv.Itoa(user.Id), nil, &out)
}

func (gc *GrafanaClient) SetUserAdmin(userId int, isAdmin bool) error {
	log.Printf("In Grafana setting user id: %v isAdmin: %v", userId, isAdmin)
	form := dtos.AdminUpdateUserPermissionsForm{
		IsGrafanaAdmin:isAdmin,
	}
	out := grafanaUserCreateResponse{}
	gurl := fmt.Sprintf("/api/admin/users/%d/permissions", userId)
	return gc.makeRequest("PUT", gurl, &form, &out)
}

type grafanaUserCreateResponse struct {
	id      int
	message string
}

