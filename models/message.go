package models

import "time"

// Role represents the role of a message sender in a conversation.
type Role string

const (
	// RoleSystem represents the system role, used for initial instructions.
	RoleSystem Role = "system"
	// RoleUser represents the user role.
	RoleUser Role = "user"
	// RoleAssistant represents the AI assistant role.
	RoleAssistant Role = "assistant"
)

// Message represents a single message in a conversation.
type Message struct {
	ID        string    `json:"id,omitempty"`
	Role      Role      `json:"role"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at,omitempty"`
}

// Conversation holds a sequence of messages forming a dialogue.
type Conversation struct {
	ID        string    `json:"id"`
	Messages  []Message `json:"messages"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ChatRequest is the payload sent to the AI provider.
type ChatRequest struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	MaxTokens   int       `json:"max_tokens,omitempty"`
	Temperature float64   `json:"temperature,omitempty"`
	Stream      bool      `json:"stream,omitempty"`
}

// ChatResponse is the response received from the AI provider.
type ChatResponse struct {
	ID      string   `json:"id"`
	Object  string   `json:"object"`
	Created int64    `json:"created"`
	Model   string   `json:"model"`
	Choices []Choice `json:"choices"`
	Usage   Usage    `json:"usage"`
}

// Choice represents a single completion choice returned by the AI.
type Choice struct {
	Index        int     `json:"index"`
	Message      Message `json:"message"`
	FinishReason string  `json:"finish_reason"`
}

// Usage contains token usage statistics for a request.
type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

// NewConversation creates a new Conversation with a system prompt.
func NewConversation(id, systemPrompt string) *Conversation {
	now := time.Now()
	c := &Conversation{
		ID:        id,
		Messages:  []Message{},
		CreatedAt: now,
		UpdatedAt: now,
	}
	if systemPrompt != "" {
		c.Messages = append(c.Messages, Message{
			Role:      RoleSystem,
			Content:   systemPrompt,
			CreatedAt: now,
		})
	}
	return c
}

// AddMessage appends a new message to the conversation and updates the timestamp.
// Both the message timestamp and the conversation UpdatedAt are set to the same
// instant so they stay consistent with each other.
func (c *Conversation) AddMessage(role Role, content string) {
	now := time.Now()
	c.Messages = append(c.Messages, Message{
		Role:      role,
		Content:   content,
		CreatedAt: now,
	})
	c.UpdatedAt = now
}

// LastAssistantMessage returns the most recent assistant message, or nil if none exists.
func (c *Conversation) LastAssistantMessage() *Message {
	for i := len(c.Messages) - 1; i >= 0; i-- {
		if c.Messages[i].Role == RoleAssistant {
			return &c.Messages[i]
		}
	}
	return nil
}
