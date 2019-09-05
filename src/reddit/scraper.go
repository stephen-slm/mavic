package reddit

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strings"
)

type Scraper struct {
	after              int
	scrapingOptions    Options
	supportedFileTypes map[string]bool
	supportedPageTypes map[string]bool
	uniqueImageIds     map[string]map[string]bool
}

// NewRedditScraper creates a instance of the reddit reddit used for taking images
// from the reddit site and downloading them into the given directory. Additionally
// sets the default options and data into the reddit reddit.
func NewScraper(options Options) Scraper {
	redditScraper := Scraper{
		after:              0,
		supportedFileTypes: map[string]bool{"jpeg": true, "png": true, "gif": true, "apng": true, "tiff": true, "pdf": true, "xcf": true},
		supportedPageTypes: map[string]bool{"hot": true, "new": true, "rising": true, "controversial": true, "top": true},
		uniqueImageIds:     map[string]map[string]bool{},
	}

	if options.ImageLimit > 100 {
		fmt.Println("Option 'limit' is currently enforced to 100 or les due ot a on going problem")
		options.ImageLimit = 100
	}

	if options.ImageLimit <= 0 || options.ImageLimit > 500 {
		redditScraper.scrapingOptions.ImageLimit = 50
	}

	if options.FrontPage {
		options.Subreddits = append(options.Subreddits, "frontpage")
	}

	redditScraper.scrapingOptions = options
	return redditScraper
}

// ProcessSubreddits starts the downloading process of all the images
// in the sub reddits
func (s Scraper) ProcessSubreddits() {
	for _, sub := range s.scrapingOptions.Subreddits {
		// if we have not already done this sub reddit before, then
		// create a new unique entry into he unique image list to keep
		// track of all the already downloaded images by id.
		if _, ok := s.uniqueImageIds[sub]; !ok {
			s.uniqueImageIds[sub] = map[string]bool{}
		}

		listings, _ := s.gatherRedditFeed(sub)
		links := parseLinksFromListings(listings)

		dir := path.Join(s.scrapingOptions.OutputDirectory, sub)

		if _, err := os.Stat(dir); os.IsNotExist(err) {
			_ = os.Mkdir(dir, os.ModePerm)
		}

		fmt.Printf("\n\nDownloading %v images from /r/%v", len(links), sub)

		for _, image := range links {
			// if the image id is already been downloaded (the post came up twice) or the image id that we managed
			// to obtain was empty, then continue since we don't have anything to work with. Skipping or attempting
			// to not download a non-existing image.
			if strings.TrimSpace(image.imageId) == "" || s.uniqueImageIds[sub][image.imageId] {
				continue
			}

			fmt.Printf("\nDownloading %20v - /r/%-20v - %v", image.imageId, image.subreddit, image.source)
			s.uniqueImageIds[sub][image.imageId] = true

			downloadImage(dir, image)
		}
	}
}

func downloadImage(outDir string, image Image) {
	// replace gif-v with mp4 for a preferred download as a gif-v file does not work really well on windows
	// machines but require additional processing. While mp4s work fine.
	if strings.HasSuffix(image.link, "gifv") {
		image.link = image.link[:len(image.link)-4] + "mp4"
	}

	// the image id again but this time containing the file type,
	// which allows us to determine the file type without having
	// to do any fancy work.
	imageIdSplit := strings.Split(image.link, "/")
	imageId := imageIdSplit[len(imageIdSplit)-1]

	// returning early if the file already exists, ensuring another check before we go and
	// attempt to download the file, reducing the chance of re-downloading already existing
	// posts.
	imagePath := path.Join(outDir, imageId)
	if _, err := os.Stat(imagePath); !os.IsNotExist(err) {
		return
	}

	out, _ := os.Create(imagePath)
	defer out.Close()

	resp, _ := http.Get(image.link)
	defer resp.Body.Close()

	_, _ = io.Copy(out, resp.Body)
}

// determineRedditUrl will take in a sub reddit that will be used to determine
// what reddit url would be used based on the scraping options, this includes
// setting and marking the image limit and what stage they are currently at.
// handling empty page type or invalid types (defaulting to hot)
func (s Scraper) determineRedditUrl(sub string) string {
	emptyPageType := strings.TrimSpace(s.scrapingOptions.PageType) == ""
	invalidType := s.supportedFileTypes[s.scrapingOptions.PageType]

	if sub == "frontpage" {
		return fmt.Sprintf("https://www.reddit.com/.json?limit=%d&after=%d", s.scrapingOptions.ImageLimit, s.after)
	}
	if emptyPageType || invalidType {
		return fmt.Sprintf("https://www.reddit.com/r/%s/.json?limit=%d&after=%d",
			sub, s.scrapingOptions.ImageLimit, s.after)
	}

	return fmt.Sprintf("https://www.reddit.com/r/%s/%s.json?limit=%d&after=%d",
		sub, s.scrapingOptions.PageType, s.scrapingOptions.ImageLimit, s.after)
}

// Downloads and parses the reddit json feed based on the sub reddit. Ensuring that
// the sub reddit is not empty and ensuring that we send a valid user-agent to ensure
// that reddit does not rate limit us
func (s Scraper) gatherRedditFeed(sub string) (Listings, error) {
	if strings.TrimSpace(sub) == "" {
		return Listings{}, errors.New("sub reddit is required for downloading")
	}

	client := &http.Client{}
	req, _ := http.NewRequest("GET", s.determineRedditUrl(sub), nil)
	req.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64)")

	resp, _ := client.Do(req)

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	return UnmarshalListing(body)
}

// parseLinksFromListings parses all the links and core information out from
// the listings into a more usable formatted listings to allow for a simpler
// image downloading process.
func parseLinksFromListings(listings Listings) []Image {
	if listings.Data == nil || len(listings.Data.Children) == 0 {
		return []Image{}
	}

	// the filtered list of all domains of imgur or that the given post hint
	// states that the given message could be a image.
	var filteredList []Child

	for _, value := range listings.Data.Children {
		if (value.Data.Domain != nil && strings.Contains(*value.Data.Domain, "imgur")) ||
			(value.Data.PostHint != nil && strings.Contains(*value.Data.PostHint, "image")) {
			filteredList = append(filteredList, value)
		}
	}

	// preallocate the direct size required to process all the images, since there is no need to let
	// the underling array double constantly when we already know the size required to process.
	returnableImages := make([]Image, len(filteredList))

	for k, v := range filteredList {
		returnableImages[k] = redditChildToImage(v)
	}

	return returnableImages
}

// redditChildToImage takes in a single reddit listings child data object and converts it to a local
// metadata object that is used to process and download the image.
func redditChildToImage(child Child) Image {
	// the image id is the last section of the source url, so this requires
	// splitting on the forward slash and then taking everything after the dot
	// of the last item and then taking that last item.
	splitUrl := strings.Split(*child.Data.URL, "/")
	imageId := strings.Split(splitUrl[len(splitUrl)-1], ".")[0]

	return Image{
		author: Author{
			link: fmt.Sprintf("https://www.reddit.com/user/%s/", *child.Data.Author),
			name: *child.Data.Author,
		},
		id:        *child.Data.ID,
		imageId:   imageId,
		postLink:  *child.Data.Permalink,
		link:      *child.Data.URL,
		title:     *child.Data.Title,
		subreddit: *child.Data.Subreddit,
		source:    *child.Data.Domain,
	}
}
