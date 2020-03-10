package httputil

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
)

// Client performs client that calls to a remote server with an optional token.
type Client struct {
	Server *url.URL
	Token  string // Optional token to be put in the Bearer HTTP header.

	UserAgent string // Optional User-Agent for each request.
	Accept    string // Optional Accept header.

	Transport http.RoundTripper
}

func (c *Client) addHeaders(h http.Header) {
	headerSetAuthToken(h, c.Token)
	setHeader(h, "User-Agent", c.UserAgent)
	setHeader(h, "Accept", c.Accept)
}

func (c *Client) makeClient() *http.Client {
	return &http.Client{Transport: c.Transport}
}

func (c *Client) do(req *http.Request) (*http.Response, error) {
	return c.makeClient().Do(req)
}

func (c *Client) req(m, p string, r io.Reader) (*http.Request, error) {
	u, err := makeURL(c.Server, p)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(m, u, r)
	if err != nil {
		return nil, err
	}

	c.addHeaders(req.Header)
	return req, nil
}

func (c *Client) reqJSON(m, p string, r io.Reader) (*http.Request, error) {
	req, err := c.req(m, p, r)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	return req, nil
}

// Put puts a stream to a path on the server.
func (c *Client) Put(p string, r io.Reader) error {
	req, err := c.req("PUT", p, r)
	if err != nil {
		return err
	}

	resp, err := c.do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if !isSuccess(resp) {
		return RespError(resp)
	}
	return resp.Body.Close()
}

// PutBytes puts bytes to a path on the server.
func (c *Client) PutBytes(p string, bs []byte) error {
	return c.Put(p, bytes.NewBuffer(bs))
}

// JSONPut puts an object in JSON encoding.
func (c *Client) JSONPut(p string, v interface{}) error {
	bs, err := json.Marshal(v)
	if err != nil {
		return err
	}
	return c.PutBytes(p, bs)
}

func (c *Client) poke(m, p string) error {
	req, err := c.req(m, p, nil)
	if err != nil {
		return err
	}

	resp, err := c.do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if !isSuccess(resp) {
		return RespError(resp)
	}
	return nil
}

// GetCode gets a response from a route and returns the
// status code.
func (c *Client) GetCode(p string) (int, error) {
	req, err := c.req("GET", p, nil)
	if err != nil {
		return 0, err
	}
	resp, err := c.do(req)
	if err != nil {
		return 0, err
	}
	code := resp.StatusCode
	resp.Body.Close()
	return code, nil
}

// Poke posts a signal to the given route on the server.
func (c *Client) Poke(p string) error {
	return c.poke("POST", p)
}

// Get gets a response from a route on the server.
func (c *Client) Get(p string) (*http.Response, error) {
	req, err := c.req("GET", p, nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.do(req)
	if err != nil {
		return nil, err
	}
	if !isSuccess(resp) {
		defer resp.Body.Close()
		return nil, RespError(resp)
	}
	return resp, nil
}

// GetString gets the string response from a route on the server.
func (c *Client) GetString(p string) (string, error) {
	resp, err := c.Get(p)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	return respString(resp)
}

// GetBytes gets the byte array from a route on the server.
func (c *Client) GetBytes(p string) ([]byte, error) {
	resp, err := c.Get(p)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

// JSONGet gets the content of a path and decodes the response
// into resp as JSON.
func (c *Client) JSONGet(p string, resp interface{}) error {
	httpResp, err := c.Get(p)
	if err != nil {
		return err
	}
	defer httpResp.Body.Close()

	dec := json.NewDecoder(httpResp.Body)
	if err := dec.Decode(resp); err != nil {
		return err
	}
	return httpResp.Body.Close()
}

func copyRespBody(resp *http.Response, w io.Writer) error {
	defer resp.Body.Close()
	if !isSuccess(resp) {
		return RespError(resp)
	}
	if w == nil {
		return nil
	}
	if _, err := io.Copy(w, resp.Body); err != nil {
		return err
	}
	return resp.Body.Close()
}

// Post posts with request body from r, and copies the response body
// to w.
func (c *Client) Post(p string, r io.Reader, w io.Writer) error {
	req, err := c.req("POST", p, ioutil.NopCloser(r))
	if err != nil {
		return err
	}
	resp, err := c.do(req)
	if err != nil {
		return err
	}
	return copyRespBody(resp, w)
}

func (c *Client) jsonPost(p string, req interface{}) (*http.Response, error) {
	bs, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	httpReq, err := c.reqJSON("POST", p, bytes.NewBuffer(bs))
	if err != nil {
		return nil, err
	}
	return c.do(httpReq)
}

// JSONPost posts a JSON object as the request body and writes the body
// into the given writer.
func (c *Client) JSONPost(p string, req interface{}, w io.Writer) error {
	resp, err := c.jsonPost(p, req)
	if err != nil {
		return err
	}
	return copyRespBody(resp, w)
}

// JSONCall performs a call with the request as a marshalled JSON object,
// and the response unmarhsalled as a JSON object.
func (c *Client) JSONCall(p string, req, resp interface{}) error {
	httpResp, err := c.jsonPost(p, req)
	if err != nil {
		return err
	}
	defer httpResp.Body.Close()

	if !isSuccess(httpResp) {
		return RespError(httpResp)
	}
	if resp == nil {
		return nil
	}
	dec := json.NewDecoder(httpResp.Body)
	if err := dec.Decode(resp); err != nil {
		return err
	}
	return httpResp.Body.Close()
}

// Delete sends a delete message to the particular path.
func (c *Client) Delete(p string) error {
	return c.poke("DELETE", p)
}
