package common

import (
	"context"
	"fmt"
	"net/url"

	"github.com/hofstadter-io/hof/lib/agent/services/environ"
)

// extractSid extracts the session ID from a URI string.
// Supported formats:
//   veg://host:port/sid:version/path...
//   oci://host:port/sid:version?path=...
// Translation from veg:// to oci:// happens in the VS Code extension.
// func extractSid(uriStr string) string {
// 	u, err := url.Parse(uriStr)
// 	if err != nil {
// 		return ""
// 	}

// 	// For both veg:// and oci://, the first segment of the path 
// 	// (after the authority) is expected to be sid:version or just sid.
// 	parts := strings.Split(strings.TrimPrefix(u.Path, "/"), "/")
// 	if len(parts) > 0 {
// 		idPart := parts[0]
// 		// The sid might be followed by a version (sid:version)
// 		if idx := strings.Index(idPart, ":"); idx != -1 {
// 			return idPart[:idx]
// 		}
// 		return idPart
// 	}
// 	return ""
// }

func FilesysWrite(ctx context.Context, ar Runtime, user, uri, path, content, sid string) (string, error) {
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

	if sid != "" {
		// Update session state currEnv
		u, _ := url.Parse(nextUri)
		envVal := u.Host + u.Path
		err := SessionStatePut(ctx, ar, user, sid, "currEnv", envVal)
		if err != nil {
			fmt.Printf("failed to update session state: %v\n", err)
		}

		// TODO: Add partial user event to session history
	}

	return nextUri, nil
}

func FilesysRead(ctx context.Context, ar Runtime, user, uri, path, sid string, diff bool) (string, error) {
	if sid != "" {
		_, err := SessionGet(ctx, ar, user, sid)
		if err != nil {
			return "", fmt.Errorf("session access denied: %w", err)
		}
	}
	return environ.Client().ReadFile(uri, path, diff)
}

func FilesysList(ctx context.Context, ar Runtime, user, uri, path, sid string, diff bool) (any, error) {
	if sid != "" {
		_, err := SessionGet(ctx, ar, user, sid)
		if err != nil {
			return nil, fmt.Errorf("session access denied: %w", err)
		}
	}
	return environ.Client().ReadDirectory(uri, path, diff)
}

func FilesysStat(ctx context.Context, ar Runtime, user, uri, path, sid string, diff bool) (any, error) {
	if sid != "" {
		_, err := SessionGet(ctx, ar, user, sid)
		if err != nil {
			return nil, fmt.Errorf("session access denied: %w", err)
		}
	}
	return environ.Client().Stat(uri, path, diff)
}

func FilesysDelete(ctx context.Context, ar Runtime, user, uri, path, sid string) (string, error) {
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

	if sid != "" {
		u, _ := url.Parse(nextUri)
		envVal := u.Host + u.Path
		SessionStatePut(ctx, ar, user, sid, "currEnv", envVal)
		// TODO: Add partial user event
	}

	return nextUri, nil
}

func FilesysMkdir(ctx context.Context, ar Runtime, user, uri, path, sid string) (string, error) {
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

	if sid != "" {
		u, _ := url.Parse(nextUri)
		envVal := u.Host + u.Path
		SessionStatePut(ctx, ar, user, sid, "currEnv", envVal)
		// TODO: Add partial user event
	}

	return nextUri, nil
}

func FilesysRename(ctx context.Context, ar Runtime, user, uri, src, dst, sid string) (string, error) {
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

	if sid != "" {
		u, _ := url.Parse(nextUri)
		envVal := u.Host + u.Path
		SessionStatePut(ctx, ar, user, sid, "currEnv", envVal)
		// TODO: Add partial user event
	}

	return nextUri, nil
}

func FilesysCopy(ctx context.Context, ar Runtime, user, uri, src, dst, sid string) (string, error) {
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

	if sid != "" {
		u, _ := url.Parse(nextUri)
		envVal := u.Host + u.Path
		SessionStatePut(ctx, ar, user, sid, "currEnv", envVal)
		// TODO: Add partial user event
	}

	return nextUri, nil
}
