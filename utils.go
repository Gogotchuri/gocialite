package gocialite

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"strings"

	"github.com/danilopolani/gocialite/structs"
	"golang.org/x/oauth2"
)

//gocialJSONable 'private' Struct to Marshal the Gocial instance into JSON, including private fields
type gocialJSONable struct {
	Driver string         `json:"driver"`
	State  string         `json:"state"`
	Scopes []string       `json:"scopes"`
	Conf   *oauth2.Config `json:"conf"`
	User   structs.User   `json:"user"`
	Token  *oauth2.Token  `json:"token"`
}

func (g *Gocial) createJSONable() *gocialJSONable {
	return &gocialJSONable{
		Driver: g.driver,
		State:  g.state,
		Scopes: g.scopes,
		Conf:   g.conf,
		User:   g.User,
		Token:  g.Token,
	}
}

func (g *Gocial) fillFromJSONable(gj *gocialJSONable) {
	g.driver = gj.Driver
	g.state = gj.State
	g.scopes = gj.Scopes
	g.conf = gj.Conf
	g.User = gj.User
	g.Token = gj.Token
}

// Generate a random token
func randToken() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.StdEncoding.EncodeToString(b)
}

// Check if a value is in a string slice
func inSlice(v string, s []string) bool {
	for _, scope := range s {
		if scope == v {
			return true
		}
	}

	return false
}

// Decode a json or return an error
func jsonDecode(js []byte) (map[string]interface{}, error) {
	var decoded map[string]interface{}
	decoder := json.NewDecoder(strings.NewReader(string(js)))
	decoder.UseNumber()

	if err := decoder.Decode(&decoded); err != nil {
		return nil, err
	}

	return decoded, nil
}

// Return the keys of a map
func keys(m map[string]string) []string {
	var keys []string
	for k := range m {
		keys = append(keys, k)
	}

	return keys
}

func complexKeys(m map[string]map[string]string) []string {
	var keys []string
	for k := range m {
		keys = append(keys, k)
	}

	return keys
}
