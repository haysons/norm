package clause

import (
	"fmt"
	"github.com/haysons/norm/resolver"
	"strconv"
	"strings"
)

type CreateTag struct {
	IfNotExist bool
	Tag        *resolver.VertexTag
}

const CreateTagName = "CREATE_TAG"

func (ct CreateTag) Name() string {
	return CreateTagName
}

func (ct CreateTag) MergeIn(clause *Clause) {
	clause.Expression = ct
}

func (ct CreateTag) Build(nGQL Builder) error {
	nGQL.WriteString("CREATE TAG ")
	if ct.IfNotExist {
		nGQL.WriteString("IF NOT EXISTS ")
	}
	nGQL.WriteString(ct.Tag.TagName)
	nGQL.WriteByte('(')
	propsLen := len(ct.Tag.GetProps())
	ttlCols := make([]string, 0, 1)
	ttlDuration := ""
	for i, prop := range ct.Tag.GetProps() {
		if prop.Name == "" || prop.DataType == "" {
			return fmt.Errorf("norm: %w, build create tag clause failed, tag prop must has name and data type", ErrInvalidClauseParams)
		}
		nGQL.WriteString(prop.Name)
		nGQL.WriteByte(' ')
		nGQL.WriteString(prop.DataType)
		if prop.NotNull {
			nGQL.WriteString(" NOT NULL")
		}
		if prop.Default != "" {
			nGQL.WriteString(" DEFAULT ")
			switch strings.ToLower(prop.DataType) {
			case "string", "fixed_string":
				nGQL.WriteString(strconv.Quote(prop.Default))
			default:
				nGQL.WriteString(prop.Default)
			}
		}
		if prop.Comment != "" {
			nGQL.WriteString(" COMMENT ")
			nGQL.WriteByte('"')
			nGQL.WriteString(prop.Comment)
			nGQL.WriteByte('"')
		}
		if prop.TTL != "" {
			ttlCols = append(ttlCols, prop.Name)
			ttlDuration = prop.TTL
		}
		if i < propsLen-1 {
			nGQL.WriteString(", ")
		}
	}
	nGQL.WriteByte(')')
	if len(ttlCols) > 1 {
		return fmt.Errorf("norm: %w, build create tag clause failed, must only one ttl col", ErrInvalidClauseParams)
	}
	if len(ttlCols) == 1 && ttlDuration != "" {
		nGQL.WriteString(" TTL_DURATION = ")
		nGQL.WriteString(ttlDuration)
		nGQL.WriteString(", TTL_COL = ")
		nGQL.WriteString(strconv.Quote(ttlCols[0]))
	}
	return nil
}
