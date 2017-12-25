# Vagrant Golang

Vagrant configuration for a simple golang setup.

## Requirements
You have to install [Vagrant](https://www.vagrantup.com) and [git](https://git-scm.com/) first.  
If you're not familiar with Vagrant you can read it's great [docs](https://www.vagrantup.com/docs/).

## Usage
Run following commands in your preferred directory.
```bash
git clone https://github.com/mbndr/vagrant-golang
cd vagrant-golang/
vagrant up
```
Now the virtual machine is running and accessible through `vagrant ssh`.

## Go version
To edit the version to install of Go you can edit the `bootstrap.sh` (first three lines).

Another method is to download the preferred Go binary release from [this](https://golang.org/dl/) page and locate it in the clone repository with the name `go.tar.gz`.  
If this file exists, nothing will be downloaded instead this file is used.  
After bootstrapping the manually downloaded `go.tar.gz` can be removed.

You can call the bootstrap script again with `vagrant provision`

## Go paths
The script is setting up Go's paths as follows:
```bash
GOROOT=$HOME/.go         # /home/vagrant/.go
GOPATH=/vagrant/gopath   # The folder shared with the host system
```
The folders `src/`, `bin/` and `pkg/` are also created in the script.  
The `GOROOT` and the `bin/` directory in the `GOPATH` are added to the `PATH` automatically.
