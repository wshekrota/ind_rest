# ind_rest

ra2.go 

A quick pass at a simple method of extending an API using another port. This is an example and not meant to be secure or full featured.
Set the host of your elasticsearch instance in the bash global RUN_HOST prior to starting this interface.
The endpoint ‘_cat/indices’ for elastic will return a report of existing indices in the database.
My endpoint ‘/distribution’ aggregates that report giving a breakdown of the indices prefixed by the name ‘log-’.
The returned data is of the form json for ease of handling.
Start ra2 in one process and you will get the API version message back. It is waiing on port 10000.

you may then 
   ‘curl -i http://host_of_ra2:10000/distribution’
   
   Additional feature
   To refuse methods passed since we are doing only GET for this.
   
Functional Description:
----------------------
Listen for our single endpoint on port 10000 in preparation to calling aggregate function using HandleFunc.

Use URL.Query to get parms passed as map. Range through the key/value pairs looking for 'X'.
If found write 404 to header and return.

When aggregate response complete convert to serial and return with Fprintf to Http.ResponseWriter.
