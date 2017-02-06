all:
	$(MAKE) -C copy/make all
install:
	cp btsoot.py /usr/local/bin/btsoot
	mkdir /etc/btsoot
	cp ./copy/copy /etc/btsoot/
uninstall:
	rm /usr/local/bin/btsoot
	rm -rf /etc/btsoot
clean:
	$(MAKE) -C copy/make clean
