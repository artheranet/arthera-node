# InfluxDB
docker run -d --network host --name influxdb -v $PWD/influxdata:/var/lib/influxdb2 influxdb:2.7

## Setup InfluxDB
docker exec influxdb influx setup \
--username arthera \
--password arthera2023 \
--org arthera \
--bucket arthera \
--retention 2d \
--force

## Create an access token
docker exec -it influxdb influx auth create -o arthera -d "Arthera node token" --read-buckets --write-buckets --read-orgs

## Add node monitoring opts
Add the following opts to your validator node:
```
--metrics
--metrics.influxdbv2
```

## Clear data
docker exec -it influxdb influx delete -b validator1 --start '1970-01-01T00:00:00Z' --stop '2070-01-01T00:00:00Z'

# Grafana
docker volume create grafana-storage
docker volume inspect grafana-storage
docker run -d --network host --name=grafana -v grafana-storage:/var/lib/grafana grafana/grafana-enterprise
