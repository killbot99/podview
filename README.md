# podview

Uses kubebuilder to make a simple CRD and controller to deploy podinfo on-demand.

## To Do

I timeboxed myself to 5 hours. These are a handful of things I would've done next.

- unit tests using envtest
- add redis container to deployment if `redis: true` is passed
- bubble deployment status.readyReplicas up to the CR status
- validation webhook (for example validate the color parameter is `#XXXXXX` format)

## diagram

```
      ┌─────────────────────────────────────┐
      │ podview-system namespace            │
      │                                     │
      │    ┌──────────────────────────┐     │
      │    │ Beholder                 │     │
      │    │ Controller               │     │
      │    │                          │     │
      │    │                          │     │
      │    │                          │     │
      │    │                          │     │
      │    │                          │     │
      │    └──────────────────────────┘     │
      │                                     │
      └─────────────────────────────────────┘


┌────────────────────────────────────────────────────────────────────┐
│ some namespace                                                     │
│                                                                    │
│                                                                    │
│                                                                    │
│                   ┌──────────────────────────────┐                 │
│                   │ Beholder                     │                 │
│                   │ Custom Resource              │                 │
│                   │                              │                 │
│                   │                              │                 │
│                   │                              │                 │
│                   │                              │                 │
│                   │                              │                 │
│                   │                              │                 │
│                   │                              │                 │
│                   │                              │                 │
│                   └────┬───────────────────────┬─┘                 │
│                        │                       │                   │
│                        │                       │                   │
│                        │                       │                   │
│                        │                       │                   │
│                        │                       │                   │
│         ┌──────────────▼──┐                    │                   │
│         │                 │                  ┌─▼──────────────┐    │
│         │podview          │                  │                │    │
│         │deployment       │                  │podview         │    │
│         │                 │                  │service         │    │
│         │                 │                  │                │    │
│         │                 │                  │                │    │
│         └─────────────────┘                  └────────────────┘    │
│                                                                    │
└────────────────────────────────────────────────────────────────────┘
```

## Instructions

```
# install CRD
make install

# deploy beholder controller from dockerhub image
make deploy IMG=davidhoman99/podview:latest

# create a "podview-1" namespace and an example beholder CR in that namespace
kubectl apply -f example/beholder.yaml
```
