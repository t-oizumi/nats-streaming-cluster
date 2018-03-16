#!/bin/bash
sleep $1
rm -rf /go/src/store
cp /go/src/web/server.go /go/src/github.com/nats-io/nats-streaming-server/server/server.go
cd /go/src/github.com/nats-io/nats-streaming-server/
go build
### 5台構成
./nats-streaming-server -p 4222 -cluster_sync=false -file_sync=false -store file -dir /go/src/store/nats -cluster_log_path /go/src/store/nats_cluster -cluster_node_id d -cluster_peers a,c,b,e -clustered -cluster nats://nats_d:6222 -routes nats://nats_c:6222,nats://nats_a:6222,nats://nats_b:6222,nats://nats_e:6222
