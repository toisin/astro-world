Get all wines:
curl -i -X GET http://localhost:3000/wines

Get wine with _id value of 5069b47aa892630aae000007 (use a value that exists in your database):
curl -i -X GET http://localhost:3000/wines/5069b47aa892630aae000007

Delete wine with _id value of 5069b47aa892630aae000007:
curl -i -X DELETE http://localhost:3000/wines/5069b47aa892630aae000007

Add a new wine:
curl -i -X POST -H 'Content-Type: application/json' -d '{"name": "New Wine", "year": "2009"}' http://localhost:3000/wines

Modify wine with _id value of 5069b47aa892630aae000007:
curl -i -X PUT -H 'Content-Type: application/json' -d '{"name": "New Wine", "year": "2010"}' http://localhost:3000/wines/5069b47aa892630aae000007