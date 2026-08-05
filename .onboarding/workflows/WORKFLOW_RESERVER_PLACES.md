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
Permettre à un code appelant d'enregistrer la réservation de `n` places sur un créneau d'activité, en retournant un nouveau `Slot` reflétant l'état après réservation. Dans la logique métier visée (réservation d'activités nautiques en Polynésie française), c'est l'opération terminale du cycle de vie d'un créneau disponible : on consulte la disponibilité (`IsAvailable`), puis on réserve (`Book`). À ce stade du pilote, `Book` est la seule opération qui modifie l'état d'un créneau — sans persistance ni validation de capacité.

## Acteurs
- **Appelant Go** : tout code qui importe `github.com/arthurfromtahiti/shift-pilot-go/internal/booking` et passe un `Slot` et un entier `n` à `Book`. Aucun acteur humain ni système externe dans le code actuel.

## Points d'entrée
- `Book(s Slot, n int) Slot` — `internal/booking/booking.go:25`

## Étapes principales

1. **L'appelant construit ou obtient un `Slot` et choisit `n`** — le `Slot` doit porter `Capacity` et `Booked` actuels. `n` est le nombre de places à réserver. Aucune validation sur `n` (peut être 0, négatif, ou supérieur à la capacité restante).

2. **`Book` incrémente `Booked` de `n`** — `s.Booked += n` (`internal/booking/booking.go:26`). La fonction opère sur une **copie locale** du `Slot` reçu par valeur — le `Slot` original de l'appelant n'est pas modifié.

3. **`Book` retourne le `Slot` mis à jour** — `return s` (`internal/booking/booking.go:27`). Le créneau retourné est une nouvelle valeur avec `Booked` incrémenté ; les autres champs (`ID`, `Activity`, `Start`, `Capacity`) sont inchangés.

4. **L'appelant est responsable de la suite** — persistance (absente dans le pilote), vérification de disponibilité préalable (non effectuée par `Book`), propagation à d'autres systèmes (non prévue dans le code actuel). Tout ce qui suit l'appel à `Book` est hors périmètre du code existant.

## Règles métier

- **`Book` n'effectue aucune vérification de disponibilité** : l'incrément `s.Booked += n` est exécuté sans contrôler `IsAvailable(s)` ni `Remaining(s)` (`internal/booking/booking.go:25-28`). Un créneau complet ou déjà sur-réservé peut être « réservé » davantage sans erreur.
- **La réservation est non destructive pour le `Slot` source** : `Book` reçoit `s Slot` par valeur (copie) et retourne une nouvelle valeur — pas de mutation en place, pas de pointeur. L'appelant doit utiliser la valeur de retour pour obtenir l'état mis à jour.
- **`n` peut être négatif (libération implicite)** : `n` est un `int` sans validation. `Book(s, -2)` décrémentera `Booked` de 2 — simulant une annulation. Ce n'est pas documenté dans le code ni testé (`internal/booking/booking_test.go` teste uniquement `n=2` positif).

## Données

- **`Slot`** (`internal/booking/booking.go:6-12`) : créneau d'activité, seule entité de données.
  - **En entrée** : `Booked int` (état courant) + `Capacity int` (non lu par `Book`, mais nécessaire pour que `Remaining`/`IsAvailable` soient cohérents après l'appel)
  - **En sortie** : nouveau `Slot` avec `Booked` incrémenté de `n` ; tous les autres champs inchangés
- **`n int`** : nombre de places à réserver (paramètre libre, non validé)

## Intégrations

Aucune intégration externe explicite visible. Pas de persistance, pas de message broker, pas d'appel réseau dans le code actuel. Si `Book` doit à terme déclencher une confirmation, enregistrer en base, ou notifier, ces connexions ne sont pas présentes.

## Risques

- **Sur-réservation silencieuse** : `Book` n'appelle pas `IsAvailable` avant d'incrémenter `Booked` (`internal/booking/booking.go:26`). Un appelant qui appelle `Book` sans vérification préalable produit un `Slot` avec `Booked > Capacity` — `Remaining` retourne alors un entier négatif, `IsAvailable` retourne `false`, sans qu'aucune erreur ne soit levée. Impact direct : intégrité des données compromise sans signal. Scénario concret : `Book(Slot{Capacity:10, Booked:10}, 1)` retourne `Slot{Booked:11}` sans erreur.
- **`n` négatif accepté silencieusement** : `Book(s, -3)` retourne un `Slot` avec `Booked` décrémenté, pouvant passer en dessous de 0. Comportement non documenté, non testé, potentiellement trompeur si utilisé par inadvertance. Impact : `Remaining` retournerait une valeur supérieure à `Capacity`.
- **Aucune persistance** : la valeur retournée par `Book` est une valeur Go en mémoire — si l'appelant ne la conserve pas ou ne la persiste pas, la réservation est perdue à la fin du scope. Impact futur critique lors de l'ajout d'un serveur : oublier de sauvegarder le `Slot` retourné équivaut à ne pas réserver.
- **Absence d'atomicité lecture-puis-écriture** : le schéma `IsAvailable(s)` → `Book(s, n)` n'est pas atomique. Entre les deux appels, un autre appelant concurrent pourrait modifier l'état du `Slot` (si un serveur ou goroutine est introduit à l'avenir). Impact futur : race condition TOCTOU classique.

## Questions ouvertes

- La validation de disponibilité (`IsAvailable` avant `Book`) est-elle attendue côté appelant, ou sera-t-elle intégrée directement à `Book` dans une prochaine itération ?
- `n` négatif est-il un cas d'usage volontaire (annulation/libération de places) ou une omission du pilote ? Aucun test ni commentaire n'en précise le rôle.
- Quelle couche sera responsable de la persistance du `Slot` retourné par `Book` ? Pas visible dans le code actuel.
- Le choix du passage par valeur (plutôt que pointeur) est-il un invariant architectural du projet ou une simplification du pilote ?

## Preuves

- `internal/booking/booking.go` — lu intégralement (`VÉRIFIÉ_CODE`) : type `Slot` (lignes 6-12), `Book` (lignes 25-28)
- `internal/booking/booking_test.go` — lu intégralement (`VÉRIFIÉ_CODE`) : `TestBook` (lignes 24-29) — assertions lues mais **non exécutées** (statut d'exécution `INCONNU`)
- `go.mod` — lu (`VÉRIFIÉ_CODE`) : module `github.com/arthurfromtahiti/shift-pilot-go`, Go 1.21, aucune dépendance externe
