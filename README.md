# Mnemosyne [![Build Status](https://travis-ci.com/arekn/mnemosyne.svg?branch=master)](https://travis-ci.com/arekn/mnemosyne)

Mnemosyne is an application parsing `proc/meminfo` and calculating amount of free memory on linux device. 

Measurements are saved as .csv every minute. Each day is saved into separate file.

Why? I want to see if my 16Gb of RAM is enough.

### To do next

1. Enable to run this as systemd service
1. Measure amount of memory used by various processes to see what is consuming most of memory (I got my eyes on you Chrome)
1. Generate weekly reports out of gathered data
1. Generate plots
