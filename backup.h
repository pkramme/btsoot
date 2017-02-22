#ifndef BACKUP_H_INCLUDED
#define BACKUP_H_INCLUDED

#include<stdio.h>
#define _XOPEN_SOURCE 500
#include<ftw.h>

#include"btsoot.h"
#include"sqlite3.h"

int backup(job *job_import);

#endif

