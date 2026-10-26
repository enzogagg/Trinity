# Installation de GitLab et du runner

## Architecture des VM

| VM | Adresse | Ressources | Rôle |
| --- | --- | --- | --- |
| `gitlab` | `192.168.252.2` | 2 vCPU, 4 Go RAM, 25 Go | GitLab CE et Container Registry |
| `runner` | `192.168.252.3` | 2 vCPU, 3 Go RAM, 20 Go | GitLab Runner, Docker et déploiement |

Les deux VM utilisent Ubuntu Server 24.04 sans interface graphique sur un réseau privé Multipass. Cette séparation évite que la compilation des images Docker consomme les ressources du serveur GitLab.

## Création des VM

```bash
multipass launch 24.04 --name gitlab --cpus 2 --memory 4G --disk 25G
multipass launch 24.04 --name runner --cpus 2 --memory 3G --disk 20G
multipass list
```

## Installation de GitLab CE

Dans la VM `gitlab` :

```bash
sudo apt update
sudo apt install -y curl ca-certificates openssh-server
sudo fallocate -l 2G /swapfile
sudo chmod 600 /swapfile
sudo mkswap /swapfile
sudo swapon /swapfile
echo '/swapfile none swap sw 0 0' | sudo tee -a /etc/fstab

curl -fsSL https://packages.gitlab.com/install/repositories/gitlab/gitlab-ce/script.deb.sh \
  -o /tmp/gitlab-install.sh
sudo bash /tmp/gitlab-install.sh
sudo EXTERNAL_URL="http://192.168.252.2" apt install -y gitlab-ce
```

Ajouter à `/etc/gitlab/gitlab.rb` :

```ruby
puma['worker_processes'] = 0
sidekiq['concurrency'] = 5
prometheus_monitoring['enable'] = false
registry_external_url 'http://192.168.252.2:5050'
```

Appliquer et vérifier la configuration :

```bash
sudo gitlab-ctl reconfigure
sudo gitlab-ctl status
```

Le mot de passe administrateur initial se trouve temporairement dans `/etc/gitlab/initial_root_password`. Il ne doit jamais être ajouté au dépôt ou à un log partagé.

## Installation du runner

Dans la VM `runner` :

```bash
sudo apt update
sudo apt install -y docker.io docker-compose-v2 curl ca-certificates
sudo systemctl enable --now docker

curl -fsSL https://packages.gitlab.com/install/repositories/runner/gitlab-runner/script.deb.sh \
  -o /tmp/runner-install.sh
sudo bash /tmp/runner-install.sh
sudo apt install -y gitlab-runner
```

Le registry utilise HTTP uniquement sur le réseau privé du laboratoire. Docker doit donc l'autoriser dans `/etc/docker/daemon.json` :

```json
{
  "insecure-registries": ["192.168.252.2:5050"]
}
```

```bash
sudo systemctl restart docker
```

## Enregistrement du runner

Créer un runner d'instance dans l'interface GitLab, puis saisir son token sans l'enregistrer dans l'historique du terminal :

```bash
read -rs RUNNER_TOKEN
sudo gitlab-runner register --non-interactive \
  --url http://192.168.252.2 \
  --token "$RUNNER_TOKEN" \
  --executor docker \
  --docker-image alpine:3.20 \
  --docker-volumes /var/run/docker.sock:/var/run/docker.sock \
  --docker-volumes /cache
unset RUNNER_TOKEN
```

Le bloc `[[runners]]` de `/etc/gitlab-runner/config.toml` contient notamment :

```toml
request_concurrency = 2

[runners.docker]
  image = "alpine:3.20"
  volumes = ["/var/run/docker.sock:/var/run/docker.sock", "/cache"]
```

Le montage du socket Docker permet aux jobs de construire et déployer des conteneurs. Il donne aussi un accès équivalent à `root` sur la VM : seuls les projets et utilisateurs de confiance doivent pouvoir utiliser ce runner.

## Vérifications

```bash
sudo gitlab-runner verify
sudo systemctl status gitlab-runner --no-pager
curl -I http://192.168.252.2:5050/v2/
```

Le runner doit être `valid` et `active (running)`. Le registry doit répondre `401 Unauthorized` avec l'en-tête `Docker-Distribution-Api-Version: registry/2.0`, ce qui prouve qu'il est joignable et protégé par authentification.

Dans `Settings > CI/CD > Runners`, activer les runners d'instance pour le projet et autoriser les jobs sans tag, ou ajouter le tag du runner aux jobs.

## Variables CI/CD

Créer les variables suivantes dans `Settings > CI/CD > Variables` :

| Clé | Visibilité | Protection |
| --- | --- | --- |
| `POSTGRES_DB` | Visible | Protégée |
| `POSTGRES_USER` | Visible | Protégée |
| `POSTGRES_PASSWORD` | Masquée et cachée | Protégée |

GitLab fournit automatiquement les variables temporaires `CI_REGISTRY`, `CI_REGISTRY_USER` et `CI_REGISTRY_PASSWORD`. Aucun identifiant du registry n'est stocké dans le dépôt.
