# Secrets

Generate a secret manifest
```
kubectl -n default create secret generic my-secret-name \
    --from-literal=key=value \
    --dry-run=client \
    -o yaml > secret.yaml
```

Encrypt secret manifest using kubeseal
```
kubeseal --format=yaml --cert=pub-sealed-secrets.pem \
< secret.yaml > sealed-secret.yaml
```


Clean up old unencrypted secret
```
rm secret.yaml
```