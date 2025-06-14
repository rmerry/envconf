/*
Package goenvconf provides functionality for populating structs with values
defined in environment variables.

The following basic types are supported:

  - string
  - bool
  - int
  - int8
  - int16
  - int32
  - int64
  - uint
  - uint8
  - uint16
  - uint32
  - uint64
  - float32
  - float64
  - complex64
  - complex128

Usage:

	type Config struct {
		AppName string  `env:"APP_NAME,required"`
		Port    int     `env:"PORT,default=8080"`
		Debug   bool    `env:"DEBUG"`
		Timeout float64 `env:"TIMEOUT,default=5.5"`
	}

	func main() {
		var cfg Config goenvconf.Process(&cfg)
		// ...
	}

Supported Tag Attributes:

  - default=VALUE - use VALUE when environment variable not set.

  - required - panic if environment variable not set.

    Note: If both `required` and `default` are
    provided the `required` tag is ignored.
*/
package goenvconf

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
)

const (
	tagKey = "env"

	tagAttrAssignmentSymbol = "="
	tagAttrDefault          = "default"
	tagAttrRequired         = "required"
)

// Makes unit testing easier.
var getEnvFunc func(string) string = os.Getenv

// Process populates the fields of a struct based on environment variables
// defined in struct tags.
//
// The input `v` must be a pointer to a struct. The function will panic if `v`
// is not a pointer or does not point to a struct type.
//
// All exported fields of the struct (and any exported nested structs) are
// processed. Unexported fields are ignored. The function recurses into nested
// structs, whether they are embedded by value or by pointer.
//
// Each relevant field should have a struct tag that specifies the environment
// variable name and optional attributes (such as whether the variable is
// required or a default value). The function retrieves the value from the
// environment, attempts to convert it to the field's type, and assigns it. The
// struct is modified in place.
//
// This function will panic under the following conditions: - A required
// environment variable is not set and no default value is provided. - A value
// retrieved from the environment cannot be converted to the field's type (e.g.,
// non-numeric string for an int).
func Process(v any) {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Pointer || rv.Elem().Kind() != reflect.Struct {
		panic("expected pointer to struct")
	}

	processFields(rv)
}

// processFields takes a `[]reflect.StructField` a `reflect.Value` and iterates
// through the `StructField` collection processing only the fields that are both
// exported and decorated with an appropriate tag (see `tagKey`).
//
// This function is recursive and will also iterate through all levels of struct
// nesting (struct embedding) so long as the structs are exported. Fields that
// are unexported or that do not contain a valid tag are skipped. This function
// will panic if a required environment variable is not supplied.
func processFields(v reflect.Value) {
	for _, field := range reflect.VisibleFields(v.Elem().Type()) {
		if !field.IsExported() {
			continue
		}
		// Recurse into structs and struct pointers.
		var (
			isStruct    = field.Type.Kind() == reflect.Struct
			isStructPtr = field.Type.Kind() == reflect.Pointer &&
				field.Type.Elem().Kind() == reflect.Struct
		)
		if isStruct || isStructPtr {
			fV := v.Elem().FieldByIndex(field.Index)

			if isStructPtr {
				if fV.IsNil() {
					fV.Set(reflect.New(fV.Type().Elem()))
				}
				fV = fV.Elem()
			}

			processFields(fV.Addr())
			continue
		}

		key, required, defaultVal := parseTag(field.Tag)
		if key == "" {
			continue // Ignore any field with no tag.
		}

		val := getEnvFunc(key)
		if val == "" && defaultVal != "" {
			val = defaultVal
		} else if val == "" && required {
			panic(fmt.Sprintf("env var %q not set", key))
		} else if val == "" {
			continue
		}

		var (
			err      error
			fieldPtr = v.Elem().FieldByIndex(field.Index)
		)
		switch field.Type.Kind() {
		case reflect.String:
			fieldPtr.SetString(val)

		case reflect.Int:
			var (
				bitSize int = strconv.IntSize
				i       int64
			)
			i, err = strconv.ParseInt(val, 10, bitSize)
			fieldPtr.SetInt(int64(i))
		case reflect.Uint:
			var (
				bitSize int = strconv.IntSize
				i       uint64
			)
			i, err = strconv.ParseUint(val, 10, bitSize)
			fieldPtr.SetUint(i)

		case reflect.Int8:
			var i int64
			i, err = strconv.ParseInt(val, 10, 8)
			fieldPtr.SetInt(i)
		case reflect.Int16:
			var i int64
			i, err = strconv.ParseInt(val, 10, 16)
			fieldPtr.SetInt(i)
		case reflect.Int32:
			var i int64
			i, err = strconv.ParseInt(val, 10, 32)
			fieldPtr.SetInt(i)
		case reflect.Int64:
			var i int64
			i, err = strconv.ParseInt(val, 10, 64)
			fieldPtr.SetInt(i)
		case reflect.Uint8:
			var i uint64
			i, err = strconv.ParseUint(val, 10, 8)
			fieldPtr.SetUint(i)
		case reflect.Uint16:
			var i uint64
			i, err = strconv.ParseUint(val, 10, 16)
			fieldPtr.SetUint(i)
		case reflect.Uint32:
			var i uint64
			i, err = strconv.ParseUint(val, 10, 32)
			fieldPtr.SetUint(i)
		case reflect.Uint64:
			var i uint64
			i, err = strconv.ParseUint(val, 10, 64)
			fieldPtr.SetUint(i)
		case reflect.Float32:
			var f float64
			f, err = strconv.ParseFloat(val, 32)
			fieldPtr.SetFloat(f)
		case reflect.Float64:
			var f float64
			f, err = strconv.ParseFloat(val, 64)
			fieldPtr.SetFloat(f)
		case reflect.Bool:
			var b bool
			b, err = strconv.ParseBool(val)
			fieldPtr.SetBool(b)
		case reflect.Complex64:
			var v complex128
			v, err = strconv.ParseComplex(val, 64)
			fieldPtr.SetComplex(v)
		case reflect.Complex128:
			var v complex128
			v, err = strconv.ParseComplex(val, 128)
			fieldPtr.SetComplex(v)
		}
		if err != nil {
			panic(fmt.Sprintf("invalid %s value supplied: %q",
				field.Type.Kind().String(), val))
		}
	}
}

// parseTag takes a `reflect.StructTag` and parses it for the presence of
// `tagKey`. The function returns a 3-tuple of (key, required, default value).
//
// If `tagKey` is not present `key` will be an empty string. If an invalid tag
// attribute is provided the function will panic.
func parseTag(st reflect.StructTag) (string, bool, string) {
	var (
		key        string
		required   bool
		defaultVal string
	)

	val := st.Get(tagKey)
	// Tag does not contain `tagKey`.
	if val == "" {
		return key, required, defaultVal
	}

	splits := strings.Split(val, ",")
	key = splits[0]

	// Only key is supplied in tag (i.e., no additional attributes).
	if len(splits) == 1 {
		return key, required, defaultVal
	}

	// Extract and process all tag attributes.
	for _, attr := range splits[1:] {
		if attr == tagAttrRequired {
			required = true
		} else if strings.HasPrefix(attr,
			tagAttrDefault+tagAttrAssignmentSymbol) {
			defaultVal = strings.TrimPrefix(attr,
				tagAttrDefault+tagAttrAssignmentSymbol)
		} else {
			panic(fmt.Sprintf("unrecognised struct tag attribute: %q", attr))
		}
	}

	return key, required, defaultVal
}
