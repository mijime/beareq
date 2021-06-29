package attrs

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/itchyny/gojo"
)

type JSONObject struct {
	Args []string
}

func NewJSONObject() JSONObject {
	return JSONObject{Args: make([]string, 0)}
}

func (x *JSONObject) String() string {
	if x == nil || x.Args == nil {
		return ""
	}

	return strings.Join(x.Args, " ")
}

func (x *JSONObject) Set(v string) error {
	if len(v) == 0 {
		return errors.New("require value")
	}

	x.Args = append(x.Args, v)

	return nil
}

func (x JSONObject) Exists() bool {
	return len(x.Args) > 0
}

func (x JSONObject) Parse() (io.Reader, error) {
	buf := bytes.NewBuffer([]byte{})

	if err := gojo.New(gojo.Output(buf), gojo.Args(x.Args)).Run(); err != nil {
		return nil, fmt.Errorf("failed to create buffer use gojo: %w", err)
	}

	return buf, nil
}
