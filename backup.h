#ifndef BACKUP_H_INCLUDED
#define BACKUP_H_INCLUDED

#define _XOPEN_SOURCE 500
#include<stdio.h>
#include<ftw.h>
#include<stdint.h>

#include"btsoot.h"
#include"sqlite3.h"

int backup(job *job_import);

#endif

