# Architecture — Audit

> Confiance : medium — l'intégralité du code source a été lue (`VÉRIFIÉ_CODE` sur 4 fichiers), mais le contexte d'exécution (pipeline CI/CD, environnement de déploiement, intentions futures) est `INCONNU` faute de configuration visible. La matière est volontairement pauvre : c'est un pilote de 28 lignes de code métier.

## Compréhension globale

`shift-pilot-go` est une bibliothèque Go sans exécutable. L'unique package métier (`internal/booking`) expose trois fonctions pures sur un seul type valeur (`Slot`). Il n'existe aucune couche de présentation, aucune couche de persistance, aucune injection de dépendances, aucune interface — c'est un noyau de logique à l'état de germe, posé dans une structure de module Go standard.

## Résumé exécutif

Le projet n'a qu'un seul niveau architectural : un package `internal/booking` isolé correctement par la directive `internal/`, portant trois fonctions pures et un type valeur. L'architecture est trop embryonnaire pour être qualifiée de bonne ou mauvaise — elle ne pose pas encore les décisions structurantes d'une vraie application (transport, persistance, interfaces, gestion d'erreurs). Ce qui existe est propre dans ses conventions Go mais silencieux sur tout le reste.

Le risque architectural immédiat n'est pas dans ce qui est là — c'est dans ce qui n'y est pas : aucune interface ne définit le contrat de `booking`, aucun type d'erreur ne structure les signaux de défaut, et le README mentionne un déploiement vers "staging" sans qu'aucune configuration de build ni de CI/CD ne soit visible dans le dépôt. Si quelqu'un doit faire évoluer ce pilote en service réel, il part d'un canevas vierge sur tous les axes non-métier.

## Constats détaillés

**`VÉRIFIÉ_CODE` — Structure de module conforme aux idiomes Go.** Le fichier `go.mod` (`go.mod:1`) déclare le module `github.com/arthurfromtahiti/shift-pilot-go`, Go 1.21, sans aucune dépendance externe. Le package `internal/booking` (`internal/booking/booking.go:1`) est placé sous `internal/` : en Go, ce chemin empêche tout module externe d'importer ce package directement, ce qui est la convention correcte pour du code non encore stable ou destiné à rester privé.

**`VÉRIFIÉ_CODE` — Aucun exécutable, aucun serveur, aucune route.** Il n'existe pas de `func main` dans le dépôt (`find` sur `*.go` : un seul fichier de prod, `booking.go`, et un fichier de tests). Le projet est une bibliothèque. Il n'y a ni handler HTTP, ni CLI, ni goroutine de service. Le README (`README.md:12`) mentionne `git → staging` comme instruction de déploiement, mais aucun `Dockerfile`, `Makefile`, `.github/`, ni configuration CI n'est présent dans le dépôt. Ce que "staging" désigne est `INCONNU`.

**`VÉRIFIÉ_CODE` — Fonctions pures, passage par valeur.** Les trois fonctions de `booking.go` reçoivent `Slot` par valeur (copie) et retournent soit un scalaire, soit un nouveau `Slot`. Il n'y a pas de pointeur, pas de mutation d'état partagé, pas d'effet de bord observable. C'est structurellement correct pour des calculs purs et rend chaque fonction testable en isolation totale. En revanche, `Book` (`booking.go:25-28`) retourne un `Slot` modifié que l'appelant doit impérativement utiliser — sinon la réservation est perdue. Ce pattern est idiomatique en Go pour des types valeur, mais il exige une discipline d'usage que rien dans le code n'enforces.

**`VÉRIFIÉ_CODE` — Aucune interface, aucun type d'erreur.** Aucun `interface` n'est défini dans le package. Les fonctions retournent `int`, `bool`, et `Slot` — jamais `error`. Il n'existe pas de sentinel error, pas de type personnalisé d'erreur, pas de `errors.New`. Le package n'a donc aucun mécanisme pour signaler qu'un appel est sémantiquement invalide (ex. `Book` sur un créneau plein).

**`HYPOTHÈSE` — La croissance vers un service va nécessiter des décisions architecturales majeures.** Si le pilote doit devenir un service (HTTP, gRPC, ou autre), il manque : une couche de transport, une couche de persistance, un modèle d'erreur, une couche d'interfaces pour testabilité et découplage, et probablement un schéma de concurrence (goroutines, canaux ou mutex). Aucune de ces décisions n'est amorcée. Ce n'est pas une critique du pilote tel qu'il est — c'est un signal pour le planificateur.

## Forces

- `VÉRIFIÉ_CODE` : La convention `internal/` protège le package de toute importation accidentelle par un module externe (`internal/booking/booking.go:1`).
- `VÉRIFIÉ_CODE` : Zéro dépendance externe (absence de bloc `require` dans `go.mod`) — pas de supply-chain risk, pas de gestion de versions transitives, build reproductible.
- `VÉRIFIÉ_CODE` : Les trois fonctions sont pures (pas d'état global, pas de pointeur partagé) : elles sont composables, prévisibles, et testables sans mock ni setup (`booking.go:15-28`).

## Dettes techniques

- `VÉRIFIÉ_CODE` : Absence de type d'erreur — `Book` ne peut pas signaler une tentative de sur-réservation (`booking.go:25-28`). Cette dette est légère aujourd'hui (c'est une lib de pilote), mais bloquante dès qu'une couche appelante doit distinguer succès et échec.
- `VÉRIFIÉ_CODE` : Absence d'interface — l'appelant de `booking` est couplé aux types concrets. Si le comportement de `Book` ou `Remaining` évolue, tous les appelants doivent être mis à jour sans filet (`booking.go:1-28`).
- `HYPOTHÈSE` : Absence de configuration de build et de CI/CD — si "staging" est un vrai environnement, le chemin de build → déploiement est soit documenté ailleurs, soit informel. Dans les deux cas, il n'est pas versionné avec le code.

## Zones critiques

- `VÉRIFIÉ_CODE` : `internal/booking/booking.go` — le seul fichier de production. Toute évolution du projet passe par là. Le risque de sur-réservation silencieuse a été résorbé : `Book` valide désormais `n > 0` et `n ≤ Remaining(s)`, retourne `(Slot, error)` (`booking.go:36-45`).

## Risques

- `VÉRIFIÉ_CODE` (impact futur, plutôt que défaut) : Lorsqu'un transport (HTTP, etc.) sera ajouté, le passage de `Slot` par valeur sans persistance rendra chaque réservation éphémère — la valeur retournée par `Book` doit être explicitement stockée ou propagée par l'appelant (`booking.go:43-44`). Oublier de le faire est un bug silencieux. Note : la validation dans `Book` elle-même garantit désormais que seuls les états valides sont retournés.
- `HYPOTHÈSE` : Si le même slot est manipulé par plusieurs goroutines (futur serveur concurrent), le schéma `IsAvailable → Book` n'est pas atomique : une race condition TOCTOU est possible. Le code actuel n'en souffre pas (mono-goroutine, pas de serveur), mais aucun avertissement ni mécanisme de verrouillage n'est en place pour signaler ce risque lors de l'évolution.

## Recommandations priorisées

1. **Définir une interface `BookingService`** avant d'ajouter une couche de transport — permet de substituer l'implémentation, de mocker en test, et de découpler — `internal/booking/booking.go`.
2. **Introduire un type `error` sur `Book`** (`func Book(s Slot, n int) (Slot, error)`) pour signaler la sur-réservation — sans ça, toute couche appelante doit dupliquer la logique de garde — `internal/booking/booking.go:25`.
3. **Ajouter une configuration CI (GitHub Actions ou équivalent)** qui exécute `go test ./...` à chaque push — garantit que le code reste vert même après ajout de dépendances ou de fonctions — (fichier à créer : `.github/workflows/ci.yml`).
4. **Clarifier la cible du pilote** (lib vs service) dans `README.md` — le "Git → staging" est ambigu et peut mener à des déploiements non bornés.

## Questions ouvertes

- Qu'est-ce que "staging" dans le README — un vrai environnement qui tourne le package, ou une note de processus ? Aucun artefact de déploiement visible.
- Le projet est-il censé rester une bibliothèque Go importable, ou doit-il évoluer en service HTTP autonome ? Cette décision détermine toute l'architecture future.
- Y a-t-il d'autres dépôts (frontend, service orchestrateur) qui importent ou dépendent de ce package ? Non visible dans ce dépôt.
