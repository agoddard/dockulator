dockulator
==========

dockulator


user hits web server
user enters a calculation
calculation is stored in redis as follows:
GET /calculations
GET /calculations/new
POST /calculations
{
  "calculation": "1 + 4",
  "id": 1,
  "os": null,
  "language": null,
  "answer": null
}



PUT /calculations/:id


docker middleware polls redis for calculations with no answer
docker middleware receives a calculation, spins up random docker OS/language combination
docker receives request for OS/language, runs request, responds via HTTP PUT to ?
 PUT data-->Mongo unique ID, request ID, metadata, calculation response


 

