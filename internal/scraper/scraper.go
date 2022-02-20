package scraper

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/mambadev/mavic/internal/reddit"
	"github.com/schollz/progressbar/v3"
)

// The progress bar of the downloading progress that os currently happening
// instead of just happening update notifications.
var progressBar *progressbar.ProgressBar

// DownloadState is the outcome constant of the download process. Used
// to determine the message to be generated and shown to the user.
type DownloadState int

const (
	DOWNLOADING DownloadState = iota
	SUCCESS                   = iota
	SKIPPED                   = iota
	FAILED                    = iota
)

// updateState is used to determine how a downloading progress has occurred and on
// what subreddit that this happened.
type updateState struct {
	// The image metadata that was used to download the given image, this will
	// be used to correctly format a message that will be displayed briefly
	// during the downloading process.
	image reddit.Image
	// The state that the downloading is currently in.
	state DownloadState
}

// Scraper is the type that will be containing all the configuration and
// data used for the parsing process. Including references to already
// downloaded ids + channels for the message and image pump.
type Scraper struct {
	// after is used for when you increase over the number of possible records, so the limit on
	// reddit is 100, so if you ask for 110 images, first we must check the first 100 and then
	// update after to 1, to see the next 10.
	after int
	// the options used for the scraping downloadRedditMetadata, this includes limits, pages, page types and
	// sub reddits to be parsed. This is the central point of truth.
	scrapingOptions Options
	// the supported page types that can be used on reddit, this is hot, new, rising, etc. if
	// the user chooses a unsupported page type, then we will just default to reddits default
	// which is currently hot.
	supportedPageTypes map[string]bool
}

// Start is exposed and called into when a new Scraper is created, this is called
// when the cli commands are parsed and the application is ready to start.
func (s Scraper) Start() {
	// setup the progress bar on start with the rendering of the blank empty state
	// otherwise the loading bar could be displayed before the contents are being
	// parsed.
	progressBar = progressbar.NewOptions(1, progressbar.OptionSetRenderBlankState(s.scrapingOptions.DisplayLoading))

	done := make(chan interface{})
	defer close(done)

	imageStream := s.downloadMetadata(done, progressBar, s.scrapingOptions.Subreddits)

	// The downloaded images once download will pump a message to this channel
	// which will log back out to the user the information they are expecting
	// to be notified that they have been downloaded.
	downloadedMessagePumpChannel := s.downloadImages(done, imageStream)
	var downloaded, failed, skipped int

	for msg := range downloadedMessagePumpChannel {
		var downloadState string
		var addingAmount = 1

		switch msg.state {
		case DOWNLOADING:
			downloadState = "Downloading"
			addingAmount = 0
			break
		case SUCCESS:
			downloadState = "Downloaded"
			downloaded += 1
			break
		case SKIPPED:
			downloadState = "Skipped"
			skipped += 1
			break
		case FAILED:
			downloadState = "Failed Downloading"
			failed += 1
			break
		}

		if s.scrapingOptions.DisplayLoading {
			progressBar.Describe(fmt.Sprintf("%s Image %s from r/%s...", downloadState, msg.image.ImageId, msg.image.Subreddit))
			_ = progressBar.Add(addingAmount)
		}
	}

	if s.scrapingOptions.DisplayLoading {
		progressBar.Describe(fmt.Sprintf("%v images processed. Downloaded %v, skipped %v and failed %v.",
			progressBar.GetMax(), downloaded, skipped, failed))
		_ = progressBar.Finish()
	}
}

// NewRedditScraper creates a instance of the reddit reddit used for taking images
// from the reddit site and downloading them into the given directory. Additionally
// sets the default options and data into the reddit reddit.
func NewScraper(options Options) Scraper {
	redditScraper := Scraper{
		after: 0,
		supportedPageTypes: map[string]bool{"hot": true, "new": true, "rising": true, "best": true,
			"top-hour": true, "top-week": true, "top-month": true, "top-year": true, "top-all": true, "top": true,
			"controversial-hour": true, "controversial-week": true, "controversial-month": true,
			"controversial-year": true, "controversial-all": true, "controversial": true,
		},
	}

	// we don't want to continue to process the data if the given page
	// type is not valid. Determined it will exit earlier over
	// trying to handle it later to improve code quality.
	if !redditScraper.supportedPageTypes[options.PageType] {
		log.Fatalf("Invalid page type '%v' used, reference README for valid page types.\n", options.PageType)
	}

	if options.ImageLimit <= 0 || options.ImageLimit > 500 {
		options.ImageLimit = 50
	}

	if options.ImageLimit > 100 {
		fmt.Println("Option 'limit' is currently enforced to 100 or les due ot a on going problem")
		options.ImageLimit = 100
	}

	if options.FrontPage {
		options.Subreddits = append(options.Subreddits, "frontpage")
	}

	redditScraper.scrapingOptions = options
	return redditScraper
}

// downloads the metadata for a given sub and syncs with a sync group. This will download
// the data, parse it and pump all the images into the download image stream that will
// perform a fan out approach to download all the images.
func (s Scraper) downloadMetadata(done <-chan interface{}, progressBar *progressbar.ProgressBar, subreddit []string) <-chan reddit.Image {
	imageStream := make(chan reddit.Image)

	go func() {
		defer close(imageStream)

		for _, sub := range subreddit {
			select {
			case <-done:
				return
			default:
			}

			listings, _ := s.gatherRedditFeed(sub)
			links := parseLinksFromListings(listings)

			dir := path.Join(s.scrapingOptions.OutputDirectory, sub)

			// if we are only going into the root folder, there is no reason
			// for us to be creating any of the sub folders, just the root.
			if s.scrapingOptions.RootFolderOnly {
				dir = s.scrapingOptions.OutputDirectory
			}

			if _, err := os.Stat(dir); os.IsNotExist(err) {
				_ = os.MkdirAll(dir, os.ModePerm)
			}

			progressBar.ChangeMax(progressBar.GetMax() + len(links))

			for _, image := range links {
				// reassign the sub reddit since it could be the front page and
				// the front page folder is which we want the folder to enter into.
				image.Subreddit = sub

				select {
				case <-done:
					return
				case imageStream <- image:
				}
			}
		}
	}()

	return imageStream
}

// Iterates through the download image pump channel and constantly blocks
// and takes the images pushed to it to be downloaded. calling into the
// download image each time, until closed.
func (s Scraper) downloadImages(done <-chan interface{}, imageStream <-chan reddit.Image) <-chan updateState {
	statusStream := make(chan updateState)

	go func() {
		defer close(statusStream)

		for img := range imageStream {
			select {
			case <-done:
				return
			default:
			}

			outDir := path.Join(s.scrapingOptions.OutputDirectory, img.Subreddit)
			s.downloadImage(statusStream, outDir, img)
		}
	}()

	return statusStream

}

// downloadImage takes in the directory, image and sync group used to
// download a given reddit image to a given directory.
func (s Scraper) downloadImage(statusStream chan<- updateState, outDir string, img reddit.Image) {
	statusStream <- updateState{img, DOWNLOADING}

	// if we are just going into the root, remove everything after the last forward slash.
	if s.scrapingOptions.RootFolderOnly {
		outDir = strings.Replace(outDir, img.Subreddit, "", 1)
	}

	// replace gif-v with mp4 for a preferred download as a gif-v file does not work really well on windows
	// machines but require additional processing. While mp4s work fine.
	if strings.HasSuffix(img.Link, "gifv") {
		img.Link = img.Link[:len(img.Link)-4] + "mp4"
	}

	// the img id again but this time containing the file type,
	// which allows us to determine the file type without having
	// to do any fancy work.
	imageIdSplit := strings.Split(img.Link, "/")
	imageId := imageIdSplit[len(imageIdSplit)-1]

	// returning early if the file already exists, ensuring another check before we go and
	// attempt to download the file, reducing the chance of re-downloading already existing
	// posts.
	imagePath := path.Join(outDir, imageId)
	if _, fileErr := os.Stat(imagePath); !os.IsNotExist(fileErr) {
		statusStream <- updateState{img, SKIPPED}
		return
	}

	out, createErr := os.Create(imagePath)

	// early return if the os failed to create any of the folders, since there is
	// no reason to attempt to download the file if we don't have any where to
	// write the file to after wards.
	if createErr != nil {
		statusStream <- updateState{img, FAILED}
		return
	}

	defer Close(out)
	resp, httpErr := http.Get(img.Link)

	// early return if we failed to download the given file due to a
	// unexpected http error.
	if httpErr != nil {
		statusStream <- updateState{img, FAILED}
		return
	}

	defer Close(resp.Body)
	_, ioErr := io.Copy(out, resp.Body)

	if ioErr != nil {
		statusStream <- updateState{img, FAILED}
		return
	}

	statusStream <- updateState{img, SUCCESS}
}

// Downloads and parses the reddit json feed based on the sub reddit. Ensuring that
// the sub reddit is not empty and ensuring that we send a valid user-agent to ensure
// that reddit does not rate limit us
func (s Scraper) gatherRedditFeed(sub string) (reddit.Listings, error) {
	if strings.TrimSpace(sub) == "" {
		return reddit.Listings{}, errors.New("sub reddit is required for downloading")
	}

	client := &http.Client{}
	req, _ := http.NewRequest("GET", s.determineRedditUrl(sub), nil)
	req.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64)")

	resp, err := client.Do(req)

	if err != nil {
		log.Panic(err)
	}

	defer Close(resp.Body)
	body, _ := ioutil.ReadAll(resp.Body)

	return reddit.UnmarshalListing(body)
}

// parseLinksFromListings parses all the links and core information out from
// the listings into a more usable formatted listings to allow for a simpler
// image downloading downloadRedditMetadata.
func parseLinksFromListings(listings reddit.Listings) []reddit.Image {
	if listings.Data == nil || len(listings.Data.Children) == 0 {
		return []reddit.Image{}
	}

	// the filtered list of all domains of imgur or that the given post hint
	// states that the given message could be a image.
	var filteredList []reddit.Child

	for _, value := range listings.Data.Children {
		if (value.Data.Domain != nil && strings.Contains(*value.Data.Domain, "imgur")) ||
			(value.Data.PostHint != nil && strings.Contains(*value.Data.PostHint, "image")) {

			splitLink := strings.Split(*value.Data.URL, "/")

			// ensure that we have not got a gallery or something, making sure that
			// what we are downloading is a direct image and nothing else.
			if strings.Contains(splitLink[len(splitLink)-1], ".") {
				filteredList = append(filteredList, value)
			}
		}
	}

	// preallocate the direct size required to downloadRedditMetadata all the images, since there is no need to let
	// the underling array double constantly when we already know the size required to downloadRedditMetadata.
	returnableImages := make([]reddit.Image, len(filteredList))

	for k, v := range filteredList {
		image := reddit.RedditChildToImage(v)

		// if the image id is already been downloaded (the post came up twice) or the image id that we managed
		// to obtain was empty, then continue since we don't have anything to work with. Skipping or attempting
		// to not download a non-existing image.
		if strings.TrimSpace(image.ImageId) == "" {
			continue
		}

		returnableImages[k] = image
	}

	return returnableImages
}

// determineRedditUrl will take in a sub reddit that will be used to determine
// what reddit url would be used based on the scraping options, this includes
// setting and marking the image limit and what stage they are currently at.
// (defaulting to hot)
func (s Scraper) determineRedditUrl(sub string) string {
	pageType := s.scrapingOptions.PageType
	additional := ""

	// if a page type is a type that supports having a time span (e.g top and controversial) then
	// split out the page type and adjust the additional to contain the time span and assign the page
	// type to the correct reddit representation.
	if strings.Contains(pageType, "-") {
		pageSplit := strings.Split(pageType, "-")
		additional = fmt.Sprintf("&t=%v", pageSplit[1])
		pageType = pageSplit[0]
	}

	if sub == "frontpage" {
		return fmt.Sprintf("https://www.reddit.com/%v/.json?limit=%v&after=%v%v", pageType, s.scrapingOptions.ImageLimit, s.after, additional)
	}

	url := fmt.Sprintf("https://www.reddit.com/r/%v/%v.json?limit=%v&after=%v%v",
		sub, pageType, s.scrapingOptions.ImageLimit, s.after, additional)

	return url
}

// Close is designed to handle a defer closed on a closer. Correctly and
// fatally exiting if a error occurs on the close.
func Close(c io.Closer) {
	err := c.Close()
	if err != nil {
		log.Fatal(err)
	}
}
