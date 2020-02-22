-- https://stackoverflow.com/questions/5755576/mysql-use-environment-variables-in-script

-- remove anon users, just in case
-- DELETE FROM mysql.user WHERE User='';
-- only allow localhost for root
-- DELETE FROM mysql.user WHERE User='root' AND Host NOT IN ('localhost', '127.0.0.1', '::1');
-- remove test databases
-- DROP DATABASE IF EXISTS test;
-- DELETE FROM mysql.db WHERE Db='test' OR Db='test\\_%';

CREATE DATABASE craftcms;
CREATE USER 'craftcms'@'%' IDENTIFIED BY 'CHANGEME';
GRANT CREATE, ALTER, INDEX, LOCK TABLES, REFERENCES, UPDATE, DELETE, DROP, SELECT, INSERT ON `craftcms`.* TO 'craftcms'@'%';
FLUSH PRIVILEGES;
