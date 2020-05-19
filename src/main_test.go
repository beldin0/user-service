package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/kelseyhightower/envconfig"
	_ "github.com/lib/pq"
	dockertest "github.com/ory/dockertest/v3"
	dc "github.com/ory/dockertest/v3/docker"
	"github.com/phayes/freeport"
	"github.com/stretchr/testify/require"
	"golang.org/x/sync/errgroup"
	"gotest.tools/assert"
)

var db *sqlx.DB

func TestMain(m *testing.M) {
	shutdown, err := setup()
	if err != nil {
		log.Fatal(err)
	}
	defer shutdown()
	if db == nil {
		log.Fatal("failed to connect")
	}
	ctx, cancel := context.WithCancel(context.Background())
	g, ctx := errgroup.WithContext(ctx)
	g.Go(func() error {
		return run(ctx, db)
	})
	log.Println("Waiting for HTTP server to be ready")
	expiry := time.Now().Add(2 * time.Second)
	for {
		if time.Now().After(expiry) {
			log.Fatal("unable to connect to test HTTP server")
		}
		time.Sleep(100 * time.Millisecond)
		_, err := http.Head("http://localhost:8080/users")
		if err == nil {
			break
		}
	}
	code := m.Run()
	cancel()
	g.Wait()
	os.Exit(code)
}

func setup() (teardown func(), err error) {
	c := config{}
	envconfig.Process("", &c)
	c.Host = "localhost"
	c.Password = "testPassword"
	c.Port, err = freeport.GetFreePort()
	if err != nil {
		return nil, err
	}
	log.Println("Starting test Postgres container")
	pool, err := dockertest.NewPool("")
	if err != nil {
		return nil, err
	}
	res, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository:   "postgres",
		ExposedPorts: []string{strconv.Itoa(c.Port)},
		PortBindings: map[dc.Port][]dc.PortBinding{
			"5432/tcp": {
				{
					HostIP:   "localhost",
					HostPort: strconv.Itoa(c.Port),
				},
			},
		},
		Env: []string{
			fmt.Sprintf("POSTGRES_PASSWORD=%s", c.Password),
		},
	})
	if err != nil {
		return nil, err
	}
	expiry := uint(12 * time.Second.Seconds())
	defer res.Expire(expiry)
	purgeFunc := func() {
		if db != nil {
			db.Close()
		}
		pool.Purge(res)
	}
	waitForDb(c)
	return purgeFunc, nil
}

func waitForDb(c config) (err error) {
	log.Println("Waiting for Postgres container")
	expiry := time.Now().Add(6 * time.Second)
	for db == nil {
		if time.Now().After(expiry) {
			return err
		}
		time.Sleep(250 * time.Millisecond)
		db, err = sqlx.Connect("postgres", c.ConnString())
	}
	return nil
}

func TestDB(t *testing.T) {
	// Set up a test user
	user := map[string]interface{}{
		"firstName": "Alan",
		"lastName":  "Smith",
		"nickname":  "alan112",
		"password":  "pass",
		"email":     "alan112@faceit.com",
		"country":   "UK",
	}
	userJSON, err := json.Marshal(user)
	require.NoError(t, err)
	var id float64

	t.Run("Add user", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodPost, "http://localhost:8080/users", bytes.NewReader(userJSON))
		require.NoError(t, err)
		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, resp.StatusCode)
		defer resp.Body.Close()
		jBody := map[string]interface{}{}
		require.NoError(t, json.NewDecoder(resp.Body).Decode(&jBody))
		newID, ok := jBody["id"].(float64)
		require.Equal(t, true, ok)
		assert.Equal(t, true, newID != 0) // assert that the returned user has an ID
		id = newID
	})
	if t.Failed() {
		return
	}
	t.Run("Get user", func(t *testing.T) {
		user := user
		user["id"] = id
		req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("http://localhost:8080/users/%v", id), nil)
		require.NoError(t, err)
		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()
		jBody := map[string]interface{}{}
		require.NoError(t, json.NewDecoder(resp.Body).Decode(&jBody))
		require.Equal(t, http.StatusOK, resp.StatusCode)
		for k, v := range user {
			assert.Equal(t, v, jBody[k]) // assert that the values returned match the values that were added
		}
	})
	if t.Failed() {
		return
	}

	t.Run("Modify user", func(t *testing.T) {
		user["firstName"] = "John"
		userJSON, err := json.Marshal(user)
		require.NoError(t, err)
		req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("http://localhost:8080/users/%v", id), bytes.NewReader(userJSON))
		require.NoError(t, err)
		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, resp.StatusCode)
		defer resp.Body.Close()
		jBody := map[string]interface{}{}
		require.NoError(t, json.NewDecoder(resp.Body).Decode(&jBody))
		require.Equal(t, user, jBody) // assert that the values returned are the updated values
	})
	if t.Failed() {
		return
	}

	t.Run("Search for user - fail", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "http://localhost:8080/users?firstName=alan", nil)
		require.NoError(t, err)
		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, resp.StatusCode)
		defer resp.Body.Close()
		jBody := map[string][]interface{}{}
		require.NoError(t, json.NewDecoder(resp.Body).Decode(&jBody))
		require.Equal(t, 0, len(jBody["users"])) // assert that no results were returned
	})
	if t.Failed() {
		return
	}

	t.Run("Search for user - success", func(t *testing.T) {
		user := user
		user["id"] = id
		req, err := http.NewRequest(http.MethodGet, "http://localhost:8080/users?firstName=john", nil)
		require.NoError(t, err)
		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, resp.StatusCode)
		defer resp.Body.Close()
		jBody := map[string][]interface{}{}
		require.NoError(t, json.NewDecoder(resp.Body).Decode(&jBody))
		require.Equal(t, 1, len(jBody["users"])) // assert that one result was returned
		require.Equal(t, user, jBody["users"][0])
	})
	if t.Failed() {
		return
	}

	t.Run("Delete user", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("http://localhost:8080/users/%v", id), nil)
		require.NoError(t, err)
		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, resp.StatusCode)
	})
	if t.Failed() {
		return
	}

	t.Run("Get user - should be deleted", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("http://localhost:8080/users/%v", id), nil)
		require.NoError(t, err)
		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, resp.StatusCode)
		defer resp.Body.Close()
		jBody := map[string]interface{}{}
		require.NoError(t, json.NewDecoder(resp.Body).Decode(&jBody))
		require.Equal(t, 0, len(jBody)) // assert that no results were returned
	})
}
