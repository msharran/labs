#!/bin/bash


echo "Starting build"
SECONDS=0
sleep 2
duration=$SECONDS
echo "$(($duration / 60)) minutes and $(($duration % 60)) seconds elapsed."

echo "Sending metrics to server"
curl -i -X POST localhost:9900/metrics/pipeline/stage -d '{
  "duration_seconds": '$duration',
  "status": "SUCCESS",
  "stage_name": "build",
  "pipeline_name": "foo-prod-deployment"
}'
