# kubeadm ansible

ansible + terraform to bootstrap a 1 controller + 3 worker kubernetes cluster

You should be able to use this once you're logged on GCP with `gcloud auth login` or directly on Google Cloud Shell.

## Pre-requisites

- You must install `ansible`

```shell
pip install ansible --user
```

- Install terraform 

```shell
curl -LO https://releases.hashicorp.com/terraform/0.11.13/terraform_0.11.13_linux_amd64.zip
unzip terraform_0.11.13_linux_amd64.zip
sudo mv terraform /usr/local/bin/
```

- Install terraform-inventory

```shell
curl -LO https://github.com/adammck/terraform-inventory/releases/download/v0.8/terraform-inventory_v0.8_linux_amd64.zip
unzip terraform-inventory_v0.8_linux_amd64.zip
sudo terraform-inventory /usr/local/bin/
```
