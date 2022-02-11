package gocialite

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/gogotchuri/gocialite/drivers"
	"github.com/gogotchuri/gocialite/structs"
	"golang.org/x/oauth2"
	"gopkg.in/oleiade/reflections.v1"
)

func init() {
	drivers.InitializeDrivers(RegisterNewDriver)
}

// Gocial is the main struct of the package with json tags
type Gocial struct {
	driver string
	state  string
	scopes []string
	conf   *oauth2.Config
	User   structs.User
	Token  *oauth2.Token
}

var (
	// Set the basic information such as the endpoint and the scopes URIs
	apiMap = map[string]map[string]string{}

	// Mapping to create a valid "user" struct from providers
	userMap = map[string]map[string]string{}

	// Map correct endpoints
	endpointMap = map[string]oauth2.Endpoint{}

	// Map custom callbacks
	callbackMap = map[string]func(client *http.Client, u *structs.User){}

	// Default scopes for each driver
	defaultScopesMap = map[string][]string{}
)

//RegisterNewDriver adds a new driver to the existing set
func RegisterNewDriver(driver string, defaultscopes []string, callback func(client *http.Client, u *structs.User), endpoint oauth2.Endpoint, apimap, usermap map[string]string) {
	apiMap[driver] = apimap
	userMap[driver] = usermap
	endpointMap[driver] = endpoint
	callbackMap[driver] = callback
	defaultScopesMap[driver] = defaultscopes
}

// Driver is needed to choose the correct social
func (g *Gocial) Driver(driver string) *Gocial {
	g.driver = driver
	g.scopes = defaultScopesMap[driver]

	// BUG: sequential usage of single Gocial instance will have same CSRF token. This is serious security issue.
	// NOTE: Dispatcher eliminates this bug.
	if g.state == "" {
		g.state = randToken()
	}

	return g
}

// Scopes is used to set the oAuth scopes, for example "user", "calendar"
func (g *Gocial) Scopes(scopes []string) *Gocial {
	g.scopes = append(g.scopes, scopes...)
	return g
}

// Redirect returns an URL for the selected social oAuth login
func (g *Gocial) Redirect(clientID, clientSecret, redirectURL string) (string, error) {
	// Check if driver is valid
	if !inSlice(g.driver, complexKeys(apiMap)) {
		return "", fmt.Errorf("Driver not valid: %s", g.driver)
	}

	// Check if valid redirectURL
	_, err := url.ParseRequestURI(redirectURL)
	if err != nil {
		return "", fmt.Errorf("Redirect URL <%s> not valid: %s", redirectURL, err.Error())
	}
	if !strings.HasPrefix(redirectURL, "http://") && !strings.HasPrefix(redirectURL, "https://") {
		return "", fmt.Errorf("Redirect URL <%s> not valid: protocol not valid", redirectURL)
	}

	g.conf = &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,
		Scopes:       g.scopes,
		Endpoint:     endpointMap[g.driver],
	}
	url := g.conf.AuthCodeURL(g.state)

	return url, nil
}

// Handle callback from provider
func (g *Gocial) Handle(state, code string) error {
	// Handle the exchange code to initiate a transport.
	if g.state != state {
		return fmt.Errorf("Invalid state: %s", state)
	}

	// Check if driver is valid
	if !inSlice(g.driver, complexKeys(apiMap)) {
		return fmt.Errorf("Driver not valid: %s", g.driver)
	}

	token, err := g.conf.Exchange(oauth2.NoContext, code)
	if err != nil {
		return fmt.Errorf("oAuth exchanged failed: %s", err.Error())
	}

	client := g.conf.Client(oauth2.NoContext, token)

	// Set gocial token
	g.Token = token

	// Retrieve all from scopes
	driverAPIMap := apiMap[g.driver]
	driverUserMap := userMap[g.driver]
	userEndpoint := strings.Replace(driverAPIMap["userEndpoint"], "%ACCESS_TOKEN", token.AccessToken, -1)

	// Get user info
	req, err := client.Get(driverAPIMap["endpoint"] + userEndpoint)
	if err != nil {
		return err
	}

	defer req.Body.Close()
	res, _ := ioutil.ReadAll(req.Body)
	data, err := jsonDecode(res)
	if err != nil {
		return fmt.Errorf("Error decoding JSON: %s", err.Error())
	}

	// Scan all fields and dispatch through the mapping
	mapKeys := keys(driverUserMap)
	gUser := structs.User{}
	for k, f := range data {
		if !inSlice(k, mapKeys) { // Skip if not in the mapping
			continue
		}

		// Assign the value
		// Dirty way, but we need to convert also int/float to string
		_ = reflections.SetField(&gUser, driverUserMap[k], fmt.Sprint(f))
	}

	// Set the "raw" user interface
	gUser.Raw = data

	// Custom callback
	callbackMap[g.driver](client, &gUser)

	// Update the struct
	g.User = gUser

	return nil
}

//Marshal marshals Gocial struct to JSON
func Marshal(g *Gocial) ([]byte, error) {
	return json.Marshal(g.createJSONable())
}

//Unmarshal the JSON to Gocial
func Unmarshal(data []byte) (*Gocial, error) {
	gj := &gocialJSONable{}
	err := json.Unmarshal(data, gj)
	if err != nil {
		return nil, err
	}
	g := &Gocial{}
	g.fillFromJSONable(gj)

	return g, nil
}

//NewGocial Create Gocial instance from passed arguments
func NewGocial(driver, state string, scopes []string, user structs.User, conf *oauth2.Config, token *oauth2.Token) *Gocial {
	g := &Gocial{}
	g.driver = driver
	g.scopes = scopes
	g.state = state
	g.User = user
	g.conf = conf
	g.Token = token
	return g
}

//Equals compares two Gocial instances
func (g *Gocial) Equals(g2 *Gocial) bool {
	//Compare basic fields
	if g.User.ID != g2.User.ID {
		return false
	}
	if g.driver != g2.driver {
		return false
	}
	if g.state != g2.state {
		return false
	}
	//Compare tokens
	if g.Token != nil {
		if g2.Token == nil {
			return false
		}
		if g.Token.AccessToken != g2.Token.AccessToken {
			return false
		}
		if g.Token.Expiry != g2.Token.Expiry {
			return false
		}
		if g.Token.RefreshToken != g2.Token.RefreshToken {
			return false
		}
	}

	if g.User.ID != g2.User.ID {
		return false
	}
	if g.User.Email != g2.User.Email {
		return false
	}
	//Compare configs
	if g.conf != nil {
		if g2.conf == nil {
			return false
		}
		if g.conf.ClientID != g2.conf.ClientID {
			return false
		}
		if g.conf.ClientSecret != g2.conf.ClientSecret {
			return false
		}
		if g.conf.RedirectURL != g2.conf.RedirectURL {
			return false
		}
	}
	//Compare scopes
	if len(g.scopes) != len(g2.scopes) {
		return false
	}
	for i := range g.scopes {
		if g.scopes[i] != g2.scopes[i] {
			return false
		}
	}
	return true
}
