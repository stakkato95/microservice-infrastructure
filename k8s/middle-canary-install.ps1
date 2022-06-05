cd ../middle
helm install middle-canary helm --set app.version="canary" --set image.tag=canary
cd ../k8s