package app

import (
	"fmt"
	"time"
)

// APIError represents an error response from the Z.AI API.
// Use errors.As to extract this type from wrapped errors.
type APIError struct {
	StatusCode int
	Body       string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("API error: %d - %s", e.StatusCode, e.Body)
}

// ChatRequest represents the API request payload.
type ChatRequest struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	Stream      bool      `json:"stream"`
	Temperature float64   `json:"temperature,omitempty"`
	MaxTokens   int       `json:"max_tokens,omitempty"`
	TopP        float64   `json:"top_p,omitempty"`
	Thinking    *Thinking `json:"thinking,omitempty"`
}

// Thinking configures the thinking/reasoning mode.
// Type: "enabled" or "disabled"
type Thinking struct {
	Type string `json:"type"`
}

// Message represents a chat message.
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ChatResponse represents the API response.
type ChatResponse struct {
	ID      string   `json:"id"`
	Object  string   `json:"object"`
	Created int64    `json:"created"`
	Model   string   `json:"model"`
	Choices []Choice `json:"choices"`
	Usage   Usage    `json:"usage"`
}

// Choice represents a response choice.
type Choice struct {
	Index        int     `json:"index"`
	Message      Message `json:"message"`
	FinishReason string  `json:"finish_reason"`
}

// Usage represents token usage statistics.
type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

// ModelsResponse represents the /models API response.
type ModelsResponse struct {
	Object string  `json:"object"`
	Data   []Model `json:"data"`
}

// Model represents a single model in the list.
type Model struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	OwnedBy string `json:"owned_by"`
}

// ImageGenerationRequest represents the image generation API request.
type ImageGenerationRequest struct {
	Model   string `json:"model"` // "glm-image"
	Prompt  string `json:"prompt"`
	Quality string `json:"quality,omitempty"` // "hd" or "standard"
	Size    string `json:"size,omitempty"`    // "1024x1024"
	UserID  string `json:"user_id,omitempty"` // Optional
}

// ImageResponse represents the image generation API response.
type ImageResponse struct {
	ID      string      `json:"id"`
	Object  string      `json:"object"`
	Created int64       `json:"created"`
	Model   string      `json:"model"`
	Data    []ImageData `json:"data"`
	Usage   Usage       `json:"usage,omitempty"`
}

// ImageData represents a generated image.
type ImageData struct {
	URL           string `json:"url"`
	RevisedPrompt string `json:"revised_prompt,omitempty"`
	Width         int    `json:"width"`
	Height        int    `json:"height"`
	Format        string `json:"format"`
}

// ImageModel represents an image generation model.
type ImageModel struct {
	ID          string `json:"id"`
	Object      string `json:"object"`
	Created     int64  `json:"created"`
	OwnedBy     string `json:"owned_by"`
	Description string `json:"description,omitempty"`
}

// ImageOptions configures image generation requests.
type ImageOptions struct {
	Quality string // "hd" or "standard"
	Size    string // "widthxheight" format
	UserID  string // Optional user ID for analytics
	Model   string // Override default model
}

// WebReaderRequest represents a web reader API request.
type WebReaderRequest struct {
	URL               string `json:"url"`
	Timeout           *int   `json:"timeout,omitempty"`             // default: 20
	NoCache           *bool  `json:"no_cache,omitempty"`            // default: false
	ReturnFormat      string `json:"return_format,omitempty"`       // default: "markdown"
	RetainImages      *bool  `json:"retain_images,omitempty"`       // default: true
	NoGFM             *bool  `json:"no_gfm,omitempty"`              // default: false
	KeepImgDataURL    *bool  `json:"keep_img_data_url,omitempty"`   // default: false
	WithImagesSummary *bool  `json:"with_images_summary,omitempty"` // default: false
	WithLinksSummary  *bool  `json:"with_links_summary,omitempty"`  // default: false
}

// WebReaderResponse represents a web reader API response.
type WebReaderResponse struct {
	ID           string       `json:"id"`
	Created      int64        `json:"created"`
	ReaderResult ReaderResult `json:"reader_result"`
}

// ReaderResult contains the web reader results.
type ReaderResult struct {
	Content           string                 `json:"content"`
	Description       string                 `json:"description"`
	Title             string                 `json:"title"`
	URL               string                 `json:"url"`
	ExternalResources map[string]interface{} `json:"external_resources,omitempty"`
	Metadata          map[string]interface{} `json:"metadata,omitempty"`
}

// WebReaderOptions configures web reader requests.
type WebReaderOptions struct {
	Timeout           *int
	NoCache           *bool
	ReturnFormat      string // "markdown" or "text"
	RetainImages      *bool
	NoGFM             *bool
	KeepImgDataURL    *bool
	WithImagesSummary *bool
	WithLinksSummary  *bool
}

// ChatOptions configures chat requests.
type ChatOptions struct {
	Model       string   // Override default model
	Temperature *float64 // Override default temperature
	MaxTokens   *int     // Override default max tokens
	TopP        *float64 // Override default top_p
	Thinking    *bool    // Enable thinking mode
	WebEnabled  *bool    // Enable web content fetching
	WebTimeout  *int     // Web fetch timeout in seconds

	// Legacy fields for backward compatibility
	FilePath     string    // Optional file to include in context
	Context      []Message // Previous messages for context
	Think        bool      // Enable thinking/reasoning mode (legacy)
	SystemPrompt string    // Custom system prompt
}

// WebSearchRequest represents a web search API request.
type WebSearchRequest struct {
	SearchEngine        string  `json:"search_engine"` // "search_std", "search_pro", "search_pro_sogou", "search_pro_quark"
	SearchQuery         string  `json:"search_query"`
	Count               *int    `json:"count,omitempty"` // 1-50, default 10
	SearchDomainFilter  *string `json:"search_domain_filter,omitempty"`
	SearchRecencyFilter *string `json:"search_recency_filter,omitempty"` // oneDay/oneWeek/oneMonth/oneYear/noLimit
	RequestID           *string `json:"request_id,omitempty"`
	UserID              *string `json:"user_id,omitempty"`
	ContentSize         string  `json:"content_size,omitempty"`    // "medium" or "high"
	SearchResult        *bool   `json:"search_result,omitempty"`   // return detailed sources
	RequireSearch       *bool   `json:"require_search,omitempty"`  // force search-based response
	SearchPrompt        string  `json:"search_prompt,omitempty"`   // custom prompt for processing search results
	ResultSequence      string  `json:"result_sequence,omitempty"` // "before" or "after"
}

// SearchResult represents a single search result.
type SearchResult struct {
	Title       string `json:"title"`
	Content     string `json:"content"`
	Link        string `json:"link"`
	Media       string `json:"media"`
	Icon        string `json:"icon"`
	Refer       string `json:"refer"`
	PublishDate string `json:"publish_date"`
}

// WebSearchResponse represents a web search API response.
type WebSearchResponse struct {
	ID           string         `json:"id"`
	Created      int64          `json:"created"`
	SearchResult []SearchResult `json:"search_result"`
}

// SearchOptions configures search requests.
type SearchOptions struct {
	Count          int
	DomainFilter   string
	RecencyFilter  string
	Language       string
	RequestID      string
	UserID         string
	SearchEngine   string // "search_std", "search_pro", "search_pro_sogou", "search_pro_quark"
	ContentSize    string // "medium" or "high"
	SearchResult   bool   // return detailed sources
	RequireSearch  bool   // force search-based response
	SearchPrompt   string // custom prompt for processing search results
	ResultSequence string // "before" or "after"
}

// SearchOutputFormat represents the output format for search results.
type SearchOutputFormat string

const (
	// SearchOutputTable displays results in a formatted table.
	SearchOutputTable SearchOutputFormat = "table"
	// SearchOutputDetailed displays results with full details.
	SearchOutputDetailed SearchOutputFormat = "detailed"
	// SearchOutputJSON outputs results as JSON.
	SearchOutputJSON SearchOutputFormat = "json"
)

// SearchCacheEntry represents a cached search result.
type SearchCacheEntry struct {
	Query     string         `json:"query"`
	Results   []SearchResult `json:"results"`
	CachedAt  time.Time      `json:"cached_at"`
	ExpiresAt time.Time      `json:"expires_at"`
	Hash      string         `json:"hash"` // SHA256 of query + options
}

// RetryConfig configures retry behavior for transient failures.
type RetryConfig struct {
	MaxAttempts    int           // Maximum number of retry attempts (default: 3)
	InitialBackoff time.Duration // Initial backoff duration (default: 1s)
	MaxBackoff     time.Duration // Maximum backoff duration (default: 30s)
}

// VisionRequest represents a vision/image analysis API request.
type VisionRequest struct {
	Model       string          `json:"model"`
	Messages    []VisionMessage `json:"messages"`
	Stream      bool            `json:"stream"`
	Temperature float64         `json:"temperature,omitempty"`
	MaxTokens   int             `json:"max_tokens,omitempty"`
	TopP        float64         `json:"top_p,omitempty"`
}

// VisionMessage represents a message in vision API (supports multimodal content).
type VisionMessage struct {
	Role    string        `json:"role"`
	Content []ContentPart `json:"content"`
}

// ContentPart represents a part of multimodal content (text or image).
type ContentPart struct {
	Type     string           `json:"type"` // "text" or "image_url"
	Text     string           `json:"text,omitempty"`
	ImageURL *ImageURLContent `json:"image_url,omitempty"`
}

// ImageURLContent contains image URL or base64 data.
type ImageURLContent struct {
	URL string `json:"url"`
}

// VisionOptions configures vision/analysis requests.
type VisionOptions struct {
	Model       string   // Override default model (e.g., "glm-4.6v")
	Temperature *float64 // Override default temperature
	MaxTokens   *int     // Override default max tokens
	TopP        *float64 // Override default top_p
}

// TranscriptionResponse represents the audio transcription API response.
type TranscriptionResponse struct {
	ID        string `json:"id"`
	Created   int64  `json:"created"`
	RequestID string `json:"request_id,omitempty"`
	Model     string `json:"model"`
	Text      string `json:"text"`
}

// TranscriptionOptions configures audio transcription requests.
type TranscriptionOptions struct {
	Model     string   // Override default model (default: glm-asr-2512)
	Prompt    string   // Context from prior transcriptions (max 8000 chars)
	Hotwords  []string // Domain vocabulary (max 100 items)
	Stream    bool     // Enable streaming via Event Stream
	UserID    string   // End user ID (6-128 characters)
	RequestID string   // Client-provided unique identifier
}

// StreamChunk represents a single SSE chunk from the streaming API.
type StreamChunk struct {
	ID      string         `json:"id"`
	Object  string         `json:"object"`
	Created int64          `json:"created"`
	Model   string         `json:"model"`
	Choices []StreamChoice `json:"choices"`
	Usage   *Usage         `json:"usage,omitempty"`
}

// StreamChoice represents a choice in a streaming response chunk.
type StreamChoice struct {
	Index        int         `json:"index"`
	Delta        StreamDelta `json:"delta"`
	FinishReason *string     `json:"finish_reason"`
}

// StreamDelta represents the incremental content in a streaming chunk.
type StreamDelta struct {
	Role    string `json:"role,omitempty"`
	Content string `json:"content,omitempty"`
}

// VideoGenerationRequest represents the video generation API request.
type VideoGenerationRequest struct {
	Model     string   `json:"model"`                // "cogvideox-3"
	Prompt    string   `json:"prompt,omitempty"`     // Text description (max 512 chars)
	ImageURL  []string `json:"image_url,omitempty"`  // URL or base64, 1-2 images (first/last frame)
	Quality   string   `json:"quality,omitempty"`    // "quality" or "speed" (default: speed)
	WithAudio bool     `json:"with_audio,omitempty"` // AI sound effects (default: false)
	Size      string   `json:"size,omitempty"`       // 1280x720, 1920x1080, 3840x2160, etc.
	FPS       int      `json:"fps,omitempty"`        // 30 or 60 (default: 30)
	Duration  int      `json:"duration,omitempty"`   // 5 or 10 seconds (default: 5)
	RequestID string   `json:"request_id,omitempty"`
	UserID    string   `json:"user_id,omitempty"`
}

// VideoGenerationResponse represents the async video generation API response.
type VideoGenerationResponse struct {
	ID         string `json:"id"` // Task ID for polling
	Model      string `json:"model"`
	RequestID  string `json:"request_id"`
	TaskStatus string `json:"task_status"` // PROCESSING, SUCCESS, FAIL
}

// VideoResultResponse represents the result retrieval API response.
type VideoResultResponse struct {
	Model       string        `json:"model"`
	VideoResult []VideoResult `json:"video_result"`
	TaskStatus  string        `json:"task_status"` // PROCESSING, SUCCESS, FAIL
	RequestID   string        `json:"request_id"`
}

// VideoResult represents a generated video.
type VideoResult struct {
	URL           string `json:"url"`             // Video URL
	CoverImageURL string `json:"cover_image_url"` // Thumbnail URL
}

// VideoOptions configures video generation requests.
type VideoOptions struct {
	Model     string   // Override default model
	Quality   string   // "quality" or "speed"
	Size      string   // Resolution
	FPS       int      // Frame rate
	Duration  int      // Duration in seconds
	WithAudio bool     // Include AI sound effects
	ImageURLs []string // First/last frame images
	UserID    string   // User ID for analytics
	RequestID string   // Unique request ID
}

// TTSSpeechRequest represents a text-to-speech API request.
type TTSSpeechRequest struct {
	Model          string `json:"model"`
	Input          string `json:"input"`
	Voice          string `json:"voice"`
	Speed          *int   `json:"speed,omitempty"`
	Volume         *int   `json:"volume,omitempty"`
	ResponseFormat string `json:"response_format,omitempty"`
}

// TTSOptions configures text-to-speech requests.
type TTSOptions struct {
	Model          string
	Voice          string
	Speed          *int
	Volume         *int
	ResponseFormat string
}

// EmbeddingRequest represents an embedding API request.
type EmbeddingRequest struct {
	Model string   `json:"model"`
	Input []string `json:"input"`
}

// EmbeddingData represents a single embedding result.
type EmbeddingData struct {
	Object    string    `json:"object"`
	Embedding []float64 `json:"embedding"`
	Index     int       `json:"index"`
}

// EmbeddingResponse represents the embedding API response.
type EmbeddingResponse struct {
	Object string          `json:"object"`
	Data   []EmbeddingData `json:"data"`
	Model  string          `json:"model"`
	Usage  EmbeddingUsage  `json:"usage"`
}

// EmbeddingUsage represents token usage for embeddings.
type EmbeddingUsage struct {
	PromptTokens int `json:"prompt_tokens"`
	TotalTokens  int `json:"total_tokens"`
}

// EmbeddingOptions configures embedding requests.
type EmbeddingOptions struct {
	Model string
}
