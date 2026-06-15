package traffic

import (
	"crypto/sha256"
	"encoding/hex"
	"net"
	"strings"

	"github.com/DoraZa/mini-gateway/config"
	"github.com/gin-gonic/gin"
)

func matchIPLimit(cfg *config.Config, clientIP string) *config.TrafficRateLimit {
	for pattern, limit := range cfg.Traffic.RateLimit.IPLimits {
		if !limit.Enabled {
			continue
		}
		if matchIPPattern(pattern, clientIP) {
			copied := limit
			return &copied
		}
	}
	return nil
}

func matchRouteLimit(cfg *config.Config, path string) (*config.TrafficRateLimit, string) {
	var bestPattern string
	var bestLimit *config.TrafficRateLimit
	for pattern, limit := range cfg.Traffic.RateLimit.RouteLimits {
		if !limit.Enabled || !matchPathPattern(pattern, path) {
			continue
		}
		if len(pattern) > len(bestPattern) {
			copied := limit
			bestLimit = &copied
			bestPattern = pattern
		}
	}
	return bestLimit, bestPattern
}

func matchPathPattern(pattern, path string) bool {
	if pattern == path {
		return true
	}
	if strings.Contains(pattern, "*") {
		prefix := strings.Split(pattern, "*")[0]
		return strings.HasPrefix(path, prefix)
	}
	return false
}

func matchIPPattern(pattern, clientIP string) bool {
	if pattern == clientIP {
		return true
	}
	if _, ipNet, err := net.ParseCIDR(pattern); err == nil {
		ip := net.ParseIP(clientIP)
		return ip != nil && ipNet.Contains(ip)
	}
	return false
}

func routeLimiterKey(pattern string, limit *config.TrafficRateLimit, c *gin.Context) string {
	if limit == nil || !limit.PerUser {
		return pattern
	}
	return pattern + ":user:" + callerIdentity(c)
}

func callerIdentity(c *gin.Context) string {
	if userID := c.GetHeader("X-User-ID"); userID != "" {
		return userID
	}
	if auth := c.GetHeader("Authorization"); strings.HasPrefix(auth, "Bearer ") {
		sum := sha256.Sum256([]byte(auth))
		return hex.EncodeToString(sum[:8])
	}
	return "ip:" + c.ClientIP()
}
