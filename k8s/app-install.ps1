cd ../frontend
helm install frontend helm
echo ""

cd ../middle
helm install middle helm
echo ""

cd ../backend-1
helm install backend-1 helm

cd ../k8s