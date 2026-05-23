package helpers

import (
	"encoding/json"
	"fmt"
	"reflect"
)

func Debug(label string, v any) {
	rv := reflect.ValueOf(v)
	rv = deref(rv)

	// can error invalid reflect.Value, so just print as is it
	if rv.Kind() == reflect.Invalid {
		fmt.Printf("%s : %#v\n", label, v)
		return
	}

	fmt.Printf("%s : %#v\n", label, rv.Interface())
}

func DebugJson(label string, v any) {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		fmt.Printf("%s : <json error: %v>\n", label, err)
		return
	}
	fmt.Printf("%s : %s\n", label, string(b))
}

func deref(v reflect.Value) reflect.Value {
	for v.Kind() == reflect.Ptr && !v.IsNil() {
		v = v.Elem()
	}
	return v
}
