// Copyright 2012-2016 Apcera Inc. All rights reserved.

package main

import (
	"flag"
	"log"
	"time"
	"github.com/nats-io/go-nats-streaming"
	"fmt"
	"os"
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

	ch := make(chan bool)
	var guid string
	acb := func(lguid string, err error) {
		if err != nil {
			log.Fatalf("Error in server ack for guid %s: %v\n", lguid, err)
		}
		ch <- true
	}
	finished := make(chan bool)

	funcs := []func () {
		func() {
			for j := 0; j < 10 ; j++ {
				for i := 0; i < publishMsg; i++ {
					guid, err = sc.PublishAsync("BTC", []byte(time.Now().String()), acb)
					if err != nil {
						log.Fatalf("Error during async publish: %v\n", err)
					}
					if guid == "" {
						log.Fatal("Expected non-empty guid to be returned.")
					}
				}
			}
			finished <- true
		},
		func() {
			for j := 0; j < 10 ; j++ {
				for i := 0; i < publishMsg; i++ {
					guid, err = sc.PublishAsync("ETH", []byte(time.Now().String()), acb)
					if err != nil {
						log.Fatalf("Error during async publish: %v\n", err)
					}
					if guid == "" {
						log.Fatal("Expected non-empty guid to be returned.")
					}
				}
			}
			finished <- true
		},
		func() {
			for j := 0; j < 10 ; j++ {
				for i := 0; i < publishMsg; i++ {
					guid, err = sc.PublishAsync("XEM", []byte(time.Now().String()), acb)
					if err != nil {
						log.Fatalf("Error during async publish: %v\n", err)
					}
					if guid == "" {
						log.Fatal("Expected non-empty guid to be returned.")
					}
				}
			}
			finished <- true
		},
	}
	// 並行化する
	for _, sleep := range funcs {
		go sleep()
	}

	// 終わるまで待つ
	for i := 0; i < len(funcs); i++ {
		<-finished
	}

	log.Print("all finished.")

}
