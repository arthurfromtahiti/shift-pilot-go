# Fonctionnel — Audit

> Confiance : medium — les fonctions implémentées sont entièrement lisibles (`VÉRIFIÉ_CODE`). La cohérence fonctionnelle est évaluée par rapport à l'objectif annoncé dans le README (réservation d'activités nautiques). Tout ce qui concerne l'intention finale du pilote, les règles métier non codées, et les fonctionnalités futures est `HYPOTHÈSE` ou `INCONNU`.

## Compréhension globale

Le projet implémente le cycle de disponibilité de créneau. Les trois opérations sont correctes dans leur périmètre : `Remaining` et `IsAvailable` respectent leurs contrats (`places = capacité − réservées`, `disponible si places > 0`) ; `Book` valide désormais la disponibilité et la validité de `n` avant tout incrément, retourne `(Slot, error)`. Le reste de la fonctionnalité annoncée (catalogue d'activités, clients, paiement, annulation, notifications) reste absent.

## Résumé exécutif

Sur les 3 fonctions implémentées : les 3 sont correctes. `Book` valide désormais ses pré-conditions (`n > 0`, `n ≤ Remaining(s)`) et retourne `(Slot, error)`. Le projet couvre ~5% des fonctionnalités qu'impliquerait un système de réservation d'activités nautiques complet : aucun catalogue, aucun client, aucune persistance, aucune annulation, aucune confirmation, aucune interface utilisateur. C'est cohérent avec la désignation "pilote de test" dans le README.

## Constats détaillés

**`VÉRIFIÉ_CODE` — `Remaining` : correct.** Retourne `s.Capacity - s.Booked` (`booking.go:16`). Correct pour un créneau dans un état valide. Si `Booked > Capacity` (état atteignable via `Book`), retourne un entier négatif sans erreur — comportement arithmétiquement correct mais sémantiquement trompeur.

**`VÉRIFIÉ_CODE` — `IsAvailable` : correct dans sa définition.** Retourne `Remaining(s) > 0` (`booking.go:21`). Conforme à la règle métier "disponible si au moins une place reste". Cas limite : `Booked == Capacity` → `Remaining == 0` → `IsAvailable == false`. Ce comportement est correct.

**`VÉRIFIÉ_CODE` — `Book` : fonctionnellement complet dans son périmètre.** `Book(s Slot, n int) (Slot, error)` valide `n > 0` (retourne `ErrInvalidBookingCount` sinon) et `n ≤ Remaining(s)` (retourne `ErrCapacityExceeded` sinon) avant d'incrémenter `s.Booked` (`booking.go:36-45`). La règle métier élémentaire d'un système de réservation — "on ne peut réserver que ce qui est disponible" — est encodée. `Book(Slot{Capacity:10, Booked:10}, 1)` retourne désormais `(Slot{}, ErrCapacityExceeded)`.

**`VÉRIFIÉ_CODE` — Annulation : non implémentée.** Il n'existe pas de fonction `Cancel`, `Release`, ou `Unbook` dans le package (`booking.go:1-45`). L'annulation via `Book(s, -n)` n'est plus possible : `Book` rejette désormais `n ≤ 0` avec `ErrInvalidBookingCount`. Une fonction dédiée sera nécessaire si l'annulation doit être supportée.

**`VÉRIFIÉ_CODE` — Persistance : absente.** La valeur retournée par `Book` est une valeur Go en mémoire. Si l'appelant ne la stocke pas, la réservation est perdue à la fin du scope. Le README mentionne "staging" sans préciser si la couche de persistance existe ailleurs (`README.md:12`) — c'est `INCONNU` depuis ce dépôt.

**`HYPOTHÈSE` — Gap fonctionnel majeur entre le README et le code.** Le README annonce un projet de "réservation d'activités nautiques (Polynésie française)". Un tel système implique au minimum : un catalogue d'activités, des créneaux planifiés, des clients, un flux de réservation avec confirmation, un flux d'annulation, et probablement un paiement. Aucune de ces entités ni de ces fonctionnalités n'est présente dans le code. Ce gap est attendu d'un "pilote de test" — mais il est important de le nommer clairement pour que les planificateurs d'évolution mesurent l'écart.

**`HYPOTHÈSE` — La désignation "pilote de test SHIFT/Paperclip".** Le README qualifie ce dépôt de pilote. Rien dans le code ne confirme que `Book` sans garde est un choix délibéré du pilote (validation déléguée à l'appelant) plutôt qu'une lacune. Les deux lectures sont possibles, et trancher relève du board, pas du code.

## Forces

- `VÉRIFIÉ_CODE` : `Remaining` et `IsAvailable` implémentent correctement leurs règles métier respectives (`booking.go:24-31`).
- `VÉRIFIÉ_CODE` : `Book` implémente correctement la mécanique d'incrément avec validation des pré-conditions (`n > 0`, `n ≤ Remaining(s)`).
- `VÉRIFIÉ_CODE` : Les noms de fonctions sont sémantiquement cohérents avec le domaine métier (`Remaining`, `IsAvailable`, `Book` — anglais métier clair).

## Dettes techniques

- `VÉRIFIÉ_CODE` : Annulation non implémentée — `Book` rejette désormais `n ≤ 0`, donc `Book(s, -n)` n'est plus utilisable comme mécanisme implicite d'annulation. Une fonction `Cancel(s Slot, n int) (Slot, error)` reste à écrire si l'annulation est dans le scope.
- `VÉRIFIÉ_CODE` : Aucune persistance — toute réservation est perdue à la fin du scope appelant.

## Zones critiques

- `VÉRIFIÉ_CODE` : `internal/booking/booking.go:36-45` (`Book`) — cœur fonctionnel du projet, la seule opération de réservation. La règle métier principale (on ne réserve que si disponible) est désormais encodée.

## Risques

- ~~`VÉRIFIÉ_CODE` (impact direct, prouvé) : sur-réservation silencieuse~~ — **RÉSOLU** : `Book` retourne `ErrCapacityExceeded` si `n > Remaining(s)`. Toute couche appelante reçoit une erreur explicite.
- `HYPOTHÈSE` : Si "staging" est un vrai environnement avec des utilisateurs réels, l'absence de persistance reste un risque opérationnel potentiel.

## Recommandations priorisées

1. ~~**Ajouter la garde de capacité dans `Book`**~~ — **FAIT** (`booking.go:36-45`).
2. **Ajouter une fonction `Cancel(s Slot, n int) (Slot, error)`** pour rendre l'annulation explicite, validée (`n > 0`, `n ≤ s.Booked`) et documentée — `internal/booking/booking.go`.
3. **Clarifier dans le README** si le pilote vise à grandir en service complet ou à rester une bibliothèque de logique pure.
4. **Concevoir la couche de persistance** avant d'ajouter d'autres fonctionnalités.

## Questions ouvertes

- Y a-t-il d'autres fonctionnalités (catalogue, clients, paiement) développées dans un dépôt séparé, ou sont-elles à venir dans ce dépôt ? Non visible ici.
- Le terme "SHIFT/Paperclip" dans le README a-t-il une signification métier précise pour ce projet ? Non déductible du code.
- Quelle est la définition de "done" pour ce pilote — quelles fonctionnalités doit-il couvrir avant de passer en production ?
