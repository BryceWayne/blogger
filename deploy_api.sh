#!/bin/bash
set -e  # Exit on any error

echo "Starting deployment script."

# Take project_id as a command-line argument
project_id=$1

if [ -z "$project_id" ]; then
  echo "Error: No project_id provided."
  exit 1
fi

echo "Project ID: $project_id"

# Navigate to 'api' directory
echo "Changing to 'api' directory."
cd api

# Build the Docker image
echo "Building Docker image."
docker build -t blogger-gpt-api .

# Tag the image for Google Container Registry
echo "Tagging Docker image."
docker tag blogger-gpt-api:latest gcr.io/$project_id/blogger-gpt-api:latest

# Push the image to Google Container Registry
echo "Pushing Docker image to GCR."
docker push gcr.io/$project_id/blogger-gpt-api:latest

# Deploy the image to Google Cloud Run
echo "Deploying to Google Cloud Run."
gcloud run deploy --image gcr.io/$project_id/blogger-gpt-api:latest --platform managed

# Navigate back to project root
echo "Navigating back to project root."
cd ..

echo "Deployment completed."
