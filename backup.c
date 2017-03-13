#include"backup.h"

static sqlite3 *database = NULL;

time_t t0;

static int filewalk_info_callback(const char *fpath, const struct stat *sb, int tflag, struct FTW *ftwbuf)
{
	FILE *fp = fopen(fpath, "rb");
	if(fp == NULL)
	{
		return 0;
	}

	uint64_t h64;
	char zsql[10000];

	switch(tflag)
	{
		case FTW_F:
		{
			size_t initsize;
			if(sb->st_size < FILEBUFFER) 
			{
				initsize = sb->st_size;
			}
			else
			{
				initsize = FILEBUFFER;
			}

			int8_t buffer[initsize];

			XXH64_state_t state64;
			size_t total_read = 1;

			XXH64_reset(&state64, 0);
			while(total_read)
			{
				total_read = fread(buffer, 1, initsize, fp);
				XXH64_update(&state64, buffer, initsize);
			}
			h64 = XXH64_digest(&state64);
			break;
		}
		default:
		{
			h64 = 0;
			break;
		}
	}


	sqlite3_snprintf(sizeof(zsql), zsql,
	"INSERT INTO files (filename, path, type, size, level, hash, scantime) VALUES ('%q', '%q', %i, %i, %i, %i, %i)"
		, fpath + ftwbuf->base, fpath, tflag, sb->st_size, ftwbuf->level, h64, t0);

	char *errormessage = 0;
	sqlite3_exec(database, zsql, NULL, NULL, &errormessage);
	if(errormessage != NULL)
	{
		printf("%s\n", errormessage);
	}

	fclose(fp);
	return 0;
}

static int diff_callback(void *notused, int argc, char **argv, char **azcolname)
{
	for(int i = 0; i < argc; i++)
	{
		if(strcmp(azcolname[i], "scantime") == 0 && atoi(argv[i]) < t0)
		{
			printf("Comparing %li with %s\n", t0, argv[i]);
			return 1;
		}
	}
	return 0;
}

int backup(job_t *job_import)
{
	t0 = time(0);

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
	char *errormesgselect = 0;
	char *zsqlsel = sqlite3_mprintf("SELECT * FROM files \
					WHERE %i > scantime \
					ORDER BY scantime DESC", t0);
	sqlite3_exec(database, zsqlsel, diff_callback, NULL, &errormesgselect);
	if(errormesgselect != NULL)
	{
		printf("%s\n", errormesgselect);
	}


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

