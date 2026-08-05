# Fonctionnel — Audit

> Confiance : medium — les fonctions implémentées sont entièrement lisibles (`VÉRIFIÉ_CODE`). La cohérence fonctionnelle est évaluée par rapport à l'objectif annoncé dans le README (réservation d'activités nautiques). Tout ce qui concerne l'intention finale du pilote, les règles métier non codées, et les fonctionnalités futures est `HYPOTHÈSE` ou `INCONNU`.

## Compréhension globale

Le projet implémente partiellement un cycle de disponibilité de créneau. Deux des trois opérations sont correctes dans leur périmètre : `Remaining` et `IsAvailable` respectent leurs contrats (`places = capacité − réservées`, `disponible si places > 0`). La troisième, `Book`, est fonctionnellement incomplète : elle enregistre une réservation sans vérifier la disponibilité, ce qui viole l'invariant fondamental d'un système de réservation. Le reste de la fonctionnalité annoncée (catalogue d'activités, clients, paiement, annulation, notifications) est absent.

## Résumé exécutif

Sur les 3 fonctions implémentées : 2 sont correctes (`Remaining`, `IsAvailable`) ; 1 est incomplète (`Book` — manque la garde de capacité). Le projet couvre ~5% des fonctionnalités qu'impliquerait un système de réservation d'activités nautiques complet : aucun catalogue, aucun client, aucune persistance, aucune annulation, aucune confirmation, aucune interface utilisateur. C'est cohérent avec la désignation "pilote de test" dans le README, mais le défaut de `Book` ne relève pas du scope limité — c'est une règle métier absente dans la seule fonctionnalité implémentée.

## Constats détaillés

**`VÉRIFIÉ_CODE` — `Remaining` : correct.** Retourne `s.Capacity - s.Booked` (`booking.go:16`). Correct pour un créneau dans un état valide. Si `Booked > Capacity` (état atteignable via `Book`), retourne un entier négatif sans erreur — comportement arithmétiquement correct mais sémantiquement trompeur.

**`VÉRIFIÉ_CODE` — `IsAvailable` : correct dans sa définition.** Retourne `Remaining(s) > 0` (`booking.go:21`). Conforme à la règle métier "disponible si au moins une place reste". Cas limite : `Booked == Capacity` → `Remaining == 0` → `IsAvailable == false`. Ce comportement est correct.

**`VÉRIFIÉ_CODE` — `Book` : fonctionnellement incomplet.** `Book(s Slot, n int) Slot` incrémente `s.Booked` de `n` sans vérifier `IsAvailable(s)` ni `n > 0` (`booking.go:25-28`). La règle métier élémentaire d'un système de réservation — "on ne peut réserver que ce qui est disponible" — n'est pas encodée. En conséquence : `Book(Slot{Capacity:10, Booked:10}, 1)` retourne `Slot{Booked:11}` sans erreur. Ce n'est pas un cas d'usage non prévu — c'est le cas d'usage le plus attendu d'un système de réservation (tentative de réservation sur créneau plein) et il est traité incorrectement.

**`VÉRIFIÉ_CODE` — Annulation : non implémentée.** Il n'existe pas de fonction `Cancel`, `Release`, ou `Unbook` dans le package (`booking.go:1-28`). L'annulation est physiquement possible via `Book(s, -n)`, mais ce comportement n'est ni documenté (`booking.go:24` : commentaire de `Book` ne mentionne pas les valeurs négatives), ni testé (`booking_test.go` ne teste que `n=2`), ni une fonction dédiée.

**`VÉRIFIÉ_CODE` — Persistance : absente.** La valeur retournée par `Book` est une valeur Go en mémoire. Si l'appelant ne la stocke pas, la réservation est perdue à la fin du scope. Le README mentionne "staging" sans préciser si la couche de persistance existe ailleurs (`README.md:12`) — c'est `INCONNU` depuis ce dépôt.

**`HYPOTHÈSE` — Gap fonctionnel majeur entre le README et le code.** Le README annonce un projet de "réservation d'activités nautiques (Polynésie française)". Un tel système implique au minimum : un catalogue d'activités, des créneaux planifiés, des clients, un flux de réservation avec confirmation, un flux d'annulation, et probablement un paiement. Aucune de ces entités ni de ces fonctionnalités n'est présente dans le code. Ce gap est attendu d'un "pilote de test" — mais il est important de le nommer clairement pour que les planificateurs d'évolution mesurent l'écart.

**`HYPOTHÈSE` — La désignation "pilote de test SHIFT/Paperclip".** Le README qualifie ce dépôt de pilote. Rien dans le code ne confirme que `Book` sans garde est un choix délibéré du pilote (validation déléguée à l'appelant) plutôt qu'une lacune. Les deux lectures sont possibles, et trancher relève du board, pas du code.

## Forces

- `VÉRIFIÉ_CODE` : `Remaining` et `IsAvailable` implémentent correctement leurs règles métier respectives (`booking.go:15-22`).
- `VÉRIFIÉ_CODE` : `Book` implémente correctement la mécanique d'incrément — le problème est l'absence de pré-condition, pas un bug dans la computation elle-même.
- `VÉRIFIÉ_CODE` : Les noms de fonctions sont sémantiquement cohérents avec le domaine métier (`Remaining`, `IsAvailable`, `Book` — anglais métier clair).

## Dettes techniques

- `VÉRIFIÉ_CODE` : `Book` sans garde de disponibilité — la règle métier centrale "on réserve seulement si disponible" n'est pas encodée (`booking.go:25-28`).
- `VÉRIFIÉ_CODE` : Annulation non implémentée — `Book(s, -n)` est un mécanisme implicite non documenté, pas un contrat (`booking.go:24`).
- `VÉRIFIÉ_CODE` : Aucune persistance — toute réservation est perdue à la fin du scope appelant.

## Zones critiques

- `VÉRIFIÉ_CODE` : `internal/booking/booking.go:25-28` (`Book`) — le cœur fonctionnel du projet, la seule opération de réservation, est incomplète dans sa règle métier principale.

## Risques

- `VÉRIFIÉ_CODE` (impact direct, prouvé) : Toute couche appelante qui n'appelle pas `IsAvailable` avant `Book` produit une sur-réservation silencieuse. Dans l'état actuel (pas de serveur), le risque est théorique — mais c'est une dette à solder avant toute exposition.
- `HYPOTHÈSE` : Si "staging" est un vrai environnement avec des utilisateurs réels, la sur-réservation silencieuse de `Book` est un incident opérationnel potentiel — des places pourraient être vendues en double sans signal d'erreur.

## Recommandations priorisées

1. **Ajouter la garde de capacité dans `Book`** : vérifier `IsAvailable(s)` et `n > 0` avant d'incrémenter, retourner `(Slot, error)` avec une erreur descriptive si la condition n'est pas remplie — `internal/booking/booking.go:25`.
2. **Ajouter une fonction `Cancel(s Slot, n int) (Slot, error)`** pour rendre l'annulation explicite, validée (`n > 0`, `n ≤ s.Booked`) et documentée — `internal/booking/booking.go`.
3. **Clarifier dans le README** si le pilote vise à grandir en service complet ou à rester une bibliothèque de logique pure — le gap actuel entre l'objectif annoncé et le code existant nécessite une décision d'orientation.
4. **Concevoir la couche de persistance** avant d'ajouter d'autres fonctionnalités — sans elle, `Book` retourne des valeurs qui disparaissent, et toute fonctionnalité ajoutée sera testable mais pas déployable.

## Questions ouvertes

- La garde de disponibilité dans `Book` est-elle intentionnellement absente (délégée à l'appelant) ou une omission du pilote ? Non tranchable par le code — question à poser au board.
- Y a-t-il d'autres fonctionnalités (catalogue, clients, paiement) développées dans un dépôt séparé, ou sont-elles à venir dans ce dépôt ? Non visible ici.
- Le terme "SHIFT/Paperclip" dans le README a-t-il une signification métier précise pour ce projet ? Non déductible du code.
- Quelle est la définition de "done" pour ce pilote — quelles fonctionnalités doit-il couvrir avant de passer en production ?
