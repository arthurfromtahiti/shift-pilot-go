# Relecture — Documents de référence (shift-pilot-go)

> Artefacts relus : `PROJECT_CONTEXT.md`, `CDC_FONCTIONNEL.md`, `CARTOGRAPHIE_CODE.md`, `CAHIER_RECETTE.md`
> Matériau amont consulté : `CARTE_DES_DOMAINES.md`, `WORKFLOW_CONSULTER_DISPONIBILITE.md`, `WORKFLOW_RESERVER_PLACES.md`, `FUNCTIONAL_AUDIT.md`, `ARCHITECTURE_AUDIT.md`, `DATA_MODEL_AUDIT.md`, `SECURITY_ROBUSTNESS_AUDIT.md`, `TESTING_AUDIT.md`, `CODE_HOTSPOTS_AUDIT.md`

## Verdict global

**Acceptable avec réserves** — L'ensemble des quatre documents est fidèle à l'amont. Les claims structurants (nature de la lib, trois fonctions pures, défaut de `Book` sans garde, absence de persistance, couverture de tests lacunaire, risques TOCTOU) sont tous traçables avec citations de fichier et de ligne. Un exemple non tracé s'est glissé dans `CDC_FONCTIONNEL.md`, la table de risques de `PROJECT_CONTEXT.md` sous-estime `Remaining`, et un typo est présent dans `CAHIER_RECETTE.md`. Corrections rapides — ne requièrent pas de refonte.

---

## Problèmes bloquants

_Aucun défaut de tracé fonctionnel bloquant identifié._

---

## Problèmes mineurs

### P1 — Exemple inventé : « surf » dans `CDC_FONCTIONNEL.md`
- **Fichier** : `CDC_FONCTIONNEL.md`, section « Contexte métier » — *"créneaux d'activités nautiques (plongée, surf, etc.)"*
- **Problème** : le terme « surf » n'est tracé dans aucun artefact amont. L'activité « Plongée » est la seule nommée dans le code (`booking_test.go:9`, confirmé par `SECURITY_ROBUSTNESS_AUDIT.md`, ligne 15 : *"La seule constante textuelle est `"Plongée"` dans la fixture de test"*). Le README dit « activités nautiques » sans en lister d'autres.
- **Gravité** : faible — c'est un exemple illustratif, pas une règle métier ; mais un relecteur ne peut pas confirmer que « surf » est dans le scope du pilote.
- **Correction attendue** : remplacer « (plongée, surf, etc.) » par « (plongée, etc.) » ou simplement « (activités nautiques) ».

### P2 — Risque de `Remaining` sous-estimé dans la table de `PROJECT_CONTEXT.md`
- **Fichier** : `PROJECT_CONTEXT.md`, tableau « État courant », ligne `Remaining`, colonne Risque.
- **Problème** : la cellule indique *"Aucun dans la forme actuelle — fonction pure, pas de garde"*. Or `WORKFLOW_CONSULTER_DISPONIBILITE.md` liste explicitement ce risque : *"`Remaining` peut retourner un entier négatif sans avertissement"* (section Risques, ligne 54). `SECURITY_ROBUSTNESS_AUDIT.md` le confirme. La prose du document couvre ce point (section Risques prioritaires), mais la table — qui est lue en premier — est inexacte.
- **Gravité** : faible — le risque est bien couvert ailleurs dans le même document ; la table est juste incomplète sur ce point.
- **Correction attendue** : amender la cellule en quelque chose comme *"Peut retourner une valeur négative si l'état est invalide (Booked > Capacity) — pas de garde"*.

### P3 — Typo dans `CAHIER_RECETTE.md`
- **Fichier** : `CAHIER_RECETTE.md`, tableau récapitulatif de couverture, ligne `Book` annulation implicite (-n), colonne Recommandation.
- **Problème** : *"Formalizerr ou interdire"* — double « r » dans « Formalizerr ».
- **Correction attendue** : corriger en *"Formaliser ou interdire"*.

---

## Points vérifiés et corrects

- **Tracé de `Book` avec gardes** : `Book` est correctement décrit comme validant ses pré-conditions (`n > 0`, `n ≤ Remaining(s)`) avec références à `internal/booking/booking.go:36-45` (gardes aux lignes 37-41, incrément à ligne 43). Tracé à `FUNCTIONAL_AUDIT.md`, `WORKFLOW_RESERVER_PLACES.md`, `TESTING_AUDIT.md`.
- **Absence de persistance** : correctement traitée dans `PROJECT_CONTEXT.md`, `CDC_FONCTIONNEL.md` et `CAHIER_RECETTE.md`. Tracée à `ARCHITECTURE_AUDIT.md` et `DATA_MODEL_AUDIT.md`.
- **Risque TOCTOU (`IsAvailable → Book` non atomique)** : correctement mentionné dans `CDC_FONCTIONNEL.md` (WF2) et `CAHIER_RECETTE.md` (WS-1). Tracé à `WORKFLOW_CONSULTER_DISPONIBILITE.md` (section Risques) et `CODE_HOTSPOTS_AUDIT.md`.
- **Couverture de tests** : `CAHIER_RECETTE.md` liste fidèlement les cas non couverts (TC-1.2, TC-2.2, TC-3.3, TC-3.4) en les faisant correspondre aux lacunes identifiées dans `TESTING_AUDIT.md` et `CODE_HOTSPOTS_AUDIT.md`.
- **`n` négatif (annulation implicite)** : correctement qualifié de comportement implicite non documenté dans `CDC_FONCTIONNEL.md` et `CAHIER_RECETTE.md`. Tracé à `WORKFLOW_RESERVER_PLACES.md` et `FUNCTIONAL_AUDIT.md`.
- **Hypothèses marquées** : les mentions de futures fonctionnalités (HTTP, persistance, catalogue) sont toutes qualifiées d'`HYPOTHÈSE` ou de constat d'absence, conformément au signal `HYPOTHÈSE` des audits.
- **Limites assumées** : tous les documents signalent explicitement le statut « pilote de test » et les capacités absentes. Aucun des documents ne prétend couvrir plus que ce qui est prouvé.
- **Confiance déclarée** : les niveaux medium/high sont cohérents avec ceux déclarés dans les artefacts amont.
- **SHA et module Go** : `8122e2d` et `github.com/arthurfromtahiti/shift-pilot-go` tracés à `CARTE_DES_DOMAINES.md`.
- **Contenu de `sample()`** : `{ID:1, Activity:"Plongée", Start:time.Now(), Capacity:10, Booked:4}` dans `CARTOGRAPHIE_CODE.md` est tracé mot pour mot à `TESTING_AUDIT.md` (ligne 15).
- **Zéro dépendance externe** : correctement répercuté dans `CARTOGRAPHIE_CODE.md`. Tracé à `ARCHITECTURE_AUDIT.md`.
- **Calculs du cahier de recette** : TC-3.3 (`Book(Slot{Cap:10, Booked:4}, 7)` → `Booked=11`) : 4+7=11 ✓. TC-1.3 (`Remaining(Slot{Cap:10, Booked:15})` → `-5`) : 10-15=-5 ✓.

---

## Recommandations de correction

1. **`CDC_FONCTIONNEL.md` — Contexte métier** : remplacer `(plongée, surf, etc.)` par `(plongée, etc.)`. Une ligne à changer, aucune réécriture.
2. **`PROJECT_CONTEXT.md` — table État courant, ligne Remaining** : remplacer la cellule Risque par *"Peut retourner une valeur négative si `Booked > Capacity` — pas de garde"*. Alignement avec `WORKFLOW_CONSULTER_DISPONIBILITE.md` (section Risques).
3. **`CAHIER_RECETTE.md` — tableau de couverture** : corriger `Formalizerr` → `Formaliser`.

---

## Résultat de relecture — Cycle 2 (2026-08-05)

**APPROUVÉ** — Les trois corrections ont été vérifiées et appliquées par le rédacteur :

- ✅ P1 — `CDC_FONCTIONNEL.md` : `(plongée, surf, etc.)` → `(plongée, etc.)` — confirmé.
- ✅ P2 — `PROJECT_CONTEXT.md` : cellule Risque de `Remaining` corrigée en *"Peut retourner une valeur négative si `Booked > Capacity` — pas de garde"* — confirmé.
- ✅ P3 — `CAHIER_RECETTE.md` : `Formalizerr` → `Formaliser` — confirmé.

Aucun nouveau défaut introduit. Les documents de référence sont validés.
