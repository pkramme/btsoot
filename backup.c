#include"backup.h"

static sqlite3 *database = NULL;

static int filewalk_info_callback(const char *fpath, const struct stat *sb, int tflag, struct FTW *ftwbuf)
{
	FILE *fp = fopen(fpath, "rb");

	XXH64_state_t state64;
	char buffer[45000];
	uint64_t total_read = 1;

	if(tflag == FTW_F)
	{
		XXH64_reset(&state64, 0);
		while(total_read)
		{
			total_read = fread(buffer, 1, 45000, fp);
			XXH64_update(&state64, buffer, sizeof(buffer));
		}
		uint64_t h64 = XXH64_digest(&state64);

		/*
		fprintf(scanfile, "%-3s %2d %7jd %-40s 0x%llx\n",
			(tflag == FTW_D) ?   "d"   : (tflag == FTW_DNR) ? "dnr" :
			(tflag == FTW_DP) ?  "dp"  : (tflag == FTW_F) ?   "f" :
			(tflag == FTW_NS) ?  "ns"  : (tflag == FTW_SL) ?  "sl" :
			(tflag == FTW_SLN) ? "sln" : "???",
			ftwbuf->level, (intmax_t) sb->st_size,
			fpath, state64
		);
		*/
	}
	else
	{
		//printf("Not a file\n");
	}
	fclose(fp);
	return 0;
}


int backup(job_t *job_import)
{
	/*DATABASE CREATE*/
	int *error;
	error = (int) db_init(job_import->block_name);
	if(error != NULL)
	{
		sqlite3_free(error);
		printf("%s\n", error);
	}

	/*CURRENT DATABASE INIT*/
		/*USE CLEAR FROM CREATE*/
	sqlite3_open(job_import->block_name, &database);

	/*FILEWALKER*/
	printf("%s\n", job_import->src_path);
	if(nftw(job_import->src_path, filewalk_info_callback, 20, 0) == -1)
	{
		fprintf(stderr, "ERROR NFTW\n");
		exit(EXIT_FAILURE);
	}

	/*CRC CHECK*/

	/*EXECUTOR*/




	/**
	 * BACKUP PIPELINE
	 * 
	 * Needed functions:
	 *  - scan for files and directories, record size of files
	 *  - scan for crc values
	 *  - diff this scan with the last
	 *  - execute all necessary changes
	 */
	sqlite3_close(database);
	return 0;
}

