# Trinity

Plateforme d'une chaîne d'épiceries : API REST en Go, back-office web en React et, plus tard, application mobile. Cette première partie met en place l'environnement DevOps : conteneurs dev et prod, pipeline GitLab CI/CD et déploiement automatique sans interruption.

## Démarrage rapide

```bash
cp .env.example .env.dev
docker compose -f deploy/compose.dev.yml up --build
```

Web sur http://localhost:5173, API sur http://localhost:8080/api/health.

## Pipeline

Chaque commit passe par : détection de secrets, lint, tests avec rapports JUnit et couverture, build des images. Chaque fusion sur `main` est déployée automatiquement en blue/green puis vérifiée.

## Documentation

- [Architecture, technologies et pipeline CI/CD](docs/architecture.md)
- [Développement local](docs/development.md)
- [Installation de GitLab et du runner](docs/runner-installation.md)
- [Sécurité](docs/security.md)
- [Conventions de contribution](CONTRIBUTING.md)
