# Checklist de soutenance DevOps

## Avant la présentation

- [ ] Démarrer les VM `gitlab` et `runner`.
- [ ] Vérifier `http://192.168.252.2` et `http://192.168.252.3/api/health`.
- [ ] Vérifier que le runner est `Online` dans GitLab.
- [ ] Préparer un commit sans secret pour déclencher le pipeline.
- [ ] Fermer les onglets et terminaux qui affichent des mots de passe ou tokens.

## Déroulé conseillé

1. Présenter les deux VM et expliquer leur séparation.
2. Montrer le dépôt, les Dockerfiles dev/prod et les services Compose.
3. Montrer les variables GitLab masquées et protégées sans révéler leur valeur.
4. Déclencher un pipeline et expliquer Gitleaks, lint, tests, build et deploy.
5. Ouvrir les rapports JUnit et les logs lisibles de chaque job.
6. Montrer le runner connecté dans `Settings > CI/CD > Runners`.
7. Démontrer la bascule blue/green sans interruption.

## Démonstration zero downtime

Avant de déclencher le pipeline :

```bash
while true; do
  curl -s -o /dev/null -w "%{http_code} " http://192.168.252.3/api/health
  sleep 0.2
done
```

Pendant le job `deploy-prod`, montrer les messages suivants :

```text
active=green -> deploying blue
traffic switched to blue
```

Toutes les requêtes doivent rester en HTTP 200. Après la bascule :

```bash
multipass exec runner -- sudo docker ps
multipass exec runner -- sudo docker exec trinity-proxy \
  cat /etc/nginx/active/color.conf
```

## Captures à conserver

- [ ] Pipeline complet vert.
- [ ] Job Gitleaks sans fuite détectée.
- [ ] Rapports de tests JUnit visibles dans GitLab.
- [ ] Runner `Online` et `Idle`.
- [ ] Variables CI/CD avec valeurs masquées.
- [ ] API et Web en état `healthy` après déploiement.
- [ ] Série continue de réponses HTTP 200 pendant la bascule.

## Réponses techniques

**Pourquoi séparer GitLab et le runner ?**

Les builds Docker consomment CPU, mémoire et espace disque. Une VM dédiée évite qu'un build ralentisse ou arrête le serveur GitLab.

**Pourquoi utiliser des variables GitLab ?**

Les secrets ne sont ni inscrits dans les images ni versionnés. Les variables masquées empêchent leur affichage accidentel et les variables protégées limitent leur usage aux branches protégées.

**Pourquoi le socket Docker est-il risqué ?**

Un job ayant accès à `/var/run/docker.sock` contrôle le démon Docker et possède pratiquement les droits `root` sur la VM. Le runner doit donc être réservé aux projets et utilisateurs de confiance.

**Comment la bascule évite-t-elle une panne ?**

La nouvelle couleur démarre en parallèle. Le script attend que les healthchecks API et Web réussissent avant de recharger nginx. L'ancienne couleur reste disponible jusqu'à cette bascule.

**Pourquoi le registry répond-il HTTP 401 pendant le test ?**

Le service est joignable mais exige une authentification. Le pipeline utilise les identifiants temporaires fournis automatiquement par GitLab.

## Fin de présentation

- Résumer les protections : Gitleaks, variables masquées, branche protégée et images minimales.
- Rappeler que le déploiement est entièrement automatique sur `main`.
- Montrer la documentation reproductible du runner et du pipeline.
