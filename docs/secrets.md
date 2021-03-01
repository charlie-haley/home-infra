# Secrets

## Task

```
task generate-secret
```
Edit the `secret.yaml` file that's generated with the values you want to store in the secret.

```
task seal-secret
```

### Manually

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

### Issues
kubseal sometimes fails to get the public cert, you can fetch the cert like this instead
```
task fetch-cert
```