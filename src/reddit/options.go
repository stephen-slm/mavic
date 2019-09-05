package reddit

type Options struct {
	//  The directory in which we will be downloading all the images into, based on the folder name
	//  of the given sub-reddit.
	OutputDirectory string

	// The total number of images to download max per sub-reddit before we continue to the next one.
	ImageLimit int

	// If set to true, the tool will scrape the front page of reddit for the current most active
	// sub-reddits and then scrape all the links directly from them sub-reddits.
	FrontPage bool

	// You can change this to adjust on what kind of images you get from Reddits filtering options
	// (hot, new, rising, controversial, top), hot is the default by reddit while also the default
	// in the tool.
	PageType string

	// What subreddits are going to be scrapped for downloading of sad images.
	Subreddits []string
}
