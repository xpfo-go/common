package config

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
	"unicode"
)

func ApplyDefaults(cfg any) error {
	return walkStruct(cfg, func(field reflect.Value, sf reflect.StructField, envKey string) error {
		def := sf.Tag.Get("default")
		if def == "" || !field.IsZero() {
			return nil
		}
		if err := setValueFromString(field, def); err != nil {
			return fmt.Errorf("apply default for %s: %w", sf.Name, err)
		}
		return nil
	})
}

func LoadFromEnv(prefix string, cfg any) error {
	return walkStruct(cfg, func(field reflect.Value, sf reflect.StructField, envKey string) error {
		fullKey := strings.ToUpper(envKey)
		if prefix != "" {
			fullKey = strings.ToUpper(prefix) + "_" + fullKey
		}
		raw, ok := os.LookupEnv(fullKey)
		if !ok {
			return nil
		}
		if err := setValueFromString(field, raw); err != nil {
			return fmt.Errorf("load env %s for %s: %w", fullKey, sf.Name, err)
		}
		return nil
	})
}

func ValidateRequired(cfg any) error {
	return walkStruct(cfg, func(field reflect.Value, sf reflect.StructField, envKey string) error {
		if sf.Tag.Get("required") != "true" {
			return nil
		}
		if field.IsZero() {
			return fmt.Errorf("required field missing: %s", envKey)
		}
		return nil
	})
}

func walkStruct(cfg any, fn func(field reflect.Value, sf reflect.StructField, envKey string) error) error {
	v := reflect.ValueOf(cfg)
	if v.Kind() != reflect.Ptr || v.IsNil() {
		return fmt.Errorf("config must be a non-nil pointer")
	}
	if v.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("config must point to a struct")
	}
	return walk(v.Elem(), nil, fn)
}

func walk(v reflect.Value, path []string, fn func(field reflect.Value, sf reflect.StructField, envKey string) error) error {
	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		sf := t.Field(i)
		if sf.PkgPath != "" {
			continue
		}
		key := sf.Tag.Get("config")
		if key == "-" {
			continue
		}
		if key == "" {
			key = toSnakeCase(sf.Name)
		}
		field := v.Field(i)

		if sf.Type.Kind() == reflect.Ptr && sf.Type.Elem().Kind() == reflect.Struct && sf.Type.Elem() != reflect.TypeOf(time.Duration(0)) {
			if field.IsNil() {
				field.Set(reflect.New(sf.Type.Elem()))
			}
			if err := walk(field.Elem(), append(path, key), fn); err != nil {
				return err
			}
			continue
		}

		if sf.Type.Kind() == reflect.Struct && sf.Type != reflect.TypeOf(time.Duration(0)) {
			if err := walk(field, append(path, key), fn); err != nil {
				return err
			}
			continue
		}

		envKey := strings.Join(append(path, key), "_")
		if err := fn(field, sf, envKey); err != nil {
			return err
		}
	}
	return nil
}

func setValueFromString(field reflect.Value, raw string) error {
	if field.Kind() == reflect.Ptr {
		if field.IsNil() {
			field.Set(reflect.New(field.Type().Elem()))
		}
		return setValueFromString(field.Elem(), raw)
	}

	if field.Type() == reflect.TypeOf(time.Duration(0)) {
		d, err := time.ParseDuration(strings.TrimSpace(raw))
		if err != nil {
			return err
		}
		field.SetInt(int64(d))
		return nil
	}

	raw = strings.TrimSpace(raw)
	switch field.Kind() {
	case reflect.String:
		field.SetString(raw)
	case reflect.Bool:
		b, err := strconv.ParseBool(raw)
		if err != nil {
			return err
		}
		field.SetBool(b)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		n, err := strconv.ParseInt(raw, 10, field.Type().Bits())
		if err != nil {
			return err
		}
		field.SetInt(n)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		n, err := strconv.ParseUint(raw, 10, field.Type().Bits())
		if err != nil {
			return err
		}
		field.SetUint(n)
	case reflect.Float32, reflect.Float64:
		f, err := strconv.ParseFloat(raw, field.Type().Bits())
		if err != nil {
			return err
		}
		field.SetFloat(f)
	case reflect.Slice:
		if field.Type().Elem().Kind() != reflect.String {
			return fmt.Errorf("unsupported slice type: %s", field.Type().String())
		}
		if raw == "" {
			field.Set(reflect.MakeSlice(field.Type(), 0, 0))
			return nil
		}
		parts := strings.Split(raw, ",")
		vals := make([]string, 0, len(parts))
		for _, p := range parts {
			vals = append(vals, strings.TrimSpace(p))
		}
		field.Set(reflect.ValueOf(vals))
	default:
		return fmt.Errorf("unsupported field type: %s", field.Type().String())
	}
	return nil
}

func toSnakeCase(in string) string {
	if in == "" {
		return ""
	}
	runes := []rune(in)
	var b strings.Builder
	for i, r := range runes {
		if unicode.IsUpper(r) {
			if i > 0 && (unicode.IsLower(runes[i-1]) ||
				(unicode.IsUpper(runes[i-1]) && i+1 < len(runes) && unicode.IsLower(runes[i+1]))) {
				b.WriteByte('_')
			}
			b.WriteRune(unicode.ToLower(r))
			continue
		}
		b.WriteRune(r)
	}
	return b.String()
}
