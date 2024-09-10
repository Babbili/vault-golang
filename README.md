# Vault secret management on Kubernetes with Go

This repository showcases installing a Hashicorp Vault server on Kubernetes for centralized secret management and encryption. It demonstrates how to retrieve a secret from Vault *KV v2* secret engine in a Go application that accesses Vault using *approle auth*. This setup provides a secure method for managing, versioning, and encrypting sensitive data across all applications, with fine-grained access control, all centralized in a Vault server & UI.


## Install Vault server on Kubernetes

install Hashicorp Vault helm chart on Kubernetes

```bash
helm -n vault install vault hashicorp/vault \
  --set "server.resources.requests.memory=256Mi" \
  --set "server.resources.requests.cpu=250m" \
  --set "server.resources.limits.memory=256Mi" \
  --set "server.resources.limits.cpu=250m" \
  --set "server.readinessProbe.enabled=true" \
  --set "server.readinessProbe.path='/v1/sys/health?standbyok=true&sealedcode=204&uninitcode=204'" \
  --set "server.livenessProbe.enabled=true" \
  --set "server.livenessProbe.path='/v1/sys/health? standbyok=true'" \
  --set "server.livenessProbe.initialDelaySeconds=60" \
  --set "ui.enabled=true"

# output
NAME: vault
LAST DEPLOYED: Sat Aug 10 13:11:39 2024
NAMESPACE: vault
STATUS: deployed
REVISION: 1
NOTES:
Thank you for installing HashiCorp Vault!

# pull helm chart locally
helm pull hashicorp/vault --untardir="."
```

list deployed Vault resources
```bash
kubectl -n vault get all
NAME                                       READY   STATUS    RESTARTS   AGE
pod/vault-0                                1/1     Running   0          3m50s
pod/vault-agent-injector-ff58f5d77-4n6rh   1/1     Running   0          3m51s

NAME                               TYPE        CLUSTER-IP      EXTERNAL-IP   PORT(S)             AGE
service/vault                      ClusterIP   10.106.251.84   <none>        8200/TCP,8201/TCP   3m51s
service/vault-agent-injector-svc   ClusterIP   10.109.200.74   <none>        443/TCP             3m51s
service/vault-internal             ClusterIP   None            <none>        8200/TCP,8201/TCP   3m51s
service/vault-ui                   ClusterIP   10.97.149.224   <none>        8200/TCP            3m51s

NAME                                   READY   UP-TO-DATE   AVAILABLE   AGE
deployment.apps/vault-agent-injector   1/1     1            1           3m51s

NAME                                             DESIRED   CURRENT   READY   AGE
replicaset.apps/vault-agent-injector-ff58f5d77   1         1         1       3m51s

NAME                     READY   AGE
statefulset.apps/vault   1/1     3m51s
```

create an approle auth and a `go-app` access policy
```bash
kubectl -n vault exec -it vault-0 -- sh

/ $ vault auth enable approle
Success! Enabled approle auth method at: approle/
/ $ vault write auth/approle/role/go-app token_policies="go-app" token_ttl=24h token_max_ttl=48h
Success! Data written to: auth/approle/role/go-app

vault policy write go-app - <<EOF
path "go-app/secret/*" {
  capabilities = ["read", "list"]
}
EOF
Success! Uploaded policy: go-app
```

port-forward vault-ui service, and login to access the UI
```bash
kubectl -n vault port-forward svc/vault-ui 8200

```

We can access various secret engines in Vault. I used the **KV** (Key-Value) engine at the path `go-app/secret`, other secret engines, such as the *Transit engine* provides encryption as a service. This can be useful for tasks like encrypting database entries and decrypting them in applications, or even encrypting Kubernetes secrets

<img src="./vault-ui/Screenshot%20from%202024-08-10%2014-24-30.png" width="100%" height="auto" alt="vault-secret-engines">

<img src="./vault-ui/Screenshot%20from%202024-08-10%2014-26-27.png" width="100%" height="auto" alt="create-kv-secret">

<img src="./vault-ui/Screenshot%20from%202024-08-10%2014-26-48.png" width="100%" height="auto" alt="list-secret">


Auth methods such as JWT, Kubernetes, username & password, ..., we're using **approle** for `go-app` to access Vault

<img src="./vault-ui/Screenshot%20from%202024-08-10%2014-30-52.png" width="100%" height="auto" alt="vault-auth-methods">


## deploy go-app to Kubernetes

apply go-app Kubernetes manifests in *kubernetes-manifests/go-app.yaml*

```bash
kubectl apply -f kubernetes-manifests/go-app.yaml
namespace/apps unchanged
serviceaccount/go-app created
secret/go-app-svca-token created
deployment.apps/go-app created

kubectl -n apps get all
NAME                          READY   STATUS              RESTARTS   AGE
pod/go-app-868c97f47c-cf64v   1/1     Running   0          7m34s
pod/go-app-868c97f47c-jfvs8   1/1     Running   0          7m34s
pod/go-app-868c97f47c-z4xjm   1/1     Running   0          7m34s

NAME                     READY   UP-TO-DATE   AVAILABLE   AGE
deployment.apps/go-app   3/3     3            0           7m34s

NAME                                DESIRED   CURRENT   READY   AGE
replicaset.apps/go-app-868c97f47c   3         3         0       7m34s

```
```bash
kubectl -n apps logs pod/go-app-868c97f47c-cf64v
# username is babbili
```
