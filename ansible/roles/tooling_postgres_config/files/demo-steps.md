# Demo

## Video

[![asciicast](https://asciinema.org/a/273133.svg)](https://asciinema.org/a/273133)

## Steps With Liquibase

1. `titan clone s3://titan-data-demo/hello-world/postgres titan-pg`
2. `titan ls`
3. `source .sourceme` (if you want to persist/inspect the pushes to the ssh server, execute `source .sourceme PERSIST`)
4. `type liquibase`
5. `liquibase generateChangeLog`
6. `cat changelog/changelog.yaml`
7. `liquibase changeLogSync`
8. `psql postgres://postgres:postgres@localhost/postgres -t -c 'SELECT * FROM DATABASECHANGELOG;'`
9. `titan commit -m "baselined with Liquibase" titan-pg`
10. `cat newtable.yaml >> changelog/changelog.yaml`
11. `cat changelog/changelog.yaml`
12. `liquibase validate`
13. `liquibase updateSQL`
14. `liquibase update`
15. `psql postgres://postgres:postgres@localhost/postgres -t -c '\d newtable;'`
16. `titan commit -m "Created table newtable with Liquibase." titan-pg`
17. `psql postgres://postgres:postgres@localhost/postgres -t -c '\dt;'`
18. `psql postgres://postgres:postgres@localhost/postgres -t -c 'DROP SCHEMA public CASCADE;'`
19. `psql postgres://postgres:postgres@localhost/postgres -t -c '\dt;'`
20. `titan log titan-pg`
21. `titan log titan-pg | sed '/commit/h; $!d; x'| awk '{print $2}'`
22. `titan checkout -c $(titan log titan-pg | sed '/commit/h; $!d; x'| awk '{print $2}') titan-pg`
23. `psql postgres://postgres:postgres@localhost/postgres -t -c '\dt;'`
24. `titan remote ls titan-pg`
25. `titan remote log titan-pg`
26. `titan remote add -r ssh ssh://root:root@${SSHSERVER}:22/data titan-pg`
27. `titan remote ls titan-pg`
28. `titan push -r ssh titan-pg`
29. `titan remote log titan-pg`
30. `titan rm -f titan-pg`

## Steps Titan Only

1. `titan run -- --name titan-mongo -p 27017:27017 -d mongo:latest`
2. `mongo --eval 'db.employees.insert({firstName: "Adam", lastName: "Bowen"})'`
3. `titan commit -m "first employee" titan-mongo`
4. `mongo --eval 'db.employees.insert({firstName: "Sanjeev", lastName: "Sharma"})'`
5. `mongo --eval 'db.employees.find()'`
6. `titan checkout -c $(titan log titan-mongo | sed '/commit/h; $!d; x'| awk '{print $2}') titan-mongo`
7. `mongo --eval 'db.employees.find()'`
8. `titan clone s3://titan-data-demo/hello-world/dynamodb titan-dynamo`
9. `aws dynamodb scan --endpoint http://localhost:8000 --table-name messages | jq -r ".Items[0].message.S"`
10. `titan ls`
11. `titan rm -f titan-mongo`
12. `titan rm -f titan-dynamo`
13. `titan ls`
14. `docker rm -f titan-sshd`
