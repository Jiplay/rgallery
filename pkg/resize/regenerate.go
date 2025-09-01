package resize

import (
	"path/filepath"
	"sync"

	"github.com/robbymilo/rgallery/pkg/config"
	"github.com/robbymilo/rgallery/pkg/sizes"
)

// HandleResize coordinates the parallel (go routine) generation of all non-existing thumbnails.
func HandleResize(regenThumb bool, media Media, c Conf) (int, error) {
	// build a map of sizes for the thumbnail
	var s []int
	final := false

	for _, size := range sizes.GetSizes() {
		if size <= media.Width {
			s = append(s, size)
		} else if !final {
			final = true
			s = append(s, media.Width)
		}
	}

	if len(s) > 0 && c.PreGenerateThumb && regenThumb {
		var wg sync.WaitGroup
		errChan := make(chan error, len(s))

		for _, size := range s {
			wg.Add(1)
			go AddImageThumb(media, size, c, &wg, errChan)
		}

		// Start a goroutine to close error channel when all workers are done
		go func() {
			wg.Wait()
			close(errChan)
		}()

		// Collect any errors from the goroutines
		var errors []error
		for err := range errChan {
			if err != nil {
				errors = append(errors, err)
			}
		}

		// If any errors occurred, return the first one
		if len(errors) > 0 {
			return len(s), errors[0]
		}
	}

	return len(s), nil
}

// AddImageThumb creates a request to generate a single thumbnail from the original image.
func AddImageThumb(media Media, size int, c Conf, wg *sync.WaitGroup, errChan chan<- error) {
	defer wg.Done()

	path := filepath.Join(config.MediaPath(c), media.Path)
	_, err := GenerateSingleThumb(path, media, size, c)
	if err != nil {
		c.Logger.Error("error generating thumbnail",
			"path", media.Path,
			"size", size,
			"error", err)
		errChan <- err
		return
	}
	errChan <- nil
}
