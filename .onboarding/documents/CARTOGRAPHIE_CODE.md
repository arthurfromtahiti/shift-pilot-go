# Cartographie du code — shift-pilot-go

> **Confiance : high**. Corpus entièrement lu (57 lignes totales : 28 production + 29 tests). Aucun ambiguïté sur la structure. La pauvreté du livrable rend l'analyse exhaustive facile, mais elle signifie aussi que beaucoup d'architecture reste à concevoir.

## Structure générale

```
shift-pilot-go/
├── .git/                     # Historique Git
├── .onboarding/              # Artefacts d'analyse (domaines, workflows, audits, documents)
├── .claude/                  # Configuration Claude Code (si applicable)
├── internal/
│   └── booking/              # Domaine unique : réservation & disponibilité
│       ├── booking.go        # Production : type Slot + 3 fonctions pures
│       └── booking_test.go   # Tests : 3 fonctions de test (couverture nominale seulement)
├── go.mod                    # Déclaration du module (aucune dépendance externe)
├── go.sum                    # Checksum des dépendances (vide)
└── README.md                 # Documentation d'introduction
```

**Taille totale** : 28 lignes de production, 29 lignes de tests. Rien d'autre.

## Domaines et fichiers

### Domaine : Réservation & disponibilité (`reservation-creneaux`)

**Priorité** : cœur — seul domaine métier présent.
**Catégorie** : métier.
**Confiance** : medium.

**Fichiers**
- **`internal/booking/booking.go`** (28 lignes, `VÉRIFIÉ_CODE`)
  - Type `Slot` : créneau d'activité (5 champs).
  - Fonction `Remaining(Slot) int` : places restantes.
  - Fonction `IsAvailable(Slot) bool` : disponibilité.
  - Fonction `Book(Slot, int) Slot` : enregistrement de réservation.

- **`internal/booking/booking_test.go`** (29 lignes, `VÉRIFIÉ_CODE`, exécution `INCONNU`)
  - Fonction `sample() Slot` : fixture de test.
  - Fonction `TestRemaining` : vérifie `Remaining(Slot{Capacity:10, Booked:4}) == 6`.
  - Fonction `TestIsAvailable` : vérifie `IsAvailable(Slot{Capacity:10, Booked:4}) == true`.
  - Fonction `TestBook` : vérifie `Book(Slot{Capacity:10, Booked:4}, 2).Booked == 6`.

## Points critiques du code

| Fichier | Ligne | Fonction | Rôle | Risque |
|---|---|---|---|---|
| `booking.go` | 6–12 | Type `Slot` | Entité unique ; cœur du modèle | Invariants non enforced (Booked > Capacity possible) |
| `booking.go` | 15–17 | `Remaining` | Calcule places libres (Capacity − Booked) | Peut retourner négatif ; pas de garde |
| `booking.go` | 20–22 | `IsAvailable` | Prédicat disponibilité (Remaining > 0) | Correct dans sa forme, mais passif sur état invalide |
| `booking.go` | 25–28 | **`Book`** | **Seule fonction qui mutate l'état** | **CRITIQUE** : pas de vérification de capacité, accepte tout `n` |
| `booking_test.go` | 8–10 | Fixture `sample()` | Bloc de test partagé | `time.Now()` rend non-déterministe si `Start` devient fonctionnel |
| `booking_test.go` | 12–29 | Trois tests | Couverture nominale | Aucun cas limite ; sur-réservation non testée |

## Flux d'appels

```
Appelant externe (code client)
    ↓
[importe github.com/arthurfromtahiti/shift-pilot-go/internal/booking]
    ↓
Fonctions publiques:
  ├─ Remaining(Slot) → int
  ├─ IsAvailable(Slot) → bool
  └─ Book(Slot, n) → Slot
```

Aucun point d'entrée supplémentaire : pas de serveur HTTP, pas de CLI, pas de fonction `main`, pas de package `cmd`.

## Dépendances externes

**Bloc `require` dans `go.mod`** : vide — aucune dépendance externe.

**Imports** :
- `time` (bibliothèque standard Go) — seul import dans `booking.go`.
- `testing` (bibliothèque standard Go, utilisé par les tests).

**Supply-chain** : zéro risque de CVE transitif. Build reproductible sans gestion de versions externes.

## Architecture actuelle

### Caractéristiques

- **Monocouche** : un seul package métier (`booking`), zéro abstraction (pas d'interface, pas de contrat).
- **Fonctions pures** : pas d'état global, pas de pointeur partagé, pas de mutation destructive.
- **Typographie simple** : un seul type (`Slot`), cinq champs, pas de générique, pas de constructeur.
- **Pas de configuration** : aucun fichier `.env`, aucun YAML, aucun paramètre de build.
- **Pas d'erreurs** : aucune fonction ne retourne `error` ; états invalides retournés comme valeurs nominales.

### Forces

1. `Slot` est un type valeur — prévisible, testable, no-copy semantics simple.
2. Zéro dépendance — maintenance facile, pas de breaking change transitive.
3. Code lisible — chaque fonction tient en 3–4 lignes, sans branchement.

### Dettes immédiates

1. **Pas d'interface** — tout appelant est couplé aux types concrets. Si `Remaining` ou `Book` change de signature, tous les appelants doivent être mis à jour.
2. **Pas de type d'erreur** — impossible de distinguer un résultat valide d'un appel sémantiquement invalide.
3. **Absence de validation de garde** — `Book` n'enforces aucun invariant.
4. **Pas de configuration de build/CI** — aucun workflow GitHub Actions, aucun Makefile, aucune automatisation visible.

## Tests

**Couverture nominale** : 3 tests passants en mode lecture (`VÉRIFIÉ_CODE`), exécution réelle `INCONNU` (toolchain Go absente).

**Structure** :
- White-box : tests dans le même package `booking` (accès aux non-exportés si existants).
- Fixture unique : `sample()` retourne toujours `Slot{ID:1, Activity:"Plongée", Start:time.Now(), Capacity:10, Booked:4}`.
- Assertions simples : chaque test = 1 appel + 1 vérification.

**Cas non testés** (tous critiques pour la robustesse) :
- `Book(s, 0)` — no-op silencieux.
- `Book(s, n)` où `n > Remaining(s)` — sur-réservation.
- `Book(s, -n)` — annulation implicite.
- `Remaining(s)` où `Booked > Capacity` — valeur négative.
- Immuabilité de `Slot` après `Book` — non testée explicitement.

**Absence de CI** : aucun workflow GitHub Actions, aucune exécution automatique à chaque push.

## Zones à concevoir avant l'évolution

| Composant | État | Impact si absent |
|---|---|---|
| Couche HTTP/transport | Absent | Impossible d'exposer en tant que service REST |
| Couche de persistance | Absent | Toute réservation disparaît en fin de scope |
| Gestion d'erreurs | Absent | Impossible de signaler une sur-réservation ou un appel invalide |
| Interfaces/contrats | Absent | Tout appelant est couplé à `Slot` et aux signatures actuelles |
| CI/CD | Absent | Tests ne s'exécutent qu'à la main ; pas d'automatisation |
| Configuration | Absent | Aucun paramètre externalizable (DB, logs, feature flags) |
| Validation | Absent | Invariants de `Slot` non enforced (`Capacity > 0`, `Booked ≤ Capacity`) |

## Recommandations pour l'architecture future

1. **Avant d'ajouter un serveur HTTP** : introduire une interface `type BookingService interface { ... }` pour découpler transport et logique métier.
2. **Avant d'ajouter une base de données** : concevoir le schéma d'intégrité (constraints CHECK, migrations) et les propriétés de transactions (atomicité de lecture-puis-écriture).
3. **Avant d'exposer publiquement** : implémenter la validation des invariants en amont et la gestion d'erreurs — `Book` doit retourner `(Slot, error)`, pas juste `Slot`.
4. **En parallèle** : mettre en place un workflow CI qui exécute `go test ./... -race` à chaque push et PR.
5. **Pour la testabilité** : remplacer `time.Now()` par une heure fixe dans les fixtures, utiliser des tests table-driven pour couvrir les cas limites.
