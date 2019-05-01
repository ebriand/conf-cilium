# Local development

Local development is done with [Kind](https://kind.sigs.k8s.io/) and [Skaffold](https://github.com/GoogleContainerTools/skaffold).
It won't work with Cilium for the moment, but it's sufficient to test the components.

## Usage

Create a k8s cluster (here with kind, but you might use something else)

```shell
kind create cluster
```

Launch skaffold for local development specifying your repository for pushing images

```shell
skaffold dev --default-repo docker.io/pyaillet
```
