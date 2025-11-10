# Access the id_github file from Secret Manager, and setup SSH
steps:
  - name: 'gcr.io/cloud-builders/git'
    secretEnv: ['SSH_KEY']
    entrypoint: 'bash'
    args:
      - "-c"
      - |
        echo "$$SSH_KEY" >> /root/.ssh/id_rsa
        chmod 400 /root/.ssh/id_rsa
        ssh-keyscan github.com > /root/.ssh/known_hosts
    volumes:
      - name: "ssh"
        path: /root/.ssh

  - name: golang:1.18
    entrypoint: bash
    args:
      - "-c"
      - |
        git config --global url."support@hubbuy.com:hubbuy/".insteadOf "https://github.com/hubbuy/"
        chmod 0600 /root/.ssh/id_rsa
        go env -w GOPRIVATE=github.com/hubbuy
        go mod vendor
        ls -al ./vendor
    volumes:
      - name: "ssh"
        path: /root/.ssh
  # Build the container image
  - name: 'gcr.io/cloud-builders/docker'
    args: ['build', '-t', 'gcr.io/$PROJECT_ID/$_SERVICE_NAME:$COMMIT_SHA', '.']
  # Push the container image to Container Registry
  - name: 'gcr.io/cloud-builders/docker'
    args: ['push', 'gcr.io/$PROJECT_ID/$_SERVICE_NAME:$COMMIT_SHA']
  # Deploy container image to Cloud Run
  - name: 'gcr.io/google.com/cloudsdktool/cloud-sdk'
    entrypoint: gcloud
    args:
      - 'run'
      - 'deploy'
      - '$_SERVICE_NAME'
      - '--image'
      - 'gcr.io/$PROJECT_ID/$_SERVICE_NAME:$COMMIT_SHA'
      - '--region'
      - '$_DEPLOY_REGION'
images:
  - 'gcr.io/$PROJECT_ID/$_SERVICE_NAME:$COMMIT_SHA'
availableSecrets:
  secretManager:
    - versionName: projects/$PROJECT_ID/secrets/ssh-key-github/versions/latest
      env: 'SSH_KEY'