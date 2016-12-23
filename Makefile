clean:
	rm btsoot.conf
fullclean:
	rm -rf btsoot.conf compare datatransfer
install:
	git clone https://git.paukra.com/libs/compare.git
	git clone https://git.paukra.com/libs/datatransfer.git
update:
	rm -rf compare
	git clone https://git.paukra.com/libs/compare.git
	rm -rf datatransfer-lib
	git clone https://git.paukra.com/libs/datatransfer.git
