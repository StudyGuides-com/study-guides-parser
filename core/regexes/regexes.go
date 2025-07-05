package regexes

import "regexp"

// ListItemPrefixRegex is a compiled regular expression that matches list item prefixes
// in the format of numbered items (e.g., "1.") or bullet points (e.g., "*" or "-")
var ListItemPrefixRegex = regexp.MustCompile(`^(\d+\.|\*|\-)\s+`)
