package config

import (
	"os"
	"runtime"
	"strings"
)

func (c *Config) ResolvedMappings(local *LocalConfig) (map[string]string, map[string]string) {
	vars := c.resolveVars(local)
	merged := make(map[string]string)
	for k, v := range c.Mappings {
		merged[k] = v
	}
	switch runtime.GOOS {
	case "windows":
		for k, v := range c.Windows {
			merged[k] = v
		}
	case "darwin":
		for k, v := range c.Macos {
			merged[k] = v
		}
	case "linux":
		for k, v := range c.Linux {
			merged[k] = v
		}
	}
	resolved := make(map[string]string)
	missingVars := make(map[string]string)
	for src, dest := range merged {
		out, ok := substituteVars(dest, vars)
		if ok {
			resolved[src] = out
		} else {
			missingVars[src] = dest
		}
	}
	return resolved, missingVars
}

func (c *Config) resolveVars(local *LocalConfig) map[string]string {
	vars := make(map[string]string)
	for _, name := range c.Variables {
		if v, ok := local.Values[name]; ok {
			vars[name] = v
			continue
		}
		if v := os.Getenv(name); v != "" {
			vars[name] = v
		}
	}
	return vars
}

func substituteVars(template string, vars map[string]string) (string, bool) {
	out := template
	for name, val := range vars {
		key := "$" + name
		out = strings.ReplaceAll(out, key, val)
	}
	if strings.Contains(out, "$") {
		return "", false
	}
	return out, true
}

func (c *Config) MissingVariables(local *LocalConfig) []string {
	var missing []string
	for _, name := range c.Variables {
		if _, ok := local.Values[name]; ok {
			continue
		}
		if os.Getenv(name) != "" {
			continue
		}
		missing = append(missing, name)
	}
	return missing
}
