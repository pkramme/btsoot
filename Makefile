all:
	cd ./copy/ && $(MAKE) all
install:
	cp btsoot.py /usr/local/bin/btsoot
	mkdir /etc/btsoot
	cp ./copy/copy /usr/local/bin/copy
uninstall:
	rm /usr/local/bin/btsoot
	rm -rf /etc/btsoot
clean:
	cd ./copy/ && $(MAKE) clean
