package mutlibars

import (
	"fmt"

	"github.com/schollz/progressbar"
)

var multiProgressBarCache map[string]*progressbar.ProgressBar

// GetOrCreateProgressbarByName will attempt to first gather the progress bar by
// its given reference name, if this errors than the progressbar will be created with
// the given max being the number of iterations before being 100%. Returning the newly
// created progress bar if not existing.
func GetOrCreateProgressbarByName(ref string, max int) *progressbar.ProgressBar {
	existingBar, err := GetProgressbarByName(ref)

	if err != nil {
		return existingBar
	}

	return CreateProgressbar(ref, max)
}

// CreateProgressbar will go and create a new progress bar and store it into
// the cache based on the reference name, allowing the given progress bar
// to be gathered later again by the ref name. Returning the progressbar
// pointer at the end.
func CreateProgressbar(ref string, max int) *progressbar.ProgressBar {
	bar := progressbar.New(max)

	multiProgressBarCache[ref] = bar
	return bar
}

// GetProgressbarByName returns a existing progress bar pointer by the name
// that is given, if the progress bar does not exist at the time of being called
// then a error is returned in response.
func GetProgressbarByName(ref string) (*progressbar.ProgressBar, error) {
	if bar, ok := multiProgressBarCache[ref]; ok {
		return bar, nil
	} else {
		return nil, fmt.Errorf("progress bar does not exist by name %v", ref)
	}
}
