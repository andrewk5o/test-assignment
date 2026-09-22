package main

import (
	"reflect"
	"testing"
)

func TestCitiesStartingWith(t *testing.T) {
	cities := []string{"Shanghai", "Cairo", "Chongqing", "Rio de Janeiro", "Chengdu", "São Paulo", "chicago"}

	tests := []struct {
		name   string
		letter string
		want   []string
	}{
		{"case-insensitive, preserves input order", "c", []string{"Cairo", "Chongqing", "Chengdu", "chicago"}},
		{"upper-case input", "C", []string{"Cairo", "Chongqing", "Chengdu", "chicago"}},
		{"single match", "r", []string{"Rio de Janeiro"}},
		{"unicode name", "s", []string{"Shanghai", "São Paulo"}},
		{"no match", "z", []string{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CitiesStartingWith(cities, tt.letter)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}
