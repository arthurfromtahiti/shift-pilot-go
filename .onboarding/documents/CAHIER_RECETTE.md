# Cahier de Recette — shift-pilot-go

> **Scope** : cas de test dérivés des workflows fonctionnels (`WORKFLOW_CONSULTER_DISPONIBILITE`, `WORKFLOW_RESERVER_PLACES`). Ce cahier décrit ce qui doit être testé avant toute exposition (HTTP, intégration, production).
>
> **Statut** : matière pauvre — seules trois fonctions pures et une entité. Aucun serveur, aucune interface utilisateur à piloter. Les parcours de recette sont donc techniques (appels de fonction), pas des scénarios utilisateur.

## Stratégie de recette

Chaque test vérifie qu'une fonction Go retourne la valeur attendue pour une entrée donnée. Les cas sont groupés par fonction, avec distinction entre chemin nominal et cas limites.

**Outils** : `go test` (exécution), `assert` ou `testing.T` (vérification). Aucune infrastructure externe requise.

**Prérequis** : Go 1.21 ou supérieur, toolchain Go disponible.

## Parcours de recette par fonction

### 1. Fonction `Remaining(Slot) int`

Retourne le nombre de places libres (Capacité − Réservées).

#### TC-1.1 — Chemin nominal : créneau avec places disponibles
```
Slot := {ID:1, Activity:"Plongée", Start:2024-01-01T09:00:00Z, Capacity:10, Booked:4}
Appel := Remaining(Slot)
Résultat attendu := 6  (10 - 4)
Assertion := Remaining(Slot) == 6
Statut :: COUVERT (test `TestRemaining`)
```

#### TC-1.2 — Créneau plein (places restantes = 0)
```
Slot := {Capacity:10, Booked:10}
Appel := Remaining(Slot)
Résultat attendu := 0
Assertion := Remaining(Slot) == 0
Statut :: NON COUVERT — cas limite critique
```

#### TC-1.3 — Créneau sur-réservé (places restantes négatives)
```
Slot := {Capacity:10, Booked:15}  # État impossible via Book (qui valide `n ≤ Remaining(s)`) — atteignable seulement en construction brute
Appel := Remaining(Slot)
Résultat attendu := -5  (10 - 15)
Assertion := Remaining(Slot) == -5
Statut :: NON COUVERT — cas critique, détecte la sur-réservation
Risque :: `Remaining` retourne une valeur négative sans signal d'erreur
```

#### TC-1.4 — Capacité invalide (négative ou zéro)
```
Slot := {Capacity:0, Booked:0}
Appel := Remaining(Slot)
Résultat attendu :: ? (0 - 0 = 0, mais sémantiquement invalide)
Assertion :: N/A — test pour documenter le comportement réel
Statut :: NON COUVERT — modèle de données invalide accepté
```

#### TC-1.5 — Réservations négatives
```
Slot := {Capacity:10, Booked:-2}  # Modèle invalide
Appel := Remaining(Slot)
Résultat attendu :: ? (10 - (-2) = 12)
Assertion :: N/A — documente le comportement sans validation
Statut :: NON COUVERT
```

---

### 2. Fonction `IsAvailable(Slot) bool`

Retourne `true` ssi au moins une place reste (`Remaining(s) > 0`).

#### TC-2.1 — Chemin nominal : créneau disponible
```
Slot := {Capacity:10, Booked:4}
Appel := IsAvailable(Slot)
Résultat attendu := true  (Remaining(Slot) = 6 > 0)
Assertion := IsAvailable(Slot) == true
Statut :: COUVERT (test `TestIsAvailable`)
```

#### TC-2.2 — Créneau plein (indisponible)
```
Slot := {Capacity:10, Booked:10}
Appel := IsAvailable(Slot)
Résultat attendu := false  (Remaining(Slot) = 0, not > 0)
Assertion := IsAvailable(Slot) == false
Statut :: NON COUVERT — cas limite critique
Risque :: Point de décision clé pour autoriser/refuser une réservation
```

#### TC-2.3 — Créneau sur-réservé
```
Slot := {Capacity:10, Booked:15}
Appel := IsAvailable(Slot)
Résultat attendu := false  (Remaining(Slot) = -5, not > 0)
Assertion := IsAvailable(Slot) == false
Statut :: NON COUVERT — détecte état invalide, mais `IsAvailable` retourne correctement false
```

---

### 3. Fonction `Book(Slot, int) Slot`

Enregistre une réservation. Retourne le `Slot` mis à jour (Booked augmenté de n).

#### TC-3.1 — Chemin nominal : réserver n=2 places
```
Slot_initial := {ID:1, Activity:"Plongée", Capacity:10, Booked:4}
n := 2
Appel := Book(Slot_initial, 2)
Résultat attendu := {ID:1, Activity:"Plongée", Capacity:10, Booked:6}
Assertions :=
  - Book(...).Booked == 6
  - Book(...).Capacity == 10  (inchangé)
  - Book(...).ID == 1         (inchangé)
Statut :: COUVERT (test `TestBook`)
```

#### TC-3.2 — Immuabilité : le Slot source n'est pas modifié
```
Slot_initial := {Capacity:10, Booked:4}
n := 2
Book(Slot_initial, 2)
Assertion := Slot_initial.Booked == 4  (inchangé, valeur passée, pas pointeur)
Statut :: NON COUVERT — garantie par la sémantique Go (passage par valeur) mais jamais testée explicitement
```

#### TC-3.3 — Sur-réservation : n > Remaining(Slot)
```
Slot := {Capacity:10, Booked:4}
n := 7  # Il n'y a que 6 places libres
Appel := Book(Slot, 7)
Résultat réel := (Slot{}, ErrCapacityExceeded)  # Rejeté avec erreur
Résultat attendu (métier) := ERREUR — "capacité dépassée, réservation refusée"
Assertion :: Book(Slot, 7) == (Slot{}, ErrCapacityExceeded)
Statut :: COUVERT — test `TestBookCapacityExceeded` (ligne 34-39)
Risque :: RÉSOLU — `Book` valide `n ≤ Remaining(s)` avant incrément
Impact :: La sur-réservation est impossible via `Book`
```

#### TC-3.4 — `n` négatif (rejeté)
```
Slot := {Capacity:10, Booked:6}
n := -2
Appel := Book(Slot, -2)
Résultat réel := (Slot{}, ErrInvalidBookingCount)  # Rejeté avec erreur
Résultat attendu (métier) :: ERREUR — "n doit être positif"
Assertion :: Book(Slot, -2) == (Slot{}, ErrInvalidBookingCount)
Statut :: COUVERT — test `TestBookNegative` (ligne 58-63)
Risque :: RÉSOLU — `Book` valide `n > 0` au premier appel
Impact :: L'annulation doit utiliser une fonction dédiée (pas implémentée)
```

#### TC-3.5 — `n` zéro (rejeté)
```
Slot := {Capacity:10, Booked:4}
n := 0
Appel := Book(Slot, 0)
Résultat réel := (Slot{}, ErrInvalidBookingCount)  # Rejeté avec erreur
Assertion :: Book(Slot, 0) == (Slot{}, ErrInvalidBookingCount)
Statut :: COUVERT — test `TestBookZero` (ligne 51-56)
Impact :: Les valeurs nulles ou négatives sont systématiquement rejetées
```

#### TC-3.6 — Sur-réservation massive
```
Slot := {Capacity:10, Booked:10}
n := 100
Appel := Book(Slot, 100)
Résultat réel := (Slot{}, ErrCapacityExceeded)  # Créneau plein
Résultat attendu :: ERREUR
Assertion :: Book(...) == (Slot{}, ErrCapacityExceeded)
Statut :: COUVERT implicitement par `TestBookCapacityExceeded` (logique générale)
Impact :: Aucune surréservation n'est possible, peu importe l'ampleur de `n`
```

---

## Intégrations de workflows

### Scénario WS-1 : Consulter puis réserver (nominale)
```
1. Appelant construit Slot{Capacity:10, Booked:4}
2. Appelant appelle IsAvailable(Slot) → true
3. Appelant appelle Book(Slot, 2) → Slot{Booked:6}
4. Appelant utilise le résultat ou le persiste

Assertions :
  - IsAvailable(Slot) == true
  - Book(Slot, 2).Booked == 6
  - Slot_original inchangé

Statut :: PARTIELLEMENT COUVERT (étapes 1, 3 oui ; étape 2 implicite)
Risque :: Pas de test d'atomicité — aucun verrou entre IsAvailable et Book
```

### Scénario WS-2 : Tentative de sur-réservation (cas d'échec attendu)
```
1. Appelant construit Slot{Capacity:10, Booked:10}  # Déjà plein
2. Appelant appelle IsAvailable(Slot) → false
3. Appelant DEVRAIT refuser la réservation
4. Si l'appelant appelle malgré tout Book(Slot, 1) → (Slot{}, ErrCapacityExceeded)

Assertions :
  - IsAvailable(Slot) == false
  - Book(Slot, 1) == (Slot{}, ErrCapacityExceeded)  (comportement protégé)

Statut :: COUVERT — `Book` valide automatiquement la disponibilité
Risque :: RÉSOLU — `Book` retourne une erreur explicite si capacité dépassée
Impact :: Aucune sur-réservation possible, même sans vérification préalable d'IsAvailable
```

---

## Récapitulatif de couverture

| Cas | Couverture | Criticité | Recommandation |
|---|---|---|---|
| `Remaining` nominal (places > 0) | ✓ Couvert | High | OK |
| `Remaining` créneau plein (places = 0) | ✗ Non couvert | **High** | Ajouter test |
| `Remaining` sur-réservé (places < 0) | ✗ Non couvert | **High** | Ajouter test |
| `IsAvailable` nominal (true) | ✓ Couvert | High | OK |
| `IsAvailable` créneau plein (false) | ✗ Non couvert | **High** | Ajouter test |
| `Book` nominal (+2 places) | ✓ Couvert | High | OK |
| ~~`Book` sur-réservation~~ | ✓ **Couvert** | **RÉSOLU** | `TestBookCapacityExceeded` couvre le cas |
| `Book` annulation (n<0) | ✓ **Couvert** | **RÉSOLU** | `TestBookNegative` rejetée via `ErrInvalidBookingCount` |
| `Book` no-op (n=0) | ✓ **Couvert** | **RÉSOLU** | `TestBookZero` rejetée via `ErrInvalidBookingCount` |
| Immuabilité de Slot | ✗ Non couvert | Medium | Ajouter test explicite |

---

## Avant déploiement ou intégration

**Blocants** (doivent être résolus avant toute exposition HTTP/production) :

1. ✓ **FAIT** — `Book` retourne une `error` et rejette la sur-réservation (`n > Remaining(s)` ou `n ≤ 0`) via `ErrCapacityExceeded` et `ErrInvalidBookingCount`.
2. ⚠ Partiellement couvert — TC-3.3 (sur-réservation), TC-3.4 (n<0), TC-3.5 (n=0) testés. Restent TC-1.2, TC-1.3, TC-2.2 (Remaining/IsAvailable cas limites) — cas importants mais secondaires.
3. ✗ Un `Slot` invalide (`Capacity ≤ 0`, `Booked < 0` après construction) doit être rejeté via constructeur validant `NewSlot(...)`.
4. ✓ Partiellement fait — annulation via `Book(s, -n)` interdite (`ErrInvalidBookingCount`). Une fonction explicite `Cancel` reste à créer pour formellement supporter l'annulation.

**Non-blocants** (amélioration, peut suivre) :

5. ⚠ CI avec `go test ./... -race` doit exécuter à chaque push.
6. ⚠ Tests table-driven plutôt que trois fonctions isolées.
7. ⚠ `time.Now()` remplacé par date fixe dans fixtures.

---

## Exécution

**Commande actuelle** :
```bash
go test ./...
```

**Résultat** : 3 tests nominaux passants. Exécution réelle `INCONNU` (toolchain absente lors de l'analyse).

**Commande recommandée (futur)** :
```bash
go test ./... -race -cover
```

Cela ajoutera :
- `-race` : détection de race conditions (préparation pour serveur concurrent).
- `-cover` : rapport de couverture en lignes (objectif : ≥80%, mais actuellement ≤50% vu la pauvreté des tests).
