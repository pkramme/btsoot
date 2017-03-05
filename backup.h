#ifndef BACKUP_H_INCLUDED
#define BACKUP_H_INCLUDED

#define _XOPEN_SOURCE 500
#ifndef XXHASH_C_2097394837
#define XXHASH_C_2097394837
#include<stdio.h>
#include<ftw.h>
#include<stdint.h>
#include<inttypes.h>

#include<unistd.h>

#include"sqlite3.h"
#include"btsoot.h"
#include"crc.h"

#define XXH_STATIC_LINKING_ONLY
#include"xxhash.h"
#include"db.h"

int backup(job_t *job_import);

#endif
#endif

