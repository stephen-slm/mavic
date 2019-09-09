package main

import (
	"log"
	"os"
	"strings"

	"github.com/tehstun/mavic/src/reddit"
	"gopkg.in/urfave/cli.v2"
)

var app = &cli.App{}
var options = reddit.Options{}

func setupApplicationInformation() {
	app.Name = "Mavic"
	app.Description = "Mavic is a CLI application designed to download direct images found on selected reddit subreddits."
	app.Usage = ".\\mavic.exe -l 100 --output ./pictures -f cute pics memes"
	app.Authors = []*cli.Author{{Name: "Stephen Lineker-Miller", Email: "slinekermiller@gmail.com"}}
	app.Version = "0.1.0"
}

func setupApplicationFlags() {
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:        "output",
			Aliases:     []string{"o"},
			Usage:       "The output directory to store the images.",
			Value:       "./",
			Destination: &options.OutputDirectory,
		}, &cli.IntFlag{
			Name:        "limit",
			Aliases:     []string{"l"},
			Usage:       "The total number of posts max per sub-reddit",
			Value:       50,
			Destination: &options.ImageLimit,
		},
		&cli.BoolFlag{
			Name:        "frontpage",
			Aliases:     []string{"f"},
			Usage:       "If the front page should be scrapped or not.",
			Destination: &options.FrontPage,
		},
		&cli.StringFlag{
			Name:        "type",
			Aliases:     []string{"t"},
			Usage:       "What kind of page type should reddit be during the scrapping process. e.g hot, new. top.",
			Value:       "hot",
			Destination: &options.PageType,
		},
		&cli.BoolFlag{
			Name:        "root",
			Aliases:     []string{"r"},
			Usage:       "If specified, downloads the images directly into the root, not the subreddit folder.",
			Destination: &options.RootFolderOnly,
		},
		&cli.IntFlag{
			Name:        "concurrentCount",
			Aliases:     []string{"c"},
			Usage:       "The number of images that can be downloaded at the same time.",
			Value:       25,
			Destination: &options.MaxConcurrentDownloads,
		},
	}
}

// processSubreddits takes in a ring of possible sub reddits and splits them into
// a slice of the sub reddits to be processed, there is currently a bug with the
// cli tools which is resulting in the funky processing and its best to just
// process it as a string for the time being.
func processSubreddits(arguments []string) []string {
	// since it only seems to parse the first element, even though more was selected
	// so we push it here and then go and grab the remaining.
	var processed []string

	for i := 0; i < len(arguments); i++ {
		value := arguments[i]

		// if we have hit the next command, then we must breakout since we
		// no longer have any more subs.
		if strings.HasPrefix(value, "-") {
			break
		}

		processed = append(processed, value)
	}

	return processed
}

// start is called by the cli control when the cli controls are parsed, setting up
// and building a context around the cli application. This is the time the sub
// reddits are parsed since the cli tools don't support binding stringSlices.
func start(c *cli.Context) error {
	options.Subreddits = processSubreddits(c.Args().Slice())

	// if it equals nil, and no sub reddits was given, then just set them
	// as s empty slice, letting the scraper handle the empty case as
	// it should.
	if options.Subreddits == nil {
		options.Subreddits = []string{}
	}

	// create a new reddit scraper and process through all the sub reddits
	// downloading the images in the output folder / sub reddit / image.
	scraper := reddit.NewScraper(options)
	scraper.Start()
	return nil
}

func main() {
	setupApplicationInformation()
	setupApplicationFlags()

	app.Action = start
	err := app.Run(os.Args)

	if err != nil {
		log.Fatal(err)
	}
}
