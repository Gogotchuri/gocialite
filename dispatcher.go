package gocialite

import (
	"fmt"
	"sync"

	"github.com/gogotchuri/gocialite/structs"
	"golang.org/x/oauth2"
)

// Dispatcher allows to safely issue concurrent Gocials
type Dispatcher struct {
	mu      sync.RWMutex
	storage GocialStorage
}

// NewDispatcher creates new Dispatcher
func NewDispatcher(storage GocialStorage) *Dispatcher {
	return &Dispatcher{storage: storage}
}

// New Gocial instance
func (d *Dispatcher) New() *Gocial {
	d.mu.Lock()
	defer d.mu.Unlock()
	state := randToken()
	g := &Gocial{state: state}
	d.storage.Set(state, g)
	return g
}

type StateConf struct {
	Driver      string
	ClientID    string
	Secret      string
	RedirectURL string
	Scopes      []string
}

//GenerateRedirectURL creates url for authorization with new state. basically shorthand for New().Driver().Scopes().Redirect()
func (d *Dispatcher) GenerateRedirectURL(sc StateConf) (string, error) {
	d.mu.Lock()
	defer d.mu.Unlock()
	state := randToken()
	g := &Gocial{state: state}
	//Sets driver
	g.Driver(sc.Driver)
	//Sets scopes
	if sc.Scopes != nil && len(sc.Scopes) > 0 {
		g.Scopes(sc.Scopes)
	}
	//Generate redirection url
	url, err := g.Redirect(sc.ClientID, sc.Secret, sc.RedirectURL)
	if err != nil {
		return "", err
	}
	//Store gocial instance id url is generated
	d.storage.Set(state, g)
	return url, err
}

//Update Gocial instance, should be called after every update on gocial instance
func (d *Dispatcher) Update(g *Gocial) error {
	d.mu.Lock()
	defer d.mu.Unlock()
	return d.storage.Set(g.state, g)
}

// Handle callback. Can be called only once for given state.
func (d *Dispatcher) Handle(state, code string) (*structs.User, *oauth2.Token, error) {
	d.mu.RLock()
	g, err := d.storage.Get(state)
	d.mu.RUnlock()
	if err != nil {
		return nil, nil, fmt.Errorf("invalid CSRF token: %s \n %s \n", state, err.Error())
	}
	err = g.Handle(state, code)
	d.mu.Lock()
	d.storage.Delete(state)
	d.mu.Unlock()
	return &g.User, g.Token, err
}
