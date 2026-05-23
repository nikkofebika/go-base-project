package helpers

import "encoding/json"

func AddToMapNullableValidator[T any](m map[string]any, key string, val Nullable[T]) {
	if val.Set {
		m[key] = val.Value
	}
}

func AddToMapIfNotNil[T any](m map[string]any, key string, value *T) {
	if value != nil {
		m[key] = *value
	}
}

func ToPointer[T any](v T) *T {
	var i any = v
	switch val := i.(type) {
	case string:
		if val == "" || val == "0001-01-01 00:00:00" {
			return nil
		}
	case uint:
		if val == 0 {
			return nil
		}
	}

	return &v
}

type Nullable[T any] struct {
	Set   bool
	Value *T
}

func (n *Nullable[T]) UnmarshalJSON(b []byte) error {
	n.Set = true

	if string(b) == "null" {
		n.Value = nil
		return nil
	}

	var v T
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}

	n.Value = &v
	return nil
}
