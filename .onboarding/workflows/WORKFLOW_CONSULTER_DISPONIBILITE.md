# WORKFLOW_CONSULTER_DISPONIBILITE — Consulter la disponibilité d'un créneau

## Classification
- **Type** : `technical_flow`
- **Sous-type** : calcul de disponibilité — API de package Go (bibliothèque, pas de serveur)
- **Visibilité** : `technical` — aucun point d'entrée utilisateur ou HTTP ; exposé uniquement comme API de package Go
- **Acteur principal** : code appelant (autre package Go)
- **Acteurs** : appelant Go (seul acteur connu — pas de couche HTTP, pas d'utilisateur humain)
- **Criticité** : Haute — condition de garde nécessaire avant toute réservation ; toute erreur produit des sur-réservations silencieuses
- **Confiance** : medium — les trois fichiers du package ont été lus intégralement ; les tests ont été lus (`VÉRIFIÉ_CODE`) mais **non exécutés** (toolchain Go absente, statut d'exécution `INCONNU`)
- **Justification** : Les deux fonctions (`Remaining`, `IsAvailable`) sont entièrement visibles dans `internal/booking/booking.go`. Le type `technical_flow` est retenu car il n'existe aucun point d'entrée HTTP, aucune route, aucune commande CLI, et aucun utilisateur humain qui déclenche directement ces fonctions dans le code actuel — elles sont une API de bibliothèque Go.

## Objectif
Permettre à un code appelant de connaître le nombre de places restantes sur un créneau d'activité, et de savoir si ce créneau peut encore accepter des réservations. C'est la condition préalable à toute opération de réservation : on consulte la disponibilité avant de tenter de réserver. À ce stade du pilote, ces fonctions constituent la seule capacité métier observable du côté de la lisibilité de l'état d'un créneau.

## Acteurs
- **Appelant Go** : tout code qui importe `github.com/arthurfromtahiti/shift-pilot-go/internal/booking` et passe un `Slot` à `Remaining` ou `IsAvailable`. Aucun acteur humain ni système externe identifié dans le code actuel.

## Points d'entrée
- `Remaining(s Slot) int` — `internal/booking/booking.go:15`
- `IsAvailable(s Slot) bool` — `internal/booking/booking.go:20`

Ces deux fonctions sont les seuls points d'entrée pour consulter la disponibilité. `IsAvailable` délègue à `Remaining` (appel interne direct).

## Étapes principales

1. **L'appelant construit ou obtient un `Slot`** — le `Slot` doit porter des valeurs de `Capacity` (capacité totale) et `Booked` (places déjà réservées) cohérentes. Aucune validation d'entrée n'est effectuée par les fonctions (`internal/booking/booking.go:15-17, 20-22`).

2. **Calcul des places restantes via `Remaining`** — retourne `s.Capacity - s.Booked` (`internal/booking/booking.go:16`). Opération arithmétique pure, sans effet de bord.

3. *(Optionnel, selon le besoin de l'appelant)* **Test de disponibilité via `IsAvailable`** — retourne `Remaining(s) > 0` (`internal/booking/booking.go:21`). `IsAvailable` n'est qu'un prédicat sur le résultat de `Remaining` ; appeler l'une n'empêche pas d'appeler l'autre.

4. **Retour du résultat à l'appelant** — `Remaining` retourne un `int` (nombre de places restantes, peut être négatif si sur-réservé) ; `IsAvailable` retourne un `bool`. Aucun état n'est modifié.

## Règles métier

- **Places restantes = Capacité − Réservées** : `Remaining(s) = s.Capacity - s.Booked` (`internal/booking/booking.go:16`). La règle est symétrique : un créneau vide (`Booked = 0`) retourne `Capacity` ; un créneau plein retourne `0`.
- **Disponible si et seulement si il reste au moins une place** : `IsAvailable(s) = Remaining(s) > 0` (`internal/booking/booking.go:21`). Un créneau est indisponible dès que `Remaining ≤ 0`.
- **Résultat potentiellement négatif non bloqué** : si `Booked > Capacity` (sur-réservation possible via `Book` — voir `WORKFLOW_RESERVER_PLACES.md`), `Remaining` retourne une valeur négative et `IsAvailable` retourne `false`. Aucune garde contre cette situation dans ces deux fonctions.

## Données

- **`Slot`** (`internal/booking/booking.go:6-12`) : créneau d'activité, support de toutes les opérations. Champs utilisés par ce workflow :
  - `Capacity int` — capacité totale du créneau (nombre de places offertes)
  - `Booked int` — nombre de places déjà réservées
  - `ID int`, `Activity string`, `Start time.Time` — présents sur le `Slot` mais non lus par `Remaining` ni `IsAvailable` ; portés passivement

## Intégrations

Aucune intégration externe explicite visible. Le package ne dépend que de la bibliothèque standard Go (`time`). Aucun appel réseau, base de données, ou système tiers.

## Risques

- **`Remaining` peut retourner un entier négatif sans avertissement** : si l'appelant a sur-réservé via `Book` (voir `WORKFLOW_RESERVER_PLACES.md`), `Capacity - Booked < 0`. Aucune garde ni `panic` dans `Remaining` (`internal/booking/booking.go:16`). Impact : l'appelant reçoit un entier négatif sans signal d'erreur — comportement silencieusement incorrect si non anticipé.
- **Aucune validation des champs `Capacity` ou `Booked`** : un `Slot` avec `Capacity = -1` ou `Booked = -5` est accepté sans erreur. Ces cas ne sont pas testés (`internal/booking/booking_test.go` ne couvre que le cas nominal `Capacity=10, Booked=4`). Impact : résultats arithmétiquement faux mais silencieux.
- **Absence de concurrence** : les fonctions sont pures (pas de mutation, pas de pointeur), donc sûres en lecture concurrente au sens Go. En revanche, le schéma « lire `IsAvailable` puis appeler `Book` » n'est pas atomique : aucun mécanisme de verrouillage entre ces deux appels (`internal/booking/booking.go:20-28`). Impact futur : si un serveur ou goroutine concurrente est introduit, une race condition classique TOCTOU (vérification-puis-action) devient possible.

## Questions ouvertes

- `Capacity` ou `Booked` peuvent-ils être négatifs en pratique ? Existe-t-il une validation amont (côté appelant ou couche future) qui l'interdit ? Non visible dans le code.
- `IsAvailable` est-il destiné à être appelé systématiquement avant `Book`, ou `Book` sera-t-il enrichi d'une garde interne à l'avenir ?
- La structure `Slot` est actuellement immuable (valeur, pas pointeur) : est-ce un choix architectural délibéré ou une simplification du pilote ?

## Preuves

- `internal/booking/booking.go` — lu intégralement (`VÉRIFIÉ_CODE`) : type `Slot` (lignes 6-12), `Remaining` (lignes 15-17), `IsAvailable` (lignes 20-22)
- `internal/booking/booking_test.go` — lu intégralement (`VÉRIFIÉ_CODE`) : `TestRemaining` (lignes 12-16), `TestIsAvailable` (lignes 18-22) — assertions lues mais **non exécutées** (statut d'exécution `INCONNU`)
- `go.mod` — lu (`VÉRIFIÉ_CODE`) : module `github.com/arthurfromtahiti/shift-pilot-go`, Go 1.21, aucune dépendance externe
