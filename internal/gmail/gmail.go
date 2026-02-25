package gmail

import (
	"context"
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/anderskaae/spam-time-waster/internal/auth"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

// Message represents the data needed to eventually answer an email.
type Message struct {
	ID       string
	ThreadID string
	From     string
	Subject  string
	Body     string
}

// Service wraps the gmail.Service and provides methods for interacting with Gmail.
type Service struct {
	*gmail.Service
	UserEmail string
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

	profile, err := srv.Users.GetProfile("me").Do()
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve user profile: %v", err)
	}

	return &Service{
		Service:   srv,
		UserEmail: profile.EmailAddress,
	}, nil
}

// GetMessagesByLabel retrieves messages with the specified label ID and returns them as Message objects.
// It only returns threads where the LATEST message is NOT from the authenticated user.
func (s *Service) GetMessagesByLabel(userId, labelId string) ([]*Message, error) {
	r, err := s.Users.Messages.List(userId).LabelIds(labelId).MaxResults(10).Do()
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve messages: %v", err)
	}

	var messages []*Message
	processedThreads := make(map[string]bool)

	for _, m := range r.Messages {
		if processedThreads[m.ThreadId] {
			continue
		}

		thread, err := s.Users.Threads.Get(userId, m.ThreadId).Do()
		if err != nil {
			return nil, fmt.Errorf("unable to retrieve thread %s: %v", m.ThreadId, err)
		}

		if len(thread.Messages) == 0 {
			continue
		}

		// Get the latest message in the thread
		latestMsg := thread.Messages[len(thread.Messages)-1]
		
		var fromEmail string
		for _, header := range latestMsg.Payload.Headers {
			if header.Name == "From" {
				fromEmail = header.Value
				break
			}
		}

		// If the latest message is from the user, skip this thread
		if strings.Contains(strings.ToLower(fromEmail), strings.ToLower(s.UserEmail)) {
			fmt.Printf("Debug: Skipping thread %s because latest sender (%s) matches user (%s)\n", m.ThreadId, fromEmail, s.UserEmail)
			processedThreads[m.ThreadId] = true
			continue
		}

		fmt.Printf("Debug: Including thread %s because latest sender (%s) is NOT user (%s)\n", m.ThreadId, fromEmail, s.UserEmail)

		// Map to our internal Message object
		msg := &Message{
			ID:       latestMsg.Id,
			ThreadID: latestMsg.ThreadId,
		}

		for _, header := range latestMsg.Payload.Headers {
			if header.Name == "From" {
				msg.From = header.Value
			}
			if header.Name == "Subject" {
				msg.Subject = header.Value
			}
		}

		msg.Body = getBody(latestMsg.Payload)
		messages = append(messages, msg)
		processedThreads[m.ThreadId] = true
	}

	return messages, nil
}

// getBody extracts the body from the message payload.
func getBody(payload *gmail.MessagePart) string {
	if payload.Body.Data != "" {
		data, err := base64.URLEncoding.DecodeString(payload.Body.Data)
		if err != nil {
			return ""
		}
		return string(data)
	}

	for _, part := range payload.Parts {
		body := getBody(part)
		if body != "" {
			return body
		}
	}

	return ""
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
