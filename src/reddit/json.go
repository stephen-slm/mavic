package reddit

import "encoding/json"

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
