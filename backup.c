#include"backup.h"

static int filewalk_info_callback(const char *fpath, const struct stat *sb, int tflag, struct FTW *ftwbuf)
{
	FILE *scanfile = fopen("scan.scan", "a");
	fprintf(scanfile, "%-3s %2d %7jd   %-40s %d %s\n",
		(tflag == FTW_D) ?   "d"   : (tflag == FTW_DNR) ? "dnr" :
		(tflag == FTW_DP) ?  "dp"  : (tflag == FTW_F) ?   "f" :
		(tflag == FTW_NS) ?  "ns"  : (tflag == FTW_SL) ?  "sl" :
		(tflag == FTW_SLN) ? "sln" : "???",
		ftwbuf->level, (intmax_t) sb->st_size,
		fpath, ftwbuf->base, fpath + ftwbuf->base);
	fclose(scanfile);
	printf("%d", ftwbuf->base);
	return 0;
}

static int sqlite_callback(void *notused, int argc, char **argv, char **azcolumnname)
{
	int i;
	for(i = 0; i < argc; i++)
	{
		printf("%s = %s\n", azcolumnname[i], argv[i] ? argv[i] : "NULL");
	}
	printf("\n");
	return 0;
}

int backup(job *job_import)
{
	/*DATABASE INIT*/
	sqlite3 *database;
	char *errormessage = 0;
	int recall;

	char db_name[256];
	strcpy(db_name, job_import->block_name); 
	strcat(db_name, ".dat");
	recall = sqlite3_open(db_name, &database);

	/*FILEWALKER*/
	printf("%s\n", job_import->src_path);
	if(nftw(job_import->src_path, filewalk_info_callback, 20, 0) == -1)
	{
		fprintf(stderr, "ERROR NFTW\n");
		exit(EXIT_FAILURE);
	}

	/*CRC CHECK*/

	/*EXECUTOR*/




	/* BACKUP PIPELINE
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

