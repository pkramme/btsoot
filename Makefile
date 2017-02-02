all:
	make ./copy/
install:
	cp btsoot.py /usr/local/bin/btsoot
	cp ./copy/copy /etc/btsoot/
uninstall:
	rm /usr/local/bin/btsoot
	rm -rf /etc/btsoot
