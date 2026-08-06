# Relecture — CARTE_DES_DOMAINES.md

> Relecteur : **Relecteur de domaines** (agent `2f9907e0-47bf-4657-a2ca-4e0182da3704`)
> Artefact relu : `.onboarding/domaines/CARTE_DES_DOMAINES.md`
> SHA de tête analysé : `8122e2ddaa200f5e30b4a50364c0d1758fcec72f`
> Date de relecture : 2026-08-05

## Verdict global

**Bon** — exploitable sans réserve bloquante.

La carte reflète fidèlement ce qui existe dans le dépôt : un seul domaine prouvé (`reservation-creneaux`), toutes les lignes citées exactes, aucun domaine inventé, aucun pan oublié. Les incertitudes sont honnêtement documentées. La carte peut servir de fondation à l'étape Workflows.

## Problèmes bloquants

Aucun.

## Problèmes mineurs

**1. Confiance du domaine "medium" légèrement conservatrice.**
Le code est transparent : 3 fonctions pures sans ambiguïté, 1 struct bien nommée, nommage cohérent partout. `VÉRIFIÉ_CODE` sur 100 % des affirmations. La confiance "high" serait plus précise. "Medium" est défendable (le contour peut évoluer à la prochaine itération du pilote), mais elle communique une incertitude que le code actuel ne justifie pas.
→ Observation seulement, non bloquant. Le producteur peut corriger à sa discrétion.

## Points vérifiés et corrects

1. **Preuves du domaine `reservation-creneaux` — toutes vérifiées par lecture directe des fichiers.**
   - `Slot struct` à `booking.go:6` → confirmé (`type Slot struct {` ligne 6). `VÉRIFIÉ_CODE`.
   - Champs `ID`, `Activity`, `Start time.Time`, `Capacity`, `Booked` → confirmés lignes 7-11. `VÉRIFIÉ_CODE`.
   - `Remaining(s Slot) int` à `booking.go:15` → confirmé. `VÉRIFIÉ_CODE`.
   - `IsAvailable(s Slot) bool` à `booking.go:20` → confirmé. `VÉRIFIÉ_CODE`.
   - `Book(s Slot, n int) (Slot, error)` à `booking.go:36` → confirmé. `VÉRIFIÉ_CODE`.
   - `Book` retourne par valeur (pas de mutation en place) → confirmé : `s.Booked += n; return s` (`booking.go:26-27`). `VÉRIFIÉ_CODE`.

2. **Tests cités présents.**
   - `TestRemaining`, `TestIsAvailable`, `TestBook` → tous présents dans `booking_test.go`. `VÉRIFIÉ_CODE`.
   - Assertion `TestRemaining` : `Remaining(sample()) == 6` pour `Capacity:10, Booked:4` → mathématiquement correct. `VÉRIFIÉ_CODE`.
   - Assertion `TestBook` : après `Book(sample(), 2)`, `Booked == 6` → correct (4+2). `VÉRIFIÉ_CODE`.

3. **Indices de rattachement corrects et non-envahissants.**
   - Pattern `internal/booking/*.go` testé par grep : matche exactement `booking.go` et `booking_test.go`, rien d'autre dans le repo. Les identifiants `Capacity`, `Booked`, `Remaining`, `IsAvailable`, `Book`, `Activity` n'apparaissent que dans ces deux fichiers. `VÉRIFIÉ_CODE`.

4. **Aucun domaine inventé.**
   - Le repo ne contient, hors `.git` et `.claude/`, que : `internal/booking/booking.go`, `internal/booking/booking_test.go`, `go.mod`, `README.md`. Tout est couvert par le seul domaine ou placé en incertitudes. `VÉRIFIÉ_CODE`.

5. **Aucun pan oublié.**
   - Exploration indépendante du repo (4 fichiers de code + go.mod + README) : aucune entité, route, contrôleur, job, intégration externe hors périmètre de la carte. `VÉRIFIÉ_CODE`.

6. **Champ "Dépend de la base" honnête.**
   - `non` : aucune persistance, aucun signal du §6 (`content`, `layout`, `blocks`, `config`, renderer récursif). `go.mod` ne liste aucune dépendance externe. `VÉRIFIÉ_CODE`.

7. **Confiance globale "low" justifiée.**
   - Aucun `func main`, aucun serveur HTTP, aucune route, aucune couche de stockage. Le README confirme « pilote de test ». La confiance basse sur la stabilité de la carte est correcte pour ce stade du projet. `VÉRIFIÉ_CODE` + `OBSERVÉ` (README).

8. **Granularité correcte dans le contexte.**
   - La règle des 4-12 domaines s'applique à des projets de production. Ce pilote de 28 lignes de logique métier ne peut pas être fragmenté sans inventer des domaines. Un seul domaine est la seule réponse correcte et honnête. `VÉRIFIÉ_CODE`.

9. **Observation fonctionnelle `Book` sans garde — correctement remontée en incertitude.**
   - `Book` valide `n > 0` et `n ≤ Remaining(s)` avant incrément de `Booked` (`booking.go:37-43`) : la sur-réservation est impossible. La carte documente ce comportement correctement. `VÉRIFIÉ_CODE`.

## Recommandations de correction

- (Optionnel) Relever la confiance du domaine de "medium" à "high" pour refléter la transparence totale du code actuel.

Aucune autre correction requise.
