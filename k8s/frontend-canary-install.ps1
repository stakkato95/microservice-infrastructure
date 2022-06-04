cd ../frontend
helm install frontend-canary helm --set app.version="new" --set image.tag=canary
cd ../k8s