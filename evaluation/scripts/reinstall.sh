export PATH=$PATH:/usr/local/go/bin:$HOME/go/bin
export GOROOT=/usr/local/go
export GOPATH=$HOME/go

cd ~/go/src/rbbc
git pull origin master
make install

cd ~/go/src/ether-rbbc
git pull origin master
make install
