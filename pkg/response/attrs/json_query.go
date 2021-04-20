package attrs

import (
	"fmt"

	"github.com/itchyny/gojq"
)

type JSONQuery struct {
	*gojq.Query
}

func NewJSONQuery() JSONQuery {
	return JSONQuery{Query: nil}
}

func (q *JSONQuery) Exists() bool {
	return q != nil && q.Query != nil
}

func (q *JSONQuery) String() string {
	if q == nil || q.Query == nil {
		return ""
	}

	return q.Query.String()
}

func (q *JSONQuery) Set(v string) error {
	res, err := gojq.Parse(v)
	if err != nil {
		return fmt.Errorf("failed to parse jq: %w", err)
	}

	q.Query = res

	return nil
}
