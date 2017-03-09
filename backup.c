#include"backup.h"

static sqlite3 *database = NULL;

static int filewalk_info_callback(const char *fpath, const struct stat *sb, int tflag, struct FTW *ftwbuf)
{
	FILE *fp = fopen(fpath, "rb");
	if(fp == NULL)
	{
		return 1;
	}
	XXH64_state_t state64;
	size_t total_read = 1;
	uint64_t h64;
	char zsql[10000];
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

	switch(tflag)
	{
		case FTW_F:
		{
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
	"INSERT INTO files (filename, path, type, size, level, hash) VALUES ('%q', '%q', %i, %i, %i, %i)"
		, fpath + ftwbuf->base, fpath, tflag, sb->st_size, ftwbuf->level, h64);

	char *errormessage = 0;
	sqlite3_exec(database, zsql, NULL, NULL, &errormessage);
	if(errormessage != NULL)
	{
		printf("%s\n", errormessage);
	}

	fclose(fp);
	return 0;
}

int backup(job_t *job_import)
{
	/*DATABASE CREATE*/
	db_init(job_import->block_name);
	/*CURRENT DATABASE INIT*/
		/*USE CLEAR FROM CREATE*/
	sqlite3_open(job_import->block_name, &database);

	//BEGIN SQLITE TRANSACTION AND SPEED HACKS
	sqlite3_exec(database, "BEGIN TRANSACTION", NULL, NULL, NULL);
	sqlite3_exec(database, "PRAGMA synchronous = off", NULL, NULL, NULL);
	sqlite3_exec(database, "PRAGMA journal_mode = MEMORY", NULL, NULL, NULL);

	//Execute filewalk
	if(nftw(job_import->src_path, filewalk_info_callback, 20, 0) == -1)
	{
		fprintf(stderr, "ERROR NFTW\n");
		exit(EXIT_FAILURE);
	}

	//CLOSE SQLITE TRANSACTION
	sqlite3_exec(database, "END TRANSACTION", NULL, NULL, NULL);

	/*CRC CHECK*/

	/*EXECUTOR*/

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
