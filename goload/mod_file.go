package goload

import (
	"bufio"
	"bytes"
	"errors"
	"io/ioutil"
	"strconv"
	"strings"
)

type modFile struct {
	name string
}

var (
	errInvalidModFile = errors.New("invalid go.mod file")
)

func parseModFile(f string) (*modFile, error) {
	bs, err := ioutil.ReadFile(f)
	if err != nil {
		return nil, err
	}
	name, err := modulePath(bs)
	if err != nil {
		return nil, err
	}
	return &modFile{name: name}, nil
}

// modulePath returns the module path from the gomod file text.
// If it cannot find a module path, it returns an empty string.
// It is tolerant of unrelated problems in the go.mod file.
func modulePath(bs []byte) (string, error) {
	s := bufio.NewScanner(bytes.NewReader(bs))

	for s.Scan() {
		line := strings.TrimSpace(s.Text())
		if line == "" || strings.HasPrefix(line, "//") {
			continue
		}
		if !strings.HasPrefix(line, "module") {
			continue
		}

		line = strings.TrimSpace(strings.TrimPrefix(line, "module"))

		// TODO: this is incorrect for quoted module path
		if pos := strings.Index(line, "//"); pos >= 0 {
			line = strings.TrimSpace(line[:pos])
			if line == "" {
				return "", errInvalidModFile
			}
		}

		if line == "" {
			return "", errInvalidModFile
		}
		if line[0] == '"' || line[0] == '`' {
			p, err := strconv.Unquote(line)
			if err != nil || p == "" {
				return "", errInvalidModFile
			}
			return p, nil
		}

		return line, nil
	}

	return "", errInvalidModFile
}

func isValidModPath(p, modPath string) bool {
	if modPath == p {
		return false
	}

	prefix := p + "/v"
	if !strings.HasPrefix(modPath, prefix) {
		return false
	}

	ver := strings.TrimPrefix(modPath, prefix)
	if _, err := strconv.Atoi(ver); err != nil {
		return false
	}

	return true
}
