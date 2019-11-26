package httputil

import (
	"testing"

	"io/ioutil"
)

func TestClientGetCode(t *testing.T) {
	s := newHelloServer()
	c := NewClient(s.URL)
	got, err := c.GetCode("/")
	if err != nil {
		t.Fatal(err)
	}
	if got != 200 {
		t.Errorf("want 200, got %d", got)
	}

	got, err = c.GetCode("/secret")
	if err != nil {
		t.Fatal(err)
	}
	if got != 403 {
		t.Errorf("want 403, got %d", got)
	}
}

func TestClientGet(t *testing.T) {
	s := newHelloServer()
	c := NewClient(s.URL)
	resp, err := c.Get("/")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	got, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != testHelloMessage {
		t.Errorf("got %q, want %q", string(got), testHelloMessage)
	}
}
