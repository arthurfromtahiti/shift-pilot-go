package booking

import "time"

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
func Book(s Slot, n int) Slot {
	s.Booked += n
	return s
}
