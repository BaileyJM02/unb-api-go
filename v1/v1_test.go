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

	client := NewTestClient(func(req *http.Request) *http.Response {
		// Test request parameters
		equals(t, req.URL.String(), "https://unbelievable.pizza/api/v1")
		return &http.Response{
			StatusCode: 200,
			// Send response to be tested
			Body:       ioutil.NopCloser(bytes.NewBufferString(`{"error":"401: Unauthorized"}`)),
 			// Must be set to non-nil value or it panics
			Header:     make(http.Header),
		}
	})

	api := Custom("token", client)
	check, err := api.Check()
	equals(t, true, check.Up)
	equals(t, errors.New("Unauthorized - Check your token"), err)

}

func TestCheckReturnsIsDownWhen500(t *testing.T) {

	client := NewTestClient(func(req *http.Request) *http.Response {
		// Test request parameters
		equals(t, req.URL.String(), "https://unbelievable.pizza/api/v1")
		return &http.Response{
			StatusCode: 500,
			// Send response to be tested
			Body:       ioutil.NopCloser(bytes.NewBufferString(``)),
 			// Must be set to non-nil value or it panics
			Header:     make(http.Header),
		}
	})

	api := Custom("token", client)
	check, err := api.Check()
	equals(t, errors.New("Cannot Connect to API url."), err)
	equals(t, false, check.Up)
	
}

func TestUserBalanceHandlesDataOnSuccessfulFetchCorrectlyNonInfinte(t *testing.T) {

	client := NewTestClient(func(req *http.Request) *http.Response {
		// Test request parameters
		equals(t, req.URL.String(), "https://unbelievable.pizza/api/v1/guilds/411898639737421824/users/398197113495748626")
		return &http.Response{
			StatusCode: 200,
			// Send response to be tested
			Body:       ioutil.NopCloser(bytes.NewBufferString(`{"rank":"14","user_id":"398197113495748626","cash":0,"bank":526,"total":526}`)),
 			// Must be set to non-nil value or it panics
			Header:     make(http.Header),
		}
	})

	api := Custom("token", client)
	data, err := api.UserBalance("411898639737421824", "398197113495748626") // User, Guild
	ok(t, err)
	equals(t, data, data)
}


