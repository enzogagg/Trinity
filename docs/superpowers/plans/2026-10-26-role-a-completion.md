# Role A Completion Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Terminer les livrables infrastructure, sécurité, déploiement et soutenance du rôle A pour Trinity Part One.

**Architecture:** GitLab CE tourne sur la VM `gitlab` (`192.168.252.2`) et GitLab Runner avec Docker sur la VM `runner` (`192.168.252.3`). Le pipeline contrôle les secrets, le lint et les tests, construit les images, puis effectue un déploiement blue/green derrière nginx.

**Tech Stack:** GitLab CE 19.5, GitLab Runner 19.5, Multipass, Docker Compose, nginx, PostgreSQL, Gitleaks.

**Spec:** `docs/README.md`

## Global Constraints

- Aucun secret réel ne doit être ajouté au dépôt ou à son historique.
- `main` reste protégée et le force push reste désactivé.
- Chaque changement passe par une branche et une merge request GitLab.
- Les validations doivent être reproductibles pendant la soutenance.

---

### Task 1: Prouver la stratégie de branches

**Files:**
- Modify: `docs/README.md`

**Interfaces:**
- Consumes: dépôt GitLab auto-hébergé et branche protégée `main`
- Produces: une merge request visible et un historique de validation

- [ ] **Step 1: Créer une branche pour la documentation infrastructure**

```bash
git switch -c docs/infrastructure
```

- [ ] **Step 2: Pousser la branche vers GitLab**

```bash
git push -u gitlab docs/infrastructure
```

- [ ] **Step 3: Vérifier les règles de `main`**

Dans GitLab, ouvrir `Settings > Repository > Protected branches` et vérifier que `main` est protégée et que le force push est désactivé.

- [ ] **Step 4: Utiliser cette branche pour les tâches 2 et 3**

Les documents suivants sont commités sur `docs/infrastructure`, puis intégrés par merge request.

---

### Task 2: Documenter GitLab et le runner

**Files:**
- Create: `docs/runner-installation.md`
- Modify: `docs/README.md`

**Interfaces:**
- Consumes: VM GitLab `192.168.252.2`, VM runner `192.168.252.3`
- Produces: procédure reproductible d'installation et d'enregistrement du runner

- [ ] **Step 1: Documenter les ressources et le réseau**

Écrire les caractéristiques utilisées : Ubuntu Server 24.04, GitLab avec 2 vCPU/4 Go/25 Go, runner avec 2 vCPU/3 Go/20 Go, et réseau Multipass privé.

- [ ] **Step 2: Documenter l'installation GitLab**

Inclure l'ajout du dépôt GitLab CE, `EXTERNAL_URL=http://192.168.252.2`, le profil léger, le registry `http://192.168.252.2:5050` et `gitlab-ctl reconfigure`.

- [ ] **Step 3: Documenter l'installation du runner**

Inclure Docker, le dépôt GitLab Runner, l'enregistrement avec l'exécuteur Docker et le montage suivant :

```toml
volumes = ["/var/run/docker.sock:/var/run/docker.sock", "/cache"]
```

- [ ] **Step 4: Documenter les vérifications**

```bash
sudo gitlab-runner verify
sudo systemctl status gitlab-runner --no-pager
curl -I http://192.168.252.2:5050/v2/
```

Résultats attendus : runner valide et actif, registry répondant `401 Unauthorized` avec `Docker-Distribution-Api-Version: registry/2.0`.

- [ ] **Step 5: Ajouter le document au sommaire**

Ajouter le lien `runner-installation.md` dans `docs/README.md`.

---

### Task 3: Documenter l'architecture et le pipeline

**Files:**
- Create: `docs/architecture.md`
- Modify: `docs/README.md`

**Interfaces:**
- Consumes: `.gitlab-ci.yml`, `deploy/compose.dev.yml`, `deploy/compose.infra.yml`, `deploy/compose.app.yml`, `deploy/deploy.sh`
- Produces: justification des technologies, différences dev/prod et diagramme CI/CD complet

- [ ] **Step 1: Décrire les technologies et leurs rôles**

Documenter Go pour l'API, React/Vite pour le web, PostgreSQL pour les données, nginx pour le proxy et GitLab Runner pour l'automatisation.

- [ ] **Step 2: Expliquer les différences dev/prod**

Décrire le hot reload et les volumes en développement, puis les builds multi-stage, images minimales, variables GitLab et redémarrage automatique en production.

- [ ] **Step 3: Ajouter le diagramme Mermaid**

```mermaid
flowchart LR
  Push[Push ou merge sur GitLab] --> Secrets[Gitleaks]
  Secrets --> Lint[Lint API et Web]
  Lint --> Tests[Tests API et Web + JUnit]
  Tests --> Build[Build et push des images]
  Build --> Deploy[Déploiement automatique]
  Deploy --> New[Blue ou Green + healthchecks]
  New --> Proxy[Bascule nginx]
  Proxy --> Stop[Arrêt de l'ancienne version]
```

- [ ] **Step 4: Ajouter le document au sommaire**

Ajouter le lien `architecture.md` dans `docs/README.md`.

- [ ] **Step 5: Vérifier les documents et commiter**

```bash
git diff --check
git add docs
git commit -m "docs: document infrastructure and CI/CD architecture"
git push gitlab docs/infrastructure
```

- [ ] **Step 6: Créer et fusionner la merge request**

Créer une merge request `docs/infrastructure` vers `main`, demander la relecture du binôme, attendre un pipeline vert, puis fusionner sans force push.

---

### Task 4: Valider les environnements après redémarrage

**Files:**
- Modify: `docs/deployment-evidence.md`

**Interfaces:**
- Consumes: deux VM et conteneurs avec `restart: unless-stopped`
- Produces: preuve de résilience du système

- [ ] **Step 1: Redémarrer les VM**

```bash
multipass restart gitlab runner
```

- [ ] **Step 2: Vérifier GitLab et le runner**

```bash
curl -I http://192.168.252.2
multipass exec runner -- sudo gitlab-runner verify
```

- [ ] **Step 3: Vérifier l'application**

```bash
curl --fail http://192.168.252.3/api/health
multipass exec runner -- sudo docker ps
```

Résultats attendus : redirection HTTP GitLab, runner valide, API `{"status":"ok"}` et conteneurs en état `healthy`.

- [ ] **Step 4: Consigner le test**

Ajouter la date et le résultat du redémarrage à `docs/deployment-evidence.md` via une nouvelle branche et une merge request.

---

### Task 5: Préparer la démonstration du rôle A

**Files:**
- Create: `docs/keynote-checklist.md`

**Interfaces:**
- Consumes: pipeline vert, runner en ligne et documents précédents
- Produces: déroulé reproductible pour la soutenance

- [ ] **Step 1: Préparer les captures**

Capturer le runner `Online`, le pipeline complet vert, le rapport JUnit, le job Gitleaks et les variables GitLab dont les valeurs sont masquées.

- [ ] **Step 2: Préparer la démonstration blue/green**

Lancer la surveillance suivante avant un nouveau pipeline :

```bash
while true; do
  curl -s -o /dev/null -w "%{http_code} " http://192.168.252.3/api/health
  sleep 0.2
done
```

Montrer `active=green -> deploying blue`, `traffic switched to blue` et uniquement des réponses HTTP 200.

- [ ] **Step 3: Préparer quatre réponses techniques**

Expliquer pourquoi GitLab et le runner sont séparés, pourquoi les secrets vivent dans GitLab, pourquoi le socket Docker est puissant et risqué, et comment les healthchecks empêchent la bascule vers une version défectueuse.

- [ ] **Step 4: Faire une répétition chronométrée**

Présenter l'infrastructure, lancer un pipeline, montrer les logs et conclure sur la sécurité et le zero downtime en moins de dix minutes.
