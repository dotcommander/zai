package app

import (
	"strings"

	"github.com/dotcommander/zai/internal/app/fileutil"
)

// EnableSessionFileReuse scopes local file reuse to this client and one chat file.
// Call before starting chat requests; ordinary one-shot clients remain uncached.
func (c *Client) EnableSessionFileReuse(path string) {
	if path != "" && !strings.HasPrefix(path, "https://") && !strings.HasPrefix(path, "http://") {
		c.fileReader = fileutil.NewSessionFileReader(c.fileReader, path)
	}
}
