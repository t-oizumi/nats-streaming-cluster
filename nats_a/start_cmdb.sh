#!/bin/bash
sleep $1
cd /go/src/github.com/nats-io/nats-streaming-server/
### 3台構成
./nats-streaming-server -ns nats://nats_a:4222 -cluster_sync=false -file_sync=false -store memory -dir /go/src/store/nats -cluster_log_path /go/src/store/nats_cluster -cluster_node_id a -cluster_peers b,c -clustered -cluster nats://nats_a:6222 -routes nats://nats_c:6222,nats://nats_b:6222
