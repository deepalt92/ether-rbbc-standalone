#!/usr/bin/env bash

mkdir -p $HOME/go/src/

# Install go
wget https://dl.google.com/go/go1.13.1.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.13.1.linux-amd64.tar.gz

# Install make
sudo apt-get -y install build-essential

# configure .profile
echo "export PATH=$PATH:/usr/local/go/bin:$HOME/go/bin
export GOROOT=/usr/local/go
export GOPATH=$HOME/go" >> .profile

# Re-configure shell
. ~/.profile

cd $HOME/go/src/
git clone https://hades.it.usyd.edu.au/yhua/ether-rbbc.git
git clone https://hades.it.usyd.edu.au/yhua/go-rbbc.git rbbc

cd ~/go/src/rbbc
make install

cd ~/go/src/ether-rbbc
make install

# Install Geth
# sudo apt-get install software-properties-common
# sudo add-apt-repository -y ppa:ethereum/ethereum
# sudo apt-get update
# sudo apt-get install ethereum
