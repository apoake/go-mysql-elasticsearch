package river

import (
	"github.com/siddontang/go-mysql/schema"
	"strings"
)

// If you want to sync MySQL data into elasticsearch, you must set a rule to let use know how to do it.
// The mapping rule may thi: schema + table <-> index + document type.
// schema and table is for MySQL, index and document type is for Elasticsearch.
type Rule struct {
	Schema string `toml:"schema"`
	Table  string `toml:"table"`
	Index  string `toml:"index"`
	Type   string `toml:"type"`
	Parent []string `toml:"parent"`
	ID []string `toml:"id"`

	// Default, a MySQL table field name is mapped to Elasticsearch field name.
	// Sometimes, you want to use different name, e.g, the MySQL file name is title,
	// but in Elasticsearch, you want to name it my_title.
	FieldMapping map[string]string `toml:"field"`

	// MySQL table information
	TableInfo *schema.Table

	//only MySQL fields in fileter will be synced , default sync all fields
	Fileter []string `toml:"filter"`

	Extension	string `toml:"extension"`
	Extensions 	map[string]string
}

func newDefaultRule(schema string, table string) *Rule {
	r := new(Rule)

	r.Schema = schema
	r.Table = table
	r.Index = table
	r.Type = table
	r.FieldMapping = make(map[string]string)
	return r
}

func (r *Rule) prepare() error {
	if r.FieldMapping == nil {
		r.FieldMapping = make(map[string]string)
	}

	if len(r.Index) == 0 {
		r.Index = r.Table
	}

	if len(r.Type) == 0 {
		r.Type = r.Index
	}

	r.initExtensions()

	return nil
}

// add
func (r *Rule) initExtensions() {
	r.Extensions = make(map[string]string)
	if r.Extension != "" {
		strs := strings.Split(r.Extension, ",")
		for _, str := range  strs {
			if len(str) == 0 {
				continue
			}
			subStrs := strings.Split(str, ":")
			r.Extensions[strings.TrimSpace(subStrs[0])] = strings.TrimSpace(subStrs[1])
		}
	}
}

func (r *Rule) CheckFilter(field string) bool {
	if r.Fileter == nil {
		return true
	}

	for _, f := range r.Fileter {
		if f == field {
			return true
		}
	}
	return false
}
