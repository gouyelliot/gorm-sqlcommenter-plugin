package sqlcommenter

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const GormClauseName = "plugin:sql-commenter"

type SqlCommenterPlugin struct{}

func (e *SqlCommenterPlugin) Name() string {
	return "SqlCommenterPlugin"
}

// Initialize register BuildClauses used to generate the comment
func (e *SqlCommenterPlugin) Initialize(db *gorm.DB) error {
	db.ClauseBuilders[GormClauseName] = func(c clause.Clause, builder clause.Builder) {
		if sqlComment, ok := c.Expression.(SqlComment); ok {
			sqlComment.Build(builder)
		}
	}

	db.Callback().Query().Clauses = append(db.Callback().Query().Clauses, GormClauseName)
	db.Callback().Create().Clauses = append(db.Callback().Create().Clauses, GormClauseName)
	db.Callback().Update().Clauses = append(db.Callback().Update().Clauses, GormClauseName)
	db.Callback().Delete().Clauses = append(db.Callback().Delete().Clauses, GormClauseName)
	db.Callback().Raw().Clauses = append(db.Callback().Raw().Clauses, GormClauseName)
	db.Callback().Row().Clauses = append(db.Callback().Row().Clauses, GormClauseName)
	return nil
}

// New create a new SqlCommenterPlugin
func New() *SqlCommenterPlugin {
	return &SqlCommenterPlugin{}
}
