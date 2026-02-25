package gmail

import (
	"context"
	"fmt"

	"github.com/anderskaae/spam-time-waster/internal/auth"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

// Service wraps the gmail.Service and provides methods for interacting with Gmail.
type Service struct {
	*gmail.Service
}

// NewService creates a new Gmail service using the provided credentials.
func NewService(ctx context.Context, clientID, clientSecret, tokenFile string) (*Service, error) {
	client, err := auth.GetClient(ctx, clientID, clientSecret, tokenFile)
	if err != nil {
		return nil, fmt.Errorf("failed to get oauth client: %w", err)
	}

	srv, err := gmail.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve Gmail client: %v", err)
	}

	return &Service{srv}, nil
}

// ListMessages returns a list of messages.
func (s *Service) ListMessages(userId string) ([]*gmail.Message, error) {
	r, err := s.Users.Messages.List(userId).MaxResults(10).Do()
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve messages: %v", err)
	}
	return r.Messages, nil
}

// GetLabelByName searches for a label by name and returns it.
func (s *Service) GetLabelByName(userId, name string) (*gmail.Label, error) {
	res, err := s.Users.Labels.List(userId).Do()
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve labels: %v", err)
	}

	for _, label := range res.Labels {
		if label.Name == name {
			return label, nil
		}
	}
	return nil, nil // Not found
}
