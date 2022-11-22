package sqlcommenter

import (
	"gorm.io/gorm/clause"
	"net/url"
	"sort"
)

// SqlComment contains tags for the request
type SqlComment struct {
	Tags map[string]string
}

func (sqlComment SqlComment) Name() string {
	return GormClauseName
}

// Build builds the comment clause
func (sqlComment SqlComment) Build(builder clause.Builder) {
	if len(sqlComment.Tags) > 0 {
		builder.WriteString("/*")

		keys := make([]string, 0, len(sqlComment.Tags))
		for key := range sqlComment.Tags {
			keys = append(keys, key)
		}
		sort.Strings(keys)

		isFirstValue := true
		for _, key := range keys {
			if !isFirstValue {
				builder.WriteString(",")
			}
			builder.WriteString(escapeKey(key))
			builder.WriteString("=")
			builder.WriteString(escapeValue(sqlComment.Tags[key]))
			isFirstValue = false
		}

		builder.WriteString("*/")
	}
}

// MergeClause merge SqlCommenter clauses
func (sqlComment SqlComment) MergeClause(mergeClause *clause.Clause) {
	if s, ok := mergeClause.Expression.(SqlComment); ok {
		tags := make(map[string]string, len(s.Tags)+len(sqlComment.Tags))
		for key, value := range s.Tags {
			tags[key] = value
		}
		for key, value := range sqlComment.Tags {
			tags[key] = value
		}
		sqlComment.Tags = tags
	}
	mergeClause.Expression = sqlComment
}

func escapeKey(key string) string {
	return url.PathEscape(key)
}

func escapeValue(value string) string {
	return "'" + escapeKey(value) + "'"
}

func NewTag(key string, value string) SqlComment {
	return SqlComment{Tags: map[string]string{
		key: value,
	}}
}

func NewTags(tags map[string]string) SqlComment {
	return SqlComment{Tags: tags}
}
