#!/bin/bash
set e
echo "[INFO........] k8 cluster setup is already done... during VM provisioning- For more refer ./lima-vm-config/k8infraprovision.sh"
echo "[INFO........] Cloning k8 deployment yaml from github"

mkdir project
cd project
git clone https://github.com/hiabhishek1888/ms-demo.git
cd ms-demo

echo "[INFO........] applying dep yaml and starting app"
kubectl apply -f k8s/


echo "[INFO........] run `kubectl get nodes -o wide` to see nodes"
echo "[INFO........] run `kubectl get pods -n kube-system` to see if all k8 components are up"
echo "[INFO........] run `kubectl get pods -o wide` to see pods"
echo "[INFO........] run `kubectl get svc -o wide` to see services"
echo "[INFO........] run `curl  http://localhost:30080/api/users/1` to check if everything works..."


read -p "Do you want to start PROXY SERVER so that you can get request from your host machine (you laptop)): "
if [[ "$run_script2" =~ ^[Yy](es)?$ ]]; then
    echo "Selected to Run proxy"
    echo "[INFO........] starting golang `user-space` proxy server"
    echo "[INFO........] for more understanding- refer ~/project/ms-demo/golang-reverse-proxy.go"
    echo "[INFO........] installing golang"
    sudo snap install go  --classic
    echo "[INFO........] running golang proxy server in background (not to block current process)"
    sh ./proxy_start.sh & 
else
    echo "Skipping running proxy.sh"
fi

# echo "[INFO........] starting golang `user-space` proxy server if you want to get response over host machine.. (you laptop)"
# echo "[INFO........] for more- refer ~/project/ms-demo/golang-reverse-proxy.go"
# echo "[INFO........] installing golang"
# sudo snap install go  --classic
# go run golang-reverse-proxy.go 