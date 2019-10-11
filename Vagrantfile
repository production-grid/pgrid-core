Vagrant.configure(2) do |config|
	config.vm.box = "centos/7"

	config.vm.network "forwarded_port", guest: 5432, host: 5432

	config.vm.synced_folder ".", "/vagrant", disabled: true
	config.vm.provision :shell, :path => "vagrant-init.sh"
	config.vm.provision :shell, :path => "vagrant-always.sh", run: 'always'

	config.vm.provider "virtualbox" do |vb|
		vb.customize ["modifyvm", :id, "--cpuexecutioncap", "50", "--cpus", "2"]
		vb.memory = 3048
	end

end
