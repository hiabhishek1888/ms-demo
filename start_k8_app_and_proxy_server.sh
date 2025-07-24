#!/bin/bash
set e
echo "[INFO........] k8 cluster setup is already done... during VM provisioning- For more refer ./lima-vm-config/k8infraprovision.sh"
echo "[INFO........] Cloning k8 deployment yaml from github"

mkdir project
cd project
git clone https://github.com/hiabhishek1888/ms-demo.git
cd ms-demo

echo "[INFO] applying dep yaml and starting app"
kubectl apply -f k8s/

echo "-"
echo "[INFO] k8 app started"
echo " "
echo "[code........] run    kubectl get nodes -o wide                   to see nodes"
echo "[code........] run    kubectl get pods -n kube-system             to see if all k8 components are up"
echo "[code........] run    kubectl get pods -o wide                    to see pods"
echo "[code........] run    kubectl get svc -o wide                     to see services"
echo "[code........] run    curl http://localhost:30080/api/users/1     to check if everything works..."

echo "-"
echo "-"
echo "-"
echo "-"
echo "[INFO] GOLANG USER SPACE PROXY SERVER "
echo "[INFO] FOR MORE UNDERSTANDING - refer ~/project/ms-demo/golang-reverse-proxy.go"
read -p "Do you want to start PROXY SERVER so that you can get request from your host machine (you laptop) [y/n] ): " choice
if [[ "$choice" =~ ^[Yy](es)?$ ]]; then
    echo "Selected to Run proxy"
    echo "[INFO] starting golang user-space proxy server SETUP"
    echo "[INFO] installing golang"
    sudo snap install go  --classic
    echo "[INFO] running golang proxy server in background (not to block current process)"
    sh ./proxy_start.sh & 
    echo "-"
    echo "-"
    echo "[INFO] Proxy listening on :29999 → forwarding to :30080"
    echo "[code........] run    curl http://localhost:29999/api/users/1     to check if you can recieve response from k8 app on your device "
else
    echo "Skipping running proxy.sh"
fi

echo "end of script"
echo "-"
echo "-"
# echo "[INFO........] starting golang `user-space` proxy server if you want to get response over host machine.. (you laptop)"
# echo "[INFO........] for more- refer ~/project/ms-demo/golang-reverse-proxy.go"
# echo "[INFO........] installing golang"
# sudo snap install go  --classic
# go run golang-reverse-proxy.go 