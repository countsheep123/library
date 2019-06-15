package obj

type Book struct {
	ID        string   `json:"id,omitempty"`
	CreatedAt string   `json:"created_at,omitempty"`
	UpdatedAt string   `json:"updated_at,omitempty"`
	Title     *string  `json:"title,omitempty"`
	Publisher *string  `json:"publisher,omitempty"`
	Pubdate   *string  `json:"pubdate,omitempty"`
	Authors   []string `json:"authors,omitempty"`
	ISBN      *string  `json:"isbn,omitempty"`
	Cover     *string  `json:"cover,omitempty"`

	Labels         []string `json:"labels,omitempty"`
	RecommenderIDs []string `json:"recommender_ids,omitempty"`
}

func (b *Book) Validate() error {

	if b.Title == nil {
		return InvalidRequest{
			Msg: "title is required",
		}
	}
	if b.Publisher == nil {
		return InvalidRequest{
			Msg: "publisher is required",
		}
	}
	if b.Pubdate == nil {
		return InvalidRequest{
			Msg: "pubdate is required",
		}
	}
	if b.Authors == nil {
		return InvalidRequest{
			Msg: "authors is required",
		}
	}
	if b.ISBN == nil {
		return InvalidRequest{
			Msg: "isbn is required",
		}
	}
	if b.Cover == nil {
		return InvalidRequest{
			Msg: "cover is required",
		}
	}
	if b.Labels == nil {
		return InvalidRequest{
			Msg: "labels is required",
		}
	}
	// RecommenderIDs is not required

	return nil
}
