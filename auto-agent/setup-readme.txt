#### The Astro-world auto-agent prototype
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

## Setup GitHub .gitignore to not sync react generated files
*.pid

/auto-agent/static/astro-world/js/*
/auto-agent/static/astro-world/js/.*
/auto-agent/static/astro-world/js/*.*

/exploration/data/*
/exploration/data/.*
/exploration/data/*.*

/exploration/public/js/*
/exploration/public/js/.*
/exploration/public/js/*.*






## Dev setup procedure
# Setup google appengine locally
# 
# Start goappengine locally
cd ~/src/github.com/toisin/astro-world/auto-agent
goapp serve

# Start React Transformation
cd ~/src/github.com/toisin/astro-world/auto-agent/static/astro-world
jsx -w -x jsx jsx js







## Deploy local server to cloud
cd ~/src/github.com/toisin/astro-world
# replace <app-name> with actual app name e.g. "premium-cipher-661"
# e.g goapp deploy -oauth -application premium-cipher-661 auto-agent
goapp deploy -oauth -application <app-name> auto-agent

# app url
# replace <app-name> with actual app name
# e.g. http://premium-cipher-661.appspot.com/
http://<app-name>.appspot.com/

# app dashboard
# replace <app-name> with actual app name
# e.g. https://console.developers.google.com/project/apps~premium-cipher-661/appengine
https://console.developers.google.com/project/apps~<app-name>/appengine
