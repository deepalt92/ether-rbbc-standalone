export PATH=$PATH:/usr/local/go/bin:$HOME/go/bin
export GOROOT=/usr/local/go
export GOPATH=$HOME/go

rm -rf ${HOME}/.lightchain_standalone

# Initialize rbbc
#rm -rf /home/ubuntu/.lightchain_standalone/database/chaindata
rbbc init --datadir="${HOME}/.lightchain_standalone" --standalone
#############################
rm -rf /home/deepal/.lightchain_standalone/database/chaindata
rm -rf /home/deepal/.lightchain_standalone/database/nodes
cp genesis.json /home/deepal/.lightchain_standalone/database/
#mv /home/deepal/.lightchain_standalone/database/genesis-new.json /home/deepal/.lightchain_standalone/database/genesis.json
###########################
# Run EVM
rm -f log
#rbbc run --threshold=1500 --timeout=2000 --datadir="${HOME}/.lightchain_standalone" --rpc --rpcaddr=0.0.0.0 --rpcport=8545 --rpcapi eth,net,web3,personal,admin --dbftDir=$HOME/rbcore >> log


rbbc run --threshold=1500 --timeout=2000 --datadir="${HOME}/.lightchain_standalone" --http --http.addr=0.0.0.0 --http.port=8545 --http.api eth,net,web3,personal,admin --ws --ws.addr=0.0.0.0 --ws.port=8546 --ws.origins="*" --ws.api eth,net,web3,personal,admin 

#--dbftDir=$HOME/rbcore >> log

#--v5disc --bootnodes "enode://7a12b4525656301a12e307b752319b94f95a8c8621adeacc20ded1237d60381e934a33101e36e24c5948a8356d2013dbb652c9bd3e11d96e26e6cec6a4d315b8@127.0.0.1:0?discport=30305" --networkid 161 >> log

#rm -rf /home/deepal/.lightchain_standalone/database/LOCK
