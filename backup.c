#include"backup.h"

#define MAX_THREADS 1

static sqlite3 *database = NULL;

static time_t t0;
static time_t tsearched;
static size_t total_size = 0;
static size_t max_thread_size = 0;
static size_t curr_size = 0; //Stores current size for thread filling
static int8_t thread_number = 0;

static int filewalk_info_callback(const char *fpath, const struct stat *sb, int tflag, struct FTW *ftwbuf)
{
	total_size += sb->st_size;

	char zsql[10000];

	sqlite3_snprintf(sizeof(zsql), zsql,
	"INSERT INTO files (filename, path, type, size, level, scantime) VALUES ('%q', '%q', %i, %lli, %i, %i)"
		, fpath + ftwbuf->base, fpath, tflag, sb->st_size, ftwbuf->level, t0);

	char *errormessage = 0;
	sqlite3_exec(database, zsql, NULL, NULL, &errormessage);
	if(errormessage != NULL)
	{
		printf("%s\n", errormessage);
	}

	return 0;
}

static int time_callback(void *notused, int argc, char **argv, char **azcolname)
{
	for(int i = 0; i < argc; i++)
	{
		if(strcmp(azcolname[i], "scantime") == 0 && atoi(argv[i]) < t0)
		{
			tsearched = atoi(argv[i]);
			return 1;
		}
	}
	return 0;
}

static int sql_thread_calc(void *notused, int argc, char **argv, char **azcolname)
{
	char path[4096];
	for(int i = 0; i < argc; i++)
	{	
		//printf("%s = %s\n", azcolname[i], argv[i] ? argv[i] : "NULL");
		if(strcmp(azcolname[i], "size") == 0)
		{
			char *zsql = sqlite3_mprintf(
				"UPDATE files SET thread = %i WHERE path = '%s'", thread_number, path);
			sqlite3_exec(database, zsql, NULL, NULL, NULL);
			if(curr_size <= max_thread_size)
			{
				curr_size += atoi(argv[i]);
			}
			else
			{
				curr_size = 0;
				thread_number += 1;
			}
		}
		else
		{
			strcpy(path, argv[i]);
		}
	}
	return 0;
}

static int hash(void *notused, int argc, char **argv, char **azcolname)
{

	char sql[10000];
	FILE *fp = fopen(argv[0], "rb");
	if(fp == NULL)
	{
		return 0;
	}
	uint64_t h64;

	int8_t buffer[FILEBUFFER];
	XXH64_state_t state64;	
	size_t total_read = 1;
		
	XXH64_reset(&state64, 0);
	while(total_read)
	{
		total_read = fread(buffer, 1, FILEBUFFER, fp);	
		XXH64_update(&state64, buffer, FILEBUFFER);
	}
	h64 = XXH64_digest(&state64);
	sqlite3_snprintf(sizeof(sql), sql, "UPDATE files SET hash = %llu WHERE path = '%s'", h64, argv[0]);
	sqlite3_exec(database, sql, NULL, NULL, NULL);
	
	return 0;
}

static int sql_hash(int threadsql)
{
	char *zsql = sqlite3_mprintf("SELECT path FROM files WHERE thread = %i and scantime = %i", threadsql, t0);
	sqlite3_exec(database, zsql, hash, NULL, NULL);
	return 0;
}

int backup(job_t *job_import)
{
	t0 = time(0);
	char *errormessage = 0;

	db_init(job_import->db_path);	//create and open database
	sqlite3_open(job_import->db_path, &database);

	sqlite3_exec(database, "BEGIN TRANSACTION", NULL, NULL, NULL);
	sqlite3_exec(database, "PRAGMA synchronous = off", NULL, NULL, NULL);
	sqlite3_exec(database, "PRAGMA journal_mode = MEMORY", NULL, NULL, NULL);

	//Execute filewalker
	if(nftw(job_import->src_path, filewalk_info_callback, 20, 0) == -1)
	{
		fprintf(stderr, "ERROR NFTW\n");
		exit(EXIT_FAILURE);
	}

	puts("Searching for previous scan");
	char *zsqlsel = sqlite3_mprintf("SELECT * FROM files WHERE %i > scantime ORDER BY scantime DESC", t0);
	sqlite3_exec(database, zsqlsel, time_callback, NULL, &errormessage);
	if(errormessage != NULL)
	{
		if(strcmp(errormessage, "callback requested query abort") != 0)
		{
			printf("%s\n", errormessage);
		}
	}
	sqlite3_free(zsqlsel);

	max_thread_size = total_size / MAX_THREADS;

	char *zsql_thread_calc = sqlite3_mprintf("SELECT path, size FROM files WHERE scantime = %i AND type = 0", t0);
	sqlite3_exec(database, zsql_thread_calc, sql_thread_calc, NULL, &errormessage);

	sql_hash(0);


	sqlite3_exec(database, "END TRANSACTION", NULL, NULL, NULL);

	/**
	 * BACKUP PIPELINE
	 * 
	 * Needed functions:
	 *  - scan for files and directories, record size of files - done
	 *  - scan for hash values - done
	 *  - diff this scan with the last
	 *  - execute all necessary changes
	 */
	sqlite3_close(database);
	return 0;
}

