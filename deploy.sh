#!/bin/bash

SERVICE_NAME=manydocs

PROJECT_ID="matt-welke-manydocs"
VERSION=$(git rev-parse HEAD)
IMAGE_TAG="gcr.io/$PROJECT_ID/manydocs"
REGION="us-central1"

gcloud builds submit --tag "$IMAGE_TAG"

gcloud run deploy --image "$IMAGE_TAG" --platform managed $SERVICE_NAME --region $REGION --allow-unauthenticated
