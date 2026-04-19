package urlfetch

import "net/url"

func redactURL(u *url.URL) string {
	if u == nil {
		return ""
	}
	return u.Redacted()
}
