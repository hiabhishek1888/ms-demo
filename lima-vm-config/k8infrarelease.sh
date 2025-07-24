#!/bin/bash
set -euo pipefail


limactl stop master
limactl stop worker1
limactl stop worker2

limactl delete master worker1 worker2
