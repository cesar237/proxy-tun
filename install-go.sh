#! /usr/bin/bash

VERSION=1.23.1

wget https://go.dev/dl/go${VERSION}.linux-amd64.tar.gz
sudo tar -xvf go${VERSION}.linux-amd64.tar.gz
sudo mv go /usr/local
sudo ln -s /usr/local/go/bin/* /usr/bin
rm go${VERSION}.linux-amd64.tar.gz

echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
 