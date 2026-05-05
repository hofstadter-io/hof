package common

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/hofstadter-io/hof/lib/agent/services/environ"
)

/*
	Common Information

	uri-format:
*/

// parseUriPath splits an input URI (e.g. oci://host/img:tag?path=./foo)
// into a clean OCI reference and the relative file path.
// Input uri: The full OCI reference with optional query params.
// Input path: Optional path override. If empty, uses 'path' from uri query.
// Output: (cleanUri, filePath, error)
func parseUriPath(uri, path string) (string, string, error) {
	u, err := url.Parse(uri)
	if err != nil {
		return "", "", err
	}
	if path == "" {
		path = u.Query().Get("path")
	}
	// Strip query params to get a valid OCI reference (e.g. host/img:tag)
	// We use the original string to avoid Go's url package re-encoding characters (like :) in the path segment
	cleanUri := uri
	cleanUri = strings.TrimPrefix(cleanUri, "oci://")
	if idx := strings.Index(cleanUri, "?"); idx != -1 {
		cleanUri = cleanUri[:idx]
	}
	return cleanUri, path, nil
}

func updateSessionEnv(ctx context.Context, ar Runtime, user, sid, nextUri string) error {
	if sid == "" {
		return nil
	}
	// Update session state currEnv
	parseUri := nextUri
	if !strings.Contains(parseUri, "://") {
		parseUri = "oci://" + parseUri
	}
	u, _ := url.Parse(parseUri)
	envVal := u.Host + u.Path
	err := SessionStatePut(ctx, ar, user, sid, "currEnv", envVal)
	if err != nil {
		return fmt.Errorf("failed to update session state: %w", err)
	}

	// TODO: Add partial user event to session history
	return nil
}

func FilesysWrite(ctx context.Context, ar Runtime, user, uri, path, content, sid string) (string, error) {
	var err error
	// Input uri: oci://host/img:tag?path=...
	// Input path: relative path in container
	uri, path, err = parseUriPath(uri, path)
	if err != nil {
		return "", err
	}
	// Result: uri is clean OCI ref, path is target file path

	if sid != "" {
		// Verify session access
		_, err := SessionGet(ctx, ar, user, sid)
		if err != nil {
			return "", fmt.Errorf("session access denied: %w", err)
		}
	}

	nextUri, err := environ.Client().WriteFile(uri, path, content)
	if err != nil {
		return "", err
	}

	err = updateSessionEnv(ctx, ar, user, sid, nextUri)
	if err != nil {
		return nextUri, err
	}

	return nextUri, nil
}

func FilesysRead(ctx context.Context, ar Runtime, user, uri, path, sid string, diff bool) (string, error) {
	var err error
	// Input uri: oci://host/img:tag?path=...
	uri, path, err = parseUriPath(uri, path)
	if err != nil {
		return "", err
	}

	if sid != "" {
		_, err := SessionGet(ctx, ar, user, sid)
		if err != nil {
			return "", fmt.Errorf("session access denied: %w", err)
		}
	}
	return environ.Client().ReadFile(uri, path, diff)
}

func FilesysList(ctx context.Context, ar Runtime, user, uri, path, sid string, diff bool) (any, error) {
	var err error
	// Input uri: oci://host/img:tag?path=...
	uri, path, err = parseUriPath(uri, path)
	if err != nil {
		return nil, err
	}

	if sid != "" {
		_, err := SessionGet(ctx, ar, user, sid)
		if err != nil {
			return nil, fmt.Errorf("session access denied: %w", err)
		}
	}
	return environ.Client().ReadDirectory(uri, path, diff)
}

func FilesysStat(ctx context.Context, ar Runtime, user, uri, path, sid string, diff bool) (any, error) {
	var err error
	// Input uri: oci://host/img:tag?path=...
	uri, path, err = parseUriPath(uri, path)
	if err != nil {
		return nil, err
	}

	if sid != "" {
		_, err := SessionGet(ctx, ar, user, sid)
		if err != nil {
			return nil, fmt.Errorf("session access denied: %w", err)
		}
	}
	return environ.Client().Stat(uri, path, diff)
}

func FilesysDelete(ctx context.Context, ar Runtime, user, uri, path, sid string) (string, error) {
	var err error
	// Input uri: oci://host/img:tag?path=...
	uri, path, err = parseUriPath(uri, path)
	if err != nil {
		return "", err
	}

	if sid != "" {
		_, err := SessionGet(ctx, ar, user, sid)
		if err != nil {
			return "", fmt.Errorf("session access denied: %w", err)
		}
	}

	nextUri, err := environ.Client().Delete(uri, path, true)
	if err != nil {
		return "", err
	}

	err = updateSessionEnv(ctx, ar, user, sid, nextUri)
	if err != nil {
		return nextUri, err
	}

	return nextUri, nil
}

func FilesysMkdir(ctx context.Context, ar Runtime, user, uri, path, sid string) (string, error) {
	var err error
	// Input uri: oci://host/img:tag?path=...
	uri, path, err = parseUriPath(uri, path)
	if err != nil {
		return "", err
	}

	if sid != "" {
		_, err := SessionGet(ctx, ar, user, sid)
		if err != nil {
			return "", fmt.Errorf("session access denied: %w", err)
		}
	}

	nextUri, err := environ.Client().CreateDirectory(uri, path)
	if err != nil {
		return "", err
	}

	err = updateSessionEnv(ctx, ar, user, sid, nextUri)
	if err != nil {
		return nextUri, err
	}

	return nextUri, nil
}

func FilesysRename(ctx context.Context, ar Runtime, user, uri, src, dst, sid string) (string, error) {
	var err error
	// Input uri: oci://host/img:tag (query params ignored for rename base)
	uri, src, err = parseUriPath(uri, src)
	if err != nil {
		return "", err
	}

	if sid != "" {
		_, err := SessionGet(ctx, ar, user, sid)
		if err != nil {
			return "", fmt.Errorf("session access denied: %w", err)
		}
	}

	nextUri, err := environ.Client().Move(uri, src, dst, true)
	if err != nil {
		return "", err
	}

	err = updateSessionEnv(ctx, ar, user, sid, nextUri)
	if err != nil {
		return nextUri, err
	}

	return nextUri, nil
}

func FilesysCopy(ctx context.Context, ar Runtime, user, uri, src, dst, sid string) (string, error) {
	var err error
	// Input uri: oci://host/img:tag (query params ignored for copy base)
	uri, src, err = parseUriPath(uri, src)
	if err != nil {
		return "", err
	}

	if sid != "" {
		_, err := SessionGet(ctx, ar, user, sid)
		if err != nil {
			return "", fmt.Errorf("session access denied: %w", err)
		}
	}

	nextUri, err := environ.Client().Copy(uri, src, dst, true)
	if err != nil {
		return "", err
	}

	err = updateSessionEnv(ctx, ar, user, sid, nextUri)
	if err != nil {
		return nextUri, err
	}

	return nextUri, nil
}

func FilesysDiff(ctx context.Context, ar Runtime, user, prevUri, nextUri, sid string) (any, error) {
	var err error
	// Input prevUri/nextUri: oci://host/img:tag (query params stripped)
	prevUri, _, err = parseUriPath(prevUri, "")
	if err != nil {
		return nil, err
	}
	nextUri, _, err = parseUriPath(nextUri, "")
	if err != nil {
		return nil, err
	}

	if sid != "" {
		_, err := SessionGet(ctx, ar, user, sid)
		if err != nil {
			return nil, fmt.Errorf("session access denied: %w", err)
		}
	}
	return environ.Client().DiffDirectory(prevUri, nextUri)
}
