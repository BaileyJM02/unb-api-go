package v1

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"
	"fmt"
	"path/filepath"
	"runtime"
	"reflect"
	"errors"

)



// assert fails the test if the condition is false.
func assert(tb testing.TB, condition bool, msg string, v ...interface{}) {
	if !condition {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d: "+msg+"\033[39m\n\n", append([]interface{}{filepath.Base(file), line}, v...)...)
		tb.FailNow()
	}
}

// ok fails the test if an err is not nil.
func ok(tb testing.TB, err error) {
	if err != nil {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d: unexpected error: %s\033[39m\n\n", filepath.Base(file), line, err.Error())
		tb.FailNow()
	}
}

// equals fails the test if exp is not equal to act.
func equals(tb testing.TB, exp, act interface{}) {
	if !reflect.DeepEqual(exp, act) {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d:\n\n\texp: %#v\n\n\tgot: %#v\033[39m\n\n", filepath.Base(file), line, exp, act)
		tb.FailNow()
	}
}

// Easy set-up
func setClient(code int, url, data string) *http.Client {
    client := NewTestClient(func(req *http.Request) *http.Response {
		// Test request parameters
//		equals(t, req.URL.String(), "https://unbelievable.pizza/api/v1"+url)
		return &http.Response{
			StatusCode: code,
			// Send response to be tested
			Body:       ioutil.NopCloser(bytes.NewBufferString(data)),
 			// Must be set to non-nil value or it panics
			Header:     make(http.Header),
		}
	})
	return client
}

// Test that it works or nothing else will.
func TestClientReturnsCorrectVal(t *testing.T) {
    setClient(200, "/ping", `{"test":"data"}`)
    client := NewTestClient(func(req *http.Request) *http.Response {
		// Test request parameters
		equals(t, req.URL.String(), "https://unbelievable.pizza/api/v1/ping")
		return &http.Response{
			StatusCode: 200,
			// Send response to be tested
			Body:       ioutil.NopCloser(bytes.NewBufferString(`{"test":"data"}`)),
 			// Must be set to non-nil value or it panics
			Header:     make(http.Header),
		}
	})
	
	equals(t, setClient(200, "/ping", `{"test":"data"}`).CheckRedirect,client.CheckRedirect)
}


// RoundTripFunc .
type RoundTripFunc func(req *http.Request) *http.Response

// RoundTrip .
func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}

//NewTestClient returns *http.Client with Transport replaced to avoid making real calls
func NewTestClient(fn RoundTripFunc) *http.Client {
	return &http.Client{
		Transport: RoundTripFunc(fn),
	}
}

func TestCheckReturnsIsUpOn404(t *testing.T) {
	client := NewTestClient(func(req *http.Request) *http.Response {
		// Test request parameters
		equals(t, req.URL.String(), "https://unbelievable.pizza/api/v1")
		return &http.Response{
			StatusCode: 200,
			// Send response to be tested
			Body:       ioutil.NopCloser(bytes.NewBufferString(`{"error":"404: Not found"}`)),
 			// Must be set to non-nil value or it panics
			Header:     make(http.Header),
		}
	})

	api := Custom("token", client)
	check, err := api.Check()
	ok(t, err)
	equals(t, true, check.Up)

}

func TestCheckReturnsIsUpButErrorCannotConnectOn401(t *testing.T) {
	client := setClient(200, "", `{"error":"401: Unauthorized"}`)

	api := Custom("token", client)
	check, err := api.Check()
	equals(t, true, check.Up)
	equals(t, errors.New("401 Unauthorized (Check your token)"), err)

}

func TestCheckReturnsIsDownWhen500(t *testing.T) {
	client := setClient(500, "", ``)

	api := Custom("token", client)
	check, err := api.Check()
	equals(t, errors.New("Cannot Connect to API url."), err)
	equals(t, false, check.Up)
	
}

func TestUserBalanceHandlesDataOnSuccessfulFetchCorrectlyNonInfinite(t *testing.T) {
	client := setClient(200, "/guilds/411898639737421824/users/398197113495748626", `{"rank":"14","user_id":"398197113495748626","cash":25,"bank":200,"total":526}`)

	api := Custom("token", client)
	data, err := api.UserBalance("411898639737421824", "398197113495748626") // Guild, User
	ok(t, err)
	equals(t, userBalance{14,"398197113495748626",25,false,false,200,false,false,526,false,false}, data)
}

func TestUserBalanceHandlesDataOnSuccessfulFetchCorrectlyWithNoRank(t *testing.T) {
    client := setClient(200, "/guilds/411898639737421824/users/398197113495748626", `{"user_id":"398197113495748626","cash":25,"bank":200,"total":225}`)

	api := Custom("token", client)
	data, err := api.UserBalance("411898639737421824", "398197113495748626") // Guild, User
	ok(t, err)
	equals(t, userBalance{0,"398197113495748626",25,false,false,200,false,false,225,false,false}, data)
}

func TestUserBalanceHandlesDataOnSuccessfulFetchCorrectlyWithNoRankWithInfiniteCash(t *testing.T) {
    client := setClient(200, "/guilds/411898639737421824/users/398197113495748626", `{"user_id":"398197113495748626","cash":"Infinity","bank":200,"total":"Infinity"}`)

	api := Custom("token", client)
	data, err := api.UserBalance("411898639737421824", "398197113495748626") // Guild, User
	ok(t, err)
	equals(t, userBalance{0,"398197113495748626",0,true,false,200,false,false,0,true,false}, data)
}

func TestUserBalanceHandlesDataOnSuccessfulFetchCorrectlyWithNoRankWithInfiniteBank(t *testing.T) {
    client := setClient(200, "/guilds/411898639737421824/users/398197113495748626", `{"user_id":"398197113495748626","cash":25,"bank":"Infinity","total":"Infinity"}`)

	api := Custom("token", client)
	data, err := api.UserBalance("411898639737421824", "398197113495748626") // Guild, User
	ok(t, err)
	equals(t, userBalance{0,"398197113495748626",25,false,false,0,true,false,0,true,false}, data)
}

func TestUserBalanceHandlesDataOnSuccessfulFetchCorrectlyWithNoRankWithInfiniteBankAndCash(t *testing.T) {
    client := setClient(200, "/guilds/411898639737421824/users/398197113495748626", `{"user_id":"398197113495748626","cash":"Infinity","bank":"Infinity","total":"Infinity"}`)

	api := Custom("token", client)
	data, err := api.UserBalance("411898639737421824", "398197113495748626") // Guild, User
	ok(t, err)
	equals(t, userBalance{0,"398197113495748626",0,true,false,0,true,false,0,true,false}, data)
}

func TestUserBalanceHandlesDataOnSuccessfulFetchCorrectlyWithNoRankWithNegitiveInfiniteCash(t *testing.T) {
    client := setClient(200, "/guilds/411898639737421824/users/398197113495748626", `{"user_id":"398197113495748626","cash":"-Infinity","bank":200,"total":"-Infinity"}`)

	api := Custom("token", client)
	data, err := api.UserBalance("411898639737421824", "398197113495748626") // Guild, User
	ok(t, err)
	equals(t, userBalance{0,"398197113495748626",0,false,true,200,false,false,0,false,true}, data)
}

func TestUserBalanceHandlesDataOnSuccessfulFetchCorrectlyWithNoRankWithNegitiveInfiniteBank(t *testing.T) {
    client := setClient(200, "/guilds/411898639737421824/users/398197113495748626", `{"user_id":"398197113495748626","cash":25,"bank":"-Infinity","total":"-Infinity"}`)

	api := Custom("token", client)
	data, err := api.UserBalance("411898639737421824", "398197113495748626") // Guild, User
	ok(t, err)
	equals(t, userBalance{0,"398197113495748626",25,false,false,0,false,true,0,false,true}, data)
}

func TestUserBalanceHandlesDataOnSuccessfulFetchCorrectlyWithNoRankWithNegitiveInfiniteBankAndCash(t *testing.T) {
    client := setClient(200, "/guilds/411898639737421824/users/398197113495748626", `{"user_id":"398197113495748626","cash":"-Infinity","bank":"-Infinity","total":"-Infinity"}`)

	api := Custom("token", client)
	data, err := api.UserBalance("411898639737421824", "398197113495748626") // Guild, User
	ok(t, err)
	equals(t, userBalance{0,"398197113495748626",0,false,true,0,false,true,0,false,true}, data)
}

func TestUserBalanceHandlesDataOnSuccessfulFetchCorrectlyWithNoRankWithNegitiveInfiniteBankAndInfiniteCash(t *testing.T) {
    client := setClient(200, "/guilds/411898639737421824/users/398197113495748626", `{"user_id":"398197113495748626","cash":"Infinity","bank":"-Infinity","total":0}`)

	api := Custom("token", client)
	data, err := api.UserBalance("411898639737421824", "398197113495748626") // Guild, User
	ok(t, err)
	equals(t, userBalance{0,"398197113495748626",0,true,false,0,false,true,0,false,false}, data)
}

func TestUserBalanceHandlesDataOnSuccessfulFetchCorrectlyWithNoRankWithInfiniteBankAndNegitiveInfiniteCash(t *testing.T) {
    client := setClient(200, "/guilds/411898639737421824/users/398197113495748626", `{"user_id":"398197113495748626","cash":"-Infinity","bank":"Infinity","total":0}`)

	api := Custom("token", client)
	data, err := api.UserBalance("411898639737421824", "398197113495748626") // Guild, User
	ok(t, err)
	equals(t, userBalance{0,"398197113495748626",0,false,true,0,true,false,0,false,false}, data)
}

func TestUserBalanceHandlesDataOnUnsuccessfulFetchCorrectlyWithIncorrectGuild(t *testing.T) {
    client := setClient(200, "/guilds/000000000000000000/users/398197113495748626", `{"error":"404: Not found","message":"Unknown guild"}`)

	api := Custom("token", client)
	data, err := api.UserBalance("000000000000000000", "398197113495748626") // Guild, User
	equals(t, userBalance{}, data)
	equals(t, errors.New("404: Not found (Unknown guild)"), err)
}

func TestUserBalanceHandlesDataOnUnsuccessfulFetchCorrectlyWithIncorrectUser(t *testing.T) {
    client := setClient(200, "/guilds/411898639737421824/users/000000000000000000", `{"error":"404: Not found","message":"Unknown user"}`)

	api := Custom("token", client)
	data, err := api.UserBalance("411898639737421824", "000000000000000000") // Guild, User
	equals(t, userBalance{}, data)
	equals(t, errors.New("404: Not found (Unknown user)"), err)
}

