cd ../frontend
helm install frontend-canary helm --set app.version="0.1.1" --set image.tag=canary
cd ../k8s