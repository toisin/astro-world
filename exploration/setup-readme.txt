#### The COV dialogue prototype
# Server - Node JS (https://nodejs.org/)
# UI - written in React (https://facebook.github.io/react/index.html)
# Database - MongoDB (https://www.mongodb.org/)
#          - Accessed through http requests handled by Express (http://expressjs.com/4x/api.html)

## Setup GitHub .gitignore to not sync react generated files
/exploration/public/js/*

## Dev setup procedure
# Start DB
~/Web/mongodb/bin/mongod --dbpath ~/GitHub/cov/exploration/data/

# Start webserver
cd ~/GitHub/cov/exploration
node server.js

# Start React Transformation
cd ~/GitHub/cov/exploration
jsx -w -x jsx public/jsx public/js

# Program running on URL http://localhost:3000/