# WORKFLOW_RESERVER_PLACES — Réserver des places sur un créneau

## Classification
- **Type** : `technical_flow`
- **Sous-type** : mutation d'état — API de package Go (bibliothèque, pas de serveur)
- **Visibilité** : `technical` — aucun point d'entrée utilisateur ou HTTP ; exposé uniquement comme API de package Go
- **Acteur principal** : code appelant (autre package Go)
- **Acteurs** : appelant Go (seul acteur connu — pas de couche HTTP, pas d'utilisateur humain)
- **Criticité** : Haute — cœur de la raison d'être du projet (réservation d'activités nautiques) ; absence de garde de capacité expose au risque de sur-réservation silencieuse
- **Confiance** : medium — fichiers lus intégralement (`VÉRIFIÉ_CODE`) ; tests lus mais non exécutés (toolchain Go absente, statut d'exécution `INCONNU`)
- **Justification** : `Book` est entièrement visible dans `internal/booking/booking.go` (4 lignes). Type `technical_flow` retenu pour les mêmes raisons que `WORKFLOW_CONSULTER_DISPONIBILITE` : aucun point d'entrée HTTP, route ou interface utilisateur dans le code actuel.

## Objectif
Permettre à un code appelant d'enregistrer la réservation de `n` places sur un créneau d'activité, en retournant un nouveau `Slot` reflétant l'état après réservation. Dans la logique métier visée (réservation d'activités nautiques en Polynésie française), c'est l'opération terminale du cycle de vie d'un créneau disponible : on consulte la disponibilité (`IsAvailable`), puis on réserve (`Book`). À ce stade du pilote, `Book` est la seule opération qui modifie l'état d'un créneau — sans persistance, mais avec validation de capacité et de `n`.

## Acteurs
- **Appelant Go** : tout code qui importe `github.com/arthurfromtahiti/shift-pilot-go/internal/booking` et passe un `Slot` et un entier `n` à `Book`. Aucun acteur humain ni système externe dans le code actuel.

## Points d'entrée
- `Book(s Slot, n int) (Slot, error)` — `internal/booking/booking.go:36`

## Étapes principales

1. **L'appelant construit ou obtient un `Slot` et choisit `n`** — le `Slot` doit porter `Capacity` et `Booked` actuels. `n` est le nombre de places à réserver. `n` doit être positif et ≤ `Remaining(s)`.

2. **`Book` valide les pré-conditions** — `n ≤ 0` → retourne `(s, ErrInvalidBookingCount)` ; `n > Remaining(s)` → retourne `(s, ErrCapacityExceeded)` (`internal/booking/booking.go:37-41`). La fonction opère sur une **copie locale** du `Slot` reçu par valeur — le `Slot` original de l'appelant n'est pas modifié.

3. **`Book` incrémente `Booked` de `n` et retourne `(Slot, nil)`** — `s.Booked += n; return s, nil` (`internal/booking/booking.go:42-44`). Le créneau retourné est une nouvelle valeur avec `Booked` incrémenté ; les autres champs (`ID`, `Activity`, `Start`, `Capacity`) sont inchangés.

4. **L'appelant est responsable de la suite** — persistance (absente dans le pilote), propagation à d'autres systèmes (non prévue dans le code actuel). Tout ce qui suit l'appel à `Book` est hors périmètre du code existant.

## Règles métier

- **`Book` valide ses pré-conditions** : `n ≤ 0` → `ErrInvalidBookingCount` ; `n > Remaining(s)` → `ErrCapacityExceeded` (`internal/booking/booking.go:37-41`). Un créneau complet ou une valeur `n` invalide provoque un retour d'erreur explicite — la sur-réservation silencieuse est impossible.
- **La réservation est non destructive pour le `Slot` source** : `Book` reçoit `s Slot` par valeur (copie) et retourne une nouvelle valeur — pas de mutation en place, pas de pointeur. L'appelant doit utiliser la valeur de retour pour obtenir l'état mis à jour.
- **`n` négatif ou nul est rejeté** : `Book(s, -2)` et `Book(s, 0)` retournent `ErrInvalidBookingCount`. L'annulation nécessitera une fonction dédiée.

## Données

- **`Slot`** (`internal/booking/booking.go:6-12`) : créneau d'activité, seule entité de données.
  - **En entrée** : `Booked int` (état courant) + `Capacity int` (non lu par `Book`, mais nécessaire pour que `Remaining`/`IsAvailable` soient cohérents après l'appel)
  - **En sortie** : `(Slot, error)` — nouveau `Slot` avec `Booked` incrémenté de `n` si succès, erreur non nulle sinon ; tous les autres champs du `Slot` inchangés
- **`n int`** : nombre de places à réserver (doit être > 0 et ≤ `Remaining(s)`)

## Intégrations

Aucune intégration externe explicite visible. Pas de persistance, pas de message broker, pas d'appel réseau dans le code actuel. Si `Book` doit à terme déclencher une confirmation, enregistrer en base, ou notifier, ces connexions ne sont pas présentes.

## Risques

- ~~**Sur-réservation silencieuse**~~ — **RÉSOLU** : `Book(Slot{Capacity:10, Booked:10}, 1)` retourne désormais `(Slot{}, ErrCapacityExceeded)`. L'intégrité des données est garantie par `Book` elle-même.
- ~~**`n` négatif accepté silencieusement**~~ — **RÉSOLU** : `Book(s, -3)` retourne `ErrInvalidBookingCount`.
- **Aucune persistance** : la valeur retournée par `Book` est une valeur Go en mémoire — si l'appelant ne la conserve pas ou ne la persiste pas, la réservation est perdue à la fin du scope. Impact futur critique lors de l'ajout d'un serveur : oublier de sauvegarder le `Slot` retourné équivaut à ne pas réserver.
- **Absence d'atomicité lecture-puis-écriture** : le schéma `IsAvailable(s)` → `Book(s, n)` n'est pas atomique. Entre les deux appels, un autre appelant concurrent pourrait modifier l'état du `Slot` (si un serveur ou goroutine est introduit à l'avenir). Impact futur : race condition TOCTOU classique.

## Questions ouvertes

- `Book` rejette désormais `n ≤ 0` — une fonction `Cancel` dédiée sera-t-elle ajoutée pour l'annulation de places ?
- Quelle couche sera responsable de la persistance du `Slot` retourné par `Book` ? Pas visible dans le code actuel.
- Le choix du passage par valeur (plutôt que pointeur) est-il un invariant architectural du projet ou une simplification du pilote ?

## Preuves

- `internal/booking/booking.go` — lu intégralement (`VÉRIFIÉ_CODE`) : type `Slot` (lignes 15-21), erreurs sentinelles (lignes 9, 12), `Book` (lignes 36-45)
- `internal/booking/booking_test.go` — lu intégralement (`VÉRIFIÉ_CODE`) : `TestBook` (lignes 24-32), `TestBookCapacityExceeded` (34-39), `TestBookExactCapacity` (41-49), `TestBookZero` (51-56), `TestBookNegative` (58-63)
- `go.mod` — lu (`VÉRIFIÉ_CODE`) : module `github.com/arthurfromtahiti/shift-pilot-go`, Go 1.21, aucune dépendance externe
