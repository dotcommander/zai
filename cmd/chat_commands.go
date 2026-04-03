package cmd

import (
	"context"

	tea "charm.land/bubbletea/v2"

	"github.com/dotcommander/zai/internal/app"
)

// Message types for bubbletea.

type tokenMsg string

type streamDoneMsg struct {
	input string
	usage app.Usage
}

type streamErrMsg struct{ err error }

type searchResultMsg struct {
	results *app.WebSearchResponse
	query   string
	err     error
}

type webContentMsg struct {
	resp *app.WebReaderResponse
	url  string
	err  error
}

// streamChans holds the two channels needed for streaming: tokens for
// incremental display, and done for metadata after completion.
type streamChans struct {
	tokens chan string
	done   chan streamDoneMsg
}

func newStreamChans() streamChans {
	return streamChans{
		tokens: make(chan string, 64),
		done:   make(chan streamDoneMsg, 1),
	}
}

// launchStream starts the streaming API call. Tokens go to sc.tokens;
// when complete, streamDoneMsg (with input and usage) goes to sc.done.
func launchStream(ctx context.Context, client *app.Client, input string, opts app.ChatOptions, sc streamChans) tea.Cmd {
	return func() tea.Msg {
		reader, err := client.StreamChat(ctx, input, opts)
		if err != nil {
			close(sc.tokens)
			return streamErrMsg{err: err}
		}
		drainStream(ctx, reader, sc, input)
		return nil
	}
}

// drainStream pumps tokens from reader into sc.tokens in a goroutine,
// then sends the done message. Caller must not close sc.tokens.
func drainStream(ctx context.Context, reader *app.StreamReader, sc streamChans, input string) {
	go func() {
		defer close(sc.tokens)
		defer reader.Close() //nolint:errcheck // best-effort
		for {
			token, err := reader.Next()
			if err != nil {
				break
			}
			select {
			case sc.tokens <- token:
			case <-ctx.Done():
				return
			}
		}
		sc.done <- streamDoneMsg{input: input, usage: reader.StreamUsage()}
	}()
}

// waitForToken reads one token from the token channel. When the channel
// closes (stream complete), it reads the done channel to get metadata.
func waitForToken(sc streamChans) tea.Cmd {
	return func() tea.Msg {
		tok, ok := <-sc.tokens
		if !ok {
			return <-sc.done
		}
		return tokenMsg(tok)
	}
}

func launchSearchAugmentedStream(ctx context.Context, client *app.Client, input string, opts app.ChatOptions, sc streamChans) tea.Cmd {
	return func() tea.Msg {
		messageToSend := augmentWithSearch(ctx, client, input, true)

		reader, err := client.StreamChat(ctx, messageToSend, opts)
		if err != nil {
			close(sc.tokens)
			return streamErrMsg{err: err}
		}
		drainStream(ctx, reader, sc, messageToSend)
		return nil
	}
}

// doSearch performs a web search asynchronously.
func doSearch(ctx context.Context, client *app.Client, query string, opts app.SearchOptions) tea.Cmd {
	return func() tea.Msg {
		results, err := client.SearchWeb(ctx, query, opts)
		return searchResultMsg{results: results, query: query, err: err}
	}
}

// doWebFetch fetches web content asynchronously.
func doWebFetch(ctx context.Context, client *app.Client, url string) tea.Cmd {
	return func() tea.Msg {
		webOpts := &app.WebReaderOptions{ReturnFormat: "markdown"}
		resp, err := client.FetchWebContent(ctx, url, webOpts)
		return webContentMsg{resp: resp, url: url, err: err}
	}
}
