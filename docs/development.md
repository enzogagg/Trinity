# Développement local

## Prérequis

| Outil | Version | Usage |
| --- | --- | --- |
| Docker et Docker Compose | Docker 27 ou plus | Environnement complet (obligatoire) |
| Go | 1.27 | Tests et lint de l'API hors conteneur (facultatif) |
| Node.js | 22.12 ou plus | Tests et lint du Web hors conteneur (facultatif) |

Sous Windows, le dépôt force les fins de ligne LF (`.gitattributes`) : les scripts shell et Dockerfiles restent exécutables dans les conteneurs Linux.

## Lancer l'environnement de développement

```bash
cp .env.example .env.dev        # valeurs locales, jamais commitées
docker compose -f deploy/compose.dev.yml up --build
```

| URL | Service |
| --- | --- |
| http://localhost:5173 | Web (Vite, rechargement à chaud) |
| http://localhost:5173/api/health | API via le proxy Vite |
| http://localhost:8080/api/health | API en direct |
| http://localhost:8080/api/health/db | État de la connexion PostgreSQL |

Les sources `api/` et `web/` sont montées dans les conteneurs : une modification du Web est visible immédiatement ; pour l'API, redémarrer le service :

```bash
docker compose -f deploy/compose.dev.yml restart api
```

Arrêt, avec ou sans suppression des données de développement :

```bash
docker compose -f deploy/compose.dev.yml down       # conserve la base
docker compose -f deploy/compose.dev.yml down -v    # supprime aussi le volume pgdata-dev
```

## Lint et tests

Les mêmes commandes que la pipeline :

```bash
# API
cd api
go vet ./...
go test -v -cover ./...
docker run --rm -v "$PWD:/app" -w /app golangci/golangci-lint:v2.15.0 golangci-lint run ./...

# Web
cd web
npm ci
npm run lint                     # tsc --noEmit puis ESLint
npm test                         # Vitest
npx vitest run --coverage        # avec rapport de couverture
npm run build
```

## Tester les images de production en local

```bash
docker build -f api/Dockerfile.prod -t trinity/api:local api
docker build -f web/Dockerfile.prod -t trinity/web:local web
```

Pour reproduire la production (proxy, base et une couleur) :

```bash
cd deploy
export CI_REGISTRY_IMAGE=trinity TAG=local COLOR=blue \
       POSTGRES_DB=trinity POSTGRES_USER=trinity POSTGRES_PASSWORD=localpassword
docker compose -f compose.infra.yml up -d --wait
docker compose -p trinity-blue -f compose.app.yml up -d --wait
docker exec trinity-proxy sh -c "echo 'set \$color blue;' > /etc/nginx/active/color.conf && nginx -s reload"
curl http://localhost/api/health/db

# nettoyage
docker compose -p trinity-blue -f compose.app.yml down
docker compose -f compose.infra.yml down -v
```

## Avant de pousser

```bash
sh scripts/check-conventions.sh   # nom de branche et messages de commit
```

Les conventions de branches, de commits et de merge requests sont décrites dans [CONTRIBUTING.md](../CONTRIBUTING.md).
