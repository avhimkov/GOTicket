#!/bin/bash

# Golang installation variables
VERSION="1.9"
OS="linux"
ARCH="amd64"

# Home of the vagrant user, not the root which calls this script
HOMEPATH="/home/vagrant"

# Updating and installing stuff
sudo apt-get update
sudo apt-get install -y git curl

if [ ! -e "/vagrant/go.tar.gz" ]; then
	# No given go binary
	# Download golang
	FILE="go$VERSION.$OS-$ARCH.tar.gz"
	URL="https://storage.googleapis.com/golang/$FILE"

	echo "Downloading $FILE ..."
	curl --silent $URL -o "$HOMEPATH/go.tar.gz"
else
	# Go binary given
	echo "Using given binary ..."
	cp "/vagrant/go.tar.gz" "$HOMEPATH/go.tar.gz"
fi;

echo "Extracting ..."
tar -C "$HOMEPATH" -xzf "$HOMEPATH/go.tar.gz"
mv "$HOMEPATH/go" "$HOMEPATH/.go"
rm "$HOMEPATH/go.tar.gz"

# Create go folder structure
GP="/vagrant/gopath"
mkdir -p "$GP/src"
mkdir -p "$GP/pkg"
mkdir -p "$GP/bin"

# Write environment variables, other prompt and automatic cd into /vagrant in the bashrc
echo "Editing .bashrc ..."
touch "$HOMEPATH/.bashrc"
{
	echo '# Prompt'
	echo 'export PROMPT_COMMAND=_prompt'
	echo '_prompt() {'
	echo '    local ec=$?'
	echo '    local code=""'
	echo '    if [ $ec -ne 0 ]; then'
	echo '        code="\[\e[0;31m\][${ec}]\[\e[0m\] "'
	echo '    fi'
	echo '    PS1="${code}\[\e[0;32m\][\u] \W\[\e[0m\] $ "'
	echo '}'

    echo '# Golang environments'
    echo 'export GOROOT=$HOME/.go'
    echo 'export PATH=$PATH:$GOROOT/bin'
    echo 'export GOPATH=/vagrant/gopath'
    echo 'export PATH=$PATH:$GOPATH/bin'

	echo '# Automatically change to the vagrant dir'
	echo 'cd /vagrant'
} >> "$HOMEPATH/.bashrc"

# Add Vendoring in Vagrant
go get -u github.com/golang/dep/cmd/dep
dep init