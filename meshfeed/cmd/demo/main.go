// Command demo runs the exact scenario this blueprint exists to prove:
// three simulated nodes, a subscriber that gets nothing while out of
// range, then syncs automatically — through an intermediate relay that
// never subscribed to anything — once in range. No real radio is
// available in this environment, so it runs against simtransport; see
// BLUEPRINT.md for what a real LoRa Transport would replace here.
package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"meshfeed/internal/api"
	"meshfeed/internal/domain"
	"meshfeed/internal/mesh"
	"meshfeed/internal/simtransport"
)

func main() {
	net := simtransport.NewNetwork()
	publisher := mesh.NewNode("publisher", net.NewTransport("publisher"))
	relay := mesh.NewNode("relay", net.NewTransport("relay"))
	subscriber := mesh.NewNode("subscriber", net.NewTransport("subscriber"))

	for _, n := range []*mesh.Node{publisher, relay, subscriber} {
		n := n
		n.Log = func(format string, args ...any) { fmt.Printf("[%s] "+format+"\n", append([]any{n.ID}, args...)...) }
		go n.Run()
	}

	feed := domain.NewFeedID("garden-club")
	subscriber.Subscribe(feed)

	fmt.Println("=== nobody is in range of anybody yet ===")
	if _, err := publisher.Publish(feed, []byte("compost bin needs turning"), uint32(time.Now().Unix())); err != nil {
		panic(err)
	}
	must(publisher.AdvertiseAll())
	time.Sleep(200 * time.Millisecond)
	fmt.Printf("subscriber latest seq: %d (expect 0 — nothing could reach it)\n\n", subscriber.Store.LatestSeq(feed))

	fmt.Println("=== publisher <-> relay <-> subscriber move into range (a chain, not direct) ===")
	net.SetInRange("publisher", "relay", true)
	net.SetInRange("relay", "subscriber", true)
	must(publisher.AdvertiseAll())
	time.Sleep(300 * time.Millisecond)

	items := subscriber.Store.ItemsSince(feed, 0)
	fmt.Printf("subscriber now holds %d item(s), received via the relay it never configured a subscription with:\n", len(items))
	for _, it := range items {
		fmt.Printf("  #%d: %s\n", it.Seq, it.Content)
	}

	fmt.Println("\n=== publisher posts a second update; subscriber picks it up on the next advertise ===")
	must2(publisher.Publish(feed, []byte("tomato seedlings ready for pickup"), uint32(time.Now().Unix())))
	must(publisher.AdvertiseAll())
	time.Sleep(300 * time.Millisecond)

	items = subscriber.Store.ItemsSince(feed, 0)
	fmt.Printf("subscriber now holds %d item(s) total:\n", len(items))
	for _, it := range items {
		fmt.Printf("  #%d: %s\n", it.Seq, it.Content)
	}

	fmt.Println("\n=== serving the subscriber's feed store at http://localhost:8080 (Ctrl-C to stop) ===")
	handler := api.NewHandler(subscriber.Store)
	log.Fatal(http.ListenAndServe(":8080", handler.Routes()))
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}

func must2(_ domain.Item, err error) {
	must(err)
}
