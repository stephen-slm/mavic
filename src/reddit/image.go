package reddit

// A author of a given reddit post.
type Author struct {
	// The name of the reddit user, commonly the reddit username.
	name string
	// The link directly to the reddit users profile page /u/username
	link string
}

// Image is a given image with basic metadata about said
// image posting on reddit.
type Image struct {
	// The id of the given image post.
	id string
	// The image id, the ending part of the link.
	imageId string
	// The author of the given post.
	author Author
	// The source post link directory to reddit
	postLink string
	// The link to the source image (e.g imgur.com)
	link string
	// The title of the given post
	title string
	// The sub reddit that the image was posted too.
	subreddit string
	//  The source in which the image is hosted. e.g imgur, reddit
	source string
}
