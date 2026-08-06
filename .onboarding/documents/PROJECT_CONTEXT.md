# Project Context — shift-pilot-go

> **Confiance : medium**. Codebase entièrement lue et vérifiée. Contexte d'exécution (CI/CD, déploiement staging, intentions futures) `INCONNU` : aucune configuration ne figure dans le dépôt.

## Nature du projet

**shift-pilot-go** est un **pilote de test** (énoncé dans le README : "pilote de test SHIFT/Paperclip") — une bibliothèque Go qui implémente le noyau logique d'un système de réservation d'activités nautiques en Polynésie française. Le projet a pour ambition d'évoluer en application complète, mais il en couvre aujourd'hui une fraction minimal : une seule entité (`Slot`), trois fonctions pures, aucun serveur, aucune persistance, aucune interface utilisateur.

**Caractéristiques techniques :**
- **Type** : bibliothèque Go standard (module `github.com/arthurfromtahiti/shift-pilot-go`), Go 1.21.
- **Dépendances** : aucune (bibliothèque standard uniquement — `time` seul import).
- **Empreinte** : 28 lignes de code de production (`internal/booking/booking.go`), 29 lignes de tests (`internal/booking/booking_test.go`).
- **Branche maître** : `main` (SHA initial analysé : `8122e2d`).

## État courant — matière très restreinte

Le dépôt contient un seul domaine métier **prouvé** : **Réservation & disponibilité de créneaux** (`internal/booking`). Ce domaine expose trois fonctions sur un type unique (`Slot`) :

| Fonction | Rôle | Risque |
|---|---|---|
| `Remaining(Slot) int` | Calcule le nombre de places libres (Capacité − Réservées) | Peut retourner une valeur négative si `Booked > Capacity` — pas de garde |
| `IsAvailable(Slot) bool` | Prédicat : au moins une place reste-t-elle ? | Aucun — dérive correctement de `Remaining` |
| `Book(Slot, n int) (Slot, error)` | Enregistre n places réservées, retourne le créneau mis à jour | Valide `n > 0` (`ErrInvalidBookingCount`) et `n ≤ Remaining(s)` (`ErrCapacityExceeded`) ; sur-réservation silencieuse impossible |

**Domaines absents du code** : Catalogue d'activités, comptes clients, paiement, planning/agenda, notifications, annulation. Ces fonctionnalités, logiquement attendues pour un système de réservation complet, ne sont représentées nulle part — ce qui est cohérent avec la nature de "pilote de test" mais doit être explicite pour planifier l'évolution.

## Acteurs et usages

**Acteur unique identifié** : code Go appelant, qui importe le package `internal/booking` et utilise ses trois fonctions. Aucun point d'entrée HTTP, CLI, ou interface utilisateur n'existe dans le code actuel. Les workflows énoncés dans l'audit sont donc purement techniques (appels entre fonctions), pas des scénarios utilisateur.

## Risques prioritaires

### ~~Défaut fonctionnel — `Book` sans garde de capacité~~ — RÉSOLU
`Book` valide désormais ses pré-conditions : `n > 0` (sinon `ErrInvalidBookingCount`) et `n ≤ Remaining(s)` (sinon `ErrCapacityExceeded`). `Book(Slot{Capacity:10, Booked:10}, 1)` retourne maintenant `(Slot{}, ErrCapacityExceeded)` — l'état invalide `Booked > Capacity` ne peut plus être produit par `Book`.

### Absence de persistance
Les valeurs retournées par `Book` sont des valeurs Go en mémoire. Elles disparaissent à la fin du scope appelant si elles ne sont pas explicitement conservées. Le README mentionne un déploiement vers "staging" sans préciser où les réservations sont stockées — ce mécanisme est `INCONNU`.

### Suite de tests — cas nominaux et critiques couverts
Sept tests couvrent les chemins nominaux et les cas limites de `Book` : `n=0`, `n<0`, sur-réservation, créneau exactement plein. Restent non couverts : `Remaining` quand `Booked > Capacity`, `Capacity ≤ 0`, immuabilité du `Slot` source.

## Priorités prochaines

1. **Clarifier la cible d'évolution** : le projet reste-t-il une bibliothèque importable, ou doit-il évoluer en service HTTP autonome ? Cette décision structure toute l'architecture future (transport, persistance, interfaces, gestion d'erreurs).
2. ~~**Ajouter la garde de capacité à `Book`**~~ — **FAIT** (`booking.go:36-45`).
3. **Concevoir et implémenter la couche de persistance** avant d'ajouter d'autres fonctionnalités — sans elle, aucune réservation n'est durable.
4. **Compléter la suite de tests** : `Remaining` quand `Booked > Capacity`, `Capacity ≤ 0`, immuabilité du `Slot` source.

## Questions non tranchées

- Le terme "SHIFT/Paperclip" a-t-il une signification métier précise ou une affectation interne ? Non déductible du code.
- "Staging" dans le README est-il un environnement réel tourne-t-il ce code, ou une note de processus ? Pas d'artefact de déploiement visible dans le dépôt.
- `Book` rejette désormais `n ≤ 0` — une fonction `Cancel` dédiée sera-t-elle ajoutée pour l'annulation ?
