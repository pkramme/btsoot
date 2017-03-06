#ifndef DB_H_INCLUDED
#define DB_H_INCLUDED

#include<stdio.h>
#include<string.h>

#include"sqlite3.h"

int db_init(char blockname[256]);

#endif

