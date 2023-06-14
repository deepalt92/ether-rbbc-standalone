# Smart Red Belly Blockchain - SRBB VM code

SRBB VM is an optimized EVM ported from Geth1.10.18 and is the VM component of the Smart Red Belly Blockchain.
Smart Red Belly Blockchain was accepted at the 37th IEEE International Parallel & Distributed Processing Symposium.
The proceedings will be made available shortly.

Please find the link to our manuscript at: https://gramoli.github.io/pubs/IPDPS23-SmartRedbelly.pdf


# Installation

## Requirements

### Go >= 1.16

#### Install Go

```bash
mkdir -p $HOME/go/src/
wget https://dl.google.com/go/go1.16.1.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.16.1.linux-amd64.tar.gz
```

#### Configure go

```bash
echo "export PATH=$PATH:/usr/local/go/bin:$HOME/go/bin
export GOROOT=/usr/local/go
export GOPATH=$HOME/go" >> .profile
```

### Dep

```bash
curl https://raw.githubusercontent.com/golang/dep/master/install.sh | bash
```

### Make

```bash
sudo apt-get update --fix-missing && sudo apt-get -y install build-essential
```

## Install SRBBC VM

```bash
cd $HOME/go/src/
git clone <this-URL>
cd ~/go/src/ether-rbbc-standalone
make install
```

# Command

## Initialization

This command initializes SRBB.
The target folder can be customized by changing values of the parameter --datadir.

```bash
rbbc init --datadir="${HOME}/.lightchain_standalone" --standalone
```

## Run EVM

This command runs the SRBB VM state machine with the following customizable configurations
* Threshold: 1500
* Timeout: 2000 (in milliseconds)
* The data directory: ${HOME}/.lightchain_standalone
* HTTP RPC port: 8545
* WS port: 8546

```bash
rbbc run --threshold=1500 --timeout=2000 --datadir="${HOME}/.lightchain_standalone" --http --http.addr=0.0.0.0 --http.port=8545 --http.api eth,net,web3,personal,admin --ws --ws.addr=0.0.0.0 --ws.port=8546 --ws.origins="*" --ws.api eth,net,web3,personal,admin
```
*Although you can decide to simply execute (observe.sh added to this repo)

Parameters:
* The Threshold is the maximum number of transactions received before making a proposal.
* The Timeout is the maximum duration between two proposals.
* The RPC and WS port are the ports on which the Ethereum services are hosted -- you can send transactions to these ports.
* Please refer to the Ethereum docs for further information on these
* Recommend reading web3js docs on how to send transactions to an Ethereum node
