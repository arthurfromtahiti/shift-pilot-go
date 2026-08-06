# Points chauds du code — Audit

> Confiance : high — le corpus est de 28 lignes de code de production et 29 lignes de tests, intégralement lus. Il n'y a pas de point chaud au sens classique (fichier volumineux, couplage fort, chaîne d'appels profonde) ; les risques sont qualitatifs, pas volumétriques. Le niveau de confiance `high` reflète l'exhaustivité de la lecture, pas la gravité des constats.

## Compréhension globale

Le dépôt ne contient qu'un seul fichier de production : `internal/booking/booking.go` (45 lignes). Il n'y a pas de point chaud par volume ou complexité cyclomatique — chaque fonction tient en quelques lignes. `Book` est la seule fonction qui modifie l'état ; elle valide désormais ses paramètres et retourne `(Slot, error)`. La suite de tests couvre 7 cas dont les cas limites critiques.

## Résumé exécutif

`Book` (`booking.go:36-45`) est le seul point chaud du projet — non par sa taille, mais par sa responsabilité. Elle est la seule fonction à retourner un état modifié ; elle valide désormais `n > 0` et `n ≤ Remaining(s)`, retourne `(Slot, error)`, et expose deux sentinelles d'erreur (`ErrInvalidBookingCount`, `ErrCapacityExceeded`). Les deux autres fonctions (`Remaining`, `IsAvailable`) sont sans risque. Le fichier de tests (`booking_test.go`) couvre 7 cas incluant les cas limites critiques (sur-réservation, n=0, n négatif, créneau exactement plein).

## Constats détaillés

**`VÉRIFIÉ_CODE` — `Book` : la seule mutation, avec garde de capacité et de n.** La fonction `Book` (`booking.go:36-45`) est la seule à produire un état différent de son entrée. Elle retourne `(Slot, error)` et valide ses paramètres avant tout incrément :

```go
func Book(s Slot, n int) (Slot, error) {
    if n <= 0 {
        return s, ErrInvalidBookingCount
    }
    if n > Remaining(s) {
        return s, ErrCapacityExceeded
    }
    s.Booked += n
    return s, nil
}
```

La complexité cyclomatique est 3 (deux branches de garde explicites). Un senior regarderait ici en premier avant tout refactoring ou ajout d'un transport.

**`VÉRIFIÉ_CODE` — `Remaining` et `IsAvailable` : fonctions sans risque dans l'état actuel.** `Remaining(s Slot) int` retourne `s.Capacity - s.Booked` (`booking.go:24-26`), `IsAvailable(s Slot) bool` retourne `Remaining(s) > 0` (`booking.go:29-31`). Les deux sont pures, déterministes, sans effet de bord. Leur seul risque est passif : elles acceptent et calculent sur des valeurs incohérentes (Booked > Capacity → Remaining négatif) sans jamais le signaler — mais `Book` ne peut plus produire cet état puisqu'elle valide `n ≤ Remaining(s)`.

**`VÉRIFIÉ_CODE` — `booking_test.go` : couverture nominale et cas critiques.** Le fichier compte 7 tests (`booking_test.go:12-63`) utilisant la même fixture `sample()` : `Slot{Capacity: 10, Booked: 4}`. Cas couverts :

| Cas | Résultat attendu | Couvert |
|---|---|---|
| `Book(s, 2)` | `Booked=6`, err=nil | Oui (`TestBook`) |
| `Book(s, 0)` | `ErrInvalidBookingCount` | Oui (`TestBookZero`) |
| `Book(s, -3)` | `ErrInvalidBookingCount` | Oui (`TestBookNegative`) |
| `Book(s, 7)` sur `Remaining=6` | `ErrCapacityExceeded` | Oui (`TestBookCapacityExceeded`) |
| `Book(s, 6)` — créneau exactement plein | `Booked==Capacity`, err=nil | Oui (`TestBookExactCapacity`) |
| `Remaining` quand `Booked > Capacity` | Valeur négative | Non |
| `Slot{Capacity: 0}` | Comportement indéfini | Non |

**`VÉRIFIÉ_CODE` — Absence de goroutines et de concurrence.** Aucun `go func`, aucun `sync.Mutex`, aucun `chan` dans le package (`booking.go:1-45`). Les fonctions sont sûres en lecture concurrente (valeurs, pas de pointeurs partagés). En revanche, le schéma `if IsAvailable(s) { Book(s, n) }` n'est pas atomique — c'est une race condition TOCTOU latente dès qu'un serveur concurrent est ajouté.

**`VÉRIFIÉ_CODE` — Aucun `panic`, aucun `recover`, aucun `log`.** Le package ne panique jamais, ne trace jamais, ne loggue jamais (`booking.go:1-45`). C'est correct pour une bibliothèque. Les erreurs sont désormais retournées explicitement (`ErrInvalidBookingCount`, `ErrCapacityExceeded`).

## Forces

- `VÉRIFIÉ_CODE` : `Book` valide ses pré-conditions (`n > 0`, `n ≤ Remaining(s)`) et retourne `(Slot, error)` — aucun état invalide ne peut être produit par `Book` (`booking.go:36-45`).
- `VÉRIFIÉ_CODE` : Aucun couplage entre fonctions sauf `IsAvailable → Remaining` (`booking.go:29-31`) et `Book → Remaining` (`booking.go:40`), couplages intentionnels et cohérents.
- `VÉRIFIÉ_CODE` : Zéro état global — aucune variable de package hors sentinelles d'erreur, aucun `init()`. Impossible de "polluer" l'état entre deux appels.
- `VÉRIFIÉ_CODE` : 7 tests couvrant les cas nominaux et critiques (`booking_test.go:12-63`).

## Dettes techniques

- `VÉRIFIÉ_CODE` : `booking_test.go` — fixture unique avec `time.Now()` (`booking_test.go:9`). Si le champ `Start` devient fonctionnel dans les comparaisons, les tests deviendront non déterministes sans modification.
- `VÉRIFIÉ_CODE` : `Slot` sans validation à la construction — un `Slot{Capacity: -1}` ou `Slot{Booked: 999}` est accepté sans erreur. `Book` protège contre la sur-réservation mais pas contre un état initial incohérent.

## Zones critiques

- `VÉRIFIÉ_CODE` : `internal/booking/booking.go:36-45` (`Book`) — cœur de la logique de mutation, avec garde. La fonction est désormais sûre vis-à-vis de la sur-réservation et des valeurs `n` invalides.
- `VÉRIFIÉ_CODE` : `internal/booking/booking_test.go:8-10` (`sample()`) — fixture unique avec `time.Now()`. Si le champ `Start` devient fonctionnel dans les comparaisons, ces tests deviendront non déterministes sans modification.

## Risques

- `VÉRIFIÉ_CODE` (résolu) : `Book(Slot{Capacity: 10, Booked: 10}, 1)` retourne désormais `(Slot{}, ErrCapacityExceeded)` (`booking.go:40-41`). La sur-réservation silencieuse n'est plus possible via `Book`.
- `HYPOTHÈSE` : Dès qu'un serveur concurrent est introduit, le schéma `IsAvailable → Book` sans verrou devient une race condition TOCTOU. Ce pattern est probable dans tout handler HTTP REST basique.

## Recommandations priorisées

1. ~~**Ajouter une garde de capacité dans `Book`**~~ — **FAIT** : `Book` retourne `(Slot, error)`, rejette `n ≤ 0` et `n > Remaining(s)` (`booking.go:36-45`).
2. ~~**Écrire des tests pour les cas limites de `Book`**~~ — **FAIT** : 4 nouveaux tests couvrent `n=0`, `n<0`, sur-réservation, créneau exactement plein (`booking_test.go:34-63`).
3. **Remplacer `time.Now()` dans `sample()`** par un instant fixe (`time.Date(2024, 1, 1, 9, 0, 0, 0, time.UTC)`) pour garantir la reproductibilité — `internal/booking/booking_test.go:9`.
4. **Valider `Slot` à la construction** — ajouter une fonction `NewSlot(...)` ou une validation dans un constructeur pour rejeter `Capacity ≤ 0` ou `Booked < 0` dès la création.

## Questions ouvertes

- Une annulation (libération de places) est-elle dans le scope ? Si oui, `Book(s, -n)` est-il le mécanisme prévu, ou une fonction `Cancel(s Slot, n int) (Slot, error)` sera-t-elle ajoutée ? Non déductible du code actuel.
- Y aura-t-il un mécanisme de verrouillage (mutex, transaction DB) pour rendre `IsAvailable → Book` atomique lors de l'introduction d'un serveur concurrent ?
