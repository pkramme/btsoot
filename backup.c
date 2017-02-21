#include"backup.h"

int backup(job *job_import)
{
	puts(job_import->block_name);
	/* BACKUP PIPELINE
	 * 
	 * Needed functions:
	 *  - scan for files and directories, record size of files
	 *  - scan for crc values
	 *  - diff this scan with the last
	 *  - execute all necessary changes
	 */
	return 0;
}
