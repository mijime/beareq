package attrs

import (
	"bytes"
	"encoding/json"
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
	m, err := gojo.Map(x.Args)
	if err != nil {
		return nil, fmt.Errorf("failed to create buffer use gojo: %w", err)
	}

	buf := bytes.NewBuffer([]byte{})
	if err := json.NewEncoder(buf).Encode(m); err != nil {
		return nil, fmt.Errorf("failed to create buffer use gojo: %w", err)
	}

	return buf, nil
}
