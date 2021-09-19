package data

import (
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type Movie struct {
	gorm.Model
	Title   string         `json:"title" gorm:"index:,type:gin,expression:to_tsvector('english'\\,title)"`
	Year    int32          `json:"year"`
	Runtime Runtime        `json:"runtime" gorm:"type:int"`
	Genres  pq.StringArray `json:"genres" gorm:"type:text[];index:,type:gin"`
	Version int32          `json:"version" gorm:"default:1"`
}

type MovieModel struct {
	DB *gorm.DB
}

// Add a placeholder method for inserting a new record in the movies table.
func (m MovieModel) Insert(movie *Movie) error {
	return m.DB.Create(&movie).Error
}

// Add a placeholder method for fetching a specific record from the movies table.
func (m MovieModel) Get(id int64) (*Movie, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}
	var movie Movie
	if err := m.DB.First(&movie, id).Error; err != nil {
		return nil, err
	}
	return &movie, nil
}

// Add a placeholder method for updating a specific record in the movies table.
func (m MovieModel) Update(movie *Movie) error {
	return m.DB.Save(movie).Error
}

// Add a placeholder method for deleting a specific record from the movies table.
func (m MovieModel) Delete(id int64) error {
	return m.DB.Delete(&Movie{}, id).Error
}

func (m MovieModel) GetAll(title string, genres []string, filters Filters) ([]*Movie, Metadata, error) {
	var totalRecords int64
	movies := []*Movie{}
	query := m.DB

	if title != "" {
		query = query.Where(
			"to_tsvector('english', title) @@ plainto_tsquery('english', ?)",
			title,
		)
	}

	if len(genres) > 0 {
		query = query.Where(
			"genres::text[] @> ?",
			pq.Array(genres),
		)
	}

	if filters.Sort != "" {
		query = query.Order(filters.sortColumn() + " " + filters.sortDirection())
	}

	err := query.Limit(filters.limit()).Offset(filters.offset()).Order("id ASC").Find(&movies).Count(&totalRecords)
	if err.Error != nil {
		return nil, Metadata{}, err.Error
	}

	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)

	return movies, metadata, nil
}
