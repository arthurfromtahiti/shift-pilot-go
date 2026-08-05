# Relecture — WORKFLOW_RESERVER_PLACES

## Verdict global

**Bon** — L'analyse est exacte, bien sourcée et honnête sur ses limites. La fonction `Book` est décrite conformément à son implémentation ; les règles métier (passage par valeur, `n` non validé, absence de garde de capacité) sont toutes vérifiées dans le code. Les risques sont concrets et accompagnés de scénarios précis. Aucune correction requise.

## Problèmes bloquants

Aucun.

## Problèmes mineurs

Aucun.

## Points vérifiés et corrects

- **Fichiers cités existants** : `internal/booking/booking.go`, `internal/booking/booking_test.go`, `go.mod` — tous lus et correspondant aux descriptions. (Preuve : lecture directe des trois fichiers.)
- **Point d'entrée `Book(s Slot, n int) Slot`** : à la ligne 25 de `booking.go` — exact. (Preuve : `booking.go:25`.)
- **`s.Booked += n` à la ligne 26** : exact. (Preuve : `booking.go:26`.)
- **`return s` à la ligne 27** : exact. (Preuve : `booking.go:27`.)
- **Passage par valeur (copie)** : `Book` reçoit `s Slot` par valeur, opère sur une copie locale, et retourne la copie modifiée — conformément à la signature Go `func Book(s Slot, n int) Slot`. L'appelant doit utiliser la valeur de retour. (Preuve : `booking.go:25-28`.)
- **Aucune garde de disponibilité dans `Book`** : l'incrément `s.Booked += n` est exécuté sans appel à `IsAvailable` ni `Remaining` — exact. (Preuve : corps de `Book`, `booking.go:25-28`.)
- **`n` peut être négatif** : paramètre `int` sans validation, comportement de décrément non documenté ni testé — correctement signalé. (Preuve : signature `n int` à `booking.go:25` ; `booking_test.go:24-29` ne teste que `n=2`.)
- **`TestBook` à lignes 24-29** : exact. (Preuve : `booking_test.go:24-29`, ligne 25 `s := Book(sample(), 2)`, ligne 26 `if s.Booked != 6`.)
- **Risque sur-réservation** : scénario concret `Book(Slot{Capacity:10, Booked:10}, 1)` → `Slot{Booked:11}` sans erreur — vérifiable dans le code. (Preuve : `booking.go:26`, aucune garde.)
- **Absence de persistance** : aucun appel réseau, base de données, ou système tiers dans le package — confirmé (`go.mod` sans dépendance externe).
- **Classification `technical_flow`** : justifiée pour les mêmes raisons que `WORKFLOW_CONSULTER_DISPONIBILITE` — aucun point d'entrée HTTP ni CLI. (Preuve : seuls fichiers Go présents.)
- **Niveau de confiance `medium`** : honnête — tests lus mais non exécutés (toolchain absente).

## Recommandations de correction

Aucune.
