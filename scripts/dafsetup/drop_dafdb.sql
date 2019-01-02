UPDATE pg_database SET datallowconn = 'false' WHERE datname = 'dafdb';
ALTER DATABASE dafdb CONNECTION LIMIT 0;
SELECT pg_terminate_backend(pid)
FROM pg_stat_activity
WHERE datname = 'dafdb';
drop database dafdb;