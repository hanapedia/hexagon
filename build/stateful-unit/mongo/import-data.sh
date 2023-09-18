#!/bin/bash
set -e

# Wait for MongoDB to start
until mongo --eval "print(\"waited for connection\")"; do
  sleep 1
done

# Import the data into the database and collection of your choice
mongoimport --jsonArray --db mongo --collection small --file /data/small.json
mongoimport --jsonArray --db mongo --collection medium --file /data/medium.json
mongoimport --jsonArray --db mongo --collection large --file /data/large.json

# Create the default user
mongo admin --eval 'db.createUser({user: "root", pwd: "password", roles: [{role: "readWrite", db: "mongo"}]})'

echo "Data import completed."
