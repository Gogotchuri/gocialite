package storages_test

import (
	"context"
	"fmt"
	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
	"testing"

	"github.com/gogotchuri/gocialite"
	"github.com/gogotchuri/gocialite/storages"
	"github.com/gogotchuri/gocialite/structs"
	"github.com/stretchr/testify/assert"
	"golang.org/x/oauth2"
)

// TestMemoryStorage1 tests MemoryStorage by creating instance, setting getting and deleting a Gocialite struct
func TestMemoryStorage1(t *testing.T) {
	storage := storages.NewMemoryStorage()
	state, gocial := createStateAndGocial()
	//set gocial
	err := storage.Set(state, gocial)
	//Check error is nil
	assert.Nil(t, err)
	//get gocial
	gocial2, err := storage.Get(state)
	//Check error is nil
	assert.Nil(t, err)
	//Check gocial equals gocial2
	assert.True(t, gocial.Equals(gocial2))
	//delete gocial
	err = storage.Delete(state)
	//Check error is nil
	assert.Nil(t, err)
	//try getting after delete
	gocial, err = storage.Get(state)
	//Check error is not nil
	assert.NotNil(t, err)
	assert.Nil(t, gocial)
	fmt.Println("TestMemoryStorage1 passed")
}

func TestRedisStorage1(t *testing.T) {
	redisClient := createRedisClient()
	storage := storages.NewRedisStorage(redisClient)
	state, gocial := createStateAndGocial()
	//set gocial
	err := storage.Set(state, gocial)
	//Check error is nil
	assert.Nil(t, err)
	//get gocial
	gocial2, err := storage.Get(state)
	//Check error is nil
	assert.Nil(t, err)
	//Check gocial equals gocial2
	assert.True(t, gocial.Equals(gocial2))
	//delete gocial
	err = storage.Delete(state)
	//Check error is nil
	assert.Nil(t, err)
	//try getting after delete
	gocial, err = storage.Get(state)
	//Check error is not nil
	assert.NotNil(t, err)
	assert.Nil(t, gocial)
	fmt.Println("TestRedisStorage1 passed")
}

func createStateAndGocial() (string, *gocialite.Gocial) {
	state := "ASD123as-33vw"
	//create gocial
	gocial := gocialite.NewGocial("google", state, []string{"email", "profile"}, structs.User{ID: "1as3-412bz", Email: "dooo@dodoo.do"}, &oauth2.Config{ClientID: "clientzz"}, &oauth2.Token{AccessToken: "asdasdasd", RefreshToken: "asdasdasd"})
	return state, gocial
}

func createRedisClient() *redis.Client {
	//Initializing redis
	s, err := miniredis.Run()
	client := redis.NewClient(&redis.Options{
		Addr: s.Addr(),
	})
	_, err = client.Ping(context.Background()).Result()
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}
	return client
}
