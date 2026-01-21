package mappers

import (
	"database/sql"
	"path/filepath"

	"github.com/ristep/smanzy_backend/internal/db"
	"github.com/ristep/smanzy_backend/internal/models"
)

// UserRowToModel converts a database user row to a User model.
// Since sqlc generates distinct struct types for each query (even if fields are identical),
// this function uses a type switch to handle all known user-related database row types.
func UserRowToModel(row interface{}) models.User {
	// Handle different row types from sqlc
	switch r := row.(type) {
	case db.GetUserByIDRow:
		return models.User{
			ID:            uint(r.ID),
			Email:         r.Email,
			Name:          r.Name,
			Tel:           r.Tel,
			Age:           int(r.Age),
			Gender:        r.Gender,
			Address:       r.Address,
			City:          r.City,
			Country:       r.Country,
			EmailVerified: r.EmailVerified,
			CreatedAt:     r.CreatedAt,
			UpdatedAt:     r.UpdatedAt,
		}
	case db.GetUserByEmailRow:
		return models.User{
			ID:            uint(r.ID),
			Email:         r.Email,
			Name:          r.Name,
			Tel:           r.Tel,
			Age:           int(r.Age),
			Gender:        r.Gender,
			Address:       r.Address,
			City:          r.City,
			Country:       r.Country,
			EmailVerified: r.EmailVerified,
			CreatedAt:     r.CreatedAt,
			UpdatedAt:     r.UpdatedAt,
		}
	case db.ListUsersRow:
		user := models.User{
			ID:            uint(r.ID),
			Email:         r.Email,
			Name:          r.Name,
			Tel:           r.Tel,
			Age:           int(r.Age),
			Gender:        r.Gender,
			Address:       r.Address,
			City:          r.City,
			Country:       r.Country,
			EmailVerified: r.EmailVerified,
			CreatedAt:     r.CreatedAt,
			UpdatedAt:     r.UpdatedAt,
		}
		// ListUsersRow includes soft delete information which needs to be mapped optionally
		if r.DeletedAt.Valid {
			user.DeletedAt = &r.DeletedAt.Time
		}
		return user
	case db.CreateUserRow:
		return models.User{
			ID:            uint(r.ID),
			Email:         r.Email,
			Name:          r.Name,
			Tel:           r.Tel,
			Age:           int(r.Age),
			Gender:        r.Gender,
			Address:       r.Address,
			City:          r.City,
			Country:       r.Country,
			EmailVerified: r.EmailVerified,
			CreatedAt:     r.CreatedAt,
			UpdatedAt:     r.UpdatedAt,
		}
	case db.UpdateUserRow:
		return models.User{
			ID:            uint(r.ID),
			Email:         r.Email,
			Name:          r.Name,
			Tel:           r.Tel,
			Age:           int(r.Age),
			Gender:        r.Gender,
			Address:       r.Address,
			City:          r.City,
			Country:       r.Country,
			EmailVerified: r.EmailVerified,
			CreatedAt:     r.CreatedAt,
			UpdatedAt:     r.UpdatedAt,
		}
	default:
		// Return empty user if type not recognized to avoid panics
		return models.User{}
	}
}

// ListUsersRowsToModels converts multiple user rows to User models
func ListUsersRowsToModels(rows []db.ListUsersRow) []models.User {
	users := make([]models.User, len(rows))
	for i, row := range rows {
		users[i] = UserRowToModel(row)
	}
	return users
}

// MediaRowToModel converts a database media row to a Media model.
// It also handles the generation of full URLs for files and thumbnails based on the API paths.
func MediaRowToModel(row any) models.Media {
	switch r := row.(type) {
	case db.GetMediaByIDRow:
		return models.Media{
			ID:         uint(r.ID),
			Filename:   r.Filename,
			StoredName: r.StoredName,
			// Construct public URLs for the media file and its thumbnail
			URL:          "/api/media/files/" + r.StoredName,
			ThumbnailURL: "/api/media/thumbnails/320x200/" + r.StoredName[:len(r.StoredName)-len(filepath.Ext(r.StoredName))] + ".jpg",
			Type:         r.Type,
			MimeType:     r.MimeType,
			Size:         r.Size,
			UserID:       uint(r.UserID),
			CreatedAt:    r.CreatedAt,
			UpdatedAt:    r.UpdatedAt,
		}
	case db.ListPublicMediaRow:
		return models.Media{
			ID:           uint(r.ID),
			Filename:     r.Filename,
			StoredName:   r.StoredName,
			URL:          "/api/media/files/" + r.StoredName,
			ThumbnailURL: "/api/media/thumbnails/320x200/" + r.StoredName[:len(r.StoredName)-len(filepath.Ext(r.StoredName))] + ".jpg",
			Type:         r.Type,
			MimeType:     r.MimeType,
			Size:         r.Size,
			UserID:       uint(r.UserID),
			UserName:     r.UserName,
			UserTel:      r.UserTel.String,
			UserEmail:    r.UserEmail, // Add this field
			CreatedAt:    r.CreatedAt,
			UpdatedAt:    r.UpdatedAt,
		}
	case db.ListUserMediaRow:
		return models.Media{
			ID:           uint(r.ID),
			Filename:     r.Filename,
			StoredName:   r.StoredName,
			URL:          "/api/media/files/" + r.StoredName,
			ThumbnailURL: "/api/media/thumbnails/320x200/" + r.StoredName[:len(r.StoredName)-len(filepath.Ext(r.StoredName))] + ".jpg",
			Type:         r.Type,
			MimeType:     r.MimeType,
			Size:         r.Size,
			UserID:       uint(r.UserID),
			CreatedAt:    r.CreatedAt,
			UpdatedAt:    r.UpdatedAt,
		}
	case db.CreateMediaRow:
		return models.Media{
			ID:           uint(r.ID),
			Filename:     r.Filename,
			StoredName:   r.StoredName,
			URL:          "/api/media/files/" + r.StoredName,
			ThumbnailURL: "/api/media/thumbnails/320x200/" + r.StoredName[:len(r.StoredName)-len(filepath.Ext(r.StoredName))] + ".jpg",
			Type:         r.Type,
			MimeType:     r.MimeType,
			Size:         r.Size,
			UserID:       uint(r.UserID),
			CreatedAt:    r.CreatedAt,
			UpdatedAt:    r.UpdatedAt,
		}
	case db.UpdateMediaRow:
		return models.Media{
			ID:           uint(r.ID),
			Filename:     r.Filename,
			StoredName:   r.StoredName,
			URL:          "/api/media/files/" + r.StoredName,
			ThumbnailURL: "/api/media/thumbnails/320x200/" + r.StoredName[:len(r.StoredName)-len(filepath.Ext(r.StoredName))] + ".jpg",
			Type:         r.Type,
			MimeType:     r.MimeType,
			Size:         r.Size,
			UserID:       uint(r.UserID),
			CreatedAt:    r.CreatedAt,
			UpdatedAt:    r.UpdatedAt,
		}
	case db.Medium:
		return models.Media{
			ID:           uint(r.ID),
			Filename:     r.Filename,
			StoredName:   r.StoredName,
			URL:          "/api/media/files/" + r.StoredName,
			ThumbnailURL: "/api/media/thumbnails/320x200/" + r.StoredName[:len(r.StoredName)-len(filepath.Ext(r.StoredName))] + ".jpg",
			Type:         r.Type.String,
			MimeType:     r.MimeType.String,
			Size:         r.Size,
			UserID:       uint(r.UserID),
			CreatedAt:    r.CreatedAt,
			UpdatedAt:    r.UpdatedAt,
		}
	default:
		return models.Media{}
	}
}

// ListPublicMediaRowsToModels converts multiple media rows to Media models
func ListPublicMediaRowsToModels(rows []db.ListPublicMediaRow) []models.Media {
	medias := make([]models.Media, len(rows))
	for i, row := range rows {
		medias[i] = MediaRowToModel(row)
	}
	return medias
}

// AlbumRowToModel converts a database album row to an Album model
func AlbumRowToModel(row interface{}) models.Album {
	switch r := row.(type) {
	case db.Album:
		return models.Album{
			ID:          uint(r.ID),
			Title:       r.Title,
			Description: r.Description.String,
			UserID:      uint(r.UserID),
			IsPublic:    r.IsPublic.Bool,
			IsShared:    r.IsShared.Bool,
			CreatedAt:   r.CreatedAt,
			UpdatedAt:   r.UpdatedAt,
		}
	case db.ListUserAlbumsRow:
		return models.Album{
			ID:          uint(r.ID),
			Title:       r.Title,
			Description: r.Description.String,
			UserID:      uint(r.UserID),
			UserName:    r.UserName,
			IsPublic:    r.IsPublic.Bool,
			IsShared:    r.IsShared.Bool,
			CreatedAt:   r.CreatedAt,
			UpdatedAt:   r.UpdatedAt,
		}
	case db.ListAllAlbumsRow:
		return models.Album{
			ID:          uint(r.ID),
			Title:       r.Title,
			Description: r.Description.String,
			UserID:      uint(r.UserID),
			UserName:    r.UserName,
			IsPublic:    r.IsPublic.Bool,
			IsShared:    r.IsShared.Bool,
			CreatedAt:   r.CreatedAt,
			UpdatedAt:   r.UpdatedAt,
		}
	default:
		return models.Album{}
	}
}

// ListUserAlbumsRowsToModels converts multiple album rows to Album models
func ListUserAlbumsRowsToModels(rows []db.ListUserAlbumsRow) []models.Album {
	albums := make([]models.Album, len(rows))
	for i, row := range rows {
		albums[i] = AlbumRowToModel(row)
	}
	return albums
}

// ListAllAlbumsRowsToModels converts multiple album rows to Album models
func ListAllAlbumsRowsToModels(rows []db.ListAllAlbumsRow) []models.Album {
	albums := make([]models.Album, len(rows))
	for i, row := range rows {
		albums[i] = AlbumRowToModel(row)
	}
	return albums
}

// NullStringToString safely converts sql.NullString to string.
// Returns empty string if invalid.
func NullStringToString(ns sql.NullString) string {
	if ns.Valid {
		return ns.String
	}
	return ""
}

// NullInt64ToInt safely converts sql.NullInt64 to int.
// Returns 0 if invalid.
func NullInt64ToInt(ni sql.NullInt64) int {
	if ni.Valid {
		return int(ni.Int64)
	}
	return 0
}

// NullBoolToBool safely converts sql.NullBool to bool.
// Returns false if invalid.
func NullBoolToBool(nb sql.NullBool) bool {
	if nb.Valid {
		return nb.Bool
	}
	return false
}
