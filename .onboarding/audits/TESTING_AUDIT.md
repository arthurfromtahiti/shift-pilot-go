# Tests — Audit

> Confiance : medium — les fichiers de test ont été intégralement lus (`VÉRIFIÉ_CODE`). Les tests n'ont pas pu être exécutés : la toolchain Go est absente de l'environnement d'analyse (statut d'exécution `INCONNU`). L'audit porte sur la structure et la couverture telles qu'elles sont lisibles dans le source, pas sur les résultats d'exécution.

## Compréhension globale

Le projet comporte un unique fichier de tests : `internal/booking/booking_test.go`, 29 lignes, 3 fonctions de test, une fixture partagée. Les tests couvrent le chemin nominal de chacune des trois fonctions du package. Aucun cas limite n'est testé. Il n'y a pas de CI, pas de benchmark, pas de test d'intégration, pas de test de concurrence. La couverture est formellement présente mais fonctionnellement incomplète.

## Résumé exécutif

Trois tests, un seul fixture, chemin nominal uniquement. L'outil `go test ./...` suffit à les lancer (README), mais aucun fichier de CI n'automatise cette exécution dans le dépôt. Les cas critiques pour la robustesse du package — sur-réservation, `n` négatif, créneau exactement plein, `Capacity` invalide — sont tous non testés. Un contributeur qui s'appuie sur "3 tests passants" comme signal de santé prend un risque réel : les cas qui importent le plus ne sont pas couverts.

## Constats détaillés

**`VÉRIFIÉ_CODE` — Structure des tests.** `booking_test.go` est dans le package `booking` (test blanc — accès aux membres non exportés si existants, bien qu'il n'y en ait aucun ici). Il importe `"testing"` et `"time"`. La fixture `sample()` (`booking_test.go:8-10`) retourne `Slot{ID: 1, Activity: "Plongée", Start: time.Now(), Capacity: 10, Booked: 4}`. Toutes les fonctions de test l'utilisent sans variation.

**`VÉRIFIÉ_CODE` — `TestRemaining`** (`booking_test.go:12-16`) : appelle `Remaining(sample())` et vérifie que le résultat vaut `6` (`10 - 4 = 6`). Un seul assert, chemin nominal.

**`VÉRIFIÉ_CODE` — `TestIsAvailable`** (`booking_test.go:18-22`) : appelle `IsAvailable(sample())` et vérifie qu'il retourne `true`. Un seul assert, chemin nominal.

**`VÉRIFIÉ_CODE` — `TestBook`** (`booking_test.go:24-29`) : appelle `Book(sample(), 2)` et vérifie que `s.Booked == 6` (`4 + 2`). Ne vérifie pas que `Capacity` est inchangée, ne teste pas que le `Slot` original est non muté, ne teste pas le cas de sur-réservation.

**`VÉRIFIÉ_CODE` — Cas non testés (lus dans le source, confirmés absents des tests) :**

| Cas | Fonction concernée | Impact si non couvert |
|---|---|---|
| `n = 0` | `Book` | No-op silencieux, jamais détecté |
| `n < 0` | `Book` | Annulation implicite non documentée |
| `n > Remaining(s)` | `Book` | Sur-réservation silencieuse — le risque principal |
| `Booked == Capacity` avant appel | `Book`, `IsAvailable` | Frontière exacte non couverte |
| `Remaining(s)` quand `Booked > Capacity` | `Remaining` | Valeur négative jamais testée |
| `IsAvailable(s)` quand `Booked == Capacity` | `IsAvailable` | Retourne `false` mais jamais vérifié |
| `Capacity = 0` | toutes | Comportement sur créneau sans capacité |
| `Capacity < 0` ou `Booked < 0` | toutes | Valeurs absurdes acceptées silencieusement |
| Immuabilité : le `Slot` source n'est pas muté | `Book` | Pas de test de la sémantique valeur |

**`VÉRIFIÉ_CODE` — Aucun test table-driven.** Go dispose d'un pattern idiomatique (`t.Run`, tranches de cas) pour tester plusieurs entrées sans duplication. Aucun test du package ne l'utilise (`booking_test.go:1-29`).

**`VÉRIFIÉ_CODE` — `time.Now()` dans la fixture.** `sample()` inclut `Start: time.Now()` (`booking_test.go:9`). Aucune fonction de production ne lit le champ `Start` dans le code actuel — ce n'est pas un bug aujourd'hui. Mais si `Start` est utilisé dans une comparaison future (tri, filtrage), les tests deviendront non déterministes et potentiellement en échec selon l'instant d'exécution.

**`VÉRIFIÉ_CODE` — Aucune configuration CI visible.** Le dépôt ne contient pas de `.github/workflows/`, de `Makefile` avec cible `test`, de `Taskfile`, ni de tout autre fichier d'automatisation. Le README mentionne `go test ./...` (`README.md:8`) mais rien n'exécute cette commande automatiquement à chaque push ou PR.

**`INCONNU` — Résultats d'exécution des tests.** La toolchain Go n'est pas disponible dans l'environnement d'analyse. Les assertions ont été lues et sont arithmétiquement cohérentes avec l'implémentation (ex. `10 - 4 = 6`, `4 + 2 = 6`), mais le statut vert/rouge réel des tests est `INCONNU`. Il n'existe pas de rapport de test (fichier `coverage.out`, badge CI) dans le dépôt.

## Forces

- `VÉRIFIÉ_CODE` : Les 3 tests existent et couvrent chacune des 3 fonctions au moins une fois (`booking_test.go:12-29`) — point de départ existant.
- `VÉRIFIÉ_CODE` : Package de test en boîte blanche (même package `booking`) — pas besoin de rendre des fonctions internes accessibles pour les tester si besoin.
- `VÉRIFIÉ_CODE` : Le README documente la commande d'exécution (`go test ./...`) (`README.md:8`).

## Dettes techniques

- `VÉRIFIÉ_CODE` : Aucun cas limite testé — la couverture nominale masque les comportements défectueux de `Book` en sur-réservation (`booking_test.go:24-29`).
- `VÉRIFIÉ_CODE` : `time.Now()` dans la fixture — fragilité latente si `Start` devient fonctionnel (`booking_test.go:9`).
- `VÉRIFIÉ_CODE` : Aucune CI — les tests ne s'exécutent que manuellement.

## Zones critiques

- `VÉRIFIÉ_CODE` : `internal/booking/booking_test.go:24-29` (`TestBook`) — le test le plus important est aussi le moins complet : il valide uniquement l'incrément nominal, pas la garde de capacité qui est le risque numéro 1 du package.

## Risques

- `VÉRIFIÉ_CODE` (impact prouvé par absence) : `Book(Slot{Capacity:10, Booked:10}, 1)` retourne `Slot{Booked:11}` sans erreur, et aucun test ne détecte ce comportement. Un futur contributeur qui ajoute ce cas d'usage sans test le manquera.
- `HYPOTHÈSE` : Sans CI, un contributeur peut casser les 3 tests existants et pousser sans s'en rendre compte — particulièrement risqué si ce code est utilisé comme dépendance par un autre module.

## Recommandations priorisées

1. **Convertir `TestBook` en test table-driven** couvrant : `n=2` (nominal), `n=0` (no-op), `n=-2` (décrémentation), `n=7` sur `Remaining=6` (sur-réservation) — `internal/booking/booking_test.go`.
2. **Ajouter des tests pour `TestRemaining` et `TestIsAvailable`** sur les cas limites : `Booked == Capacity` (Remaining=0, IsAvailable=false), `Booked > Capacity` (Remaining négatif) — `internal/booking/booking_test.go`.
3. **Remplacer `time.Now()` dans `sample()` par une date fixe** (`time.Date(2024, 1, 1, 9, 0, 0, 0, time.UTC)`) — `internal/booking/booking_test.go:9`.
4. **Créer un workflow CI** (ex. `.github/workflows/ci.yml`) exécutant `go test ./... -race` à chaque push sur `main` et sur toute PR — dépôt racine.
5. **Ajouter un test `TestBook_SlotsAreImmutable`** vérifiant que le `Slot` source n'est pas muté après `Book` — propriété garantie par le passage par valeur mais jamais testée explicitement.

## Questions ouvertes

- Les tests ont-ils déjà été exécutés dans un environnement CI ? Aucun badge ni rapport visible dans le dépôt.
- Le taux de couverture cible est-il défini (ex. 80% de lignes) ? Non visible.
- Y aura-t-il des tests d'intégration lorsqu'une couche de persistance sera ajoutée ? Non prévu dans le code actuel.
