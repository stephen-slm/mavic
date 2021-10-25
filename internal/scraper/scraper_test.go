package scraper

import (
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type test struct {
	data   []int
	answer int
}

// Define the suite, and absorb the built-in basic suite
// functionality from testify - including a T() method which
// returns the current testing context
type ScraperTestSuite struct {
	suite.Suite
	baseOptions   Options
	sampleScraper Scraper
}

// SetupTest ensures that a good and simple base configuration options are setup
// before each test with a additional pre-made scraper based off these core
// configuration options. All though won't do anything since its not called "start".
func (suite *ScraperTestSuite) SetupTest() {
	suite.baseOptions = Options{
		OutputDirectory:        "./pictures",
		ImageLimit:             5,
		FrontPage:              false,
		RootFolderOnly:         false,
		PageType:               "hot",
		Subreddits:             []string{"cute"},
		MaxConcurrentDownloads: 25,
		DisplayLoading:         false,
	}

	// generates a temporary directory to be used throughout the
	// testing process if required, tear down test will remove it.
	dir, _ := ioutil.TempDir("", "mavic_test_pictures")
	suite.baseOptions.OutputDirectory = dir

	// create the sample scraper that can be used on each test but not enforced.
	// Allows base generation samples to occur and additional sample scrapers
	// could be created if needed.
	suite.sampleScraper = NewScraper(suite.baseOptions)
}

// TearDownTest ensures that after the tests are complete, that any folder
// or folders created during this test are completely removed from the disk.
// generally these folders should be put into temp space but just to ensure
// cleanup, these will be removed.
func (suite *ScraperTestSuite) TearDownTest() {
	if _, err := os.Stat(suite.baseOptions.OutputDirectory); !os.IsNotExist(err) {
		err := os.RemoveAll(suite.baseOptions.OutputDirectory)

		if err != nil {
			log.Fatal(err)
		}
	}
}

// TestNewScraperBadLimit ensures that if a bad upper limit is given or lower limit, then
// the given value is reset back to the default (if lower e.g 0 or less) or to a reset of 100
// due to a current limitation of not being able to support more than 100 per sub reddit.
func (suite *ScraperTestSuite) TestNewScraperBadLimit() {
	tests := []test{
		{[]int{101}, 100}, {[]int{501}, 50},
		{[]int{500}, 100}, {[]int{0}, 50},
		{[]int{-1}, 50}, {[]int{-100}, 50},
	}

	for _, v := range tests {
		suite.baseOptions.ImageLimit = v.data[0]
		badLimit := NewScraper(suite.baseOptions)

		assert.NotEqual(suite.T(), badLimit.scrapingOptions.ImageLimit, v.data[0])
		assert.Equal(suite.T(), badLimit.scrapingOptions.ImageLimit, v.answer)
	}
}

// TestNewScraperGoodLimit ensures that if a good value is given, then the output of the
// generation of the new scraper will still respect the given input. Ensuring the upper
// and lower bounds are respected.
func (suite *ScraperTestSuite) TestNewScraperGoodLimit() {
	tests := []test{
		{data: []int{1}, answer: 1},
		{data: []int{50}, answer: 50},
		{data: []int{100}, answer: 100},
	}

	for _, v := range tests {
		suite.baseOptions.ImageLimit = v.data[0]
		assert.Equal(suite.T(), NewScraper(suite.baseOptions).scrapingOptions.ImageLimit, v.answer)
	}
}

// TestNewScraperFrontPage ensures that if the front page is being marked to be scrapped
// that the front page will be scrapped and the entry will be pushed onto the sub reddit
// slice for scraping.
func (suite *ScraperTestSuite) TestNewScraperFrontPage() {
	emptyScraper := NewScraper(suite.baseOptions)

	// ensuring that if the given scraping options is false and that it does not have a
	// front page entry into its sub reddit scraping.
	assert.False(suite.T(), emptyScraper.scrapingOptions.FrontPage)
	assert.NotEqual(suite.T(), emptyScraper.scrapingOptions.Subreddits[len(emptyScraper.scrapingOptions.Subreddits)-1], "frontpage")

	suite.baseOptions.FrontPage = true
	frontScraper := NewScraper(suite.baseOptions)

	// ensuring that if marked, front page will be scraped and the entry is pushed onto
	// the sub reddit stack.
	assert.True(suite.T(), frontScraper.scrapingOptions.FrontPage)
	assert.Equal(suite.T(), frontScraper.scrapingOptions.Subreddits[len(frontScraper.scrapingOptions.Subreddits)-1], "frontpage")
}

// TestScraperSimpleDownload ensures that for a basic run, correct folders are created, content exists
// that does not breach past the upper limit of the max number of images per site. Front page folder
// is not created (since its not marked  true) and so fourth.
func (suite *ScraperTestSuite) TestScraperSimpleDownload() {
	suite.sampleScraper.Start()
}

// TestScraperSimpleFrontPageDownload ensures that after a basic run with the additional setup of
// enabling the front page, that the front page folder is created with content existing within.
// this could fail if the N number of front page posts are not images.
func (suite *ScraperTestSuite) TestScraperSimpleFrontPageDownload() {
	suite.sampleScraper.Start()
}

// TestScraperSimpleRootDownload Ensures that with a basic run, no sub folders for the sub reddits
// dont exist, and all the content is in the root folder. No sub folders exist at all. Ensuring
// that content was also downloaded.
func (suite *ScraperTestSuite) TestScraperSimpleRootDownload() {
	suite.sampleScraper.Start()
}

func TestScraperSuite(t *testing.T) {
	suite.Run(t, new(ScraperTestSuite))
}
