# Points chauds du code — Audit

> Confiance : high — le corpus est de 28 lignes de code de production et 29 lignes de tests, intégralement lus. Il n'y a pas de point chaud au sens classique (fichier volumineux, couplage fort, chaîne d'appels profonde) ; les risques sont qualitatifs, pas volumétriques. Le niveau de confiance `high` reflète l'exhaustivité de la lecture, pas la gravité des constats.

## Compréhension globale

Le dépôt ne contient qu'un seul fichier de production : `internal/booking/booking.go` (28 lignes). Il n'y a pas de point chaud par volume ou complexité cyclomatique — chaque fonction tient en 3-4 lignes. Les points chauds ici sont qualitatifs : `Book` est la seule fonction qui modifie l'état et c'est elle qui concentre tous les risques de comportement incorrect. La suite de tests ne couvre que le chemin nominal, laissant les cas critiques sans filet.

## Résumé exécutif

`Book` (`booking.go:25-28`) est le seul point chaud du projet — non par sa taille, mais par sa responsabilité et ses lacunes : elle est la seule fonction à retourner un état modifié, elle ne valide aucun paramètre, et aucun test ne couvre ses comportements limites. Tout futur appelant (HTTP handler, CLI, autre service Go) qui invoque `Book` sans vérification préalable produit une sur-réservation silencieuse. Les deux autres fonctions (`Remaining`, `IsAvailable`) sont sans risque dans leur forme actuelle. Le fichier de tests (`booking_test.go`) constitue un second point chaud : il donne une illusion de couverture (3 tests passants) sans tester les cas qui comptent.

## Constats détaillés

**`VÉRIFIÉ_CODE` — `Book` : la seule mutation, sans garde.** La fonction `Book` (`booking.go:25-28`) est la seule à produire un état différent de son entrée. Elle incrémente `s.Booked` de `n` et retourne le `Slot` modifié :

```go
func Book(s Slot, n int) Slot {
    s.Booked += n
    return s
}
```

Aucun appel à `IsAvailable`, aucune vérification de `n > 0`, aucune vérification de `n ≤ Remaining(s)`. La complexité cyclomatique est 1 (aucune branche), ce qui est techniquement simple — mais cette simplicité est trompeuse : elle masque une absence de logique de garde, pas une logique triviale correctement implémentée. Un senior regarderait ici en premier avant tout refactoring ou ajout d'un transport.

**`VÉRIFIÉ_CODE` — `Remaining` et `IsAvailable` : fonctions sans risque dans l'état actuel.** `Remaining(s Slot) int` retourne `s.Capacity - s.Booked` (`booking.go:16`), `IsAvailable(s Slot) bool` retourne `Remaining(s) > 0` (`booking.go:21`). Les deux sont pures, déterministes, sans effet de bord. Leur seul risque est passif : elles acceptent et calculent sur des valeurs incohérentes (Booked > Capacity → Remaining négatif) sans jamais le signaler.

**`VÉRIFIÉ_CODE` — `booking_test.go` : couverture nominale uniquement.** Les trois tests (`booking_test.go:12-29`) utilisent tous la même fixture `sample()` : `Slot{Capacity: 10, Booked: 4}`. Cas couverts : `Remaining` renvoie 6, `IsAvailable` retourne `true`, `Book(s, 2)` donne `Booked=6`. Cas non couverts (lus dans le source, non présents dans les tests) :

| Cas | Risque | Couvert |
|---|---|---|
| `Book(s, 0)` | No-op silencieux | Non |
| `Book(s, -2)` | Décrémente Booked (annulation implicite) | Non |
| `Book(s, 7)` sur `Capacity=10, Booked=4` | Sur-réservation (`Booked=11 > Capacity=10`) | Non |
| `IsAvailable` quand `Booked == Capacity` | Créneau exactement plein | Non |
| `Remaining` quand `Booked > Capacity` | Valeur négative | Non |
| `Slot{Capacity: 0}` | Division implicite de sens | Non |

**`VÉRIFIÉ_CODE` — Absence de goroutines et de concurrence.** Aucun `go func`, aucun `sync.Mutex`, aucun `chan` dans le package (`booking.go:1-28`). Les fonctions sont sûres en lecture concurrente (valeurs, pas de pointeurs partagés). En revanche, le schéma `if IsAvailable(s) { Book(s, n) }` n'est pas atomique — c'est une race condition TOCTOU latente dès qu'un serveur concurrent est ajouté.

**`VÉRIFIÉ_CODE` — Aucun `panic`, aucun `recover`, aucun `log`.** Le package ne panique jamais, ne trace jamais, ne loggue jamais (`booking.go:1-28`). C'est correct pour une bibliothèque (une lib ne doit pas décider de la politique de log de son appelant), mais cela signifie que les erreurs silencieuses (sur-réservation, valeurs incohérentes) passent complètement inaperçues sans instrumentation externe.

## Forces

- `VÉRIFIÉ_CODE` : Complexité cyclomatique de 1 pour chaque fonction — aucun `if`/`for`/`switch` imbriqué, code lisible (`booking.go:15-28`).
- `VÉRIFIÉ_CODE` : Aucun couplage entre fonctions sauf `IsAvailable → Remaining` (`booking.go:21`), ce qui est un couplage intentionnel et cohérent.
- `VÉRIFIÉ_CODE` : Zéro état global — aucune variable de package, aucune variable d'initialisation `init()`. Impossible de "polluer" l'état entre deux appels.

## Dettes techniques

- `VÉRIFIÉ_CODE` : `Book` sans validation (`booking.go:25-28`) — dette fonctionnelle qui devient risque opérationnel dès que la fonction est exposée à des paramètres non filtrés.
- `VÉRIFIÉ_CODE` : `booking_test.go` — 3 tests sur un seul fixture, aucun cas limite — donne une fausse impression de couverture. Un futur contributeur qui voit "3 tests passants" peut avoir confiance à tort.

## Zones critiques

- `VÉRIFIÉ_CODE` : `internal/booking/booking.go:25-28` (`Book`) — cœur de la logique de mutation, sans garde. Si cette fonction est exposée via un handler HTTP sans validation amont, tout appel avec `n` quelconque ou sans vérification préalable de `IsAvailable` produit un état corrompu silencieusement.
- `VÉRIFIÉ_CODE` : `internal/booking/booking_test.go:8-10` (`sample()`) — fixture unique avec `time.Now()`. Si le champ `Start` devient fonctionnel dans les comparaisons, ces tests deviendront non déterministes sans modification.

## Risques

- `VÉRIFIÉ_CODE` (défaut prouvé, impact immédiat si exposé) : `Book(Slot{Capacity: 10, Booked: 10}, 1)` retourne `Slot{Booked: 11}` sans erreur (`booking.go:26`). L'état résultant est sémantiquement invalide et silencieux.
- `HYPOTHÈSE` : Dès qu'un serveur concurrent est introduit, le schéma `IsAvailable → Book` sans verrou devient une race condition TOCTOU. Ce pattern est probable dans tout handler HTTP REST basique.

## Recommandations priorisées

1. **Ajouter une garde de capacité dans `Book`** : vérifier `IsAvailable(s)` et `n > 0` avant d'incrémenter, retourner `(Slot, error)` — `internal/booking/booking.go:25`.
2. **Écrire des tests table-driven pour `Book`** couvrant : `n=0`, `n<0`, `n > Remaining(s)`, `Booked == Capacity` avant l'appel — `internal/booking/booking_test.go`.
3. **Remplacer `time.Now()` dans `sample()`** par un instant fixe (`time.Date(2024, 1, 1, 9, 0, 0, 0, time.UTC)`) pour garantir la reproductibilité — `internal/booking/booking_test.go:9`.
4. **Documenter le contrat d'appel de `Book`** (pré-conditions : `IsAvailable(s) == true`, `n > 0`) dans le commentaire de la fonction, jusqu'à ce que la validation soit intégrée — `internal/booking/booking.go:24`.

## Questions ouvertes

- Une annulation (libération de places) est-elle dans le scope ? Si oui, `Book(s, -n)` est-il le mécanisme prévu, ou une fonction `Cancel(s Slot, n int) (Slot, error)` sera-t-elle ajoutée ? Non déductible du code actuel.
- Y aura-t-il un mécanisme de verrouillage (mutex, transaction DB) pour rendre `IsAvailable → Book` atomique lors de l'introduction d'un serveur concurrent ?
