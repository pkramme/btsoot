#include"db.h"

int db_init(char path[4096])
{
	/*DATABASE INIT*/
	sqlite3 *database;
	char *errormessage = 0;

	sqlite3_open(path, &database);
	sqlite3_exec(database, 
		"CREATE TABLE IF NOT EXISTS files(filename TEXT, path TEXT, type INT, hash INT, size INT, level INT, scantime NUMERIC, thread INT)", 
		0, 
		0, 
		&errormessage
	);
	if(errormessage != NULL)
	{
		fprintf(stderr, "%s\n", errormessage);
	}

	return 0;
}
