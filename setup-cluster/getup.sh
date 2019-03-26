#!/usr/bin/env bash

set -euo pipefail

ansible-playbook create-infra.yml

ansible-playbook -i inventory.cfg setup-nodes.yml

ansible-playbook -i inventory.cfg setup-controller-0.yml

ansible-playbook -i inventory.cfg setup-workers.yml

ansible-playbook -i inventory.cfg setup-cluster.yml
