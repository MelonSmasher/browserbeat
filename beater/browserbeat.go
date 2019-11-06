package beater

import (
	"fmt"
	"time"

	"github.com/elastic/beats/libbeat/beat"
	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/logp"

	"github.com/MelonSmasher/browserbeat/config"
)

// Browserbeat configuration.
type Browserbeat struct {
	done   chan struct{}
	config config.Config
	client beat.Client
}

// New creates an instance of browserbeat.
func New(b *beat.Beat, cfg *common.Config) (beat.Beater, error) {
	c := config.DefaultConfig
	if err := cfg.Unpack(&c); err != nil {
		return nil, fmt.Errorf("Error reading config file: %v", err)
	}

	bt := &Browserbeat{
		done:   make(chan struct{}),
		config: c,
	}
	return bt, nil
}

// Run starts browserbeat.
func (bt *Browserbeat) Run(b *beat.Beat) error {
	logp.Info("browserbeat is running! Hit CTRL-C to stop it.")
	var err error
	bt.client, err = b.Publisher.Connect()
	if err != nil {
		return err
	}
	ticker := time.NewTicker(bt.config.Period)

	for {
		select {
		case <-bt.done:
			return nil
		case <-ticker.C:
		}
		browsers := getBrowserHistoryPaths()
		hn := getHostname()
		ipAddresses := getLocalIPs()
		cleanScratchDir()
		for _, browser := range []string{"chrome", "firefox", "safari"} {
			events := readBrowserData(browsers, browser, hn, ipAddresses)
			for _, data := range events {
				event := beat.Event{
					Timestamp: time.Now(),
					Fields: common.MapStr{
						"type": "browser.history",
						"data": data,
					},
				}
				bt.client.Publish(event)
				logp.Info(browser + " event sent")
			}
		}
		cleanScratchDir()
	}
}

// Stop stops browserbeat.
func (bt *Browserbeat) Stop() {
	bt.client.Close()
	close(bt.done)
}
