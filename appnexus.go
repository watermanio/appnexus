package AppNexus

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"

	"github.com/google/go-querystring/query"
)

// const defaultEndPoint = "https://api.appnexus.com"
// const defaultEndPoint = "http://hb.sand-08.adnxs.net/"

const defaultEndPoint = "http://sand.api.appnexus.com/"

// Credentials required to login
type credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Client used to make HTTP requests
type Client struct {
	client      *http.Client
	EndPoint    *url.URL
	Rate        Rate
	UserAgent   string
	token       string
	credentials credentials

	Members       *MemberService
	MemberSharing *MemberSharingService
	Segments      *SegmentService
}

// Rate contains information on the current rate limit in operation
type Rate struct {
	Reads             int     `json:"reads"`
	ReadLimit         int     `json:"read_limit"`
	ReadLimitSeconds  int     `json:"read_limit_seconds"`
	Writes            int     `json:"writes"`
	WriteLimit        int     `json:"write_limit"`
	WriteLimitSeconds int     `json:"write_limit_seconds"`
	Time              float64 `json:"time"`
}

// Response is a AppNexus API response object
type Response struct {
	*http.Response
	Obj struct {
		Status             string          `json:"status"`
		ID                 int             `json:"id,omitempty"`
		Token              string          `json:"token,omitempty"`
		Service            string          `json:"service,omitempty"`
		Method             string          `json:"method,omitempty"`
		Count              int             `json:"count,omitempty"`
		StartElement       int             `json:"start_element,omitempty"`
		NumElements        int             `json:"num_elements,omitempty"`
		MemberDataSharings []MemberSharing `json:"member_data_sharings,omitempty"`
		Member             Member          `json:"member,omitempty"`
		Segments           []Segment       `json:"segments,omitempty"`
		Rate               Rate            `json:"dbg_info"`
	} `json:"response"`
}

// ErrorResponse is a AppNexus API response to an internal error
type ErrorResponse struct {
	*http.Response
	Obj struct {
		Status           string `json:"status"`
		ErrorID          string `json:"error_id"`
		Error            string `json:"error"`
		ErrorDescription string `json:"error_description"`
		ErrorCode        string `json:"error_code"`
		Service          string `json:"service"`
		Rate             Rate   `json:"dbg_info"`
	} `json:"response"`
}

// ListOptions specifies the optional parameters to various List methods that
// support pagination.
type ListOptions struct {
	StartElement int  `url:"start_element,omitempty"`
	NumElements  int  `url:"num_elements,omitempty"`
	Active       bool `url:"active,omitempty"`
}

// NewClient returns a new AppNexus API client
// If a nil httpClient is provided, http.DefaultClient will be used.
func NewClient(httpClient *http.Client) *Client {

	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	baseURL, _ := url.Parse(defaultEndPoint)

	c := &Client{
		client:    httpClient,
		EndPoint:  baseURL,
		UserAgent: "github.com/adwww/appnexus go-appnexus-client",
	}

	c.Members = &MemberService{client: c}
	c.MemberSharing = &MemberSharingService{client: c}
	c.Segments = &SegmentService{client: c}

	return c
}

// NewRequest creates an API request using a relative URL
func (c *Client) newRequest(method, path string, body interface{}) (*http.Request, error) {
	rel, err := url.Parse(path)
	if err != nil {
		return nil, err
	}

	u := c.EndPoint.ResolveReference(rel)

	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)
		err := json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}

	req.Header.Add("User-Agent", c.UserAgent)

	if c.token != "" {
		req.Header.Add("Authorization", c.token)
	}

	return req, nil
}

// Do sends an API request and returns the API response.  The API response is
// JSON decoded and stored in the value pointed to by v, or returned as an
// error if an API error has occurred.  If v implements the io.Writer
// interface, the raw response body will be written to v, without attempting to
// first decode it.
func (c *Client) do(req *http.Request, v interface{}) (*Response, error) {
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, errors.New("client.do.do: " + err.Error())
	}

	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New("client.do.readall: " + err.Error())
	}

	err = checkResponse(resp, data)
	if err != nil {
		return nil, errors.New("client.do.checkResponse: " + err.Error())
	}

	response := &Response{Response: resp}
	c.Rate = response.Obj.Rate

	if v != nil {
		err := json.Unmarshal(data, v)
		if err != nil {
			return nil, errors.New("client.do.unmarshal: " + err.Error())
		}
	}

	return response, nil
}

// CheckResponse checks the API response for errors, and returns them if
// present.
func checkResponse(r *http.Response, data []byte) error {

	if r.StatusCode < 200 || r.StatusCode > 299 {
		return errors.New(r.Status)
	}

	if len(data) > 0 {

		resp := &ErrorResponse{Response: r}
		err := json.Unmarshal(data, resp)
		if err != nil {
			return err
		}

		if resp.Obj.ErrorID != "" || resp.Obj.Error != "" {
			str := fmt.Sprintf("AppNexus:checkResponse [%s]: %s", resp.Obj.ErrorID, resp.Obj.Error)
			return errors.New(str)
		}
	}

	return nil
}

// Login to the AppNexus API and get an authentication token
func (c *Client) Login(username string, password string) error {

	c.credentials = credentials{
		Username: username,
		Password: password,
	}

	auth := struct {
		credentials `json:"auth"`
	}{c.credentials}

	req, err := c.newRequest("POST", "auth", auth)
	if err != nil {
		return err
	}

	resp, err := c.do(req, nil)
	if err != nil {
		return err
	}

	c.token = resp.Cookies()[0].Value
	return nil
}

// addOptions adds the parameters in opt as URL query parameters to s.  opt
// must be a struct whose fields may contain "url" tags.
func addOptions(s string, opt interface{}) (string, error) {
	v := reflect.ValueOf(opt)
	if v.Kind() == reflect.Ptr && v.IsNil() {
		return s, nil
	}

	u, err := url.Parse(s)
	if err != nil {
		return s, err
	}

	qs, err := query.Values(opt)
	if err != nil {
		return s, err
	}

	u.RawQuery = qs.Encode()
	return u.String(), nil
}
