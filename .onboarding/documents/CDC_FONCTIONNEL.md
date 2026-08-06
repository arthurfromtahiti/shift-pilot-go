# Cahier des Charges Fonctionnel — shift-pilot-go

> **Confiance : medium**. Fonctionnalités implémentées entièrement lisibles (`VÉRIFIÉ_CODE`). Intentions futures relatives au contexte d'usage (interface utilisateur, paiement, notifications) sont `HYPOTHÈSE` ou `INCONNU` — déductibles uniquement du README, pas du code.

## Contexte métier

**Problème résolu** : permettre la gestion des places disponibles dans des créneaux d'activités nautiques (plongée, etc.) offerts en Polynésie française.

**Acteurs métier** : aucun acteur humain représenté dans le code actuellement. Le package n'expose qu'une API technique pour un code appelant Go. Aucun utilisateur, aucune interface de réservant, aucun administrateur d'activités n'est implémenté.

**Scope présent vs. futur** :
- ✓ **Présent** : logique pure de calcul de disponibilité et d'enregistrement de réservation sur un créneau unique.
- ✗ **Absent** : catalogue d'activités, clients/réservants, paiement, annulation explicite, persistance, interfaces utilisateur, authentification, notifications.

## Acteurs

Aucun acteur métier explicite dans le code. La seule entité de fait est un **appelant technique** — du code Go qui importe le package `internal/booking` et invoque ses fonctions. Le comportement de cet appelant n'est pas contraint par le package ; il relève de la responsabilité du code appelant.

## Parcours fonctionnels

Le code implémente deux workflows connexes. Ces workflows sont purement techniques, sans point d'entrée utilisateur.

### WF1 — Consulter la disponibilité d'un créneau

**Type** : `technical_flow` — pas de serveur, pas d'interface utilisateur, API de package Go.

**Objectif** : permettre au code appelant de connaître le nombre de places libres et si le créneau peut accepter des réservations.

**Étapes** :
1. L'appelant détient ou construit un `Slot` avec valeurs cohérentes de `Capacity` et `Booked`.
2. L'appelant appelle `Remaining(Slot) int` → obtient `Capacity - Booked`.
3. *(Optionnel)* L'appelant appelle `IsAvailable(Slot) bool` → obtient `Remaining(s) > 0`.
4. L'appelant utilise le résultat pour décider si la réservation peut procéder.

**Règles métier**
- **Places restantes = Capacité − Réservées** : `Remaining(s)` retourne `s.Capacity - s.Booked` (`internal/booking/booking.go:24-26`).
- **Disponible ssi Restantes > 0** : `IsAvailable(s)` retourne `Remaining(s) > 0` (`internal/booking/booking.go:29-31`).

**Données en entrée** : `Slot{ID, Activity, Start, Capacity, Booked}`.
**Données en sortie** : `int` (nombre de places libres, potentiellement négatif si sur-réservation), `bool` (disponible oui/non).

**Risques**
- `Remaining` peut retourner un entier négatif si `Booked > Capacity` (état atteignable via `Book` — voir WF2). Aucune garde ne le prévient.
- Les fonctions n'acceptent aucune validation d'entrée : un `Slot` avec `Capacity = -5` est traité sans erreur.
- Pas de concurrence : si une goroutine concurrente modifie l'état du `Slot` entre un appel à `IsAvailable` et l'action que l'appelant en déduit, le résultat peut être obsolète.

### WF2 — Réserver des places sur un créneau

**Type** : `technical_flow` — API de package Go, pas de serveur.

**Objectif** : enregistrer la réservation de `n` places sur un créneau, retourner le créneau mis à jour.

**Étapes**
1. L'appelant détient ou construit un `Slot` et choisit `n` (nombre de places à réserver).
2. L'appelant appelle `Book(Slot, n int) (Slot, error)`.
3. `Book` vérifie que `n > 0` (sinon retourne `ErrInvalidBookingCount`) et que `n ≤ Remaining(s)` (sinon retourne `ErrCapacityExceeded`), puis incrémente `s.Booked` de `n` (`internal/booking/booking.go:36-45`).
4. `Book` retourne un nouveau `Slot` avec `Booked` augmenté, et `nil` comme erreur.
5. L'appelant est responsable de la persistance du `Slot` retourné — sinon la réservation est perdue.

**Règles métier**
- **`Book` valide ses pré-conditions** : `n > 0` (sinon `ErrInvalidBookingCount`) et `n ≤ Remaining(s)` (sinon `ErrCapacityExceeded`). Un créneau complet ou une valeur `n` invalide provoque un retour d'erreur explicite (`internal/booking/booking.go:36-45`).
- **Non-mutatif** : `Book` reçoit `Slot` par valeur, retourne une nouvelle valeur — pas de modification en place. L'appelant doit utiliser la valeur de retour.
- **`n` négatif ou nul est rejeté** : `Book(s, -2)` et `Book(s, 0)` retournent `ErrInvalidBookingCount`. L'annulation nécessitera une fonction dédiée.

**Données en entrée** : `Slot` + `int n` (nombre de places à réserver).
**Données en sortie** : `(Slot, error)` — `Slot` modifié (`Booked` incrémenté de `n`) si succès, `error` non nulle (`ErrInvalidBookingCount` ou `ErrCapacityExceeded`) sinon.

**Risques**
- ~~**Sur-réservation silencieuse**~~ — **RÉSOLU** : `Book(Slot{Capacity:10, Booked:10}, 1)` retourne désormais `(Slot{}, ErrCapacityExceeded)`.
- ~~**`n` négatif accepté silencieusement**~~ — **RÉSOLU** : `Book(s, -3)` retourne `ErrInvalidBookingCount`.
- **Aucune persistance** : la valeur retournée est une valeur Go en mémoire. Si l'appelant oublie de la sauvegarder, la réservation disparaît.
- **Pas d'atomicité** : le schéma « vérifier la disponibilité → appeler `Book` » n'est pas atomique. Entre les deux appels, un concurrent peut modifier le `Slot`.

## Règles métier — synthèse

| Règle | Preuve | Statut |
|---|---|---|
| Places libres = Capacité − Réservées | `Remaining` ligne 24-26 | ✓ Implémentée correctement |
| Un créneau est disponible ssi places libres > 0 | `IsAvailable` ligne 29-31 | ✓ Implémentée correctement |
| On ne peut réserver que si disponible | `Book` lignes 36-45 | ✓ Implémentée — `ErrCapacityExceeded` si `n > Remaining(s)` |
| Réservation enregistrée → places réservées augmentent | `Book` ligne 43 | ✓ Implémentée avec garde |
| Annulation possible | Aucun mécanisme — `Book(s, -n)` est rejeté (`ErrInvalidBookingCount`) | ✗ **Absente** — fonction `Cancel` à créer |

## Modèle de données

**Entité unique : `Slot`** (créneau d'activité)

```
Slot:
  ID       int          — identifiant du créneau (aucune garantie d'unicité sans persistance)
  Activity string       — nom de l'activité (chaîne libre, aucune validation)
  Start    time.Time    — date-heure de début
  Capacity int          — places totales (aucune contrainte de positivité)
  Booked   int          — places réservées (peut dépasser Capacity)
```

**Invariants attendus** (non enforced par le code) :
- `Capacity > 0` — un créneau doit avoir au moins une place.
- `Booked ≥ 0` — aucune réservation ne peut être négative.
- `Booked ≤ Capacity` — places réservées ≤ capacité totale.

**Invariants actuellement violables** : deux sur trois. `Book` garantit `Booked ≤ Capacity` après appel (via `ErrCapacityExceeded`). En revanche, aucune validation à la construction n'empêche `Capacity ≤ 0` ou `Booked < 0` dans un `Slot` initial.

**Données absentes** :
- Durée ou heure de fin du créneau — `Start` seul ne suffit pas à décrire un créneau complet.
- Identité du réservant — aucun lien entre une réservation et un client.
- État du créneau (ouvert, complet, annulé) — binaire sur "disponible" seulement.
- Prix, monnaie, paiement associé.
- Historique des modifications ou audit trail.

## État de la fonctionnalité — lacunes

| Fonctionnalité attendue pour "réservation d'activités nautiques" | État | Raison |
|---|---|---|
| Catalogue d'activités | ✗ Absent | Pas d'entité Activity en base, pas de liste énumérée |
| Créneaux planifiés | ⚠ Partiel | `Slot` existe mais sans persistance, pas de requête par date/activité |
| Réservation (core) | ✓ Implémentée | `Book` valide capacité et n, retourne `(Slot, error)` |
| Annulation | ✗ Absent | `Book(s, -n)` est rejeté ; une fonction `Cancel` reste à créer |
| Paiement | ✗ Absent | Aucune entité, aucun flux |
| Authentification du réservant | ✗ Absent | Aucun concept de client/utilisateur |
| Notifications | ✗ Absent | Aucun système de message |
| Persistance | ✗ Absent | Tout existe en mémoire, disparaît à la fin du scope |

**Conclusion** : le projet couvre ~5% du périmètre annoncé. C'est attendu pour un "pilote de test" — mais cette lacune doit être explicite pour la planification.

## Questions non résolues

- **Annulation** : `Book` rejette désormais `n ≤ 0` — une fonction `Cancel(s Slot, n int) (Slot, error)` dédiée sera-t-elle ajoutée ?
- **Persistance** : où les réservations sont-elles stockées ? Base de données, fichier, in-memory cache ? Non visible ici.
- **Clients** : y aura-t-il une entité `Client` ou `User` pour lier une réservation à un réservant ? Aujourd'hui aucune notion.
- **Concurrence** : `IsAvailable → Book` doit-il être atomique ? Aucune protection actuellement.
