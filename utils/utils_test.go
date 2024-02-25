package utils

import (
	"log"
	"testing"
)

func TestRandomID(t *testing.T) {
	tests := []struct {
		name   string
		length int
	}{
		{"Test 1", 10},
		{"Test 2", 20},
		{"Test 3", 30},
		{"Test 4", 40},
		{"Test 5", 50},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := RandomID(tt.length)
			log.Print("RandomID() = ", got)
			if len(got) != tt.length {
				t.Errorf("RandomID() = %v, want %v", len(got), tt.length)
			}
		})
	}
}
