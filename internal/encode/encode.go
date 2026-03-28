package encode

import (
	"fmt"
	"net/url"
	"reflect"
	"strconv"
	"strings"
)

func Values(v any) (url.Values, error) {
	values := make(url.Values)
	if v == nil {
		return values, nil
	}

	rv := reflect.ValueOf(v)
	for rv.Kind() == reflect.Pointer {
		if rv.IsNil() {
			return values, nil
		}
		rv = rv.Elem()
	}

	switch rv.Kind() {
	case reflect.Map:
		return encodeMap(rv)
	case reflect.Struct:
		return encodeStruct(rv)
	default:
		return nil, fmt.Errorf("unsupported params type: %T", v)
	}
}

func encodeMap(rv reflect.Value) (url.Values, error) {
	if rv.Type().Key().Kind() != reflect.String {
		return nil, fmt.Errorf("map key must be string")
	}

	values := make(url.Values)
	iter := rv.MapRange()
	for iter.Next() {
		key := iter.Key().String()
		val, ok, err := stringify(iter.Value(), "")
		if err != nil {
			return nil, fmt.Errorf("encode key %q: %w", key, err)
		}
		if ok {
			values.Set(key, val)
		}
	}

	return values, nil
}

func encodeStruct(rv reflect.Value) (url.Values, error) {
	values := make(url.Values)
	rt := rv.Type()

	for i := 0; i < rv.NumField(); i++ {
		sf := rt.Field(i)
		if sf.PkgPath != "" {
			continue
		}

		tag := sf.Tag.Get("url")
		if tag == "-" {
			continue
		}

		name, opts := parseTag(tag, sf.Name)
		val, ok, err := stringify(rv.Field(i), opts)
		if err != nil {
			return nil, fmt.Errorf("encode field %q: %w", sf.Name, err)
		}
		if ok {
			values.Set(name, val)
		}
	}

	return values, nil
}

func parseTag(tag string, fallback string) (string, string) {
	if tag == "" {
		return strings.ToLower(fallback), ""
	}

	parts := strings.Split(tag, ",")
	name := parts[0]
	if name == "" {
		name = strings.ToLower(fallback)
	}

	if len(parts) == 1 {
		return name, ""
	}

	return name, strings.Join(parts[1:], ",")
}

func stringify(v reflect.Value, opts string) (string, bool, error) {
	for v.Kind() == reflect.Pointer {
		if v.IsNil() {
			return "", false, nil
		}
		v = v.Elem()
	}

	if hasOpt(opts, "omitempty") && isZero(v) {
		return "", false, nil
	}

	switch v.Kind() {
	case reflect.String:
		return v.String(), true, nil

	case reflect.Bool:
		if v.Bool() {
			return "1", true, nil
		}
		return "0", true, nil

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(v.Int(), 10), true, nil

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return strconv.FormatUint(v.Uint(), 10), true, nil

	case reflect.Float32:
		return strconv.FormatFloat(v.Float(), 'f', -1, 32), true, nil

	case reflect.Float64:
		return strconv.FormatFloat(v.Float(), 'f', -1, 64), true, nil

	case reflect.Slice, reflect.Array:
		if v.Len() == 0 && hasOpt(opts, "omitempty") {
			return "", false, nil
		}

		if hasOpt(opts, "comma") {
			parts := make([]string, 0, v.Len())
			for i := 0; i < v.Len(); i++ {
				part, ok, err := stringify(v.Index(i), "")
				if err != nil {
					return "", false, err
				}
				if ok {
					parts = append(parts, part)
				}
			}
			return strings.Join(parts, ","), true, nil
		}

		return "", false, fmt.Errorf("slice/array requires comma option")
	}

	return "", false, fmt.Errorf("unsupported kind: %s", v.Kind())
}

func hasOpt(opts string, target string) bool {
	if opts == "" {
		return false
	}
	for _, opt := range strings.Split(opts, ",") {
		if strings.TrimSpace(opt) == target {
			return true
		}
	}
	return false
}

func isZero(v reflect.Value) bool {
	return v.IsZero()
}
