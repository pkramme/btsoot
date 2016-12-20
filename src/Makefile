clean:
	rm btsoot.conf btsootdb
fullclean:
	rm -rf btsoot.conf btsootdb compare datatransfer-lib
install:
	rm -rf datatransfer-lib
	mkdir datatransfer-lib
	cd datatransfer-lib
	git clone git.paukra.com/paulkramme/datatransfer-lib.git
	cd ..
	rm -rf compare
	mkdir compare
	cd compare
	git clone git.paukra.com/paulkramme/compare.git
	cd ..