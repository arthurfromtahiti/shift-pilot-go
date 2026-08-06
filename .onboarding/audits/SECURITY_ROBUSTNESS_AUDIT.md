# Sécurité & Robustesse — Audit

> Confiance : medium — code entièrement lu (`VÉRIFIÉ_CODE`). La surface de sécurité est nulle dans le code actuel (pas de serveur, pas d'entrée utilisateur, pas de réseau), mais elle le restera seulement jusqu'à ce qu'une couche de transport soit ajoutée. Cet audit vaut donc autant comme bilan présent que comme liste de garde-fous à poser avant l'évolution.

## Compréhension globale

Le code actuel n'expose aucune surface d'attaque directe : il n'y a ni serveur, ni route HTTP, ni entrée utilisateur, ni base de données, ni appel réseau, ni secret en clair. Les fonctions de `internal/booking` ne peuvent pas être exploitées depuis l'extérieur sans passer par un appelant Go. Les risques de sécurité observables sont donc tous de nature préventive — ils concernent ce qui doit être conçu correctement avant d'exposer ce package.

## Résumé exécutif

Aucun secret dans le code (`VÉRIFIÉ_CODE`). Aucun appel réseau ni SQL (`VÉRIFIÉ_CODE`). Les risques de sécurité fonctionnels ont été résorbés : `Book` valide désormais ses paramètres et rejette `n ≤ 0` avec `ErrInvalidBookingCount` et `n > Remaining(s)` avec `ErrCapacityExceeded`. Aucun état incohérent ne peut plus être produit par `Book`. Robustesse : `Book` retourne maintenant un `error`, la suite de tests couvre les cas limites critiques (`TestBookZero`, `TestBookNegative`, `TestBookCapacityExceeded`, `TestBookExactCapacity`).

## Constats détaillés

**`VÉRIFIÉ_CODE` — Aucun secret ni credential dans le source.** Recherche effectuée sur `booking.go` et `booking_test.go` : aucune clé, token, mot de passe, URL de service tiers ou adresse IP n'est présent. La seule constante textuelle est `"Plongée"` dans la fixture de test (`booking_test.go:9`), qui est une donnée de test sans valeur de secret.

**`VÉRIFIÉ_CODE` — Aucune dépendance externe, aucun appel réseau.** Le `go.mod` n'importe que la bibliothèque standard (absence de bloc `require` dans `go.mod`). Le seul import de `booking.go` est `"time"` (`booking.go:3`). Il n'y a pas de `net/http`, `database/sql`, `os/exec`, ou `reflect` — le vecteur de supply-chain et d'injection est nul dans l'état actuel.

**`VÉRIFIÉ_CODE` — `Book` valide désormais ses entrées.** `Book(s Slot, n int) (Slot, error)` valide les pré-conditions (`booking.go:37-41`) : (1) `n ≤ 0` → retourne `ErrInvalidBookingCount` (`booking.go:37-39`), rejetant les valeurs négatives ou nulles ; (2) `n > Remaining(s)` → retourne `ErrCapacityExceeded` (`booking.go:40-41`), empêchant la sur-réservation. Tous ces cas sont couverts par des tests (`TestBookZero` ligne 51, `TestBookNegative` ligne 58, `TestBookCapacityExceeded` ligne 34, `TestBookExactCapacity` ligne 41).

**`VÉRIFIÉ_CODE` — `Book` retourne désormais un type `error` pour les défauts.** La signature `Book(s Slot, n int) (Slot, error)` (`booking.go:36`) retourne un `error` non nul pour indiquer une tentative invalide (`ErrInvalidBookingCount` si `n ≤ 0`, `ErrCapacityExceeded` si `n > Remaining(s)`). Les signatures `Remaining(s Slot) int` et `IsAvailable(s Slot) bool` demeurent sans `error` car ce sont des fonctions pures de requête, pas de mutation.

**`VÉRIFIÉ_CODE` — Aucun mécanisme de gestion de panique.** Aucun `defer`, `recover`, ou gestion d'erreur de runtime dans le package. Si un appelant fournit un `Slot` avec des valeurs extrêmes (ex. `Capacity = math.MaxInt`, `Booked = -math.MinInt`), l'arithmétique `Capacity - Booked` peut provoquer un débordement d'entier silencieux en Go (comportement défini par le compilateur, mais résultat numériquement incorrect).

**`HYPOTHÈSE` — La surface de sécurité s'agrandit dès qu'une couche de transport est ajoutée.** Si un serveur HTTP est introduit, chaque paramètre de `Slot` et chaque valeur de `n` deviendront des entrées utilisateur non validées. Sans validation préalable (n > 0, n ≤ Remaining(s), Capacity > 0, etc.), `Book` deviendra un vecteur direct de manipulation d'état côté serveur. Ce risque est aujourd'hui `HYPOTHÈSE` car il n'y a pas de serveur — mais le code ne porte aucun garde-fou préventif.

## Forces

- `VÉRIFIÉ_CODE` : Zéro secret en clair dans le code source (`booking.go`, `booking_test.go`, `go.mod`, `README.md`).
- `VÉRIFIÉ_CODE` : Zéro dépendance externe — pas de supply-chain, pas de CVE transitif possible dans l'état actuel.
- `VÉRIFIÉ_CODE` : Fonctions pures sans état global partagé — aucun état mutable accessible de l'extérieur, ce qui élimine toute une classe de bugs de concurrence en mode mono-goroutine.

## Dettes techniques

- ~~`VÉRIFIÉ_CODE` : `Book` ne valide pas `n` et ne vérifie pas la disponibilité avant d'incrémenter `Booked`~~ — **RÉSOLU** : `Book` valide `n > 0` et `n ≤ Remaining(s)` avant tout incrément, retourne `(Slot, error)` (`booking.go:36-45`).
- ~~`VÉRIFIÉ_CODE` : Aucun type `error` sur les fonctions~~ — **RÉSOLU** : `Book` expose `ErrInvalidBookingCount` et `ErrCapacityExceeded` et retourne un `error` non nul en cas de violation (`booking.go:8-12, 36-45`).
- `VÉRIFIÉ_CODE` : La fixture de test `sample()` utilise `time.Now()` (`booking_test.go:9`), ce qui rend les tests non déterministes si le champ `Start` venait à être utilisé dans une comparaison future. Ce n'est pas un problème aujourd'hui (aucune fonction ne lit `Start`), mais c'est une dette latente.

## Zones critiques

- `VÉRIFIÉ_CODE` : `internal/booking/booking.go:25-28` (`Book`) — point unique de mutation d'état sans validation. Si ce code est exposé via une API sans blindage amont, c'est la première cible à auditer.

## Risques

- ~~`VÉRIFIÉ_CODE` (impact futur, défaut prouvé) : Sans validation de `n`, une couche HTTP mal conçue exposerait directement la sur-réservation silencieuse~~ — **RÉSOLU** : `Book` valide les pré-conditions et retourne une erreur explicite. L'intégrité des données est désormais garantie par `Book` elle-même (`booking.go:36-45`).
- `HYPOTHÈSE` : Débordement d'entier possible si `Capacity` ou `Booked` sont proches de `math.MaxInt`/`math.MinInt` — Go ne panique pas sur ces cas, il retourne une valeur fausse (`booking.go:16`). Impact actuel : nul (aucun appelant externe). Impact futur : silencieusement incorrect si des valeurs non validées sont acceptées.
- `HYPOTHÈSE` : Si le README ("Git → staging") implique un déploiement réel sans layer de sécurité (auth, rate limiting, TLS), la surface d'attaque pourrait être entière dès le premier service exposé. Non vérifiable sans accès à l'environnement staging.

## Recommandations priorisées

1. **Ajouter une validation et un retour `error` à `Book`** : `func Book(s Slot, n int) (Slot, error)` retournant une erreur si `n <= 0`, si `n > Remaining(s)`, ou si `Booked < 0` — `internal/booking/booking.go:25`.
2. **Définir des constructeurs validants pour `Slot`** (ex. `func NewSlot(capacity int, ...) (Slot, error)`) pour empêcher la création d'un `Slot` avec `Capacity <= 0` — `internal/booking/booking.go`.
3. **Remplacer `time.Now()` dans la fixture de test** par une heure fixe (`time.Date(...)`) pour garantir la reproductibilité des tests si `Start` devient fonctionnel — `internal/booking/booking_test.go:9`.
4. **Documenter explicitement que `Book` ne doit jamais être appelé sans vérification préalable de `IsAvailable`** — au minimum un commentaire dans le code en attendant la validation interne — `internal/booking/booking.go:24`.

## Questions ouvertes

- L'environnement "staging" mentionné dans le README — est-il exposé sur Internet ? Avec quelle authentification ? Non visible dans ce dépôt.
- Un dépassement de `Capacity` est-il considéré comme une erreur métier bloquante (à rejeter) ou comme une règle métier acceptable (liste d'attente) ? La distinction détermine si `Book` doit retourner une erreur ou un état spécial.
- `n` négatif est-il un cas d'usage légal (annulation) ou une omission ? Si c'est un cas légal, il mérite une fonction dédiée (`Cancel`, `Release`) plutôt que d'être un sous-produit implicite de `Book`.
