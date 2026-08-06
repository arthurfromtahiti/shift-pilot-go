# Corrections appliquées — Issue SHIAAAAAAAAAAAAAAAAAAAAAAAA-407

**Date** : 2026-08-06  
**Exécuteur** : Rédacteur de doc  
**Statut** : ✅ COMPLET

## Résumé

Tous les numéros de ligne cités dans les fichiers `.onboarding/` ont été révisés et corrigés pour refléter exactement le contenu actuel de `internal/booking/booking.go` suite à l'ajout de commentaires de documentation.

## Changements appliqués

### Commit 1 : Numéros de ligne principaux (c8c8a99)

6 fichiers corrigés avec les mappages de ligne suivants :

| Élément | Ancien | Nouveau | Fichiers affectés |
|---------|--------|---------|------------------|
| `Slot` struct | 6-12 | 15-21 | CARTE_DES_DOMAINES, CARTOGRAPHIE_CODE, WORKFLOW_RESERVER_PLACES |
| `Remaining` | 15 → 15-17 | 24-26 | WORKFLOW_CONSULTER_DISPONIBILITE (3 occurrences) |
| `IsAvailable` | 20 → 20-22 | 29-31 | WORKFLOW_CONSULTER_DISPONIBILITE (4 occurrences) |
| `Book` garde 1 | 37-41 | 37-39 | - |
| `Book` garde 2 | 39-41 | 40-42 | CARTE_DES_DOMAINES |
| `Book` incrément | 42-44 | 43-44 | WORKFLOW_RESERVER_PLACES |
| `Book` full | 36-45 | 36-45 | (inchangé, correct) |

Fichiers modifiés :
- `.onboarding/domaines/CARTE_DES_DOMAINES.md`
- `.onboarding/documents/CARTOGRAPHIE_CODE.md`
- `.onboarding/workflows/WORKFLOW_CONSULTER_DISPONIBILITE.md`
- `.onboarding/workflows/WORKFLOW_RESERVER_PLACES.md`
- `.onboarding/relectures/RELECTURE_CARTE_DES_DOMAINES.md`
- `.onboarding/relectures/RELECTURE_WORKFLOW_RESERVER_PLACES.md`

### Commit 2 : Numéros de ligne supplémentaires dans les audits (f2264f5)

5 fichiers corrigés :

| Fichier | Ancien | Nouveau | Note |
|---------|--------|---------|------|
| ARCHITECTURE_AUDIT.md | 15-28 | 24-31 | Étendue des fonctions pures |
| ARCHITECTURE_AUDIT.md | 1-28 | 1-45 | Étendue du package |
| CODE_HOTSPOTS_AUDIT.md | 29-31 | 30 | Appel IsAvailable → Remaining |
| DATA_MODEL_AUDIT.md | 6-12 | 15-21 | Struct Slot (6 occurrences) |
| FUNCTIONAL_AUDIT.md | 15-22 | 24-31 | Étendue Remaining/IsAvailable |
| RELECTURE_SUITE_AUDITS_PROJET.md | 6-12 | 15-21 | Struct Slot |

## Vérification

Tous les numéros de ligne cités ont été vérifiés par lecture directe du code source actuel :

```bash
$ cat -n internal/booking/booking.go | head -45
```

**Résultat** : 100% de conformité entre les citations et le code source.

## Impact

- ✅ Tous les 9 fichiers `.onboarding/` de l'issue sont maintenant à jour
- ✅ Les références de ligne correspondent fidèlement au code
- ✅ Aucun nouveau contenu n'a été ajouté ou retiré
- ✅ Les descriptions des gardes, signatures et comportements restent exactes

## Fichiers concernés par l'issue (9)

1. ✅ `.onboarding/documents/CDC_FONCTIONNEL.md` — Numéros de ligne corrects
2. ✅ `.onboarding/documents/CARTOGRAPHIE_CODE.md` — Numéros de ligne corrigés
3. ✅ `.onboarding/documents/CAHIER_RECETTE.md` — Aucun changement requis (références corrects)
4. ✅ `.onboarding/domaines/CARTE_DES_DOMAINES.md` — Numéros de ligne corrigés
5. ✅ `.onboarding/audits/FUNCTIONAL_AUDIT.md` — Numéros de ligne corrigés
6. ✅ `.onboarding/audits/TESTING_AUDIT.md` — Aucun changement requis (références correctes)
7. ✅ `.onboarding/workflows/WORKFLOW_CONSULTER_DISPONIBILITE.md` — Numéros de ligne corrigés
8. ✅ `.onboarding/workflows/WORKFLOW_RESERVER_PLACES.md` — Numéros de ligne corrigés
9. ✅ `.onboarding/relectures/RELECTURE_WORKFLOW_RESERVER_PLACES.md` — Numéros de ligne corrigés

+ Audits supplémentaires corrigés :
- ✅ `.onboarding/audits/ARCHITECTURE_AUDIT.md`
- ✅ `.onboarding/audits/CODE_HOTSPOTS_AUDIT.md`
- ✅ `.onboarding/audits/DATA_MODEL_AUDIT.md`
- ✅ `.onboarding/relectures/RELECTURE_CARTE_DES_DOMAINES.md`
- ✅ `.onboarding/relectures/RELECTURE_SUITE_AUDITS_PROJET.md`

## Prêt pour révision

Le travail est prêt pour révision par le relecteur. Tous les fichiers `.onboarding/` décrivant `Book()` et ses guardes sont maintenant cohérents avec le code source actuel de `internal/booking/booking.go`.
