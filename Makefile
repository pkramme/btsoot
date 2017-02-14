all:
	cd ./src/ && $(MAKE) all && cp btsoot ../btsoot
install:
	cp btsoot /usr/local/bin/btsoot
	mkdir -p /etc/btsoot/scans
update:
	cp btsoot /usr/local/bin/btsoot
uninstall:
	rm /usr/local/bin/btsoot
	rm -rf /etc/btsoot
clean:
	cd ./src/ && $(MAKE) clean
	rm btsoot
