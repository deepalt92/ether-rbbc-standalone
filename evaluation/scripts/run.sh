export PATH=$PATH:/usr/local/go/bin:$HOME/go/bin
export GOROOT=/usr/local/go
export GOPATH=$HOME/go

rm -rf ${HOME}/.lightchain_standalone

# Initialize rbbc
rbbc init --datadir="${HOME}/.lightchain_standalone" --standalone

# Run RbCore
rbcore run -d $HOME/rbcore &

# Run EVM
rm log
rbbc run --datadir="${HOME}/.lightchain_standalone" --rpc --rpcaddr=0.0.0.0 --rpcport=8545 --rpcapi eth,net,web3,personal,admin --dbftDir=$HOME/rbcore --threshold=400 --timeout=1000 >> log &

