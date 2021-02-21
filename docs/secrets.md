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

## Issues
kubseal sometimes fails to get the public cert, you can fetch the cert like this instead
```
kubectl port-forward service/sealed-secrets -n flux-system 8081:8080
curl localhost:8081/v1/cert.pem > pub-sealed-secrets.pem
```