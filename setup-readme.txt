#### The COV dialogue prototype
# Engine - runs on a google appengine (https://cloud.google.com/appengine/)
# Server side logic - in Go (http://golang.org/)
# UI - in React (https://facebook.github.io/react/index.html)
# Database - App Engine Datastore (https://cloud.google.com/appengine/features/#storage)
#          - accessed through Go API (https://cloud.google.com/appengine/docs/go/datastore)

#### Helpful documentations and tutorial
# doc for appengine for Go
# https://developers.google.com/appengine/docs/go/
#
# doc for db access
# https://developers.google.com/appengine/docs/go/gettingstarted/usingdatastore
#
# Go Tutorial
# http://tour.golang.org/

## Dev setup procedure
# Setup google appengine locally
# 
# Start goappengine locally
cd ~/GitHub/cov/dialogue
goapp serve
#
# Start React Transformation
cd ~/GitHub/cov/dialogue/static/cov
jsx -w -x jsx js js

## Deploy local server to cloud
cd ~/GitHub/
# replace <app-name> with actual app name
goapp deploy -oauth -application <app-name> cov

# app url
# replace <app-name> with actual app name
http://<app-name>.appspot.com/

# app dashboard
# replace <app-name> with actual app name
https://console.developers.google.com/project/apps~<app-name./appengine
