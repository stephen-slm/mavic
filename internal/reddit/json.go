package reddit

import (
	"encoding/json"
	"fmt"
	"strings"
)

func UnmarshalListing(data []byte) (Listings, error) {
	var r Listings
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *Listings) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type Listings struct {
	Data *ListingData `json:"data,omitempty"`
}

type ListingData struct {
	After    *string `json:"after,omitempty"`
	Dist     *int64  `json:"dist,omitempty"`
	Children []Child `json:"children"`
}

type Child struct {
	Data *ChildData `json:"data,omitempty"`
}

type ChildData struct {
	Title     *string `json:"title,omitempty"`
	Domain    *string `json:"domain,omitempty"`
	ID        *string `json:"id,omitempty"`
	Author    *string `json:"author,omitempty"`
	Permalink *string `json:"permalink,omitempty"`
	PostHint  *string `json:"post_hint,omitempty"`
	URL       *string `json:"url,omitempty"`
	Subreddit *string `json:"subreddit,omitempty"`
}

// RedditChildToImage takes in a single reddit listings child data object and converts it to a local
// metadata object that is used to downloadRedditMetadata and download the image.
func RedditChildToImage(child Child) Image {
	// the image id is the last section of the source url, so this requires
	// splitting on the forward slash and then taking everything after the dot
	// of the last item and then taking that last item.
	splitUrl := strings.Split(*child.Data.URL, "/")
	imageId := strings.Split(splitUrl[len(splitUrl)-1], ".")[0]

	return Image{
		Author: Author{
			Link: fmt.Sprintf("https://www.reddit.com/user/%s/", *child.Data.Author),
			Name: *child.Data.Author,
		},
		Id:        *child.Data.ID,
		ImageId:   imageId,
		PostLink:  *child.Data.Permalink,
		Link:      *child.Data.URL,
		Title:     *child.Data.Title,
		Subreddit: *child.Data.Subreddit,
		Source:    *child.Data.Domain,
	}
}
