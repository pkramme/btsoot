#ifndef BTSOOT_H_INCLUDED
#define BTSOOT_H_INCLUDED

#include<stdio.h>
#include<stdlib.h>
#include<string.h>

#include"color.h"

#define MAX_PATH_LEN 4096

#define PIP_PURP_ID_BACKUP 1
#define PIP_PURP_ID_RESTORE 2

#define CONFIG_PATH "btsoot.conf"
#define COPY_CONFIG_PATH "btsoot.conf.temp"

struct job {
	int pip_purp_id;
	char block_name[256];
	char src_path[MAX_PATH_LEN];
	char dest_path[MAX_PATH_LEN];
	char db_path[MAX_PATH_LEN];
};

typedef struct job job;

#include"copy.h"
#include"backup.h"

#endif
