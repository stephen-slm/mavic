package reddit

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestRedditChildToImage takes a generic reddit response reddit child and calls into
// the underlining parse, this parse then is checked to validate that the given output
// is what is meant to be expected.
func TestRedditChildToImage(t *testing.T) {
	id := "d4zpeh"
	permalink := "/r/cute/comments/d4zpeh/arm_crawling_basic_training/"
	url := "https://i.redd.it/4kxuzo2zidn32.gif"
	title := "Army Crawling Basic Training"
	subreddit := "cute"
	domain := "i.redd.it"
	author := "unknown"
	postHint := "image"

	authorMerge := fmt.Sprintf("https://www.reddit.com/user/%s/", author)

	child := Child{Data: &ChildData{
		Title:     &title,
		Domain:    &domain,
		ID:        &id,
		Author:    &author,
		Permalink: &permalink,
		PostHint:  &postHint,
		URL:       &url,
		Subreddit: &subreddit,
	}}

	image := RedditChildToImage(child)

	assert.Equal(t, image.Id, id)
	assert.Equal(t, image.ImageId, strings.Split(strings.Split(url, "/")[3], ".")[0])
	assert.Equal(t, image.PostLink, permalink)
	assert.Equal(t, image.Link, url)
	assert.Equal(t, image.Title, title)
	assert.Equal(t, image.Subreddit, subreddit)
	assert.Equal(t, image.Source, domain)
	assert.Equal(t, image.Author.Name, author)
	assert.Equal(t, image.Author.Link, authorMerge)
}
