package cmd

import (
	"context"
	"errors"
	"io"

	tea "charm.land/bubbletea/v2"

	"github.com/dotcommander/zai/internal/app"
)

// Message types for bubbletea.

type tokenMsg string
type turnTokenMsg struct {
	turn  uint64
	token string
}

type streamDoneMsg struct {
	turn  uint64
	model string
	input string
	usage app.Usage
}

type streamErrMsg struct {
	turn uint64
	err  error
}

type searchResultMsg struct {
	turn    uint64
	results *app.WebSearchResponse
	query   string
	err     error
}

type webContentMsg struct {
	turn uint64
	resp *app.WebReaderResponse
	url  string
	err  error
}

// streamChans holds the two channels needed for streaming: tokens for
// incremental display, and done for metadata after completion.
type streamChans struct {
	turn   uint64
	tokens chan string
	done   chan tea.Msg
}

func newStreamChans() streamChans {
	return streamChans{
		tokens: make(chan string, 64),
		done:   make(chan tea.Msg, 1),
	}
}

// launchStream starts the streaming API call. Tokens go to sc.tokens;
// when complete, streamDoneMsg (with input and usage) goes to sc.done.
func launchStream(ctx context.Context, client *app.Client, input string, opts app.ChatOptions, sc streamChans) tea.Cmd {
	return func() tea.Msg {
		reader, err := client.StreamChat(ctx, input, opts)
		if err != nil {
			sc.done <- streamErrMsg{turn: sc.turn, err: err}
			close(sc.tokens)
			return nil
		}
		drainStream(ctx, reader, sc, input)
		return nil
	}
}

// drainStream pumps tokens from reader into sc.tokens in a goroutine,
// then sends the done message. Caller must not close sc.tokens.
func drainStream(ctx context.Context, reader *app.StreamReader, sc streamChans, input string) {
	go func() {
		// The request context or reader completion terminates this producer.
		terminal := tea.Msg(streamDoneMsg{turn: sc.turn, model: reader.Model(), input: input})
		defer func() {
			if ctx.Err() != nil {
				terminal = streamErrMsg{turn: sc.turn, err: ctx.Err()}
			}
			sc.done <- terminal
			close(sc.tokens)
		}()
		defer reader.Close() //nolint:errcheck // best-effort
		for {
			token, err := reader.Next()
			if err != nil {
				if errors.Is(err, io.EOF) {
					terminal = streamDoneMsg{turn: sc.turn, model: reader.Model(), input: input, usage: reader.StreamUsage()}
				} else {
					terminal = streamErrMsg{turn: sc.turn, err: err}
				}
				return
			}
			select {
			case sc.tokens <- token:
			case <-ctx.Done():
				return
			}
		}
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
		if sc.turn != 0 {
			return turnTokenMsg{turn: sc.turn, token: tok}
		}
		return tokenMsg(tok)
	}
}

func launchSearchAugmentedStream(ctx context.Context, client *app.Client, input string, opts app.ChatOptions, sc streamChans) tea.Cmd {
	return func() tea.Msg {
		messageToSend := augmentWithSearch(ctx, client, input, true)

		reader, err := client.StreamChat(ctx, messageToSend, opts)
		if err != nil {
			sc.done <- streamErrMsg{turn: sc.turn, err: err}
			close(sc.tokens)
			return nil
		}
		drainStream(ctx, reader, sc, messageToSend)
		return nil
	}
}

// doSearch performs a web search asynchronously.
func doSearch(ctx context.Context, client *app.Client, query string, opts app.SearchOptions, turns ...uint64) tea.Cmd {
	return func() tea.Msg {
		results, err := client.SearchWeb(ctx, query, opts)
		var turn uint64
		if len(turns) > 0 {
			turn = turns[0]
		}
		return searchResultMsg{turn: turn, results: results, query: query, err: err}
	}
}

// doWebFetch fetches web content asynchronously.
func doWebFetch(ctx context.Context, client *app.Client, url string, turns ...uint64) tea.Cmd {
	return func() tea.Msg {
		webOpts := &app.WebReaderOptions{ReturnFormat: "markdown"}
		resp, err := client.FetchWebContent(ctx, url, webOpts)
		var turn uint64
		if len(turns) > 0 {
			turn = turns[0]
		}
		return webContentMsg{turn: turn, resp: resp, url: url, err: err}
	}
}
