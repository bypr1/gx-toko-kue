package activity

import (
	"service/internal/pkg/constant"
	"service/internal/pkg/core"

	"gorm.io/gorm"
)

func InitCreate(reference ActivityModelInterface, conn *gorm.DB) UseActivity {
	return UseActivity{
		Connection:    conn,
		ReferenceID:   reference.SetReference(),
		ReferenceType: reference.TableName(),
		Action:        constant.ACTION_CREATE,
	}
}

func InitUpdate(reference ActivityModelInterface, conn *gorm.DB) UseActivity {
	return UseActivity{
		Connection:    conn,
		ReferenceID:   reference.SetReference(),
		ReferenceType: reference.TableName(),
		Action:        constant.ACTION_UPDATE,
	}
}

func InitDelete(reference ActivityModelInterface, conn *gorm.DB) UseActivity {
	return UseActivity{
		Connection:    conn,
		ReferenceID:   reference.SetReference(),
		ReferenceType: reference.TableName(),
		Action:        constant.ACTION_DELETE,
	}
}

func (ua UseActivity) ParseOldProperty(parser core.BaseActivityPropertyParserInterface, subs ...string) UseActivity {
	ua.Parser = parser
	ua.Property.Old = ua.setPropertyWithParser(ua.Action, subs...)
	return ua
}

func (ua UseActivity) ParseNewProperty(parser core.BaseActivityPropertyParserInterface, subs ...string) UseActivity {
	ua.Parser = parser
	ua.Property.New = ua.setPropertyWithParser(ua.Action, subs...)
	return ua
}
