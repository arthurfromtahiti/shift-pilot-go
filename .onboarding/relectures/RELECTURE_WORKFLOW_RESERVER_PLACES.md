# Relecture — WORKFLOW_RESERVER_PLACES

## Verdict global

**Bon** — L'analyse est exacte, bien sourcée et honnête sur ses limites. La fonction `Book` est décrite conformément à son implémentation actuelle (`Book(s Slot, n int) (Slot, error)`, deux gardes explicites). Les règles métier, risques résiduels et questions ouvertes sont vérifiés dans le code. Aucune correction requise.

## Problèmes bloquants

Aucun.

## Problèmes mineurs

Aucun.

## Points vérifiés et corrects

- **Fichiers cités existants** : `internal/booking/booking.go`, `internal/booking/booking_test.go`, `go.mod` — tous lus et correspondant aux descriptions. (Preuve : lecture directe des trois fichiers.)
- **Point d'entrée `Book(s Slot, n int) (Slot, error)`** : à la ligne 36 de `booking.go` — exact. (Preuve : `booking.go:36`.)
- **Garde `n ≤ 0` → `ErrInvalidBookingCount`** : `booking.go:37-38`. Exact. (Preuve : `if n <= 0 { return s, ErrInvalidBookingCount }`.)
- **Garde `n > Remaining(s)` → `ErrCapacityExceeded`** : `booking.go:39-41`. Exact. (Preuve : `if n > Remaining(s) { return s, ErrCapacityExceeded }`.)
- **`s.Booked += n` à la ligne 43** : exact. (Preuve : `booking.go:43`.)
- **`return s, nil` à la ligne 44** : exact. (Preuve : `booking.go:44`.)
- **Passage par valeur (copie)** : `Book` reçoit `s Slot` par valeur, opère sur une copie locale, et retourne la copie modifiée avec `nil` en cas de succès — conformément à la signature `func Book(s Slot, n int) (Slot, error)`. L'appelant doit utiliser la valeur de retour. (Preuve : `booking.go:36-45`.)
- **Sur-réservation impossible via `Book`** : `Book(Slot{Capacity:10, Booked:10}, 1)` retourne `(Slot{}, ErrCapacityExceeded)` — correctement signalé. (Preuve : garde `booking.go:39-41`.)
- **`n` négatif ou nul rejeté** : `Book(s, -3)` et `Book(s, 0)` retournent `ErrInvalidBookingCount` — correctement signalé. (Preuve : garde `booking.go:37-38`.)
- **`n` négatif n'est plus une annulation implicite** : le workflow signale qu'une fonction `Cancel` dédiée sera nécessaire — exact. (Preuve : aucune fonction `Cancel` dans `booking.go:1-45`.)
- **Tests couvrant les cas critiques** : `TestBook`, `TestBookCapacityExceeded`, `TestBookExactCapacity`, `TestBookZero`, `TestBookNegative` — vérifiés à la lecture (`booking_test.go:24-63`).
- **Absence de persistance** : aucun appel réseau, base de données, ou système tiers dans le package — confirmé (`go.mod` sans dépendance externe).
- **Classification `technical_flow`** : justifiée — aucun point d'entrée HTTP ni CLI. (Preuve : seuls fichiers Go présents.)
- **Niveau de confiance `medium`** : honnête — tests lus mais non exécutés (toolchain absente).

## Recommandations de correction

Aucune.
