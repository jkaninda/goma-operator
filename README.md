# Goma Gateway Kubernetes operator
A Kubernetes operator for managing Goma Gateway.

[![Tests](https://github.com/jkaninda/goma-gateway/actions/workflows/test.yml/badge.svg)](https://github.com/jkaninda/goma-operator/actions/workflows/test.yml)
[![GitHub Release](https://img.shields.io/github/v/release/jkaninda/goma-operator)](https://github.com/jkaninda/goma-operator/releases)
[![Go Report Card](https://goreportcard.com/badge/github.com/jkaninda/goma-operator)](https://goreportcard.com/report/github.com/jkaninda/goma-operator)
[![Go Reference](https://pkg.go.dev/badge/github.com/jkaninda/goma-operator.svg)](https://pkg.go.dev/github.com/jkaninda/goma-operator)
![Docker Image Size (latest by date)](https://img.shields.io/docker/image-size/jkaninda/goma-operator?style=flat-square)

## Links:

- [Docker Hub](https://hub.docker.com/r/jkaninda/goma-operator)
- [Github](https://github.com/jkaninda/goma-operator)

### Documentation is found at <https://jkaninda.github.io/goma-gateway/operator-manual/installation.html>


### 1. Kubernetes installation

**Install the CRDs and Operator into the cluster:**

```sh
kubectl apply -f https://raw.githubusercontent.com/jkaninda/goma-operator/main/dist/install.yaml
```



## License

Copyright 2024.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

