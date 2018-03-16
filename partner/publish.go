// Copyright 2012-2016 Apcera Inc. All rights reserved.

package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/nats-io/go-nats-streaming"
)

var usageStr = `
Usage: publish [options] <subject> <message> <subscribe_subject>

Options:
	-s, --server   <url>            NATS Streaming server URL(s)
	-c, --cluster  <cluster name>   NATS Streaming cluster name
	-id,--clientid <client ID>      NATS Streaming client ID
	-a, --async                     Asynchronous publish mode
`

// NOTE: Use tls scheme for TLS, e.g. stan-pub -s tls://demo.nats.io:4443 foo hello
func usage() {
	fmt.Printf("%s\n", usageStr)
	os.Exit(0)
}

func main() {
	var clusterID string
	var clientID string
	var async bool
	var URL string

	flag.StringVar(&URL, "s", "nats://nats_a:4222,nats://nats_b:4223", "The nats server URLs (separated by comma)")
	flag.StringVar(&URL, "server", "nats://nats_a:4222,nats://nats_b:4223", "The nats server URLs (separated by comma)")
	flag.StringVar(&clusterID, "c", "test-cluster", "The NATS Streaming cluster ID")
	flag.StringVar(&clusterID, "cluster", "test-cluster", "The NATS Streaming cluster ID")
	flag.StringVar(&clientID, "id", "stan-pub", "The NATS Streaming client ID to connect with")
	flag.StringVar(&clientID, "clientid", "stan-pub", "The NATS Streaming client ID to connect with")
	flag.BoolVar(&async, "a", false, "Publish asynchronously")
	flag.BoolVar(&async, "async", false, "Publish asynchronously")

	log.SetFlags(0)
	flag.Usage = usage
	flag.Parse()

	args := flag.Args()

	if len(args) < 1 {
		usage()
	}

	sc, err := stan.Connect(clusterID, clientID, stan.NatsURL(URL))
	if err != nil {
		log.Fatalf("Can't connect: %v.\nMake sure a NATS Streaming Server is running at: %s", err, URL)
	}
	defer sc.Close()

	subj, msg , s_subj := args[0], []byte(args[1]) , args[2]

	ch := make(chan bool)
	var glock sync.Mutex
	var guid string
	acb := func(lguid string, err error) {
		glock.Lock()
		log.Printf("Received ACK for guid %s\n", lguid)
		defer glock.Unlock()
		if err != nil {
			log.Fatalf("Error in server ack for guid %s: %v\n", lguid, err)
		}
		if lguid != guid {
			log.Fatalf("Expected a matching guid in ack callback, got %s vs %s\n", lguid, guid)
		}
		ch <- true
	}

	if !async {
		// err = sc.Publish(subj, msg)

		// partner Request & Reply
		ret, err := sc.NatsConn().Request(subj, msg, 10000*time.Millisecond)

		log.Printf("Published [%s] : '%s'\n", subj, msg)
		if err != nil {
			log.Fatalf("partner is dead : %v\n", err)
		} else {
			log.Printf("Received [%v] : '%s'\n", ret.Subject, string(ret.Data))
		}

		if string(ret.Data) == "OK" {
			err = sc.Publish(s_subj, msg)
			log.Printf("Published [%s] : '%s'\n", s_subj, msg)
			if err != nil {
				log.Fatalf("Error during publish: %v\n", err)
			}
		}
	} else {
		glock.Lock()
		guid, err = sc.PublishAsync(subj, msg, acb)
		if err != nil {
			log.Fatalf("Error during async publish: %v\n", err)
		}
		glock.Unlock()
		if guid == "" {
			log.Fatal("Expected non-empty guid to be returned.")
		}
		log.Printf("Published [%s] : '%s' [guid: %s]\n", subj, msg, guid)

		select {
		case <-ch:
			break
		case <-time.After(5 * time.Second):
			log.Fatal("timeout")
		}

	}
}
