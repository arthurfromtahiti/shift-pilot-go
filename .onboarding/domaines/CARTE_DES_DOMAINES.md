# Carte des domaines — shift-pilot-go

> **Confiance globale du livrable : low.**
> Matière très pauvre : le dépôt est un *pilote de test* d'un seul package Go
> (`internal/booking`, 3 fonctions pures + tests), sans serveur, sans route,
> sans persistance, sans authentification. La carte reflète **honnêtement** ce
> qui existe : **un seul domaine prouvé**. Le reste (paiement, catalogue
> d'activités, comptes clients, notifications…) est **absent du code** et vit en
> *Incertitudes*, pas en domaines inventés.
>
> - Dépôt : `github.com/arthurfromtahiti/shift-pilot-go`
> - Branche par défaut : `main`
> - SHA de tête analysé : `8122e2ddaa200f5e30b4a50364c0d1758fcec72f`
> - Accès base de données : non fourni (sans objet — aucune persistance dans le code).

## Nature du projet

Bibliothèque Go (module `github.com/arthurfromtahiti/shift-pilot-go`, Go 1.21,
**bibliothèque standard uniquement**) présentée par son README comme un « pilote
de test SHIFT/Paperclip — réservation d'activités nautiques (Polynésie
française) ». Le seul code métier réel est le package `internal/booking` : un
modèle `Slot` (créneau d'activité avec capacité et places réservées) et trois
fonctions pures (`Remaining`, `IsAvailable`, `Book`). Aucun point d'entrée
exécutable (`func main` absent), aucun serveur HTTP, aucune route, aucune couche
de stockage. C'est un **noyau de logique de disponibilité**, pas encore une
application : la finalité annoncée (réservation) n'est représentée qu'à l'état de
germe.

## Domaines

### Réservation & disponibilité de créneaux (`reservation-creneaux`)
- **Catégorie** : métier
- **Priorité** : cœur
- **Confiance** : medium
- **Description** : cycle de vie de la disponibilité d'un créneau d'activité —
  calcul des places restantes, test de disponibilité, et enregistrement d'une
  réservation de `n` places sur un créneau. C'est la seule capacité métier
  matérialisée dans le code, et elle porte la raison d'être annoncée du projet
  (réservation d'activités nautiques).
- **Entités** : `Slot` (champs `ID`, `Activity`, `Start time.Time`, `Capacity`,
  `Booked`) — `internal/booking/booking.go:6`.
- **Routes / points d'entrée** : *aucune* — pas de serveur ni de handler. Le
  domaine n'est exposé que comme API de package Go : `Remaining(Slot) int`
  (`booking.go:15`), `IsAvailable(Slot) bool` (`booking.go:20`),
  `Book(Slot, int) (Slot, error)` (`booking.go:36`). À noter : `Book` retourne
  un `Slot` **par valeur** et une `error`, avec deux gardes :
  `ErrInvalidBookingCount` si n ≤ 0, `ErrCapacityExceeded` si n dépasse
  les places restantes.
- **Indices de rattachement** : package `internal/booking`, type `Slot`,
  identifiants `Capacity`/`Booked`/`Remaining`/`IsAvailable`/`Book`, terme
  `Activity`. Ne matche que `internal/booking/*.go`.
- **Types de workflows attendus** : consulter la disponibilité d'un créneau,
  réserver des places, (à terme) libérer/annuler des places. Aujourd'hui limité
  aux trois opérations pures ci-dessus.
- **Preuves** : `internal/booking/booking.go` (modèle + 3 fonctions),
  `internal/booking/booking_test.go` (tests `TestRemaining`, `TestIsAvailable`,
  `TestBook` — assertions lues, `VÉRIFIÉ_CODE` ; **non exécutées** : la toolchain
  Go n'est pas disponible dans l'environnement d'analyse → statut d'exécution
  `INCONNU`).
- **Dépend de la base** : non — aucune persistance, aucun champ
  `content`/`layout`/`blocks`/`config`, aucun renderer récursif de structure
  arborescente. Aucun des trois signaux du §6 n'est présent.

## Incertitudes

- **Un seul domaine, volontairement.** Le README annonce une application de
  « réservation d'activités nautiques » ; le code n'en couvre qu'une brique. Les
  domaines qu'un tel produit finira presque certainement par avoir — **catalogue
  d'activités**, **comptes clients / réservants**, **paiement**, **planning /
  agenda des créneaux**, **notifications de confirmation** — ne sont représentés
  par **aucune entité ni route** aujourd'hui. Ils sont donc *hors carte* (règle :
  un domaine sans preuve concrète n'existe pas). À revérifier dès que du code
  arrive. Question exploitable : quelle est la portée cible du pilote — reste-t-il
  une lib de logique de disponibilité, ou doit-il grandir en service complet ?
- **Aucune persistance ni exposition.** Pas de `main`, pas de HTTP, pas de base.
  Impossible de déduire des workflows de bout en bout ou une frontière
  d'intégration à ce stade — l'Analyste de workflows aura peu de matière. À
  signaler au Chef d'Onboarding.
- ~~**`Book` sans garde de capacité.**~~ — **RÉSOLU** : `Book` valide désormais que `n > 0` (`ErrInvalidBookingCount`) et `n ≤ Remaining(s)` (`ErrCapacityExceeded`) avant incrément. Un créneau ne peut plus être sur-réservé via `Book`. Les tests `TestBookZero`, `TestBookNegative`, `TestBookCapacityExceeded`, `TestBookExactCapacity` couvrent ces cas.
- **« SHIFT/Paperclip ».** Le README qualifie le dépôt de « pilote de test
  SHIFT/Paperclip ». Rien dans le code ne relie à un système SHIFT externe ;
  terme métier non tranché ici (réservé au board).
