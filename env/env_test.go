package env

import (
	"os"
	"strings"
	"testing"
	"time"
)

func TestDuration(t *testing.T) { //nolint:paralleltest
	var tests = []struct {
		value time.Duration
	}{
		{value: 1 * time.Second},
		{value: 1 * time.Minute},
		{value: 1 * time.Hour},
	}

	for _, tt := range tests { //nolint:paralleltest
		t.Run(tt.value.String(), func(t *testing.T) {
			key := strings.ToUpper(tt.value.String())
			if err := os.Setenv(key, tt.value.String()); err != nil {
				t.Fatalf("failed to set env var %s for test: %s\n", key, err)
			}

			def := 10 * time.Minute
			if have, want := Duration(key, def), tt.value; have != want {
				t.Errorf("have %s, want %s", have, want)
			}
		})
	}

	// test default value
	def := 10 * time.Minute
	if have, want := Duration("TEST_DEFAULT", def), def; have != want {
		t.Errorf("have %s, want %s", have, want)
	}
}

func TestString(t *testing.T) { //nolint:paralleltest
	var tests = []struct {
		value string
	}{
		{value: "foo"},
		{value: "bar"},
		{value: "baz"},
	}

	for _, tt := range tests { //nolint:paralleltest
		t.Run(tt.value, func(t *testing.T) {
			key := strings.ToUpper(tt.value)
			if err := os.Setenv(key, tt.value); err != nil {
				t.Fatalf("failed to set env var %s for test: %s\n", key, err)
			}

			def := "default_value"
			if have, want := String(key, def), tt.value; have != want {
				t.Errorf("have %s, want %s", have, want)
			}
		})
	}

	// test default value
	def := "default_value"
	if have, want := String("TEST_DEFAULT", def), def; have != want {
		t.Errorf("have %s, want %s", have, want)
	}
}

func TestBool(t *testing.T) { //nolint:paralleltest
	var tests = []struct {
		env   string
		value bool
	}{
		{env: "TRUE", value: true},
		{env: "true", value: true},
		{env: "1", value: true},
		{env: "F", value: false},
		{env: "FALSE", value: false},
		{env: "false", value: false},
		{env: "0", value: false},
	}

	for _, tt := range tests { //nolint:paralleltest
		t.Run(tt.env, func(t *testing.T) {
			key := "TEST_BOOL"
			if err := os.Setenv(key, tt.env); err != nil {
				t.Fatalf("failed to set env var %s for test: %s\n", key, err)
			}

			def := false
			if have, want := Bool(key, def), tt.value; have != want {
				t.Errorf("have %v, want %v", have, want)
			}
			def = true
			if have, want := Bool(key, def), tt.value; have != want {
				t.Errorf("have %v, want %v", have, want)
			}
		})
	}

	// test default value
	def := true
	if have, want := Bool("TEST_DEFAULT", def), def; have != want {
		t.Errorf("have %v, want %v", have, want)
	}
}

func TestInt(t *testing.T) { //nolint:paralleltest
	var tests = []struct {
		env   string
		value int
	}{
		{env: "1337", value: 1337},
		{env: "1", value: 1},
		{env: "-34", value: -34},
		{env: "0", value: 0},
	}

	for _, tt := range tests { //nolint:paralleltest
		t.Run(tt.env, func(t *testing.T) {
			key := "TEST_INT"
			if err := os.Setenv(key, tt.env); err != nil {
				t.Fatalf("failed to set env var %s for test: %s\n", key, err)
			}

			if have, want := Int(key, 10), tt.value; have != want {
				t.Errorf("have %v, want %v", have, want)
			}
		})
	}

	// test default value
	def := 11
	if have, want := Int("TEST_DEFAULT", def), def; have != want {
		t.Errorf("have %v, want %v", have, want)
	}
}
