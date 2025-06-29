package parser

import (
	xtremefs "github.com/globalxtreme/go-core/v2/filesystem"
	"service/internal/pkg/model"
)

type TestingParser struct {
	Array  []model.Testing
	Object model.Testing
}

func (parser TestingParser) Get() []interface{} {
	var result []interface{}

	for _, activity := range parser.Array {
		firstParser := TestingParser{Object: activity}
		result = append(result, firstParser.First())
	}

	return result
}

func (parser TestingParser) First() interface{} {
	activity := parser.Object

	var resSubs []interface{}
	for _, sub := range activity.Subs {
		resSubs = append(resSubs, map[string]interface{}{
			"id":        sub.ID,
			"name":      sub.Name,
			"createdAt": sub.CreatedAt.Format("02/01/2006 15:04"),
		})
	}

	return map[string]interface{}{
		"id":        activity.ID,
		"name":      activity.Name,
		"createdAt": activity.CreatedAt.Format("02/01/2006 15:04"),
		"file":      xtremefs.Storage{}.GetFullPathURL("ckH2cahaAaDMNVgS2xdM1697957810885349000.png"),
		"subs":      resSubs,
	}
}

func (parser TestingParser) CreateActivity(action string) interface{} {
	activity := parser.Object

	var resSubs []interface{}
	for _, sub := range activity.Subs {
		resSubs = append(resSubs, map[string]interface{}{
			"name": sub.Name,
		})
	}

	return map[string]interface{}{
		"name": activity.Name,
		"file": xtremefs.Storage{}.GetFullPathURL("ckH2cahaAaDMNVgS2xdM1697957810885349000.png"),
		"subs": resSubs,
	}
}

func (parser TestingParser) UpdateActivity(action string) interface{} {
	return parser.CreateActivity(action)
}

func (parser TestingParser) DeleteActivity(action string) interface{} {
	return parser.CreateActivity(action)
}

func (parser TestingParser) GeneralActivity(action string) interface{} {
	if action == "onlyName" {
		activity := parser.Object

		return map[string]interface{}{
			"name": activity.Name,
		}
	}

	return parser.CreateActivity(action)
}
