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

	// A list of supported browsers
	browserNames := getSupportedBrowsers()

	for {
		select {
		case <-bt.done:
			return nil
		case <-ticker.C:
		}

		//**************************** Initial information section ****************************************************/
		//**** This runs on each iteration to pick up any changes to the environment without restarting browserbeat ***/
		// Get all all paths to browser history databases for each user on the machine
		browsers := getBrowserHistoryPaths()
		// Get the machine's hostname
		hn := getHostname()
		// Get the machine's local network IP addresses
		ipAddresses := getLocalIPs()
		//*************************************************************************************************************/

		// Clean up any copied browser databases from the last run
		cleanScratchDir()

		//********************************* Begin main loop ***********************************************************/
		// Loop for each supported browser
		for _, browser := range browserNames {
			// Read the current target browser's data base for each user on the machine
			events := readBrowserData(browsers, browser, hn, ipAddresses)
			// Loop through all of the browser history events
			for _, data := range events {
				// Store the event in a shippable format
				event := beat.Event{
					Timestamp: time.Now(),
					Fields: common.MapStr{
						"type": "browser.history",
						"data": data,
					},
				}
				// Send the data on it's way
				bt.client.Publish(event)
				// Log that we sent the data
				logp.Info(browser + " event sent")
			}
		}
		//******************************* End main loop ***************************************************************/

		// Clean up any copied browser databases from this run
		cleanScratchDir()
	}
}

// Stop stops browserbeat.
func (bt *Browserbeat) Stop() {
	bt.client.Close()
	close(bt.done)
}
