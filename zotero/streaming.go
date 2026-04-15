package zotero

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// StreamEvent represents an event received from the Zotero streaming API.
type StreamEvent struct {
	Event   string          `json:"event"`
	Topic   string          `json:"topic,omitempty"`
	Version int             `json:"version,omitempty"`
	APIKey  string          `json:"apiKey,omitempty"`
	Retry   int             `json:"retry,omitempty"`
	Raw     json.RawMessage `json:"-"`
}

// StreamSubscription represents a subscription to create on the streaming API.
type StreamSubscription struct {
	APIKey string   `json:"apiKey,omitempty"`
	Topics []string `json:"topics,omitempty"`
}

// streamMessage is the wire format for subscription commands.
type streamMessage struct {
	Action        string                `json:"action"`
	Subscriptions []StreamSubscription  `json:"subscriptions"`
}

// StreamClient manages a WebSocket connection to the Zotero streaming API.
type StreamClient struct {
	conn    *websocket.Conn
	events  chan StreamEvent
	done    chan struct{}
	mu      sync.Mutex
	closed  bool
}

// NewStreamClient connects to the Zotero streaming API and returns a StreamClient.
// The context controls the connection lifetime.
func NewStreamClient(ctx context.Context, apiKey string) (*StreamClient, error) {
	return NewStreamClientWithURL(ctx, StreamURL, apiKey)
}

// NewStreamClientWithURL connects to the given WebSocket URL.
func NewStreamClientWithURL(ctx context.Context, wsURL string, apiKey string) (*StreamClient, error) {
	dialer := websocket.DefaultDialer
	conn, _, err := dialer.DialContext(ctx, wsURL, nil)
	if err != nil {
		return nil, fmt.Errorf("zotero: stream connect: %w", err)
	}

	sc := &StreamClient{
		conn:   conn,
		events: make(chan StreamEvent, 64),
		done:   make(chan struct{}),
	}

	// Read the initial "connected" event.
	var evt StreamEvent
	if err := conn.ReadJSON(&evt); err != nil {
		conn.Close()
		return nil, fmt.Errorf("zotero: stream handshake: %w", err)
	}
	if evt.Event != "connected" {
		conn.Close()
		return nil, fmt.Errorf("zotero: unexpected initial event: %s", evt.Event)
	}

	// Start reading events in the background.
	go sc.readLoop()

	return sc, nil
}

// Events returns the channel of streaming events.
func (sc *StreamClient) Events() <-chan StreamEvent {
	return sc.events
}

// Subscribe creates subscriptions on the streaming API.
func (sc *StreamClient) Subscribe(ctx context.Context, subs []StreamSubscription) error {
	sc.mu.Lock()
	defer sc.mu.Unlock()
	if sc.closed {
		return fmt.Errorf("zotero: stream closed")
	}
	msg := streamMessage{
		Action:        "createSubscriptions",
		Subscriptions: subs,
	}
	return sc.conn.WriteJSON(msg)
}

// Unsubscribe removes subscriptions from the streaming API.
func (sc *StreamClient) Unsubscribe(ctx context.Context, subs []StreamSubscription) error {
	sc.mu.Lock()
	defer sc.mu.Unlock()
	if sc.closed {
		return fmt.Errorf("zotero: stream closed")
	}
	msg := streamMessage{
		Action:        "deleteSubscriptions",
		Subscriptions: subs,
	}
	return sc.conn.WriteJSON(msg)
}

// Close closes the streaming connection.
func (sc *StreamClient) Close() error {
	sc.mu.Lock()
	defer sc.mu.Unlock()
	if sc.closed {
		return nil
	}
	sc.closed = true
	close(sc.done)
	err := sc.conn.WriteMessage(
		websocket.CloseMessage,
		websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""),
	)
	if err != nil {
		sc.conn.Close()
		return err
	}
	return sc.conn.Close()
}

// readLoop reads events from the WebSocket and sends them to the events channel.
func (sc *StreamClient) readLoop() {
	defer close(sc.events)
	for {
		select {
		case <-sc.done:
			return
		default:
		}
		sc.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		_, data, err := sc.conn.ReadMessage()
		if err != nil {
			if !sc.closed {
				sc.events <- StreamEvent{Event: "error", Raw: json.RawMessage(fmt.Sprintf(`{"error":%q}`, err.Error()))}
			}
			return
		}
		var evt StreamEvent
		if err := json.Unmarshal(data, &evt); err != nil {
			continue
		}
		evt.Raw = data
		sc.events <- evt
	}
}
