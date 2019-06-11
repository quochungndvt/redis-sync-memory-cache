package rsmemory

import (
	"log"
	"testing"
)

func TestCalculateHashSlot(t *testing.T) {
	log.Println("TestCalculateHashSlot")
	hslt := NewHashSlotCalculator(0)
	s := hslt.CalculateHashSlot("test")
	log.Println("test__________", s)
	t.Run("message", func(t *testing.T) {

		log.Printf("recieved %q", string(s))
	})
}
