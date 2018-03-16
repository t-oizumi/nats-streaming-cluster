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

const (
	DefaultNumMsgs            = 100000
	DefaultNumPubs            = 1
	DefaultNumSubs            = 0
	DefaultAsync              = false
	DefaultMessageSize        = 128
	DefaultIgnoreOld          = false
	DefaultMaxPubAcksInflight = 1000
	DefaultClientID           = "benchmark"
	DefaultPublishMsg         = 10
)
var usageStr = `
Usage: publish [options] <subject> <subscribe_subject>

Options:
	-s, --server   <url>            NATS Streaming server URL(s)
	-c, --cluster  <cluster name>   NATS Streaming cluster name
	-id,--clientid <client ID>      NATS Streaming client ID
	-a, --async                     Asynchronous publish mode
    -ms, --messagesize <msgbyte>        byte
    -pm, --publishMessageSize
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
	var messagesize int
	var publishMsg int

	flag.StringVar(&URL, "s", "nats://nats_a:4222,nats://nats_b:4222,nats://nats_c:4222", "The nats server URLs (separated by comma)")
	flag.StringVar(&URL, "server", "nats://nats_a:4222,nats://nats_b:4222,nats://nats_c:4222", "The nats server URLs (separated by comma)")
	flag.StringVar(&clusterID, "c", "test-cluster", "The NATS Streaming cluster ID")
	flag.StringVar(&clusterID, "cluster", "test-cluster", "The NATS Streaming cluster ID")
	flag.StringVar(&clientID, "id", "stan-pub", "The NATS Streaming client ID to connect with")
	flag.StringVar(&clientID, "clientid", "stan-pub", "The NATS Streaming client ID to connect with")
	flag.BoolVar(&async, "a", false, "Publish asynchronously")
	flag.BoolVar(&async, "async", false, "Publish asynchronously")
	flag.IntVar(&messagesize ,"ms", DefaultMessageSize, "Message size in bytes.")
	flag.IntVar(&publishMsg ,"pm", DefaultPublishMsg, "Message size in bytes.")

	log.SetFlags(0)
	flag.Usage = usage
	flag.Parse()

	args := flag.Args()

	if len(args) < 1 {
		usage()
	}
	//stan.MaxPubAcksInflight(100000)
	sc, err := stan.Connect(clusterID, clientID, stan.NatsURL(URL))
	if err != nil {
		log.Fatalf("Can't connect: %v.\nMake sure a NATS Streaming Server is running at: %s", err, URL)
	}
	defer sc.Close()
	var msg []byte
	msg = make([]byte , messagesize)
	for i:= 0 ; i < messagesize ; i++ {
		msg[i] = 'a'
	}
	subj, s_subj := args[0], args[1]

	ch := make(chan bool)
	var glock sync.Mutex
	var guid string
	acb := func(lguid string, err error) {
		//glock.Lock()
		// log.Printf("Received ACK for guid %s\n", lguid)
		//defer glock.Unlock()
		if err != nil {
			log.Fatalf("Error in server ack for guid %s: %v\n", lguid, err)
		}
		//if lguid != guid {
		//	log.Fatalf("Expected a matching guid in ack callback, got %s vs %s\n", lguid, guid)
		//}
		ch <- true
	}


	if !async {
		//for j := 0; j < 10 ; j++ {
			//sTime := time.Now();
			//log.Printf("Published [%s] start : %s\n", s_subj, time.Now())
			for i := 0; i < publishMsg; i++ {
				time.Sleep(3 * time.Millisecond)
				//ret, err := sc.NatsConn().Request(subj, msg, 1000*time.Millisecond)
				//startTime := time.Now()
				//ret, err := sc.NatsConn().Request(subj, []byte(time.Now().String()), 1000*time.Millisecond)
				//log.Printf("Published [%s] : '%s'\n", i + publishMsg * j, time.Now().Sub(startTime))


				//if err != nil {
				//	log.Fatalf("partner is dead : %v\n", err)
				//}
				//if string(ret.Data) == "OK" {
				//	log.Printf("start:%s\n",  time.Now())
					startTime := time.Now()
					//err = sc.NatsConn().Publish(s_subj, []byte(time.Now().String()))
					err = sc.Publish(s_subj, []byte(time.Now().String()))
					//err = sc.Publish(s_subj, msg)
					log.Printf("end:%s\n", time.Now().Sub(startTime))

					//if err != nil {
					//	log.Fatalf("Error during publish: %v\n", err)
					//}

					//glock.Lock()
					//guid, err = sc.PublishAsync(s_subj, msg, acb)
					//guid, err = sc.PublishAsync("BTC", []byte(time.Now().String()), acb)
					//if err != nil {
					//	log.Fatalf("Error during async publish: %v\n", err)
					//}
					////glock.Unlock()
					//if guid == "" {
					//	log.Fatal("Expected non-empty guid to be returned.")
					//}
					//guid, err = sc.PublishAsync("ETH", []byte(time.Now().String()), acb)
					//if err != nil {
					//	log.Fatalf("Error during async publish: %v\n", err)
					//}
					//////glock.Unlock()
					//if guid == "" {
					//	log.Fatal("Expected non-empty guid to be returned.")
					//}
					//
					//guid, err = sc.PublishAsync("XEM", []byte(time.Now().String()), acb)
					//if err != nil {
					//	log.Fatalf("Error during async publish: %v\n", err)
					//}
					////glock.Unlock()
					//if guid == "" {
					//	log.Fatal("Expected non-empty guid to be returned.")
					//}
					//log.Printf("Published [%s] : '%s' [guid: %s]\n", subj, time.Now(), guid)
				//}
			}
			//eTime := time.Now()
			//log.Printf("Published [%s] end  : %s\n", s_subj, time.Now())
			//log.Printf("Published [%s]      : %s\n", s_subj, eTime.Sub(sTime))
		//}
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
