// Package strmangle is used exclusively by the templates in sqlboiler.
// There are many helper functions to deal with bdb.* values as well
// as string manipulation. Because it is focused on pipelining inside templates
// you will see some odd parameter ordering.
package strmangle

import (
	"fmt"
	"strings"

	"github.com/jinzhu/inflection"
)

// Plural converts singular words to plural words (eg: person to people)
func Plural(name string) string {
	splits := strings.Split(name, "_")
	splits[len(splits)-1] = inflection.Plural(splits[len(splits)-1])
	return strings.Join(splits, "_")
}

// Singular converts plural words to singular words (eg: people to person)
func Singular(name string) string {
	splits := strings.Split(name, "_")
	splits[len(splits)-1] = inflection.Singular(splits[len(splits)-1])
	return strings.Join(splits, "_")
}

// TitleCase changes a snake-case variable name
// into a go styled object variable name of "ColumnName".
// titleCase also fully uppercases "ID" components of names, for example
// "column_name_id" to "ColumnNameID".
func TitleCase(name string) string {
	splits := strings.Split(name, "_")

	for i, split := range splits {
		if split == "id" {
			splits[i] = "ID"
			continue
		}

		splits[i] = strings.Title(split)
	}

	return strings.Join(splits, "")
}

// CamelCase takes a variable name in the format of "var_name" and converts
// it into a go styled variable name of "varName".
// camelCase also fully uppercases "ID" components of names, for example
// "var_name_id" to "varNameID".
func CamelCase(name string) string {
	splits := strings.Split(name, "_")

	for i, split := range splits {
		if split == "id" && i > 0 {
			splits[i] = "ID"
			continue
		}

		if i == 0 {
			continue
		}

		splits[i] = strings.Title(split)
	}

	return strings.Join(splits, "")
}

// StringMap maps a function over a slice of strings.
func StringMap(modifier func(string) string, strs []string) []string {
	ret := make([]string, len(strs))

	for i, str := range strs {
		ret[i] = modifier(str)
	}

	return ret
}

// MakeDBName takes a table name in the format of "table_name" and a
// column name in the format of "column_name" and returns a name used in the
// `db:""` component of an object in the format of "table_name_column_name"
func MakeDBName(tableName, colName string) string {
	return fmt.Sprintf("%s_%s", tableName, colName)
}

// HasElement checks to see if the string is found in the string slice
func HasElement(str string, slice []string) bool {
	for _, s := range slice {
		if str == s {
			return true
		}
	}

	return false
}

// PrefixStringSlice with the given str.
func PrefixStringSlice(str string, strs []string) []string {
	ret := make([]string, len(strs))

	for i, s := range strs {
		ret[i] = fmt.Sprintf("%s%s", str, s)
	}

	return ret
}

// GenerateParamFlags generates the SQL statement parameter flags
// For example, $1,$2,$3 etc. It will start counting at startAt.
func GenerateParamFlags(colCount int, startAt int) string {
	cols := make([]string, 0, colCount)

	for i := startAt; i < colCount+startAt; i++ {
		cols = append(cols, fmt.Sprintf("$%d", i))
	}

	return strings.Join(cols, ",")
}

// WhereClause returns the where clause using start as the $ flag index
// For example, if start was 2 output would be: "colthing=$2 AND colstuff=$3"
func WhereClause(cols []string, start int) string {
	if start == 0 {
		panic("0 is not a valid start number for whereClause")
	}

	ret := make([]string, len(cols))
	for i, c := range cols {
		ret[i] = fmt.Sprintf(`"%s"=$%d`, c, start+i)
	}

	return strings.Join(ret, " AND ")
}

// DriverUsesLastInsertID returns whether the database driver supports the
// sql.Result interface.
func DriverUsesLastInsertID(driverName string) bool {
	switch driverName {
	case "postgres":
		return false
	default:
		return true
	}
}

// Substring returns a substring of str starting at index start and going
// to end-1.
func Substring(start, end int, str string) string {
	return str[start:end]
}

// JoinSlices merges two string slices of equal length
func JoinSlices(sep string, a, b []string) []string {
	lna, lnb := len(a), len(b)
	if lna != lnb {
		panic("joinSlices: can only merge slices of same length")
	} else if lna == 0 {
		return nil
	}

	ret := make([]string, len(a))
	for i, elem := range a {
		ret[i] = fmt.Sprintf("%s%s%s", elem, sep, b[i])
	}

	return ret
}
