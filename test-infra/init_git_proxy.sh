set -x 

# update server and install go 
apt update
apt upgrade -y
snap install go --classic

# clone repo
git clone https://github.com/coillteoir/bramble /root/bramble
cd /root/bramble/git-proxy
git switch git-proxy

#set up go environment

mkdir ~/go
export HOME=/root
export GOPATH=/root/go
export GOCACHE=/root/go/cache

go mod tidy
go mod download
go build -o bin/git-proxy .
./bin/git-proxy -p 80

