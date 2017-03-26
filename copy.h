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
#include<fcntl.h>
#include<sys/stat.h>
#include<unistd.h>
#include<sys/sendfile.h>

int copy(char *source, char *destination);

int copy_fallback(char *source, char *destination);

