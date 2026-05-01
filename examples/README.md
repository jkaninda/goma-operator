# goma-operator examples

Runnable manifests demonstrating common deployment patterns for the operator's
`Gateway`, `Route`, and `Middleware` CRDs. Apply any example with:

```sh
kubectl apply -f examples/<file>.yaml
```

| File | What it shows |
| ---- | ------------- |
| [01-gateway-basic.yaml](01-gateway-basic.yaml) | Minimal `ClusterIP` gateway with the K8s provider sidecar enabled (default). |
| [02-gateway-ingress-lb.yaml](02-gateway-ingress-lb.yaml) | Ingress-style `LoadBalancer` gateway on ports 80/443, client-IP preservation. |
| [03-gateway-nodeport.yaml](03-gateway-nodeport.yaml) | `NodePort` exposure for bare-metal / dev clusters. |
| [04-gateway-autoscaling.yaml](04-gateway-autoscaling.yaml) | HPA driven by CPU + memory utilization. |
| [05-gateway-acme-http01.yaml](05-gateway-acme-http01.yaml) | Let's Encrypt HTTP-01 with ACME store persisted to a Secret. |
| [06-gateway-acme-dns01-cloudflare.yaml](06-gateway-acme-dns01-cloudflare.yaml) | Wildcard certs via DNS-01 (Cloudflare). |
| [07-gateway-redis-metrics.yaml](07-gateway-redis-metrics.yaml) | Redis-backed rate limiting + Prometheus metrics guarded by basic-auth. |
| [08-route-basic.yaml](08-route-basic.yaml) | Simple host-based route to a backend service. |
| [09-route-load-balanced.yaml](09-route-load-balanced.yaml) | Weighted load-balancing across multiple backends with health checks. |
| [10-route-maintenance.yaml](10-route-maintenance.yaml) | Route in maintenance mode returning a custom 503 page. |
| [11-middleware-basic-auth.yaml](11-middleware-basic-auth.yaml) | `basic` auth middleware protecting `/admin`. |
| [12-middleware-rate-limit.yaml](12-middleware-rate-limit.yaml) | Per-IP rate limiter. |
| [13-middleware-jwt.yaml](13-middleware-jwt.yaml) | JWT validation against a remote JWKS. |
| [14-middleware-forward-auth.yaml](14-middleware-forward-auth.yaml) | Delegated auth via an external `forwardAuth` endpoint. |
| [15-full-stack.yaml](15-full-stack.yaml) | End-to-end: Gateway + two Routes + JWT + rate-limit middleware. |

All examples assume the namespace `default`. Adjust `metadata.namespace`,
hostnames, and backend `target`s for your environment.
