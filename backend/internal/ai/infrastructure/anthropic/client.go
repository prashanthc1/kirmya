// Package anthropic implements the AI domain's LLM port using the official
// Anthropic Go SDK (Claude). It is the primary provider; an OpenAI adapter can
// implement the same domain.LLM interface as a fallback later.
package anthropic

import (
	"context"
	"os"
	"strings"

	sdk "github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"

	"workspace-app/internal/ai/domain"
)

// model is the default per project guidance (latest, most capable Claude).
const model = sdk.ModelClaudeOpus4_8

type Client struct {
	client sdk.Client
	ready  bool
}

// New builds the client. If ANTHROPIC_API_KEY is unset, the client reports
// Ready() == false and the AI module degrades gracefully.
func New() *Client {
	key := strings.TrimSpace(os.Getenv("ANTHROPIC_API_KEY"))
	if key == "" {
		return &Client{ready: false}
	}
	return &Client{client: sdk.NewClient(option.WithAPIKey(key)), ready: true}
}

func (c *Client) Ready() bool { return c.ready }

func (c *Client) Complete(ctx context.Context, system string, messages []domain.LLMMessage, maxTokens int) (domain.Completion, error) {
	if !c.ready {
		return domain.Completion{}, domain.ErrLLMNotReady
	}

	msgs := make([]sdk.MessageParam, 0, len(messages))
	for _, m := range messages {
		block := sdk.NewTextBlock(m.Content)
		if m.Role == domain.RoleAssistant {
			msgs = append(msgs, sdk.NewAssistantMessage(block))
		} else {
			msgs = append(msgs, sdk.NewUserMessage(block))
		}
	}

	resp, err := c.client.Messages.New(ctx, sdk.MessageNewParams{
		Model:     model,
		MaxTokens: int64(maxTokens),
		System:    []sdk.TextBlockParam{{Text: system}},
		Messages:  msgs,
	})
	if err != nil {
		return domain.Completion{}, err
	}

	var sb strings.Builder
	for _, block := range resp.Content {
		if t, ok := block.AsAny().(sdk.TextBlock); ok {
			sb.WriteString(t.Text)
		}
	}

	return domain.Completion{
		Text:         sb.String(),
		Model:        model,
		InputTokens:  int(resp.Usage.InputTokens),
		OutputTokens: int(resp.Usage.OutputTokens),
	}, nil
}
