package booking

import (
	"errors"
	"time"
)

// ErrCapacityExceeded est retourné quand une réservation dépasserait la capacité du créneau.
var ErrCapacityExceeded = errors.New("capacité du créneau dépassée")

// ErrInvalidBookingCount est retourné quand n est nul ou négatif.
var ErrInvalidBookingCount = errors.New("booking: n must be positive")

// Slot est un créneau d'activité réservable.
type Slot struct {
	ID       int
	Activity string
	Start    time.Time
	Capacity int
	Booked   int
}

// Remaining retourne le nombre de places encore disponibles sur un créneau.
func Remaining(s Slot) int {
	return s.Capacity - s.Booked
}

// IsAvailable indique si un créneau peut encore accepter une réservation.
func IsAvailable(s Slot) bool {
	return Remaining(s) > 0
}

// Book réserve n places sur un créneau et retourne le créneau mis à jour.
// Retourne ErrCapacityExceeded si n dépasse les places disponibles.
func Book(s Slot, n int) (Slot, error) {
	if n > Remaining(s) {
		return s, ErrCapacityExceeded
	}
	s.Booked += n
	return s, nil
}
