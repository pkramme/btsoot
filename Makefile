all:
	cd ./copy/ && $(MAKE) all
install:
	cp btsoot.py /usr/local/bin/btsoot
	mkdir /etc/btsoot
	cp ./copy/btsoot-copy /usr/local/bin/btsoot-copy
uninstall:
	rm /usr/local/bin/btsoot
	rm /usr/local/bin/btsoot-copy
	rm -rf /etc/btsoot
clean:
	cd ./copy/ && $(MAKE) clean
