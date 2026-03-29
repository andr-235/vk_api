package encode

import (
	"net/url"
	"testing"
	"time"
)

type testParams struct {
	String      string    `url:"string,omitempty"`
	Int         int       `url:"int,omitempty"`
	Float       float64   `url:"float,omitempty"`
	Bool        bool      `url:"bool,omitempty"`
	StringSlice []string  `url:"strings,comma,omitempty"`
	Time        time.Time `url:"time,omitempty"`
}

func BenchmarkValues(b *testing.B) {
	params := testParams{
		String:      "test",
		Int:         123,
		Float:       3.14,
		Bool:        true,
		StringSlice: []string{"a", "b", "c"},
		Time:        time.Now(),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = Values(params)
	}
}

func BenchmarkValuesLarge(b *testing.B) {
	type largeParams struct {
		Field1  string   `url:"field1,omitempty"`
		Field2  int      `url:"field2,omitempty"`
		Field3  bool     `url:"field3,omitempty"`
		Field4  []string `url:"field4,comma,omitempty"`
		Field5  string   `url:"field5,omitempty"`
		Field6  int      `url:"field6,omitempty"`
		Field7  bool     `url:"field7,omitempty"`
		Field8  []string `url:"field8,comma,omitempty"`
		Field9  string   `url:"field9,omitempty"`
		Field10 int      `url:"field10,omitempty"`
	}

	params := largeParams{
		Field1:  "test",
		Field2:  123,
		Field3:  true,
		Field4:  []string{"a", "b", "c"},
		Field5:  "test",
		Field6:  123,
		Field7:  true,
		Field8:  []string{"a", "b", "c"},
		Field9:  "test",
		Field10: 123,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = Values(params)
	}
}

func BenchmarkEncodeMap(b *testing.B) {
	m := map[string]any{
		"string": "test",
		"int":    123,
		"float":  3.14,
		"bool":   true,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = Values(m)
	}
}

func BenchmarkCachePerformance(b *testing.B) {
	params := testParams{
		String:      "test",
		Int:         123,
		Float:       3.14,
		Bool:        true,
		StringSlice: []string{"a", "b", "c"},
	}

	// Первый вызов — кэш пустой
	_, _ = Values(params)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = Values(params)
	}
}

func TestValuesTypes(t *testing.T) {
	tests := []struct {
		name   string
		input  any
		expect url.Values
	}{
		{
			name: "string",
			input: struct {
				Field string `url:"field"`
			}{Field: "test"},
			expect: url.Values{"field": []string{"test"}},
		},
		{
			name: "int",
			input: struct {
				Field int `url:"field"`
			}{Field: 123},
			expect: url.Values{"field": []string{"123"}},
		},
		{
			name: "bool true",
			input: struct {
				Field bool `url:"field"`
			}{Field: true},
			expect: url.Values{"field": []string{"1"}},
		},
		{
			name: "bool false",
			input: struct {
				Field bool `url:"field"`
			}{Field: false},
			expect: url.Values{"field": []string{"0"}},
		},
		{
			name: "slice comma",
			input: struct {
				Field []string `url:"field,comma"`
			}{Field: []string{"a", "b", "c"}},
			expect: url.Values{"field": []string{"a,b,c"}},
		},
		{
			name: "omitempty empty string",
			input: struct {
				Field string `url:"field,omitempty"`
			}{},
			expect: url.Values{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Values(tt.input)
			if err != nil {
				t.Fatalf("Values() error = %v", err)
			}
			if !valuesEqual(got, tt.expect) {
				t.Errorf("Values() = %v, want %v", got, tt.expect)
			}
		})
	}
}

func valuesEqual(a, b url.Values) bool {
	if len(a) != len(b) {
		return false
	}
	for k, v := range a {
		if len(b[k]) != len(v) {
			return false
		}
		for i, vv := range v {
			if b[k][i] != vv {
				return false
			}
		}
	}
	return true
}
