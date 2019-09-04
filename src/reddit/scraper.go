package reddit

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
	"strings"
)

type Scraper struct {
	after              int
	scrapingOptions    Options
	supportedFileTypes []string
	supportedPageTypes []string
	uniqueImageIds     map[string][]string
}

// NewRedditScraper creates a instance of the reddit reddit used for taking images
// from the reddit site and downloading them into the given directory. Additionally
// sets the default options and data into the reddit reddit.
func NewScraper(options Options) Scraper {
	redditScraper := Scraper{
		after:              0,
		supportedFileTypes: []string{"jpeg", "png", "gif", "apng", "tiff", "pdf", "xcf"},
		supportedPageTypes: []string{"hot", "new", "rising", "controversial", "top"},
		uniqueImageIds:     map[string][]string{},
	}

	if options.imageLimit > 100 {
		fmt.Println("Option 'limit' is currently enforced to 100 or les due ot a on going problem")
		options.imageLimit = 100
	}

	if options.imageLimit <= 0 || options.imageLimit > 500 {
		redditScraper.scrapingOptions.imageLimit = 50
	}

	if !options.frontPage {
		options.subreddits = append(options.subreddits, "frontpage")
	}

	redditScraper.scrapingOptions = options
	return redditScraper
}

// ProcessSubreddits starts the downloading process of all the images
// in the sub reddits
func (s Scraper) ProcessSubreddits() {
	for _, sub := range s.scrapingOptions.subreddits {
		// if we have not already done this sub reddit before, then
		// create a new unique entry into he unique image list to keep
		// track of all the already downloaded images by id.
		if _, ok := s.uniqueImageIds[sub]; !ok {
			s.uniqueImageIds[sub] = []string{}
		}

		listings, _ := s.gatherRedditFeed(sub)
		links := parseLinksFromListings(listings)

		fmt.Println(links)
	}
}

// determineRedditUrl will take in a sub reddit that will be used to determine
// what reddit url would be used based on the scraping options, this includes
// setting and marking the image limit and what stage they are currently at.
// handling empty page type or invalid types (defaulting to hot)
func (s Scraper) determineRedditUrl(sub string) string {
	emptyPageType := strings.TrimSpace(s.scrapingOptions.pageType) == ""
	invalidType := sort.SearchStrings(s.supportedPageTypes, s.scrapingOptions.pageType) == len(s.supportedPageTypes)

	if sub == "frontpage" {
		return fmt.Sprintf("https://www.reddit.com/.json?limit=%d&after=%d", s.scrapingOptions.imageLimit, s.after)
	}
	if emptyPageType || invalidType {
		return fmt.Sprintf("https://www.reddit.com/r/%s/.json?limit=%d&after=%d",
			sub, s.scrapingOptions.imageLimit, s.after)
	}

	return fmt.Sprintf("https://www.reddit.com/r/%s/%s.json?limit=%d&after=%d",
		sub, s.scrapingOptions.pageType, s.scrapingOptions.imageLimit, s.after)
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
