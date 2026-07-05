package api

import "workspace-app/internal/network/domain"

type sendRequest struct {
	ReceiverID string `json:"receiver_id"`
}

type connectionResponse struct {
	ID                string `json:"id"`
	RequesterID       string `json:"requester_id"`
	ReceiverID        string `json:"receiver_id"`
	Status            string `json:"status"`
	CreatedAt         string `json:"created_at"`
	UpdatedAt         string `json:"updated_at"`
	RequesterName     string `json:"requester_name,omitempty"`
	RequesterHeadline string `json:"requester_headline,omitempty"`
	RequesterPhotoURL string `json:"requester_photo_url,omitempty"`
	ReceiverName      string `json:"receiver_name,omitempty"`
	ReceiverHeadline  string `json:"receiver_headline,omitempty"`
	ReceiverPhotoURL  string `json:"receiver_photo_url,omitempty"`
}

type statusResponse struct {
	Status      string `json:"status"`
	RequesterID string `json:"requester_id,omitempty"`
}

func toResponse(c domain.Connection) connectionResponse {
	return connectionResponse{
		ID:                c.ID,
		RequesterID:       c.RequesterID,
		ReceiverID:        c.ReceiverID,
		Status:            string(c.Status),
		CreatedAt:         c.CreatedAt,
		UpdatedAt:         c.UpdatedAt,
		RequesterName:     c.RequesterName,
		RequesterHeadline: c.RequesterHeadline,
		RequesterPhotoURL: c.RequesterPhotoURL,
		ReceiverName:      c.ReceiverName,
		ReceiverHeadline:  c.ReceiverHeadline,
		ReceiverPhotoURL:  c.ReceiverPhotoURL,
	}
}
