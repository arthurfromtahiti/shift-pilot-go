# INDEX — .onboarding/shift-pilot-go

> Généré le 2026-08-05 lors de la publication initiale des artefacts d'onboarding.
> Dépôt : `github.com/arthurfromtahiti/shift-pilot-go` — branche : `onboarding/artifacts`

| type | domaine | workflow | dépôt | fichier | date | version SHA | niveau de preuve | titre |
|---|---|---|---|---|---|---|---|---|
| domaine | réservation-creneaux | | shift-pilot-go | domaines/CARTE_DES_DOMAINES.md | 2026-08-05 | 21301ec26fcdc913cc056e212e84fc1b8ada117a | contient une hypothèse | Carte des domaines fonctionnels — bibliothèque de réservation de créneaux nautiques |
| workflow | réservation-creneaux | CONSULTER_DISPONIBILITE | shift-pilot-go | workflows/WORKFLOW_CONSULTER_DISPONIBILITE.md | 2026-08-05 | 21301ec26fcdc913cc056e212e84fc1b8ada117a | contient une hypothèse | Consulter la disponibilité d'un créneau — API de bibliothèque Go |
| workflow | réservation-creneaux | RESERVER_PLACES | shift-pilot-go | workflows/WORKFLOW_RESERVER_PLACES.md | 2026-08-05 | 21301ec26fcdc913cc056e212e84fc1b8ada117a | contient une hypothèse | Réserver des places sur un créneau — API de bibliothèque Go |
| audit | réservation-creneaux | | shift-pilot-go | audits/ARCHITECTURE_AUDIT.md | 2026-08-05 | 21301ec26fcdc913cc056e212e84fc1b8ada117a | contient une hypothèse | Audit architecture — structure minimale d'une bibliothèque Go de réservation |
| audit | réservation-creneaux | | shift-pilot-go | audits/CODE_HOTSPOTS_AUDIT.md | 2026-08-05 | 21301ec26fcdc913cc056e212e84fc1b8ada117a | contient une hypothèse | Audit points chauds — risques et lacunes dans le noyau de disponibilité |
| audit | réservation-creneaux | | shift-pilot-go | audits/DATA_MODEL_AUDIT.md | 2026-08-05 | 21301ec26fcdc913cc056e212e84fc1b8ada117a | contient une hypothèse | Audit modèle de données — entité Slot et absence de persistance |
| audit | réservation-creneaux | | shift-pilot-go | audits/FUNCTIONAL_AUDIT.md | 2026-08-05 | 21301ec26fcdc913cc056e212e84fc1b8ada117a | contient une hypothèse | Audit fonctionnel — couverture des cas d'usage et comportements |
| audit | réservation-creneaux | | shift-pilot-go | audits/SECURITY_ROBUSTNESS_AUDIT.md | 2026-08-05 | 21301ec26fcdc913cc056e212e84fc1b8ada117a | contient une hypothèse | Audit sécurité et robustesse — bibliothèque sans exposition externe |
| audit | réservation-creneaux | | shift-pilot-go | audits/TESTING_AUDIT.md | 2026-08-05 | 21301ec26fcdc913cc056e212e84fc1b8ada117a | contient une hypothèse | Audit tests — couverture unitaire du noyau de disponibilité |
| document | réservation-creneaux | | shift-pilot-go | documents/CDC_FONCTIONNEL.md | 2026-08-05 | 21301ec26fcdc913cc056e212e84fc1b8ada117a | contient une hypothèse | Cahier des charges fonctionnel — fonctionnalités actuelles et intentions futures |
| document | | | shift-pilot-go | documents/CARTOGRAPHIE_CODE.md | 2026-08-05 | 21301ec26fcdc913cc056e212e84fc1b8ada117a | contient une hypothèse | Cartographie du code — anatomie de la bibliothèque Go (57 lignes) |
| document | | | shift-pilot-go | documents/PROJECT_CONTEXT.md | 2026-08-05 | 21301ec26fcdc913cc056e212e84fc1b8ada117a | contient une hypothèse | Contexte projet — pilote de test de réservation nautique SHIFT/Paperclip |
| document | réservation-creneaux | | shift-pilot-go | documents/CAHIER_RECETTE.md | 2026-08-05 | 21301ec26fcdc913cc056e212e84fc1b8ada117a | contient une hypothèse | Cahier de recette — cas de test dérivés des deux workflows de réservation |
| journal-fabrication | réservation-creneaux | | shift-pilot-go | relectures/RELECTURE_CARTE_DES_DOMAINES.md | 2026-08-05 | 21301ec26fcdc913cc056e212e84fc1b8ada117a | contient une hypothèse | Relecture de la carte des domaines |
| journal-fabrication | | | shift-pilot-go | relectures/RELECTURE_DOCUMENTS_REFERENCE.md | 2026-08-05 | 21301ec26fcdc913cc056e212e84fc1b8ada117a | contient une hypothèse | Relecture des documents de référence |
| journal-fabrication | réservation-creneaux | | shift-pilot-go | relectures/RELECTURE_SUITE_AUDITS_PROJET.md | 2026-08-05 | 21301ec26fcdc913cc056e212e84fc1b8ada117a | contient une hypothèse | Relecture de la suite d'audits du projet |
| journal-fabrication | réservation-creneaux | CONSULTER_DISPONIBILITE | shift-pilot-go | relectures/RELECTURE_WORKFLOW_CONSULTER_DISPONIBILITE.md | 2026-08-05 | 21301ec26fcdc913cc056e212e84fc1b8ada117a | contient une hypothèse | Relecture du workflow Consulter la disponibilité |
| journal-fabrication | réservation-creneaux | RESERVER_PLACES | shift-pilot-go | relectures/RELECTURE_WORKFLOW_RESERVER_PLACES.md | 2026-08-05 | 21301ec26fcdc913cc056e212e84fc1b8ada117a | contient une hypothèse | Relecture du workflow Réserver des places |

> **Note `niveau de preuve`** : tous les artefacts contiennent au moins une mention `HYPOTHÈSE` ou `INCONNU`. C'est attendu : le dépôt est un pilote minimal sans serveur, sans base de données, sans CI/CD — de nombreux éléments de contexte d'exécution restent non vérifiables.
>
> **Note `relectures/`** : ce dossier est un journal de fabrication (traçabilité des verdicts de relecture), pas de la connaissance produit. Les compilateurs de contexte doivent l'exclure par défaut.
