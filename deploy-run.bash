#!/bin/bash
SERVICE_NAME=test-app
REGION=us-central1

# Deploy service to Cloud Run
gcloud run deploy $SERVICE_NAME --image=docker.io/ashutosh67/cloudboard-app --region=$REGION
# Get the service account associated with the service
SERVICE_ACCOUNT=$(gcloud run services describe $SERVICE_NAME --region=$REGION | grep "Service account" | awk '{print $3}')

# loop through the secrets and give the service accoutn access to them
secrets=(access-token-private-key access-token-public-key refresh-token-private-key refresh-token-public-key redis-url db-dsn)
len=${#secrets[@]}

for ((idx=0; idx<$len; idx++))
do
    # give secretAccessor role to the cloud run instance
    SECRET=${secrets[idx]}
    gcloud secrets add-iam-policy-binding $SECRET --member="serviceAccount:${SERVICE_ACCOUNT}" --role="roles/secretmanager.secretAccessor"
done

# update the service to include secrets 
gcloud run services update $SERVICE_NAME --region=$REGION --update-secrets=REDIS_URL=redis-url:1,DB_DSN=db-dsn:1,ACCESS_TOKEN_PRIVATE_KEY=access-token-private-key:1,ACCESS_TOKEN_PUBLIC_KEY=access-token-public-key:1,REFRESH_TOKEN_PRIVATE_KEY=refresh-token-private-key:1,REFRESH_TOKEN_PUBLIC_KEY=refresh-token-public-key:1