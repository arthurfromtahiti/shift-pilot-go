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
	s, err := Book(sample(), 2)
	if err != nil {
		t.Fatalf("Book inattendu: %v", err)
	}
	if s.Booked != 6 {
		t.Fatalf("Booked = %d, attendu 6", s.Booked)
	}
}

func TestBookCapacityExceeded(t *testing.T) {
	_, err := Book(sample(), 7) // 6 places restantes seulement
	if err != ErrCapacityExceeded {
		t.Fatalf("attendu ErrCapacityExceeded, obtenu %v", err)
	}
}

func TestBookExactCapacity(t *testing.T) {
	s, err := Book(sample(), 6) // exactement 6 places restantes
	if err != nil {
		t.Fatalf("Book inattendu: %v", err)
	}
	if s.Booked != s.Capacity {
		t.Fatalf("Booked = %d, attendu %d (capacité pleine)", s.Booked, s.Capacity)
	}
}

func TestBookZero(t *testing.T) {
	_, err := Book(sample(), 0)
	if err != ErrInvalidBookingCount {
		t.Fatalf("attendu ErrInvalidBookingCount pour n=0, obtenu %v", err)
	}
}

func TestBookNegative(t *testing.T) {
	_, err := Book(sample(), -3)
	if err != ErrInvalidBookingCount {
		t.Fatalf("attendu ErrInvalidBookingCount pour n=-3, obtenu %v", err)
	}
}
