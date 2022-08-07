#!/bin/sh

mysql -u${MYSQL_USER} -p${MYSQL_PASSWORD} --local-infile=1 ${MYSQL_DATABASE} -e "LOAD DATA LOCAL INFILE '/docker-entrypoint-initdb.d/race.csv' INTO TABLE race FIELDS TERMINATED BY ',' ENCLOSED BY '\"' LINES TERMINATED BY '\n'"
