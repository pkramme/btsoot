clean:
	rm btsoot.conf
fullclean:
	rm -rf btsoot.conf compare datatransfer-lib
install:
	mkdir datatransfer-lib
	cd datatransfer-lib
	git clone git@ssh.paukra.com:paulkramme/datatransfer-lib.git
	cd ..
	mkdir compare
	cd compare
	git clone git@ssh.paukra.com:libs/compare.git
	cd ..
update:
	rm -rf datatransfer-lib
	mkdir datatransfer-lib
	cd datatransfer-lib
	git clone git@ssh.paukra.com:paulkramme/datatransfer-lib.git
	cd ..
	rm -rf compare
	mkdir compare
	cd compare
	git clone git@ssh.paukra.com:libs/compare.git
	cd ..
