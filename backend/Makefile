docker run --name postgres_kurs -p 5432:5432 --rm -d -e POSTGRES_USER=mkyong -ePOSTGRES_PASSWORD=password  -ePOSTGRES_DB=kurs -v ./db/init.sql:/docker-entrypoint-initdb.d/init.sql postgres
docker exec -it postgres_kurs psql -U mkyong -d kurs
curl -D - -X POST localhost:8081/api/user/register -H "Content-Type: application/x-www-form-urlencoded"  -d "username=qerqqqe&password=asdasdf&phone=80808080808&email=fsd@mil"
update people set role='admin' where id=4;