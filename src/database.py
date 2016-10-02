database = None
def create(path = "btsootdb")
	database = open(path, "w")
	database.write("Database")

def write_database(stufftowrite):
	database.write(stufftowrite)
