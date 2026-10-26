# Preuve du déploiement sans interruption

## Résultat observé

- Date du test : 21 octobre 2026
- Commit déclencheur : `49f7943` (`ci: verify blue-green deployment`)
- Environnement : GitLab CE et GitLab Runner auto-hébergés sur deux VM Multipass
- Version initiale : `blue`
- Version finale : `green`
- Requêtes contrôlées pendant le pipeline : 1 500
- Réponses différentes de HTTP 200 : 0
- État final de l'API : `{"status":"ok"}`

Le pipeline a démarré les conteneurs `green`, attendu leurs healthchecks, basculé nginx vers cette version, puis supprimé les anciens conteneurs `blue`. Le proxy et PostgreSQL sont restés actifs pendant toute l'opération.

## État final constaté

```text
trinity-green-api-1 Up (healthy)
trinity-green-web-1 Up (healthy)
trinity-proxy Up 0.0.0.0:80->80/tcp
trinity-infra-db-1 Up (healthy)
```

Configuration active du proxy :

```nginx
set $color green;
```

## Commandes de démonstration

Surveiller la disponibilité pendant un déploiement :

```bash
while true; do
  curl -s -o /dev/null -w "%{http_code} " http://192.168.252.3/api/health
  sleep 0.2
done
```

Vérifier l'API après la bascule :

```bash
curl http://192.168.252.3/api/health
```

Résultat attendu :

```json
{"status":"ok"}
```

Vérifier les conteneurs sur la VM runner :

```bash
multipass exec runner -- sudo docker ps
multipass exec runner -- sudo docker exec trinity-proxy \
  cat /etc/nginx/active/color.conf
```

## Preuves à présenter

- Pipeline GitLab vert avec les étapes lint, test, build et deploy.
- Runner `Online` dans les paramètres du projet.
- Log `active=blue -> deploying green`, puis `traffic switched to green`.
- Série continue de statuts HTTP 200 pendant la bascule.
- Conteneurs API et Web `green` en état `healthy`.

## Reprise après redémarrage des VM

Test effectué le 26 octobre 2026 avec `multipass restart gitlab runner`.

- GitLab a retrouvé un état opérationnel et répondu avec une redirection HTTP 302 vers la page de connexion.
- Le runner s'est reconnecté automatiquement et `gitlab-runner verify` a indiqué `is valid`.
- GitLab a redémarré Gitaly, KAS, Workhorse, nginx, PostgreSQL, Puma, Redis, le registry et Sidekiq.
- PostgreSQL, l'API et le Web ont retrouvé l'état `healthy` grâce à `restart: unless-stopped`.
- Le proxy a réexposé automatiquement le port 80.
- L'endpoint `/api/health` a répondu `{"status":"ok"}` sans intervention sur les conteneurs.
