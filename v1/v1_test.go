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

func TestClientReturnsRatelimitedWhen427(t *testing.T) {
	client := setClient(429, "", `{"message":"You are being rate limited.","retry_after":36191}`)

	api := Custom("token", client)
	check, err := api.Check()
	equals(t, errors.New("You are being rate limited. Retry after: 36.191µs"), err)
	equals(t, true, check.Up)
}

func TestGetBalanceHandlesDataOnSuccessfulFetchCorrectlyNonInfinite(t *testing.T) {
	client := setClient(200, "/guilds/411898639737421824/users/398197113495748626", `{"rank":"14","user_id":"398197113495748626","cash":25,"bank":200,"total":526}`)

	api := Custom("token", client)
	data, err := api.GetBalance("411898639737421824", "398197113495748626") // Guild, User
	ok(t, err)
	equals(t, userObj{14,"398197113495748626",25,false,false,200,false,false,526,false,false }, data)
}

func TestGetBalanceHandlesDataOnSuccessfulFetchCorrectlyWithNoRank(t *testing.T) {
    client := setClient(200, "/guilds/411898639737421824/users/398197113495748626", `{"user_id":"398197113495748626","cash":25,"bank":200,"total":225}`)

	api := Custom("token", client)
	data, err := api.GetBalance("411898639737421824", "398197113495748626") // Guild, User
	ok(t, err)
	equals(t, userObj{0,"398197113495748626",25,false,false,200,false,false,225,false,false }, data)
}

func TestGetBalanceHandlesDataOnSuccessfulFetchCorrectlyWithNoRankWithInfiniteCash(t *testing.T) {
    client := setClient(200, "/guilds/411898639737421824/users/398197113495748626", `{"user_id":"398197113495748626","cash":"Infinity","bank":200,"total":"Infinity"}`)

	api := Custom("token", client)
	data, err := api.GetBalance("411898639737421824", "398197113495748626") // Guild, User
	ok(t, err)
	equals(t, userObj{0,"398197113495748626",0,true,false,200,false,false,0,true,false }, data)
}

func TestGetBalanceHandlesDataOnSuccessfulFetchCorrectlyWithNoRankWithInfiniteBank(t *testing.T) {
    client := setClient(200, "/guilds/411898639737421824/users/398197113495748626", `{"user_id":"398197113495748626","cash":25,"bank":"Infinity","total":"Infinity"}`)

	api := Custom("token", client)
	data, err := api.GetBalance("411898639737421824", "398197113495748626") // Guild, User
	ok(t, err)
	equals(t, userObj{0,"398197113495748626",25,false,false,0,true,false,0,true,false }, data)
}

func TestGetBalanceHandlesDataOnSuccessfulFetchCorrectlyWithNoRankWithInfiniteBankAndCash(t *testing.T) {
    client := setClient(200, "/guilds/411898639737421824/users/398197113495748626", `{"user_id":"398197113495748626","cash":"Infinity","bank":"Infinity","total":"Infinity"}`)

	api := Custom("token", client)
	data, err := api.GetBalance("411898639737421824", "398197113495748626") // Guild, User
	ok(t, err)
	equals(t, userObj{0,"398197113495748626",0,true,false,0,true,false,0,true,false }, data)
}

func TestGetBalanceHandlesDataOnSuccessfulFetchCorrectlyWithNoRankWithNegitiveInfiniteCash(t *testing.T) {
    client := setClient(200, "/guilds/411898639737421824/users/398197113495748626", `{"user_id":"398197113495748626","cash":"-Infinity","bank":200,"total":"-Infinity"}`)

	api := Custom("token", client)
	data, err := api.GetBalance("411898639737421824", "398197113495748626") // Guild, User
	ok(t, err)
	equals(t, userObj{0,"398197113495748626",0,false,true,200,false,false,0,false,true }, data)
}

func TestGetBalanceHandlesDataOnSuccessfulFetchCorrectlyWithNoRankWithNegitiveInfiniteBank(t *testing.T) {
    client := setClient(200, "/guilds/411898639737421824/users/398197113495748626", `{"user_id":"398197113495748626","cash":25,"bank":"-Infinity","total":"-Infinity"}`)

	api := Custom("token", client)
	data, err := api.GetBalance("411898639737421824", "398197113495748626") // Guild, User
	ok(t, err)
	equals(t, userObj{0,"398197113495748626",25,false,false,0,false,true,0,false,true }, data)
}

func TestGetBalanceHandlesDataOnSuccessfulFetchCorrectlyWithNoRankWithNegitiveInfiniteBankAndCash(t *testing.T) {
    client := setClient(200, "/guilds/411898639737421824/users/398197113495748626", `{"user_id":"398197113495748626","cash":"-Infinity","bank":"-Infinity","total":"-Infinity"}`)

	api := Custom("token", client)
	data, err := api.GetBalance("411898639737421824", "398197113495748626") // Guild, User
	ok(t, err)
	equals(t, userObj{0,"398197113495748626",0,false,true,0,false,true,0,false,true }, data)
}

func TestGetBalanceHandlesDataOnSuccessfulFetchCorrectlyWithNoRankWithNegitiveInfiniteBankAndInfiniteCash(t *testing.T) {
    client := setClient(200, "/guilds/411898639737421824/users/398197113495748626", `{"user_id":"398197113495748626","cash":"Infinity","bank":"-Infinity","total":0}`)

	api := Custom("token", client)
	data, err := api.GetBalance("411898639737421824", "398197113495748626") // Guild, User
	ok(t, err)
	equals(t, userObj{0,"398197113495748626",0,true,false,0,false,true,0,false,false }, data)
}

func TestGetBalanceHandlesDataOnSuccessfulFetchCorrectlyWithNoRankWithInfiniteBankAndNegitiveInfiniteCash(t *testing.T) {
    client := setClient(200, "/guilds/411898639737421824/users/398197113495748626", `{"user_id":"398197113495748626","cash":"-Infinity","bank":"Infinity","total":0}`)

	api := Custom("token", client)
	data, err := api.GetBalance("411898639737421824", "398197113495748626") // Guild, User
	ok(t, err)
	equals(t, userObj{0,"398197113495748626",0,false,true,0,true,false,0,false,false }, data)
}

func TestGetBalanceHandlesDataOnUnsuccessfulFetchCorrectlyWithIncorrectGuild(t *testing.T) {
    client := setClient(200, "/guilds/000000000000000000/users/398197113495748626", `{"error":"404: Not found","message":"Unknown guild"}`)

	api := Custom("token", client)
	data, err := api.GetBalance("000000000000000000", "398197113495748626") // Guild, User
	equals(t, userObj{}, data)
	equals(t, errors.New("404: Not found (Unknown guild)"), err)
}

func TestGetBalanceHandlesDataOnUnsuccessfulFetchCorrectlyWithIncorrectUser(t *testing.T) {
    client := setClient(200, "/guilds/411898639737421824/users/000000000000000000", `{"error":"404: Not found","message":"Unknown user"}`)

	api := Custom("token", client)
	data, err := api.GetBalance("411898639737421824", "000000000000000000") // Guild, User
	equals(t, userObj{}, data)
	equals(t, errors.New("404: Not found (Unknown user)"), err)
}

func TestLeaderboardHandlesDataOnSuccessfulFetchCorrectly(t *testing.T) {
    client := setClient(200, "/guilds/411898639737421824/users", `[{"rank":"1","user_id":"116293018742554625","cash":"Infinity","bank":0,"total":"Infinity"},{"rank":"2","user_id":"398197113495748626","cash":"-Infinity","bank":"Infinity","total":0},{"rank":"3","user_id":"000000000000000000","cash":"33","bank":"Infinity","total":"Infinity"}]`)

	api := Custom("token", client)
	data, err := api.Leaderboard("411898639737421824") // Guild
	ok(t, err)
	equals(t, []userObj{userObj{1, "116293018742554625", 0, true, false, 0, false, false, 0, true, false}, userObj{2, "398197113495748626", 0, false, true, 0, true, false, 0, false, false}, userObj{3, "000000000000000000", 33, false, false, 0, true, false, 0, true, false}}, data)
	equals(t,userObj{1,"116293018742554625",0,true,false,0,false,false,0,true,false } ,data[0])
}

func TestLeaderboardHandlesErrorOnUnsuccessfulFetchCorrectly(t *testing.T) {
    client := setClient(200, "/guilds/000000000000000000/users", `{"error":"404: Not found","message":"Unknown guild"}`)

	api := Custom("token", client)
	data, err := api.Leaderboard("000000000000000000") // Guild
	equals(t, []userObj{}, data)
	equals(t, errors.New("404: Not found (Unknown guild)"), err)
}

func TestSetBalanceWithNonInfiniteData(t *testing.T) {
    client := setClient(200, "/guilds/411898639737421824/users/398197113495748626", `{"user_id":"398197113495748626","cash":50,"bank":502,"total":552,"found":true}`)
    
    api := Custom("token", client)
	data, err := api.SetBalance("411898639737421824", "398197113495748626", 50, 502, "Just testing")
	ok(t, err)
	equals(t, userObj{0,"398197113495748626",50,false,false,502,false,false,552,false,false}, data)
}

func TestSetBalanceWithOnlyCashInfiniteData(t *testing.T) {
    client := setClient(200, "/guilds/411898639737421824/users/398197113495748626", `{"user_id":"398197113495748626","cash":"Infinity","bank":502,"total":"Infinity","found":true}`)
    
    api := Custom("token", client)
	data, err := api.SetBalance("411898639737421824", "398197113495748626", "Infinity", 502, "Just testing")
	ok(t, err)
	equals(t, userObj{0,"398197113495748626",0,true,false,502,false,false,0,true,false}, data)
}

func TestSetBalanceWithOnlyBankInfiniteData(t *testing.T) {
    client := setClient(200, "/guilds/411898639737421824/users/398197113495748626", `{"user_id":"398197113495748626","cash":50,"bank":"Infinity","total":"Infinity","found":true}`)
    
    api := Custom("token", client)
	data, err := api.SetBalance("411898639737421824", "398197113495748626", 50, "Infinity", "Just testing")
	ok(t, err)
	equals(t, userObj{0,"398197113495748626",50,false,false,0,true,false,0,true,false}, data)
}

func TestSetBalanceWithOnlyCashNegitiveInfiniteData(t *testing.T) {
    client := setClient(200, "/guilds/411898639737421824/users/398197113495748626", `{"user_id":"398197113495748626","cash":"-Infinity","bank":502,"total":"-Infinity","found":true}`)
    
    api := Custom("token", client)
	data, err := api.SetBalance("411898639737421824", "398197113495748626", "Infinity", 502, "Just testing")
	ok(t, err)
	equals(t, userObj{0,"398197113495748626",0,false,true,502,false,false,0,false,true}, data)
}

func TestSetBalanceWithOnlyBankNegitiveInfiniteData(t *testing.T) {
    client := setClient(200, "/guilds/411898639737421824/users/398197113495748626", `{"user_id":"398197113495748626","cash":50,"bank":"-Infinity","total":"-Infinity","found":true}`)
    
    api := Custom("token", client)
	data, err := api.SetBalance("411898639737421824", "398197113495748626", 50, "Infinity", "Just testing")
	ok(t, err)
	equals(t, userObj{0,"398197113495748626",50,false,false,0,false,true,0,false,true}, data)
}

func TestSetBalanceWithAllInfiniteData(t *testing.T) {
    client := setClient(200, "/guilds/411898639737421824/users/398197113495748626", `{"user_id":"398197113495748626","cash":"Infinity","bank":"Infinity","total":"Infinity","found":true}`)
    
    api := Custom("token", client)
	data, err := api.SetBalance("411898639737421824", "398197113495748626", "Infinity", "Infinity", "Just testing")
	ok(t, err)
	equals(t, userObj{0,"398197113495748626",0,true,false,0,true,false,0,true,false}, data)
}

func TestSetBalanceWithAllNegitiveInfiniteData(t *testing.T) {
    client := setClient(200, "/guilds/411898639737421824/users/398197113495748626", `{"user_id":"398197113495748626","cash":"-Infinity","bank":"-Infinity","total":"-Infinity","found":true}`)
    
    api := Custom("token", client)
	data, err := api.SetBalance("411898639737421824", "398197113495748626", "-Infinity", "-Infinity", "Just testing")
	ok(t, err)
	equals(t, userObj{0,"398197113495748626",0,false,true,0,false,true,0,false,true}, data)
}

func TestSetBalanceWithCashInfiniteCashNegitiveInfinite(t *testing.T) {
    client := setClient(200, "/guilds/411898639737421824/users/398197113495748626", `{"user_id":"398197113495748626","cash":"-Infinity","bank":"Infinity","total":0,"found":true}`)
    
    api := Custom("token", client)
	data, err := api.SetBalance("411898639737421824", "398197113495748626", "-Infinity", "Infinity", "Just testing")
	ok(t, err)
	equals(t, userObj{0,"398197113495748626",0,false,true,0,true,false,0,false,false}, data)
}

func TestSetBalanceWithCashInfiniteBankNegitiveInfinite(t *testing.T) {
    client := setClient(200, "/guilds/411898639737421824/users/398197113495748626", `{"user_id":"398197113495748626","cash":"Infinity","bank":"-Infinity","total":0,"found":true}`)
    
    api := Custom("token", client)
	data, err := api.SetBalance("411898639737421824", "398197113495748626", "Infinity", "-Infinity", "Just testing")
	ok(t, err)
	equals(t, userObj{0,"398197113495748626",0,true,false,0,false,true,0,false,false}, data)
}

func TestUpdateBalanceWithCorrectData(t *testing.T) {
    client := setClient(200, "/guilds/411898639737421824/users/398197113495748626", `{"user_id":"398197113495748626","cash":50,"bank":502,"total":552,"found":true}`)
    
    api := Custom("token", client)
	data, err := api.UpdateBalance("411898639737421824", "398197113495748626", 0, 0, "Just testing")
	ok(t, err)
	equals(t, userObj{0,"398197113495748626",50,false,false,502,false,false,552,false,false}, data)
}

func TestUpdateBalanceWithNegitiveData(t *testing.T) {
    client := setClient(200, "/guilds/411898639737421824/users/398197113495748626", `{"user_id":"398197113495748626","cash":50,"bank":502,"total":552,"found":true}`)
    
    api := Custom("token", client)
	data, err := api.UpdateBalance("411898639737421824", "398197113495748626", -40, -980, "Just testing")
	ok(t, err)
	equals(t, userObj{0,"398197113495748626",50,false,false,502,false,false,552,false,false}, data)
}


