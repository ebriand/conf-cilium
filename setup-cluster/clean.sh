#!/usr/bin/env bash

set -euo pipefail

ansible-playbook clean-up.yml

rm kube_rsa kube_rsa.pub \
  *.retry inventory.cfg \
  infra/terraform.tfstate*
