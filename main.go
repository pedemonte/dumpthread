package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/nbd-wtf/go-nostr"
	"github.com/nbd-wtf/go-nostr/nip19"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Suggest: true,
		Name:    "dumpthread",
		Usage:   "save a nostr thread to local filesystem",
		Commands: []*cli.Command{
			{
				Name:    "save",
				Aliases: []string{"s"},
				Usage:   "save a thread on filesystem",
				Action: func(cCtx *cli.Context) error {
					relays := cCtx.StringSlice("relay")
					eventId := cCtx.String("eventid")
					outDir := cCtx.String("outdir")

					err := processFromAllRelay(relays, eventId, outDir)

					return err
				},
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "outdir",
						Aliases: []string{"o"},
						Usage:   "the path where save the json events to",
						Value:   ".",
					},
					&cli.StringFlag{
						Name:     "eventid",
						Aliases:  []string{"e"},
						Usage:    "the event ID of the first event of the thread",
						Required: true,
					},
					&cli.StringSliceFlag{
						Name:     "relay",
						Aliases:  []string{"r"},
						Usage:    "the relays to use (can be use multiple times)",
						Required: true,
					},
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

// saveEventOnDisk saves an event to the disk at the given path
// by using its ID as the base name and the '.json' extension.
func saveEventOnDisk(e *nostr.Event, path string) error {
	json, err := e.MarshalJSON()
	if err != nil {
		return err
	}

	filename := fmt.Sprintf("%s/%s.json", path, e.ID)

	err = os.WriteFile(filename, json, 0666)
	return err
}

// processFromAllRelay downloads the specified note and all the notes that reference
// it from the provided relays. It then recursively attempts to download the notes,
// using the same criteria, from all the relays suggested within the notes themselves.
func processFromAllRelay(relays []string, noteId string, outDir string) error {
	var wg sync.WaitGroup

	eventChan := make(chan *nostr.Event, 20)
	processedEvents := make(map[string]struct{}, 1)

	relayChan := make(chan string, 20)
	relaysVisited := make(map[string]struct{}, 1)

	ctx := context.Background()

	for _, r := range relays {
		log.Printf("adding relay: %s\n", r)
		relayChan <- r
	}

	// processes the events received from the eventChan channel.
	// It maintains a record of already processed events to ensure they are not processed more than once.
	wg.Add(1)

	go func() {
		defer wg.Done()
		for e := range eventChan {
			_, exists := processedEvents[e.ID]
			if !(exists) {
				err := saveEventOnDisk(e, outDir)
				if err != nil {
					panic(err)
				}

				// track the event as processed
				processedEvents[e.ID] = struct{}{}
			}
		}
	}()

loop:
	for {
		select {
		case relay, ok := <-relayChan:
			if !ok {
				// Channel is closed, exit the loop
				break
			}

			_, exists := relaysVisited[relay]
			if exists {
				continue
			}

			log.Printf("getting notes from: %s\n", relay)

			evs, err := getThread(ctx, relay, noteId)
			if err != nil {
				log.Println(err)
			} else {
				log.Printf("got %d note(s) from: %s\n", len(evs), relay)

				for _, e := range evs {
					eventChan <- e
				}

				suggestedRelays := getSuggestedRelays(evs)

				// send new suggested relays to the relay channel
				for _, r := range suggestedRelays {
					if len(r) > 0 {
						log.Printf("adding relay %s as suggested relays\n", r)
						relayChan <- r
					}
				}
			}

			relaysVisited[relay] = struct{}{}
		default:
			// Channel is empty, exit the loop
			break loop
		}
	}

	close(eventChan)
	wg.Wait()

	return nil
}

// getThread retrieves all the notes that reference a given note ID.
// It returns the note itself, the notes referring it and the suggested relays provided by their respective authors.
func getThread(ctx context.Context, relayUrl string, noteId string) ([]*nostr.Event, error) {
	relay, err := nostr.RelayConnect(ctx, relayUrl)
	if err != nil {
		return []*nostr.Event{}, err
	}

	ctx, cancel := context.WithTimeout(ctx, 20*time.Second)
	defer cancel()
	defer relay.Close()

	// noteId is required in hex notation
	_, noteIdHex, err := nip19.Decode(noteId)

	if err != nil {
		return []*nostr.Event{}, fmt.Errorf("error decoding note ID %s, %w", noteId, err)
	}

	noteFilter := nostr.Filter{
		Limit: 1,
		IDs:   []string{noteIdHex.(string)},
	}

	referringTag := nostr.TagMap{"e": []string{noteIdHex.(string)}}

	referringFilter := nostr.Filter{
		Tags: referringTag,
	}

	mainNote, err := relay.QuerySync(ctx, noteFilter)
	if err != nil {
		return []*nostr.Event{}, fmt.Errorf("retrieving the specified note: %w", err)
	}

	referringNotes, err := relay.QuerySync(ctx, referringFilter)
	if err != nil {
		return []*nostr.Event{}, fmt.Errorf("retrieving the referring notes: %w", err)
	}

	allNotes := append(mainNote, referringNotes...)

	return allNotes, nil
}

// getSuggestedRelays retrieves the suggested relays from a slice of events
func getSuggestedRelays(evs []*nostr.Event) []string {
	relays := make(map[string]struct{})

	for _, e := range evs {
		for _, t := range e.Tags {
			relay := strings.TrimSpace(t.Relay())
			relays[relay] = struct{}{}
		}
	}

	dedupedRelays := make([]string, len(relays))

	for k := range relays {
		dedupedRelays = append(dedupedRelays, k)
	}

	return dedupedRelays
}
