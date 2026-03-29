package encode

import (
	"fmt"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

// fieldInfo хранит информацию о поле для кодирования
type fieldInfo struct {
	name      string
	omitEmpty bool
	comma     bool
	index     []int
}

// cacheEntry содержит запись кэша с временем создания
type cacheEntry struct {
	fields    []fieldInfo
	createdAt time.Time
}

// cache хранит закэшированную информацию о полях для каждого типа
// Ограничение: максимум 1000 типов, TTL 1 час
var (
	cache        sync.Map // map[reflect.Type]cacheEntry
	cacheSize    atomic.Int32
	maxCacheSize int32 = 1000
	cacheTTL           = time.Hour
)

// cleanCache очищает старые записи кэша
func cleanCache() {
	now := time.Now()
	cache.Range(func(key, value any) bool {
		if entry, ok := value.(cacheEntry); ok {
			if now.Sub(entry.createdAt) > cacheTTL {
				cache.Delete(key)
				cacheSize.Add(-1)
			}
		}
		return true
	})
}

// getCachedFields получает или создаёт кэш полей для типа
func getCachedFields(t reflect.Type) []fieldInfo {
	// Периодическая очистка (каждые 100 запросов)
	if cacheSize.Load()%100 == 0 {
		cleanCache()
	}

	if v, ok := cache.Load(t); ok {
		return v.(cacheEntry).fields
	}

	// Проверяем лимит перед добавлением
	if cacheSize.Load() >= maxCacheSize {
		cleanCache()
		// Если всё ещё переполнен, не кэшируем
		if cacheSize.Load() >= maxCacheSize {
			return buildFields(t)
		}
	}

	fields := buildFields(t)
	cache.Store(t, cacheEntry{
		fields:    fields,
		createdAt: time.Now(),
	})
	cacheSize.Add(1)
	return fields
}

// buildFields строит информацию о полях для типа
func buildFields(t reflect.Type) []fieldInfo {
	var fields []fieldInfo
	for i := 0; i < t.NumField(); i++ {
		sf := t.Field(i)
		if sf.PkgPath != "" {
			continue // пропущенные (unexported) поля
		}

		tag := sf.Tag.Get("url")
		if tag == "-" {
			continue
		}

		name, opts := parseTag(tag, sf.Name)
		fields = append(fields, fieldInfo{
			name:      name,
			omitEmpty: hasOpt(opts, "omitempty"),
			comma:     hasOpt(opts, "comma"),
			index:     sf.Index,
		})
	}
	return fields
}

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
	fields := getCachedFields(rv.Type())

	for _, f := range fields {
		fieldVal := rv.FieldByIndex(f.index)
		val, ok, err := stringify(fieldVal, fieldOpts(f.omitEmpty, f.comma))
		if err != nil {
			return nil, fmt.Errorf("encode field %q: %w", f.name, err)
		}
		if ok {
			values.Set(f.name, val)
		}
	}

	return values, nil
}

func fieldOpts(omitEmpty, comma bool) string {
	var opts string
	if omitEmpty {
		opts = "omitempty"
	}
	if comma {
		if opts != "" {
			opts += ",comma"
		} else {
			opts = "comma"
		}
	}
	return opts
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

	if hasOpt(opts, "omitempty") && v.IsZero() {
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

	case reflect.Struct:
		if t, ok := v.Interface().(time.Time); ok {
			if t.IsZero() && hasOpt(opts, "omitempty") {
				return "", false, nil
			}
			return strconv.FormatInt(t.Unix(), 10), true, nil
		}
		return "", false, fmt.Errorf("unsupported struct kind: %s", v.Type())

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
