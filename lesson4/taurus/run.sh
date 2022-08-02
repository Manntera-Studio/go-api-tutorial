#!/bin/bash
if [ "$1" = "d" ]; then
    senario=./test-direct.yml
else
    senario=./test.yml
fi

docker run --sysctl net.core.somaxconn=2048 --sysctl net.ipv4.tcp_max_syn_backlog=2048 --sysctl net.ipv4.ip_local_port_range="1024 65535" -it --rm -v `pwd`:/bzt-configs blazemeter/taurus $senario
