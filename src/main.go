package main

import (
	"log"
	"os"

	"github.com/tehstun/mavic/src/reddit"
	"github.com/urfave/cli"
)

var app = cli.NewApp()
var options = reddit.Options{}

func setupApplicationInformation() {
	app.Name = "Mavic"
	app.Description = "Mavic is a CLI application designed to download direct images found on selected reddit subreddits."
	app.Usage = ".\\mavic.exe --subreddits cute -l 100 --output ./pictures -f"
	app.Author = "Stephen Lineker-Miller <slinekermiller@gmail.com>"
	app.Version = "0.0.1"
}

func setupApplicationFlags() {
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "output, o",
			Usage:       "The output directory to store the images.",
			Required:    false,
			Value:       "./",
			Destination: &options.OutputDirectory,
		}, cli.IntFlag{
			Name:        "limit, l",
			Usage:       "The total number of posts max per sub-reddit",
			Value:       50,
			Destination: &options.ImageLimit,
		},
		cli.BoolFlag{
			Name:        "frontpage, f",
			Usage:       "If the front page should be scrapped or not.",
			Destination: &options.FrontPage,
		},
		cli.StringFlag{
			Name:        "type, t",
			Usage:       "What kind of page type should reddit be during the scrapping process. e.g hot, new. top.",
			Value:       "hot",
			Destination: &options.PageType,
		},
		cli.StringSliceFlag{
			Name:     "subreddits, s",
			Usage:    "What subreddits are going to be scrapped for downloading images.",
			Required: true,
			Value:    &cli.StringSlice{},
		},
	}
}

func main() {
	setupApplicationInformation()
	setupApplicationFlags()

	app.Action = func(c *cli.Context) error {
		options.Subreddits = c.StringSlice("subreddits")
		scraper := reddit.NewScraper(options)
		scraper.ProcessSubreddits()
		return nil
	}

	err := app.Run(os.Args)

	if err != nil {
		log.Fatal(err)
	}
}
