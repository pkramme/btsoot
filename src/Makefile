clean:
	rm btsoot.conf
fullclean:
	rm -rf btsoot.conf compare datatransfer-lib
install:
	mkdir datatransfer-lib
	cd datatransfer-lib
	git clone https://git.paukra.com/paulkramme/datatransfer-lib.git
	cd ..
	mkdir compare
	cd compare
	git clone https://git.paukra.com/paulkramme/compare.git
	cd ..
update:
	rm -rf datatransfer-lib
	mkdir datatransfer-lib
	cd datatransfer-lib
	git clone https://git.paukra.com/paulkramme/datatransfer-lib.git
	cd ..
	rm -rf compare
	mkdir compare
	cd compare
	git clone https://git.paukra.com/paulkramme/compare.git
	cd ..