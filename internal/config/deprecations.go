package config

import (
	"fmt"
	"strings"
)

type deprecated string

var (
	// fields that are deprecated along with their messages
	deprecatedFields = map[deprecated]string{
		"tracing.jaeger.enabled":          deprecatedMsgTracingJaegerEnabled,
		"cache.memory.enabled":            deprecatedMsgCacheMemoryEnabled,
		"cache.memory.expiration":         deprecatedMsgCacheMemoryExpiration,
		"experimental.filesystem_storage": deprecatedMsgExperimentalFilesystemStorage,
	}
)

const (
	deprecatedDefaultMessage = `%q is deprecated.`

	// additional deprecation messages
	deprecatedMsgTracingJaegerEnabled          = `Please use 'tracing.enabled' and 'tracing.exporter' instead.`
	deprecatedMsgCacheMemoryEnabled            = `Please use 'cache.enabled' and 'cache.backend' instead.`
	deprecatedMsgCacheMemoryExpiration         = `Please use 'cache.ttl' instead.`
	deprecatedMsgExperimentalFilesystemStorage = `The experimental filesystem storage backend has graduated to a stable feature. Please use 'storage' instead.`
)

func (d deprecated) Message() string {
	msg, ok := deprecatedFields[d]
	if !ok {
		return strings.TrimSpace(fmt.Sprintf(deprecatedDefaultMessage, d))
	}

	msg = strings.Join([]string{deprecatedDefaultMessage, msg}, " ")
	return strings.TrimSpace(fmt.Sprintf(msg, d))
}
