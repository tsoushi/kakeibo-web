package domain

import (
	"io"
	"strconv"

	"github.com/99designs/gqlgen/graphql"
	"golang.org/x/xerrors"
)

func MarshalScalarUint(t uint) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		io.WriteString(w, strconv.FormatUint(uint64(t), 10))
	})
}

func UnmarshalScalarUint(v interface{}) (uint, error) {
	switch v := v.(type) {
	case string:
		i, err := strconv.ParseUint(v, 10, 64)
		if err != nil {
			return 0, err
		}
		return uint(i), nil
	case int:
		return uint(v), nil
	case int32:
		return uint(v), nil
	case int64:
		return uint(v), nil
	case float64:
		return uint(v), nil
	default:
		return 0, xerrors.Errorf("wrong type %T", v)
	}
}
