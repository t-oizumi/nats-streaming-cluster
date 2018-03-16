#!/bin/bash
sleep $1
cd /go/src/github.com/nats-io/nats-streaming-server/
### 3台構成
#./nats-streaming-server -p 4222 -cluster_sync=false -file_sync=false -store memory -dir /go/src/store/nats -cluster_log_path /go/src/store/nats_cluster -cluster_node_id a -cluster_peers b,c -clustered -cluster nats://nats_a:6222 -routes nats://nats_c:6222,nats://nats_b:6222
### 5台構成
#./nats-streaming-server -p 4222 -cluster_sync=false -file_sync=false -store file -dir /go/src/store/nats -cluster_log_path /go/src/store/nats_cluster -cluster_node_id a -cluster_peers b,c,d,e -clustered -cluster nats://nats_a:6222 -routes nats://nats_c:6222,nats://nats_b:6222,nats://nats_d:6222,nats://nats_e:6222
### 自前3台構成
#./nats-streaming-server -p 4222 -cluster_sync=false -file_sync=false -store memory -dir /go/src/store/nats1 -cluster_log_path /go/src/store/nats_cluster1 -cluster_node_id a -cluster_peers b,c -clustered -cluster nats://nats_a:6222 -routes nats://nats_a:6223,nats://nats_a:6224
#./nats-streaming-server -p 4223 -cluster_sync=false -file_sync=false -store memory -dir /go/src/store/nats2 -cluster_log_path /go/src/store/nats_cluster2 -cluster_node_id b -cluster_peers a,c -clustered -cluster nats://nats_a:6223 -routes nats://nats_a:6222,nats://nats_a:6224
./nats-streaming-server -p 4224 -cluster_sync=false -file_sync=false -store memory -dir /go/src/store/nats3 -cluster_log_path /go/src/store/nats_cluster3 -cluster_node_id c -cluster_peers a,b -clustered -cluster nats://nats_a:6224 -routes nats://nats_a:6222,nats://nats_a:6223
