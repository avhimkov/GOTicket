# Vagrant configuration file
# Setting the base box
# Calling the bootstrap file

Vagrant.configure("2") do |config|
 
    config.vm.box = "hashicorp/precise64"

    config.vm.provision :shell, path: "bootstrap.sh"

end