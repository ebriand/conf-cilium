# kubeadm ansible

ansible + terraform to bootstrap a 1 controller + 3 worker kubernetes cluster

You should be able to use this once you're logged on GCP with `gcloud auth login` or directly on Google Cloud Shell.
You should also define two environment variables:
```shell
export GOOGLE_APPLICATION_CREDENTIALS=<PATH_TO_YOUR_ADC_JSON_FILE>
export GOOGLE_CLOUD_PROJECT=<YOUR_PROJECT>
```

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
sudo mv terraform-inventory /usr/local/bin/
```

## Cluster creation

Just launch `./getup.sh`.
If there is an error during the deployment, just launch again the faulty playbook.

## Cluster interactions

- Grab `/etc/kubernetes/admin.conf` on `controller-0`
- Replace the server section to point to localhost: `server: https://localhost:6443`
- Use gcloud to port forward the controller to your local machine: `gcloud compute ssh --ssh-flag="-L 6443:localhost:6443" controller-0`
