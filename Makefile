all:
	cd copy && $(MAKE) all
install:
	cp btsoot.py /usr/local/bin/btsoot
	cp ./copy/copy /etc/btsoot/
uninstall:
	rm /usr/local/bin/btsoot
	rm -rf /etc/btsoot
clean:
	cd copy && $(MAKE) clean
