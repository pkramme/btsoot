#include"backup.h"

static sqlite3 *database = NULL;
int8_t buffer[45000];
static int filewalk_info_callback(const char *fpath, const struct stat *sb, int tflag, struct FTW *ftwbuf)
{
	FILE *fp = fopen(fpath, "rb");
	if(fp == NULL)
	{
		return 1;
	}
	XXH64_state_t state64;
	//int8_t buffer[45000];
	uint64_t total_read = 1;
	char type[256];
	int recall;
	uint64_t h64;

	switch(tflag)
	{
		case FTW_D:
		{
			strcpy(type, "directory");
			h64 = 0;
			break;
		}
		case FTW_F:
		{
			XXH64_reset(&state64, 0);
			while(total_read)
			{
				total_read = fread(buffer, 1, 45000, fp);
				XXH64_update(&state64, buffer, 45000);
			}
			h64 = XXH64_digest(&state64);
			strcpy(type, "file");
			break;
		}
		default:
		{
			strcpy(type, "ERROR");
			h64 = 0;
			break;
		}
	}


	char *zsql = sqlite3_mprintf(
	"INSERT INTO files (filename, path, type, size, level, crc) VALUES ('%q', '%q', '%q', '%i', '%i', '%i')"
		, fpath + ftwbuf->base, fpath, type, sb->st_size, ftwbuf->level, h64);

	char *errormessage = 0;
	recall = sqlite3_exec(database, zsql, NULL, NULL, &errormessage);
	if(recall != SQLITE_OK)
	{
		puts("ERROR");
	}
	if(errormessage != NULL)
	{
		printf("%s\n", errormessage);
	}

	sqlite3_free(zsql);
	memset(buffer, 0, 45000);
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
	
	/**
	 * FILEWALKER
	 */
	printf("%s\n", job_import->src_path);

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
