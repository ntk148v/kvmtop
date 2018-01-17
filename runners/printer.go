package runners

import (
	"sync"
	"time"

	"github.com/cha87de/kvmtop/config"
	"github.com/cha87de/kvmtop/models"
)

var collectors []string

func initializePrinter(wg *sync.WaitGroup) {
	// open configured printer
	models.Collection.Printer.Open()

	// define collectors and their order
	for collectorName := range models.Collection.Collectors {
		collectors = append(collectors, collectorName)
	}

	// start continuously printing values
	for n := -1; config.Options.Runs == -1 || n < config.Options.Runs; n++ {
		start := time.Now()
		handleRun()
		nextRun := start.Add(time.Duration(config.Options.Frequency) * time.Second)
		time.Sleep(nextRun.Sub(time.Now()))
	}

	// close configured printer
	models.Collection.Printer.Close()

	// return from runner
	wg.Done()
}

func handleRun() {
	var fields []string
	var values [][]string

	// collect fields for each collector
	fields = append(fields, "UUID", "name")
	for _, collectorName := range collectors {
		collector := models.Collection.Collectors[collectorName]
		output := collector.PrintFields()
		fields = append(fields, output[0:]...)
	}

	// collect values for each domain
	for _, domain := range models.Collection.Domains {
		var domvalues []string
		domvalues = append(domvalues, domain.UUID, domain.Name)
		for _, collectorName := range collectors {
			collector := models.Collection.Collectors[collectorName]
			output := collector.PrintValues(domain)
			domvalues = append(domvalues, output[0:]...)
		}
		values = append(values, domvalues)
	}

	models.Collection.Printer.Screen(fields, values)
}
