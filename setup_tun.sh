#! /usr/bin/bash

sudo ip addr add dev tun0 10.0.0.2
sudo ip l set dev tun0 up
