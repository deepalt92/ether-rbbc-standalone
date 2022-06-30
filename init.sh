#!/usr/bin/env bash
rm -rf $HOME/.lightchain_standalone[0-9]*

rbbc init --datadir="${HOME}/.lightchain_standalone0" --standalone
rbbc init --datadir="${HOME}/.lightchain_standalone1" --standalone
rbbc init --datadir="${HOME}/.lightchain_standalone2" --standalone
rbbc init --datadir="${HOME}/.lightchain_standalone3" --standalone

# rbbc init --datadir="${HOME}/.lightchain_standalone" --standalone

# Install Geth
# sudo apt-get install software-properties-common
# sudo add-apt-repository -y ppa:ethereum/ethereum
# sudo apt-get update
# sudo apt-get install ethereum