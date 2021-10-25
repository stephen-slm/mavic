package reddit

// A author of a given reddit post.
type Author struct {
	// The name of the reddit user, commonly the reddit username.
	Name string
	// The link directly to the reddit users profile page /u/username
	Link string
}

// Image is a given image with basic metadata about said
// image posting on reddit.
type Image struct {
	// The id of the given image post.
	Id string
	// The image id, the ending part of the link.
	ImageId string
	// The author of the given post.
	Author Author
	// The source post link directory to reddit
	PostLink string
	// The link to the source image (e.g imgur.com)
	Link string
	// The title of the given post
	Title string
	// The sub reddit that the image was posted too.
	Subreddit string
	//  The source in which the image is hosted. e.g imgur, reddit
	Source string
}
