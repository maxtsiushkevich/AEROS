package middleware

import (
	"net/http"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
)

var (
	reColon  = regexp.MustCompile(`:([A-Za-z0-9_]+)`)
	reStar   = regexp.MustCompile(`\*([A-Za-z0-9_]+)`)
	reDigits = regexp.MustCompile(`^[0-9]+$`)
	reUUID   = regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`)
	reMethod = regexp.MustCompile(`(?i)^(GET|POST|PUT|PATCH|DELETE|OPTIONS|HEAD|QUERY)[:\s]+(.+)$`)
)

// NormalizePathPattern converts various route declaration styles to a
// unified pattern where path parameters are wrapped in { }.
//
//	/users/:id    -> /users/{id}
//	/files/*name   -> /files/{name}
//	/users/123     -> /users/{param}  (heuristic fallback for numeric/uuid segments)
func NormalizePathPattern(p string) string {
	if p == "" {
		return p
	}

	if matches := reMethod.FindStringSubmatch(p); len(matches) == 3 {
		p = matches[2]
	}

	// Replace gin-style ":param" and wildcard "*param" with {param}
	p = reColon.ReplaceAllString(p, `{$1}`)
	p = reStar.ReplaceAllString(p, `{$1}`)

	parts := strings.Split(p, "/")
	for i, seg := range parts {
		if seg == "" {
			continue
		}
		if strings.HasPrefix(seg, "{") && strings.HasSuffix(seg, "}") {
			continue
		}
		if reDigits.MatchString(seg) || reUUID.MatchString(seg) {
			parts[i] = "{param}"
		}
	}

	return strings.Join(parts, "/")
}

func CasbinObjFromGin(c *gin.Context) string {
	obj := c.FullPath()
	if obj == "" {
		obj = c.Request.URL.Path
	}
	return NormalizePathPattern(obj)
}

func CasbinObjFromRequest(r *http.Request) string {
	obj := r.URL.Path
	return NormalizePathPattern(obj)
}
