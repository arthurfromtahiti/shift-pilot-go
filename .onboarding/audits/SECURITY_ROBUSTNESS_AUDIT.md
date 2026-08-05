# Sécurité & Robustesse — Audit

> Confiance : medium — code entièrement lu (`VÉRIFIÉ_CODE`). La surface de sécurité est nulle dans le code actuel (pas de serveur, pas d'entrée utilisateur, pas de réseau), mais elle le restera seulement jusqu'à ce qu'une couche de transport soit ajoutée. Cet audit vaut donc autant comme bilan présent que comme liste de garde-fous à poser avant l'évolution.

## Compréhension globale

Le code actuel n'expose aucune surface d'attaque directe : il n'y a ni serveur, ni route HTTP, ni entrée utilisateur, ni base de données, ni appel réseau, ni secret en clair. Les fonctions de `internal/booking` ne peuvent pas être exploitées depuis l'extérieur sans passer par un appelant Go. Les risques de sécurité observables sont donc tous de nature préventive — ils concernent ce qui doit être conçu correctement avant d'exposer ce package.

## Résumé exécutif

Aucun secret dans le code (`VÉRIFIÉ_CODE`). Aucun appel réseau ni SQL (`VÉRIFIÉ_CODE`). Le risque de sécurité le plus concret est fonctionnel : `Book` n'effectue aucune validation de ses paramètres et peut produire silencieusement un état incohérent (`Booked > Capacity`, `n` négatif). Ce n'est pas exploitable aujourd'hui depuis l'extérieur, mais ce sera le premier vecteur de manipulation si une couche HTTP est ajoutée sans ajouter la validation correspondante. Robustesse : les fonctions ne retournent jamais d'erreur, il n'y a aucun mécanisme de gestion de panique, et la suite de tests ne couvre aucun cas limite.

## Constats détaillés

**`VÉRIFIÉ_CODE` — Aucun secret ni credential dans le source.** Recherche effectuée sur `booking.go` et `booking_test.go` : aucune clé, token, mot de passe, URL de service tiers ou adresse IP n'est présent. La seule constante textuelle est `"Plongée"` dans la fixture de test (`booking_test.go:9`), qui est une donnée de test sans valeur de secret.

**`VÉRIFIÉ_CODE` — Aucune dépendance externe, aucun appel réseau.** Le `go.mod` n'importe que la bibliothèque standard (absence de bloc `require` dans `go.mod`). Le seul import de `booking.go` est `"time"` (`booking.go:3`). Il n'y a pas de `net/http`, `database/sql`, `os/exec`, ou `reflect` — le vecteur de supply-chain et d'injection est nul dans l'état actuel.

**`VÉRIFIÉ_CODE` — `Book` n'effectue aucune validation d'entrée.** `Book(s Slot, n int) Slot` incrémente `s.Booked` de `n` sans aucune vérification (`booking.go:26`). En conséquence : (1) `n` peut être négatif — `Book(s, -3)` décrémente `Booked` de 3, simulant une annulation non documentée ni testée ; (2) `n` peut être zéro — no-op silencieux ; (3) `n` peut dépasser la capacité restante — `Book` retourne un `Slot` avec `Booked > Capacity` sans erreur. Aucun de ces cas ne lève de `panic`, ne retourne d'erreur, et n'est couvert par un test (`booking_test.go:24-29` teste uniquement `n=2` sur `Capacity=10, Booked=4`).

**`VÉRIFIÉ_CODE` — Aucun type d'erreur — aucun signal de défaut.** Les signatures `Remaining(s Slot) int`, `IsAvailable(s Slot) bool`, `Book(s Slot, n int) Slot` ne retournent jamais `error` (`booking.go:15, 20, 25`). Un état invalide (sur-réservation, `n` négatif) est retourné comme valeur nominale, indiscernable d'un résultat correct du point de vue du type.

**`VÉRIFIÉ_CODE` — Aucun mécanisme de gestion de panique.** Aucun `defer`, `recover`, ou gestion d'erreur de runtime dans le package. Si un appelant fournit un `Slot` avec des valeurs extrêmes (ex. `Capacity = math.MaxInt`, `Booked = -math.MinInt`), l'arithmétique `Capacity - Booked` peut provoquer un débordement d'entier silencieux en Go (comportement défini par le compilateur, mais résultat numériquement incorrect).

**`HYPOTHÈSE` — La surface de sécurité s'agrandit dès qu'une couche de transport est ajoutée.** Si un serveur HTTP est introduit, chaque paramètre de `Slot` et chaque valeur de `n` deviendront des entrées utilisateur non validées. Sans validation préalable (n > 0, n ≤ Remaining(s), Capacity > 0, etc.), `Book` deviendra un vecteur direct de manipulation d'état côté serveur. Ce risque est aujourd'hui `HYPOTHÈSE` car il n'y a pas de serveur — mais le code ne porte aucun garde-fou préventif.

## Forces

- `VÉRIFIÉ_CODE` : Zéro secret en clair dans le code source (`booking.go`, `booking_test.go`, `go.mod`, `README.md`).
- `VÉRIFIÉ_CODE` : Zéro dépendance externe — pas de supply-chain, pas de CVE transitif possible dans l'état actuel.
- `VÉRIFIÉ_CODE` : Fonctions pures sans état global partagé — aucun état mutable accessible de l'extérieur, ce qui élimine toute une classe de bugs de concurrence en mode mono-goroutine.

## Dettes techniques

- `VÉRIFIÉ_CODE` : `Book` ne valide pas `n` et ne vérifie pas la disponibilité avant d'incrémenter `Booked` (`booking.go:26`) — dette fonctionnelle directement connexe à la robustesse.
- `VÉRIFIÉ_CODE` : Aucun type `error` sur les fonctions — impossible pour un appelant de distinguer un résultat valide d'un résultat issu d'un appel sémantiquement invalide.
- `VÉRIFIÉ_CODE` : La fixture de test `sample()` utilise `time.Now()` (`booking_test.go:9`), ce qui rend les tests non déterministes si le champ `Start` venait à être utilisé dans une comparaison future. Ce n'est pas un problème aujourd'hui (aucune fonction ne lit `Start`), mais c'est une dette latente.

## Zones critiques

- `VÉRIFIÉ_CODE` : `internal/booking/booking.go:25-28` (`Book`) — point unique de mutation d'état sans validation. Si ce code est exposé via une API sans blindage amont, c'est la première cible à auditer.

## Risques

- `VÉRIFIÉ_CODE` (impact futur, défaut prouvé) : Sans validation de `n`, une couche HTTP mal conçue exposerait directement la sur-réservation silencieuse à un attaquant (`booking.go:26`). Impact : intégrité des données compromise, impossibilité de distinguer un état valide d'un état manipulé.
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
