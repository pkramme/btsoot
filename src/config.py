configfile = None
def create(path = "btsoot.conf"):
	configfile = open(path, "w")
	configfile.write("something\n")
	
