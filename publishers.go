package nca

import (
	"fmt"
)

type publisherFactory func() Publisher

var publisherFactories = map[string]publisherFactory{
	"spewpublisher":   func() Publisher { return &SpewPublisher{} },
	"memorypublisher": func() Publisher { return &MemoryPublisher{} },
	"execpublisher":   func() Publisher { return &ExecPublisher{} },
}

// newPublisher returns a new Publisher of the specified type.
// The type must be listed in publisherFactory otherwise this function will panic.
func newPublisher(publisherType string) Publisher {
	publisherFactory, found := publisherFactories[publisherType]
	if !found {
		panic(fmt.Sprintf("Invalid publisher type %q specified", publisherType))
	}
	return publisherFactory()
}
