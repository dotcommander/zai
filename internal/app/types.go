package app

import "time"

// ChatRequest represents the API request payload.
type ChatRequest struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	Stream      bool      `json:"stream"` // Reserved for future streaming API support
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
	Model   string `json:"model"`            // "cogview-4-250304"
	Prompt  string `json:"prompt"`
	Quality string `json:"quality,omitempty"` // "hd" or "standard"
	Size    string `json:"size,omitempty"`    // "1024x1024"
	UserID  string `json:"user_id,omitempty"` // Optional
}

// ImageResponse represents the image generation API response.
type ImageResponse struct {
	ID      string       `json:"id"`
	Object  string       `json:"object"`
	Created int64        `json:"created"`
	Model   string       `json:"model"`
	Data    []ImageData  `json:"data"`
	Usage   Usage        `json:"usage,omitempty"`
}

// ImageData represents a generated image.
type ImageData struct {
	URL            string `json:"url"`
	RevisedPrompt  string `json:"revised_prompt,omitempty"`
	Width          int    `json:"width"`
	Height         int    `json:"height"`
	Format         string `json:"format"`
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
	Timeout           *int   `json:"timeout,omitempty"`           // default: 20
	NoCache           *bool  `json:"no_cache,omitempty"`         // default: false
	ReturnFormat      string `json:"return_format,omitempty"`    // default: "markdown"
	RetainImages      *bool  `json:"retain_images,omitempty"`    // default: true
	NoGFM             *bool  `json:"no_gfm,omitempty"`           // default: false
	KeepImgDataURL    *bool  `json:"keep_img_data_url,omitempty"` // default: false
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
	Content          string                 `json:"content"`
	Description      string                 `json:"description"`
	Title            string                 `json:"title"`
	URL              string                 `json:"url"`
	ExternalResources map[string]interface{} `json:"external_resources,omitempty"`
	Metadata         map[string]interface{} `json:"metadata,omitempty"`
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
	Model       string    // Override default model
	Temperature *float64  // Override default temperature
	MaxTokens   *int      // Override default max tokens
	TopP        *float64  // Override default top_p
	Thinking    *bool     // Enable thinking mode
	WebEnabled  *bool     // Enable web content fetching
	WebTimeout  *int      // Web fetch timeout in seconds

	// Legacy fields for backward compatibility
	FilePath     string    // Optional file to include in context
	Context      []Message // Previous messages for context
	Think        bool      // Enable thinking/reasoning mode (legacy)
}

// WebSearchRequest represents a web search API request.
type WebSearchRequest struct {
	SearchEngine          string  `json:"search_engine"`           // "search-prime"
	SearchQuery           string  `json:"search_query"`
	Count                 *int    `json:"count,omitempty"`         // 1-50, default 10
	SearchDomainFilter    *string `json:"search_domain_filter,omitempty"`
	SearchRecencyFilter   *string `json:"search_recency_filter,omitempty"` // oneDay/oneWeek/oneMonth/oneYear/noLimit
	RequestID             *string `json:"request_id,omitempty"`
	UserID                *string `json:"user_id,omitempty"`
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
	ID           string        `json:"id"`
	Created      int64         `json:"created"`
	SearchResult []SearchResult `json:"search_result"`
}

// SearchOptions configures search requests.
type SearchOptions struct {
	Count         int    // Number of results (1-50)
	DomainFilter  string // Limit to specific domain
	RecencyFilter string // Time filter: oneDay, oneWeek, oneMonth, oneYear, noLimit
	RequestID     string // Unique request ID
	UserID        string // User ID for analytics
}

// SearchOutputFormat represents the output format for search results.
type SearchOutputFormat string

const (
	SearchOutputTable    SearchOutputFormat = "table"
	SearchOutputDetailed SearchOutputFormat = "detailed"
	SearchOutputJSON     SearchOutputFormat = "json"
)

// SearchCacheEntry represents a cached search result.
type SearchCacheEntry struct {
	Query      string        `json:"query"`
	Results    []SearchResult `json:"results"`
	CachedAt   time.Time     `json:"cached_at"`
	ExpiresAt  time.Time     `json:"expires_at"`
	Hash       string        `json:"hash"` // SHA256 of query + options
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
	Stream      bool            `json:"stream"` // Reserved for future streaming API support
	Temperature float64         `json:"temperature,omitempty"`
	MaxTokens   int             `json:"max_tokens,omitempty"`
	TopP        float64         `json:"top_p,omitempty"`
}

// VisionMessage represents a message in vision API (supports multimodal content).
type VisionMessage struct {
	Role    string         `json:"role"`
	Content []ContentPart  `json:"content"`
}

// ContentPart represents a part of multimodal content (text or image).
type ContentPart struct {
	Type      string           `json:"type"` // "text" or "image_url"
	Text      string           `json:"text,omitempty"`
	ImageURL  *ImageURLContent `json:"image_url,omitempty"`
}

// ImageURLContent contains image URL or base64 data.
type ImageURLContent struct {
	URL string `json:"url"`
}

// VisionOptions configures vision/analysis requests.
type VisionOptions struct {
	Model       string  // Override default model (e.g., "glm-4.6v")
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
	Model    string   // Override default model (default: glm-asr-2512)
	Prompt   string   // Context from prior transcriptions (max 8000 chars)
	Hotwords []string // Domain vocabulary (max 100 items)
	Stream   bool     // Enable streaming via Event Stream
	UserID   string   // End user ID (6-128 characters)
	RequestID string  // Client-provided unique identifier
}