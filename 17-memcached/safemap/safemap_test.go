package safemap_test

import (
	"reflect"
	"mc/safemap"
	"testing"
)

var (
	safeMap = safemap.New()
)

func TestNew(t *testing.T) {
	tests := []struct {
		name string
		want *safemap.SafeMap
	}{
		{"int int", safeMap},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := safemap.New(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSafeMap_Insert(t *testing.T) {
	tests := []struct {
		name      string
		s         *safemap.SafeMap
		key       string
		value     []byte
		wantKey   string
		wantValue []byte
	}{
		{"insert", safeMap, "k1", []byte("v1"), "k1", []byte("v1")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.s.Set(tt.key, tt.value)
			got := tt.s.Get(tt.wantKey)
			if string(got) != string(tt.wantValue) {
				t.Errorf("SafeMap.Insert() insert key %v value %v, want value %v", tt.key, tt.value, tt.wantValue)
			}
		})
	}
}
