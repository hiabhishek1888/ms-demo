#!/bin/bash
## if we use $HOME with sudo -- it will point to the /root folder instead of the /home/users folder
set -e


# 1. Create and start master node vm via lima
limactl start --name=master --yes newmaster.yaml

# 2. Create and start worker1 node vm via lima
limactl start --name=worker1 --yes newworker.yaml

# 3. Create and start worker2 node vm via lima
limactl start --name=worker2 --yes newworker.yaml  


#4. Join worker nodes to master node
JOIN_CMD=$(limactl shell master kubeadm token create --print-join-command)
#echo "{\"stdout\": \"${JOIN_CMD}\"}"

limactl shell worker1 sudo $JOIN_CMD
limactl shell worker2 sudo $JOIN_CMD


# Copying the start script master node in start directory

limactl shell master "mkdir -p \$HOME/start"
limactl copy ../start_k8_app_and_proxy_server.sh master:$(limactl shell master 'echo $HOME')/start/start_k8_app_and_proxy_server.sh



# TODO: create the func or have an exclusive check if vm are created already then just start it
#       or if it is created and stopped, then start it.. 
