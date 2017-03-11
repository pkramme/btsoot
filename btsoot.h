#pragma once

#include<stdio.h>
#include<stdlib.h>
#include<string.h>

#include"color.h"
#include"config.h"

#define MAX_PATH_LEN 4096

#define PIP_PURP_ID_BACKUP 1
#define PIP_PURP_ID_RESTORE 2

struct job {
	int pip_purp_id;
	char block_name[256];
	char src_path[MAX_PATH_LEN];
	char dest_path[MAX_PATH_LEN];
	char db_path[MAX_PATH_LEN];
	char scantime[256];
};

typedef struct job job_t;

#include"copy.h"
#include"backup.h"

