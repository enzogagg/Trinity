# Contribuer à Trinity

Ces règles sont vérifiées automatiquement par le job `conventions` sur chaque merge request (`scripts/check-conventions.sh`).

## Branches

`main` est protégée : pas de push direct ni de force push. Toute modification passe par une branche courte et une merge request.

Format : `<type>/<description-en-kebab-case>`

| Type | Usage | Exemple |
| --- | --- | --- |
| `feat` | Nouvelle fonctionnalité | `feat/api-products` |
| `fix` | Correction de bug | `fix/docker-build-context` |
| `docs` | Documentation | `docs/pipeline-diagram` |
| `refactor` | Restructuration sans changement de comportement | `refactor/api-handlers` |
| `test` | Ajout ou correction de tests | `test/web-health` |
| `ci` | Pipeline GitLab, runners | `ci/needs-artifacts` |
| `build` | Dockerfiles, compose, dépendances | `build/upgrade-go` |
| `chore` | Maintenance (gitignore, outillage) | `chore/repo-hygiene` |
| `hotfix` | Correction urgente en production | `hotfix/db-connection` |

Une branche = un sujet. Elle est supprimée après la fusion.

## Commits

Format [Conventional Commits](https://www.conventionalcommits.org/fr/v1.0.0/) :

```
<type>(<scope>): <description>
```

- **type** : `feat`, `fix`, `docs`, `style`, `refactor`, `perf`, `test`, `build`, `ci`, `chore`, `revert`
- **scope** (facultatif) : `api`, `web`, `deploy`, `docker`, `ci`…
- **description** : en anglais, à l'impératif, sans majuscule ni point final
- `!` après le type ou le scope signale un changement cassant : `feat(api)!: rename health endpoint`

Exemples :

```
feat(api): connect to PostgreSQL and expose database health endpoint
fix(docker): add dockerignore files and use npm ci
ci: pass image digests between jobs via dotenv artifacts
docs: add CI/CD pipeline diagram
```

Un commit = un changement cohérent qui compile et passe les tests. Ne jamais commiter de secret : `.env.dev` est ignoré, les valeurs de production sont des variables CI/CD GitLab (voir `docs/security.md`).

## Merge requests

1. Créer la MR vers `main` dès que la branche est poussée (en *Draft* si elle n'est pas prête).
2. Remplir le modèle de description (`.gitlab/merge_request_templates/Default.md`).
3. La pipeline doit être verte : conventions, gitleaks, lint, tests, build.
4. **Relecture croisée obligatoire** : la MR est relue et approuvée par l'autre membre de l'équipe.
5. Fusion par l'auteur après approbation, avec suppression de la branche source.

La fusion sur `main` déclenche automatiquement le déploiement en production puis la vérification (`deploy-prod`, `verify-prod`).

## Vérifier en local avant de pousser

```bash
# API
cd api && go vet ./... && go test ./...

# Web
cd web && npm ci && npm run lint && npm test

# Conventions de la branche courante
sh scripts/check-conventions.sh
```
