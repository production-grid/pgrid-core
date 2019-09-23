#!/bin/sh

#lsof -nP -i4TCP:5432 | grep LISTEN

yum -y update

yum -y install wget

yum -y install vim

yum -y update

yum -y install postgresql-server postgresql-contrib

postgresql-setup initdb

echo "listen_addresses = '*'" >> /var/lib/pgsql/data/postgresql.conf

echo "host all all 0.0.0.0/0 md5" >> /var/lib/pgsql/data/pg_hba.conf

systemctl start postgresql

systemctl enable postgresql

sudo -u postgres -i psql -c "ALTER USER postgres WITH PASSWORD '';"

sudo -u postgres -i psql -c "CREATE ROLE pgrid SUPERUSER LOGIN PASSWORD 'pgrid';"

sudo -u postgres -i psql -c "CREATE DATABASE pgrid_core;"

sudo -u postgres -i psql -c "CREATE DATABASE pgrid_schema_test;"

sudo -u postgres -i psql -d template1 -c "GRANT ALL PRIVILEGES ON DATABASE pgrid_core to pgrid;"

sudo -u postgres -i psql -d template1 -c "GRANT ALL PRIVILEGES ON DATABASE pgrid_schema_test to pgrid;"
