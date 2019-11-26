package httputil

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"shanhu.io/misc/errcode"
)

func isSuccess(resp *http.Response) bool {
	return resp.StatusCode/100 == 2
}

// AddErrCode adds error code to an error given the http status.
func AddErrCode(statusCode int, err error) error {
	switch statusCode {
	case http.StatusNotFound:
		err = errcode.Add(errcode.NotFound, err)
	case http.StatusUnauthorized, http.StatusForbidden:
		err = errcode.Add(errcode.Unauthorized, err)
	case http.StatusBadRequest:
		err = errcode.Add(errcode.InvalidArg, err)
	}
	return err
}

// RespError returns the error from an HTTP response.
func RespError(resp *http.Response) error {
	bs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	err = fmt.Errorf("%s - %s",
		resp.Status, strings.TrimSpace(string(bs)),
	)

	return AddErrCode(resp.StatusCode, err)
}
