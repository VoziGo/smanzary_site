package services

import (
	"context"
	"database/sql"
	"errors"

	"github.com/ristep/smanzy_backend/internal/db"
	"github.com/ristep/smanzy_backend/internal/mappers"
	"github.com/ristep/smanzy_backend/internal/models"
)

// AlbumService handles business logic for album operations
type AlbumService struct {
	conn    *sql.DB
	queries *db.Queries
}

// NewAlbumService creates a new album service
func NewAlbumService(conn *sql.DB, queries *db.Queries) *AlbumService {
	return &AlbumService{
		conn:    conn,
		queries: queries,
	}
}

// CreateAlbum creates a new album for a user
func (as *AlbumService) CreateAlbum(ctx context.Context, userID uint, title, description, userName string) (*models.Album, error) {
	if title == "" {
		return nil, errors.New("album title is required")
	}

	albumRow, err := as.queries.CreateAlbum(ctx, db.CreateAlbumParams{
		Title:       title,
		Description: sql.NullString{String: description, Valid: true},
		UserID:      int64(userID),
	})

	if err != nil {
		return nil, err
	}

	return &models.Album{
		ID:          uint(albumRow.ID),
		Title:       albumRow.Title,
		Description: albumRow.Description.String,
		UserID:      uint(albumRow.UserID),
		UserName:    userName,
		CreatedAt:   albumRow.CreatedAt,
		UpdatedAt:   albumRow.UpdatedAt,
	}, nil
}

// GetAlbumByID retrieves an album by its ID
func (as *AlbumService) GetAlbumByID(ctx context.Context, albumID uint) (*models.Album, error) {
	albumRow, err := as.queries.GetAlbumByID(ctx, int64(albumID))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("album not found")
		}
		return nil, err
	}

	// Fetch username
	userRow, err := as.queries.GetUserByID(ctx, albumRow.UserID)
	if err != nil {
		return nil, errors.New("failed to fetch album owner information")
	}

	// Fetch media files
	mediaRows, err := as.queries.GetAlbumMedia(ctx, albumRow.ID)
	if err != nil {
		return nil, errors.New("failed to fetch album media")
	}

	album := &models.Album{
		ID:          uint(albumRow.ID),
		Title:       albumRow.Title,
		Description: albumRow.Description.String,
		UserID:      uint(albumRow.UserID),
		UserName:    userRow.Name,
		IsPublic:    albumRow.IsPublic.Bool,
		IsShared:    albumRow.IsShared.Bool,
		CreatedAt:   albumRow.CreatedAt,
		UpdatedAt:   albumRow.UpdatedAt,
	}

	// Map media files
	for _, m := range mediaRows {
		album.MediaFiles = append(album.MediaFiles, mappers.MediaRowToModel(m))
	}

	return album, nil
}

// GetUserAlbums retrieves all albums for a user
func (as *AlbumService) GetUserAlbums(ctx context.Context, userID uint) ([]models.Album, error) {
	albumRows, err := as.queries.ListUserAlbums(ctx, int64(userID))
	if err != nil {
		return nil, err
	}

	return mappers.ListUserAlbumsRowsToModels(albumRows), nil
}

// GetAllAlbums retrieves all albums from all users
func (as *AlbumService) GetAllAlbums(ctx context.Context) ([]models.Album, error) {
	albumRows, err := as.queries.ListAllAlbums(ctx)
	if err != nil {
		return nil, err
	}

	return mappers.ListAllAlbumsRowsToModels(albumRows), nil
}

// UpdateAlbum updates an album's title and description
func (as *AlbumService) UpdateAlbum(ctx context.Context, albumID uint, title, description string) (*models.Album, error) {
	albumRaw, err := as.queries.GetAlbumByID(ctx, int64(albumID))
	if err != nil {
		return nil, err
	}

	updatedRow, err := as.queries.UpdateAlbum(ctx, db.UpdateAlbumParams{
		ID: int64(albumID),
		Title: func() string {
			if title != "" {
				return title
			} else {
				return albumRaw.Title
			}
		}(),
		Description: sql.NullString{String: description, Valid: true},
		IsPublic:    albumRaw.IsPublic,
		IsShared:    albumRaw.IsShared,
	})

	if err != nil {
		return nil, err
	}

	return &models.Album{
		ID:          uint(updatedRow.ID),
		Title:       updatedRow.Title,
		Description: updatedRow.Description.String,
		UserID:      uint(updatedRow.UserID),
		CreatedAt:   updatedRow.CreatedAt,
		UpdatedAt:   updatedRow.UpdatedAt,
	}, nil
}

// AddMediaToAlbum adds a media file to an album
func (as *AlbumService) AddMediaToAlbum(ctx context.Context, albumID, mediaID uint) error {
	return as.queries.AddMediaToAlbum(ctx, db.AddMediaToAlbumParams{
		AlbumID: int64(albumID),
		MediaID: int64(mediaID),
	})
}

// RemoveMediaFromAlbum removes a media file from an album
func (as *AlbumService) RemoveMediaFromAlbum(ctx context.Context, albumID, mediaID uint) error {
	return as.queries.RemoveMediaFromAlbum(ctx, db.RemoveMediaFromAlbumParams{
		AlbumID: int64(albumID),
		MediaID: int64(mediaID),
	})
}

// DeleteAlbum performs a soft delete on an album
func (as *AlbumService) DeleteAlbum(ctx context.Context, albumID uint) error {
	return as.queries.SoftDeleteAlbum(ctx, int64(albumID))
}

// PermanentlyDeleteAlbum permanently deletes an album from the database
func (as *AlbumService) PermanentlyDeleteAlbum(ctx context.Context, albumID uint) error {
	_, err := as.conn.ExecContext(ctx, "DELETE FROM album WHERE id = $1", int64(albumID))
	return err
}
