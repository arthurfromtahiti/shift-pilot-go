package booking

import (
	"testing"
	"time"
)

func sample() Slot {
	return Slot{ID: 1, Activity: "Plongée", Start: time.Now(), Capacity: 10, Booked: 4}
}

func TestRemaining(t *testing.T) {
	if got := Remaining(sample()); got != 6 {
		t.Fatalf("Remaining = %d, attendu 6", got)
	}
}

func TestIsAvailable(t *testing.T) {
	if !IsAvailable(sample()) {
		t.Fatal("le créneau devrait être disponible")
	}
}

func TestBook(t *testing.T) {
	s := Book(sample(), 2)
	if s.Booked != 6 {
		t.Fatalf("Booked = %d, attendu 6", s.Booked)
	}
}
