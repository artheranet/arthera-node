## Run InfluxDB
docker run -d -p 8086:8086 --name influx2 -v $PWD/influxdata:/var/lib/influxdb2 influxdb:2.7

## Setup InfluxDB
docker exec influx2 influx setup \
--username arthera \
--password arthera2023 \
--org arthera \
--bucket arthera \
--retention 7d \
--force

## Create an access token
docker exec -it influx2 influx auth create -o arthera -d "Arthera cli token" --all-access
M1AEMLc3Y7SvsWmyh7aQatk975C2Dd1LZyBO_DjvBGstKICyz-rHhtasj9ONGpUvOPZyKDBAlHsGYv6i0Mwwug==

docker exec -it influx2 influx auth create -o arthera -d "Arthera node token" --write-bucket d575e358bc41d2f0
sCQnhSRn828eebudgE8_wgnU6glEjGIHlQbFnVO3gFeQWgr3Hj2pSP84zhHctq6028tWYpAAcDagoXb_9WIN1A==


## Clear data
docker exec -it influx2 influx delete -b arthera --start '1970-01-01T00:00:00Z' --stop '2070-01-01T00:00:00Z'
