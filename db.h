/*
 * Copyright (C) Paul Kramme 2017
 * 
 * Part of BTSOOT
 * Single folder redundancy offsite-backup utility
 * 
 * Licensed under MIT License
 */

#pragma once

#include<stdio.h>
#include<string.h>

#include"config.h"
#include"sqlite3.h"

int db_init(char blockname[256]);


