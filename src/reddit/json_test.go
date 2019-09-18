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

	image := redditChildToImage(child)

	assert.Equal(t, image.id, id)
	assert.Equal(t, image.imageId, strings.Split(strings.Split(url, "/")[3], ".")[0])
	assert.Equal(t, image.postLink, permalink)
	assert.Equal(t, image.link, url)
	assert.Equal(t, image.title, title)
	assert.Equal(t, image.subreddit, subreddit)
	assert.Equal(t, image.source, domain)
	assert.Equal(t, image.author.name, author)
	assert.Equal(t, image.author.link, authorMerge)
}
