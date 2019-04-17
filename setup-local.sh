#!/usr/bin/env sh

watch "gcloud compute ssh --ssh-flag=\"-L 6443:localhost:6443 -N\" --zone europe-west1-c controller-0" 2>&1 > gcloud.log &
echo "Forwarding distant cluster to :6443"

# Port forward events
watch "kubectl port-forward -n events $(kubectl get pods -n events -l team=events -o name | cut -d'/' -f2)" 8080:80 2>&1 > events.log &
echo "Forwarding to http://events.shield.com:8080"

# Port forward internal
watch "kubectl port-forward -n internal $(kubectl get pods -n internal -l team=internal -o name | cut -d'/' -f2)" 8081:80 2>&1 > internal.log &
echo "Forwarding to http://internal.shield.com:8081"

echo "To stop just launch: killall watch"
