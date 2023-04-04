package jsonfix

import (
	"encoding/json"
	"reflect"
	"testing"
)

func Test(t *testing.T) {
	t.Run("trailing comma in object", func(t *testing.T) {
		o(t, `{"foo": 1, "bar": 2,}`, `{"foo": 1, "bar": 2}`)
	})
	t.Run("trailing comma in multiline object", func(t *testing.T) {
		o(t, `{
			"foo": 1,
			"bar": 2,
		}`, `{"foo": 1, "bar": 2}`)
	})
	t.Run("trailing comma in multiline array", func(t *testing.T) {
		o(t, `[
			{"foo": 1},
			{"bar": 2},
		]`, `[{"foo": 1}, {"bar": 2}]`)
	})
	t.Run("comments in multiline object", func(t *testing.T) {
		o(t, `
		// grand object
		{
			// one
			"foo": 1, // fubar
			// two
			"bar": 2,
		}
		// trailing comments
		// are allowed too`, `{"foo": 1, "bar": 2}`)
	})
	t.Run("comments ignored inside strings", func(t *testing.T) {
		o(t, `
		{
			"foo": "http://example.com/", // real comment
			"bar": "http://example.com\\", // third comment
			"boz": "\"http://example.com/", // another comment
		}`, `{"foo": "http://example.com/", "bar": "http://example.com\\", "boz": "\"http://example.com/"}`)
	})
}

func o(t testing.TB, src string, expected string) {
	t.Helper()

	var expData any
	ensure(t, json.Unmarshal([]byte(expected), &expData))

	raw := Bytes([]byte(src))
	t.Logf("Fixed JSON = %s", raw)
	var actData any
	ensure(t, json.Unmarshal([]byte(raw), &actData))

	deepEqual(t, actData, expData)
}

func must[T any](t testing.TB, v T, err error) T {
	if err != nil {
		t.Helper()
		t.Fatalf("** %v", err)
	}
	return v
}

func ensure(t testing.TB, err error) {
	if err != nil {
		t.Helper()
		t.Fatalf("** %v", err)
	}
}

func eq[T comparable](t testing.TB, a, e T) {
	if a != e {
		t.Helper()
		t.Fatalf("** got %v, wanted %v", a, e)
	}
}

func deepEqual[T any](t testing.TB, a, e T) {
	if !reflect.DeepEqual(a, e) {
		t.Helper()
		t.Errorf("** got %v, wanted %v", a, e)
	}
}
