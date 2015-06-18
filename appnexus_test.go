package appnexus

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"time"
)

var (
	mux    *http.ServeMux
	client *Client
	server *httptest.Server
)

// setup sets up a test HTTP server along with a AppNexus.Client that is
// configured to talk to that test server.  Tests should register handlers on
// mux which provide mock responses for the API method being tested.
func setup() {
	// test server
	mux = http.NewServeMux()
	server = httptest.NewServer(mux)

	// appnexus client configured to use test server
	client, _ = NewClient(server.URL)
}

// teardown closes the test HTTP server.
func teardown() {
	server.Close()
}

func TestNewClient(t *testing.T) {
	c, _ := NewClient("http://sand.api.appnexus.com/")

	if actual, expected := c.EndPoint.String(), "http://sand.api.appnexus.com/"; actual != expected {
		t.Errorf("NewClient EndPoint is %v, expected %v", actual, expected)
	}

	if actual, expected := c.UserAgent, "github.com/adwww/appnexus go-appnexus-client"; actual != expected {
		t.Errorf("NewClient agent is %v, expected %v", actual, expected)
	}
}

func TestNewRequest(t *testing.T) {
	c, _ := NewClient("http://sand.api.appnexus.com/")

	inURL, outURL := "/foo", "http://sand.api.appnexus.com/foo"
	inBody, outBody := &User{FirstName: "Andy"}, `{"first_name":"Andy"}`+"\n"
	req, _ := c.newRequest("GET", inURL, inBody)

	// test that relative URL was expanded
	if actual, expected := req.URL.String(), outURL; actual != expected {
		t.Errorf("NewRequest(%q) URL is %v, expected %v", inURL, actual, expected)
	}

	// test that body was JSON encoded
	body, _ := ioutil.ReadAll(req.Body)
	if actual, expected := string(body), outBody; actual != expected {
		t.Errorf("NewRequest(%q) Body is %v, expected %v", inBody, actual, expected)
	}
}

func TestCheckResponse(t *testing.T) {

	c, _ := NewClient("http://sand.api.appnexus.com/")
	data := strings.NewReader(`{"response":{"error_id":"SYNTAX","error":"invalid service","dbg_info":{"output_term":"not_found"}}}`)

	buf := new(bytes.Buffer)
	buf.ReadFrom(data)

	res := &http.Response{
		Request:    &http.Request{},
		StatusCode: http.StatusOK,
		Body:       ioutil.NopCloser(data),
	}

	_, err := c.checkResponse(res, buf.Bytes())
	if err == nil {
		t.Errorf("Expected error response")
	}

	expected := errors.New("AppNexus:checkResponse [SYNTAX]: invalid service")
	if !reflect.DeepEqual(err, expected) {
		t.Errorf("Error = %#v, expected %#v", err, expected)
	}
}

func TestWaitForRateLimit(t *testing.T) {

	c, _ := NewClient("http://sand.api.appnexus.com/")
	c.Rate.Time = time.Now()
	c.Rate.ReadLimit = 100
	c.Rate.ReadLimitSeconds = 2
	c.Rate.Reads = 99
	c.Rate.WriteLimit = 100
	c.Rate.WriteLimitSeconds = 2
	c.Rate.Writes = 0

	if actual, expected := fmt.Sprintf("%.0f", c.waitForRateLimit("GET").Seconds()), "0"; actual != expected {
		t.Errorf("Waited %v for read rate limit, expected %v", actual, expected)
	}

	c.Rate.Reads = 100
	if actual, expected := fmt.Sprintf("%.0f", c.waitForRateLimit("GET").Seconds()), "2"; actual != expected {
		t.Errorf("Waited %v for read rate limit, expected %v", actual, expected)
	}

	if actual, expected := fmt.Sprintf("%.0f", c.waitForRateLimit("POST").Seconds()), "0"; actual != expected {
		t.Errorf("Waited %v for write rate limit, expected %v", actual, expected)
	}

	if actual, expected := fmt.Sprintf("%.0f", c.waitForRateLimit("PUT").Seconds()), "0"; actual != expected {
		t.Errorf("Waited %v for write rate limit, expected %v", actual, expected)
	}

	if actual, expected := fmt.Sprintf("%.0f", c.waitForRateLimit("DELETE").Seconds()), "0"; actual != expected {
		t.Errorf("Waited %v for write rate limit, expected %v", actual, expected)
	}
}
