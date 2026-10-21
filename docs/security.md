# Sécurité et gestion des secrets

## Mesures appliquées

- Les fichiers `.env` locaux sont exclus par `.gitignore`.
- Seul `.env.example`, sans secret réel, est versionné.
- Les secrets de production sont enregistrés dans les variables CI/CD GitLab.
- `POSTGRES_PASSWORD` est masqué et protégé dans GitLab.
- Le pipeline récupère les identifiants temporaires du Container Registry depuis les variables prédéfinies de GitLab.
- Gitleaks analyse les fichiers et l'historique Git à chaque pipeline.
- Le pipeline s'arrête avant les tests, builds et déploiements si Gitleaks détecte un secret.

## Vérification

Le job `gitleaks` utilise un clone complet grâce à `GIT_DEPTH: "0"` :

```yaml
gitleaks:
  stage: security
  variables:
    GIT_DEPTH: "0"
  script:
    - gitleaks git --redact --verbose .
```

Une exécution réussie doit terminer sans détection et avec un code de sortie nul. Les éventuelles valeurs détectées sont masquées dans les logs grâce à l'option `--redact`.

## Réaction à une fuite

Si un secret est détecté, il faut immédiatement le révoquer, créer une nouvelle valeur, nettoyer l'historique Git concerné et relancer Gitleaks avant tout nouveau déploiement.
