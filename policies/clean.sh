#!/usr/bin/env sh

kubectl delete cnp --all -n kafka
kubectl delete cnp --all -n api
kubectl delete netpol --all -n internal
