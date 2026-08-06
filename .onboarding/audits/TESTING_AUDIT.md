# Tests — Audit

> Confiance : medium — les fichiers de test ont été intégralement lus (`VÉRIFIÉ_CODE`). Les tests n'ont pas pu être exécutés : la toolchain Go est absente de l'environnement d'analyse (statut d'exécution `INCONNU`). L'audit porte sur la structure et la couverture telles qu'elles sont lisibles dans le source, pas sur les résultats d'exécution.

## Compréhension globale

Le projet comporte un unique fichier de tests : `internal/booking/booking_test.go`, 63 lignes, 7 fonctions de test, une fixture partagée. Les tests couvrent le chemin nominal de chacune des trois fonctions du package, ainsi que les cas limites critiques de `Book` (sur-réservation, n=0, n négatif, créneau exactement plein). Il n'y a pas de CI, pas de benchmark, pas de test d'intégration, pas de test de concurrence.

## Résumé exécutif

Sept tests, un seul fixture, chemins nominal et cas critiques. L'outil `go test ./...` suffit à les lancer (README), mais aucun fichier de CI n'automatise cette exécution dans le dépôt. Les cas critiques de `Book` — sur-réservation, `n` négatif, `n=0`, créneau exactement plein — sont désormais couverts. Reste non couvert : `Remaining` quand `Booked > Capacity`, `Capacity` invalide à la construction.

## Constats détaillés

**`VÉRIFIÉ_CODE` — Structure des tests.** `booking_test.go` est dans le package `booking` (test blanc — accès aux membres non exportés si existants, bien qu'il n'y en ait aucun ici). Il importe `"testing"` et `"time"`. La fixture `sample()` (`booking_test.go:8-10`) retourne `Slot{ID: 1, Activity: "Plongée", Start: time.Now(), Capacity: 10, Booked: 4}`. Toutes les fonctions de test l'utilisent sans variation.

**`VÉRIFIÉ_CODE` — `TestRemaining`** (`booking_test.go:12-16`) : appelle `Remaining(sample())` et vérifie que le résultat vaut `6` (`10 - 4 = 6`). Un seul assert, chemin nominal.

**`VÉRIFIÉ_CODE` — `TestIsAvailable`** (`booking_test.go:18-22`) : appelle `IsAvailable(sample())` et vérifie qu'il retourne `true`. Un seul assert, chemin nominal.

**`VÉRIFIÉ_CODE` — `TestBook`** (`booking_test.go:24-32`) : appelle `Book(sample(), 2)`, vérifie `err == nil` et `s.Booked == 6` (`4 + 2`). Ne vérifie pas que `Capacity` est inchangée, ne teste pas que le `Slot` original est non muté.

**`VÉRIFIÉ_CODE` — Cas couverts et non couverts :**

| Cas | Fonction concernée | Couvert | Test |
|---|---|---|---|
| `n = 0` | `Book` | Oui — `ErrInvalidBookingCount` | `TestBookZero` |
| `n < 0` | `Book` | Oui — `ErrInvalidBookingCount` | `TestBookNegative` |
| `n > Remaining(s)` | `Book` | Oui — `ErrCapacityExceeded` | `TestBookCapacityExceeded` |
| `Booked == Capacity` avant appel (`Book` réussit si n=0 restants → rejeté) | `Book` | Oui (via `TestBookExactCapacity`, créneau exactement rempli) | `TestBookExactCapacity` |
| `Remaining(s)` quand `Booked > Capacity` | `Remaining` | Non | — |
| `IsAvailable(s)` quand `Booked == Capacity` | `IsAvailable` | Non | — |
| `Capacity = 0` | toutes | Non | — |
| `Capacity < 0` ou `Booked < 0` | toutes | Non | — |
| Immuabilité : le `Slot` source n'est pas muté | `Book` | Non | — |

**`VÉRIFIÉ_CODE` — Aucun test table-driven.** Go dispose d'un pattern idiomatique (`t.Run`, tranches de cas) pour tester plusieurs entrées sans duplication. Aucun test du package ne l'utilise (`booking_test.go:1-29`).

**`VÉRIFIÉ_CODE` — `time.Now()` dans la fixture.** `sample()` inclut `Start: time.Now()` (`booking_test.go:9`). Aucune fonction de production ne lit le champ `Start` dans le code actuel — ce n'est pas un bug aujourd'hui. Mais si `Start` est utilisé dans une comparaison future (tri, filtrage), les tests deviendront non déterministes et potentiellement en échec selon l'instant d'exécution.

**`VÉRIFIÉ_CODE` — Aucune configuration CI visible.** Le dépôt ne contient pas de `.github/workflows/`, de `Makefile` avec cible `test`, de `Taskfile`, ni de tout autre fichier d'automatisation. Le README mentionne `go test ./...` (`README.md:8`) mais rien n'exécute cette commande automatiquement à chaque push ou PR.

**`INCONNU` — Résultats d'exécution des tests.** La toolchain Go n'est pas disponible dans l'environnement d'analyse. Les assertions ont été lues et sont arithmétiquement cohérentes avec l'implémentation (ex. `10 - 4 = 6`, `4 + 2 = 6`), mais le statut vert/rouge réel des tests est `INCONNU`. Il n'existe pas de rapport de test (fichier `coverage.out`, badge CI) dans le dépôt.

## Forces

- `VÉRIFIÉ_CODE` : Les 7 tests couvrent les 3 fonctions nominales et 4 cas limites de `Book` (`booking_test.go:12-63`).
- `VÉRIFIÉ_CODE` : Package de test en boîte blanche (même package `booking`) — pas besoin de rendre des fonctions internes accessibles pour les tester si besoin.
- `VÉRIFIÉ_CODE` : Le README documente la commande d'exécution (`go test ./...`) (`README.md:8`).
- `VÉRIFIÉ_CODE` : `Book` expose deux sentinelles d'erreur (`ErrCapacityExceeded`, `ErrInvalidBookingCount`) que les tests comparent directement par identité — pattern correct en Go.

## Dettes techniques

- `VÉRIFIÉ_CODE` : `time.Now()` dans la fixture — fragilité latente si `Start` devient fonctionnel (`booking_test.go:9`).
- `VÉRIFIÉ_CODE` : Aucune CI — les tests ne s'exécutent que manuellement.
- `VÉRIFIÉ_CODE` : Cas limites non encore couverts : `Remaining` quand `Booked > Capacity`, `IsAvailable` quand `Booked == Capacity`, `Capacity = 0`, immuabilité du `Slot` source après `Book`.

## Zones critiques

- `VÉRIFIÉ_CODE` : `internal/booking/booking_test.go:24-32` (`TestBook`) — test nominal ; les cas d'erreur sont couverts par des tests dédiés (`TestBookCapacityExceeded`, `TestBookZero`, `TestBookNegative`, `TestBookExactCapacity`).

## Risques

- ~~`VÉRIFIÉ_CODE` : `Book(Slot{Capacity:10, Booked:10}, 1)` retourne `Slot{Booked:11}` sans erreur~~ — **RÉSOLU** : `TestBookCapacityExceeded` couvre ce cas et `Book` retourne `ErrCapacityExceeded`.
- `HYPOTHÈSE` : Sans CI, un contributeur peut casser les 7 tests existants et pousser sans s'en rendre compte — particulièrement risqué si ce code est utilisé comme dépendance par un autre module.

## Recommandations priorisées

1. ~~**Ajouter des tests pour les cas limites de `Book`**~~ — **FAIT** (`TestBookCapacityExceeded`, `TestBookZero`, `TestBookNegative`, `TestBookExactCapacity`).
2. **Ajouter des tests pour `TestRemaining` et `TestIsAvailable`** sur les cas limites : `Booked == Capacity` (Remaining=0, IsAvailable=false), `Booked > Capacity` (Remaining négatif) — `internal/booking/booking_test.go`.
3. **Remplacer `time.Now()` dans `sample()` par une date fixe** (`time.Date(2024, 1, 1, 9, 0, 0, 0, time.UTC)`) — `internal/booking/booking_test.go:9`.
4. **Créer un workflow CI** (ex. `.github/workflows/ci.yml`) exécutant `go test ./... -race` à chaque push sur `main` et sur toute PR — dépôt racine.
5. **Ajouter un test `TestBook_SlotsAreImmutable`** vérifiant que le `Slot` source n'est pas muté après `Book` — propriété garantie par le passage par valeur mais jamais testée explicitement.

## Questions ouvertes

- Les tests ont-ils déjà été exécutés dans un environnement CI ? Aucun badge ni rapport visible dans le dépôt.
- Le taux de couverture cible est-il défini (ex. 80% de lignes) ? Non visible.
- Y aura-t-il des tests d'intégration lorsqu'une couche de persistance sera ajoutée ? Non prévu dans le code actuel.
