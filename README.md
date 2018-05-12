#Chaospott Door Status Daemon

Used to add chaospott specific logic to the SpaceApiDaemon (when a door is unlocked, the space is considered to be open)

You can use the container provided under chaospott/doorstatusdaemon

Run it with
```bash
docker run -p 80:8080 -e "API_URL=https://newstatus.chaospott.de" chaospott/doorstatusdaemon
```

Then send your requests
```bash
curl https://newstatus.chaospott.de/LOCATION/VALUE -H  "Authorization: YOURTOKEN"
```