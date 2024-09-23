#! /usr/bin/bash

sudo ip addr add dev tun0 10.0.0.2
sudo ip l set dev tun0 up
sudo ip l set dev tun0 mtu 1600 txqueuelen 1000
sudo ip r add 10.0.0.0 dev tun0