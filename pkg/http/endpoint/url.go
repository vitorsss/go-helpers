package endpoint

import (
	"net/url"
	"strings"
)

func joinURL(baseURI string, paths ...string) (string, error) {
	baseURIToValidade, err := replaceURLParams(baseURI, map[string]string{}, true)
	if err != nil {
		return "", err
	}
	parsedURL, err := url.Parse(baseURIToValidade)
	if err != nil {
		return "", err
	}
	parsedURL = parsedURL.JoinPath(paths...)
	urlStr := unescapedURLString(parsedURL)
	_, err = url.Parse(urlStr)
	return strings.ReplaceAll(urlStr, baseURIToValidade, baseURI), err
}

func unescapedURLString(u *url.URL) string {
	var buf strings.Builder
	if u.Scheme != "" {
		buf.WriteString(u.Scheme)
		buf.WriteByte(':')
	}
	if u.Scheme != "" || u.Host != "" || u.User != nil {
		if u.OmitHost && u.Host == "" && u.User == nil {
			// omit empty host
		} else {
			if u.Host != "" || u.Path != "" || u.User != nil {
				buf.WriteString("//")
			}
			if ui := u.User; ui != nil {
				buf.WriteString(ui.String())
				buf.WriteByte('@')
			}
			if h := u.Host; h != "" {
				buf.WriteString(h)
			}
		}
	}
	path := u.Path
	if path != "" && path[0] != '/' && u.Host != "" {
		buf.WriteByte('/')
	}
	if buf.Len() == 0 {
		if segment, _, _ := strings.Cut(path, "/"); strings.Contains(segment, ":") {
			buf.WriteString("./")
		}
	}
	buf.WriteString(path)
	if u.ForceQuery || u.RawQuery != "" {
		buf.WriteByte('?')
		buf.WriteString(u.RawQuery)
	}
	if u.Fragment != "" {
		buf.WriteByte('#')
		buf.WriteString(u.EscapedFragment())
	}
	return buf.String()
}
