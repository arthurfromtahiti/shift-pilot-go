# Modèle de données — Audit

> Confiance : medium — la totalité du modèle de données est lisible dans un seul fichier de 28 lignes (`VÉRIFIÉ_CODE`). Il n'y a pas de base de données, pas de migration, pas de schéma SQL. La faible confiance tient à l'absence totale de persistance : le modèle est sans doute incomplet par rapport à l'intention finale, mais ce que l'intention implique reste `INCONNU`.

## Compréhension globale

Le modèle de données de `shift-pilot-go` se réduit à un seul type Go : `Slot`. C'est un struct valeur (pas un pointeur) à cinq champs, défini dans `internal/booking/booking.go`. Il n'y a pas de base de données, pas de fichier de migration, pas d'ORM, pas de schéma JSON. Le modèle est embryonnaire et ne contraint aucune invariant — les valeurs invalides (négatifs, sur-réservation) sont représentables sans erreur.

## Résumé exécutif

Un seul type, cinq champs, zéro contrainte. Le modèle décrit correctement un créneau d'activité dans le cas nominal (capacité positive, places réservées entre 0 et capacité). Il ne protège pas contre les états invalides (`Booked > Capacity`, `Capacity ≤ 0`, `n < 0` passé à `Book`), n'a pas de persistance, et est incomplet pour un vrai système de réservation (pas de durée, pas de prix, pas d'identifiant de réservant, pas d'état de créneau). Ces lacunes sont attendues pour un pilote, mais elles devront être adressées avant toute exposition.

## Constats détaillés

**`VÉRIFIÉ_CODE` — Type unique `Slot`, cinq champs.** Le type `Slot` est déclaré en `internal/booking/booking.go:6-12` :

```go
type Slot struct {
    ID       int
    Activity string
    Start    time.Time
    Capacity int
    Booked   int
}
```

- `ID int` : identifiant entier. Aucune contrainte d'unicité n'est possible sans couche de persistance. Le type `int` exclut les UUID et les identifiants opaques.
- `Activity string` : nom de l'activité (ex. `"Plongée"` dans les tests). Chaîne libre, aucune validation de longueur ni d'appartenance à un catalogue.
- `Start time.Time` : date-heure de début du créneau. Aucun champ `End` ni `Duration` — la durée d'un créneau est `INCONNU`.
- `Capacity int` : capacité totale. Aucune contrainte de positivité — un `Slot{Capacity: -5}` est syntaxiquement valide.
- `Booked int` : places actuellement réservées. Peut dépasser `Capacity` (voir `Book`), peut être négatif (`Book(s, -n)` implicite). Aucune contrainte d'intégrité.

**`VÉRIFIÉ_CODE` — Aucune invariant enforced.** Aucun constructeur validant, aucune méthode `Validate()`, aucun tag de validation (`validate:`, `json:`, etc.) n'est défini (`booking.go:6-12`). Un `Slot{Capacity: 0, Booked: 100}` ou `Slot{ID: 0}` est un état valide du point de vue du compilateur Go.

**`VÉRIFIÉ_CODE` — Aucune persistance, aucune migration.** Le dépôt ne contient ni fichier SQL, ni ORM (`gorm`, `sqlx`, `ent`, etc.), ni migration (`goose`, `atlas`, `migrate`), ni schéma JSON ou Protobuf. Le `Slot` n'existe qu'en mémoire pendant la durée de vie d'un appel de fonction. La valeur retournée par `Book` est perdue si l'appelant ne la conserve pas explicitement.

**`VÉRIFIÉ_CODE` — Passage par valeur — sémantique de copie.** Toutes les fonctions reçoivent `Slot` par valeur (`booking.go:15, 20, 25`). `Book` retourne un nouveau `Slot` — pas une modification en place. C'est idiomatique Go pour des types sans état partagé, mais cela signifie qu'il n'existe aucun identifiant d'instance stable : le `Slot` de l'appelant et le `Slot` retourné par `Book` sont deux valeurs distinctes en mémoire.

**`HYPOTHÈSE` — Le modèle est incomplet pour un système de réservation réel.** Pour couvrir les cas d'usage mentionnés dans le README (réservation d'activités nautiques), un modèle complet devrait au minimum inclure : une durée ou un champ `End`, un identifiant de réservant (client), un statut de créneau (`ouvert`, `complet`, `annulé`), un champ `Price`, et une entité `Booking` distincte de `Slot`. Ces éléments sont totalement absents — ce qui est cohérent avec la nature de pilote, mais doit être noté pour planifier l'évolution.

## Forces

- `VÉRIFIÉ_CODE` : Modèle minimal et lisible — toute la structure tient en 7 lignes, sans ambiguïté (`booking.go:6-12`).
- `VÉRIFIÉ_CODE` : Utilisation de `time.Time` pour `Start` (type fort de la bibliothèque standard) plutôt qu'un `string` ou un `int` timestamp — correct pour la manipulation de dates.
- `VÉRIFIÉ_CODE` : Aucune dépendance externe dans le modèle — pas d'ORM, pas de bibliothèque de validation tierce. Facilite un remplacement ou une extension ultérieure.

## Dettes techniques

- ~~`VÉRIFIÉ_CODE` : `Booked` peut légalement dépasser `Capacity`~~ — **RÉSOLU** : `Book` valide `n ≤ Remaining(s)` avant d'incrémenter `Booked`, garantissant `Booked ≤ Capacity` (`booking.go:40-41`). L'invariant est enforced par `Book`, bien qu'un `Slot` construit manuellement avec `Booked > Capacity` reste syntaxiquement valide.
- `VÉRIFIÉ_CODE` : `Capacity` peut être zéro ou négatif — aucune garde à la construction (`booking.go:6-12`).
- `VÉRIFIÉ_CODE` : `ID int` est insuffisant pour un système distribué ou avec persistance — collisions possibles, pas de génération automatique, pas d'opacité.
- `VÉRIFIÉ_CODE` : Absence de champ `End` ou `Duration` — impossible de calculer la durée d'un créneau, de détecter les chevauchements, ou d'afficher un créneau complet.

## Zones critiques

- `VÉRIFIÉ_CODE` : `internal/booking/booking.go:6-12` (`Slot`) — toute évolution du modèle passe par ce struct. Sa modification impacte toutes les fonctions du package et tous leurs appelants.

## Risques

- ~~`VÉRIFIÉ_CODE` (défaut prouvé, impact connu) : `Booked > Capacity` est un état atteignable via `Book` sans erreur~~ — **RÉSOLU** : `Book` rejette `n > Remaining(s)` avec `ErrCapacityExceeded` (`booking.go:40-41`). L'invariant `Booked ≤ Capacity` est garanti par `Book`.
- `HYPOTHÈSE` : Si une base de données est ajoutée sans avoir introduit de contraintes d'intégrité (CHECK `booked >= 0`, CHECK `booked <= capacity`, etc.), le modèle actuel permettra l'insertion d'états invalides au niveau SQL, avec des répercussions difficiles à corriger a posteriori sur des données existantes.
- `HYPOTHÈSE` : L'identifiant `ID int` sans unicité enforced est une source probable de conflits dès que plusieurs créneaux coexistent en mémoire ou en base.

## Recommandations priorisées

1. **Introduire un constructeur validant** `func NewSlot(id int, activity string, start time.Time, capacity int) (Slot, error)` qui rejette `capacity <= 0` et `activity == ""` — `internal/booking/booking.go`.
2. **Ajouter un champ `End time.Time` ou `Duration time.Duration`** pour rendre un créneau exploitable (affichage, détection de chevauchement) — `internal/booking/booking.go:6-12`.
3. **Séparer l'entité `Booking` (une réservation) de l'entité `Slot` (un créneau)** dès que la persistance est introduite — un `Slot` représente l'offre, une `Booking` représente l'acte de réservation d'un réservant particulier.
4. **Remplacer `ID int` par `ID string` (UUID)** avant d'introduire la persistance pour éviter les collisions et permettre la génération côté client.

## Questions ouvertes

- Y aura-t-il une entité `Client`/`User` ? Rien ne lie aujourd'hui une réservation à un réservant.
- Le champ `Activity` est-il un libellé libre ou une clé vers un catalogue d'activités ? Non déductible du code.
- La durée d'un créneau est-elle fixe (ex. toutes les plongées durent 1h) ou variable ? Le modèle ne porte ni `End` ni `Duration`.
- Quel est le schéma de persistance envisagé — SQL, NoSQL, in-memory ? Non visible dans le code.
