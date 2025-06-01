package domain

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
)

type PageCursor struct {
	ID    string `json:"id"`
	Value string `json:"value"`
}

func NewPageCursor(id, value string) *PageCursor {
	return &PageCursor{
		ID:    id,
		Value: value,
	}
}

func (c PageCursor) Encode() (string, error) {
	jsonRaw, err := json.Marshal(c)
	if err != nil {
		return "", err
	}

	base64Raw := base64.StdEncoding.EncodeToString(jsonRaw)
	return base64Raw, nil
}

func (c *PageCursor) UnmarshalGQL(v interface{}) error {
	if str, ok := v.(string); ok {
		raw, err := base64.StdEncoding.DecodeString(str)
		if err != nil {
			return err
		}

		err = json.Unmarshal(raw, c)
		if err != nil {
			return err
		}
		return nil
	}
	return nil
}

func (c PageCursor) MarshalGQL(w io.Writer) {
	encoded, err := c.Encode()
	if err != nil {
		panic(err)
	}

	io.WriteString(w, fmt.Sprintf(`"%s"`, encoded))
}
