#!/bin/bash

# get domain list
wget -O domains.csv http://www.malwaredomainlist.com/mdlcsv.php

# get only the domains
awk -F "\"*,\"*" '{print $2}' domains.csv > domains.raw
rm domains.csv

# strip empty lines
sed '/-/d' ./domains.raw > domains.stripped
rm domains.raw

# strip forward slash
awk -F'/' '{print $1}' domains.stripped > domains.nohyphen
rm domains.stripped

# remove ports
awk -F':' '{print $1}' domains.nohyphen > domains.noport
rm domains.nohyphen

# remove empty lines
sed -i '/^$/d' domains.noport

# rename file
mv domains.noport domains
