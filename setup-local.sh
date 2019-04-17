#!/usr/bin/env sh

killall watch

watch "gcloud compute ssh --ssh-flag=\"-L 6443:localhost:6443 -N\" --zone europe-west1-c controller-0" 2>&1 > gcloud.log &
echo "Forwarding distant cluster to :6443"

# Port forward events
watch "kubectl port-forward -n events service/events-frontend" 8080:80 2>&1 > events.log &
echo "Forwarding to http://events.shield.com:8080"

# Port forward internal
watch "kubectl port-forward -n internal service/internal-frontend" 8081:80 2>&1 > internal.log &
echo "Forwarding to http://internal.shield.com:8081"

echo "To stop just launch: killall watch"
