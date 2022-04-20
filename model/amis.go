package models

import (
	"fmt"
	"reflect"
	"strings"
)

//"errors"

/*
分页基类,每个分页基本都要这些字段（针对amis定制）
*/
type PagePms struct {
	Page int `json:"page" binding:"required"`
	Size int `json:"size" binding:"required,gt=0,lt=1000000"`
	// Sort []Sorts    `json:"sort"`
	Conds    Conditions `json:"conditions"`
	OrderBy  string     `json:"orderBy"`  //orderBy: "Name"
	OrderDir string     `json:"orderDir"` //orderDir: "asc"
}

type Conditions struct {
	Conjunction string       `json:"conjunction"` //conjunction: "and"
	Children    []Conditions `json:"children"`

	Id    string        `json:"id"`
	Left  ConditionLeft `json:"left"`
	Field string        `json:"field"`
	Type  string        `json:"type"`
	Op    string        `json:"op"`
	Right interface{}   `json:"right"`
}

/*
id: "ffcf44ad0ea7"
left: {type: "field", field: "Name"}
field: "Name"
type: "field"
op: "like"
right: ""
**/
// type Condition struct {
// 	Id    string        `json:"id"`
// 	Left  ConditionLeft `json:"left"`
// 	Field string        `json:"field"`
// 	Type  string        `json:"type"`
// 	Op    string        `json:"op"`
// 	Right string        `json:"right"`
// }

type ConditionLeft struct {
	Field string `json:"field"`
	Type  string `json:"type"`
}

func (m *PagePms) GetFlts() (strs []string) {
	strs = make([]string, 2)
	if m.Page == 0 {
		m.Page = 1
	}
	orderstr := ""

	if len(m.OrderBy) > 0 && len(m.OrderDir) > 0 {
		orderstr = strings.ToLower(m.OrderBy) + " " + m.OrderDir
	}

	filterarr := m.FilterItems(m.Conds.Children, m.Conds)

	fmt.Println("解析的数组：", filterarr)
	strsValid := make([]string, 0)
	//过滤掉空格字符元素
	for _, v := range filterarr {
		if len(v) > 5 {
			strsValid = append(strsValid, v)
		}
	}
	fmt.Println("过滤后的的数组：", strsValid)
	filterstr := strings.Join(strsValid, " "+m.Conds.Conjunction+" ")
	/*
	   HasPrefix 判断字符串 s 是否以 prefix 开头：
	   strings.HasPrefix(s, prefix string) bool
	   HasSuffix 判断字符串 s 是否以 suffix 结尾：
	   strings.HasSuffix(s, suffix string) bool
	*/
	filterstr = strings.TrimSpace(filterstr)
	if len(filterstr) > 0 {

		if strings.HasPrefix(filterstr, "or") || strings.HasPrefix(filterstr, "and") {
			strs[0] = "1=1 " + filterstr
		} else {
			strs[0] = filterstr
		}

	} else {
		strs[0] = filterstr
	}

	strs[1] = orderstr

	return strs
}

func (m *PagePms) FilterItems(children []Conditions, parent Conditions) (strs []string) {

	strs = make([]string, len(children))

	for _, v := range children {

		if len(v.Children) > 0 {
			str2 := m.FilterItems(v.Children, v)
			arr2str := strings.Join(str2, " "+v.Conjunction+" ")
			if len(arr2str) > 5 {
				strs = append(strs, arr2str)
			}
		} else {
			// strs = make([]string, len(cds))
			if v.Right == nil {
				continue
			}
			v.Left.Field = strings.ToLower(v.Left.Field)

			if v.Left.Type == "text" || v.Left.Type == "field" { // string
				strv := v.Right.(string)
				if v.Op == "equal" { // equal

					cd := "(" + v.Left.Field + " = '" + strv + "')"
					strs = append(strs, cd)
				} else if v.Op == "not_equal" { //not_equal
					cd := "(" + v.Left.Field + " <> '" + strv + "')"
					strs = append(strs, cd)
				} else if v.Op == "like" { // like
					cd := "(" + v.Left.Field + " like '%" + strv + "%')"
					strs = append(strs, cd)
				} else if v.Op == "starts_with" { // starts_with
					cd := "(" + v.Left.Field + " like '" + strv + "%')"
					strs = append(strs, cd)
				} else if v.Op == "ends_with" { //ends_with
					cd := "(" + v.Left.Field + " like '%" + strv + "')"
					strs = append(strs, cd)
				} else if v.Op == "not_like" { //not_like
					cd := "(" + v.Left.Field + " not like '%" + strv + "')"
					strs = append(strs, cd)
				} else if v.Op == "is_empty" { //is_empty
					cd := "(" + v.Left.Field + " is null)"
					strs = append(strs, cd)
				} else if v.Op == "is_not_empty" { //is_not_empty
					cd := "(" + v.Left.Field + " is not null)"
					strs = append(strs, cd)
				}
			} else if v.Left.Type == "boolean" { //bool
				bv := v.Right.(bool)
				strbool := "0"
				if bv {
					strbool = "1"
				}
				if v.Op == "equal" {
					cd := "(" + v.Left.Field + " = " + strbool + ")"
					strs = append(strs, cd)
				} else if v.Op == "not_equal" { //Lt
					cd := "(" + v.Left.Field + " <> " + strbool + ")"
					strs = append(strs, cd)
				}
			} else if v.Left.Type == "number" { //number
				// intv := v.Right.(float64)
				if v.Op == "equal" { // equal
					intv := v.Right.(string)
					cd := "(" + v.Left.Field + " = " + intv + ")"
					strs = append(strs, cd)
				} else if v.Op == "not_equal" { //not_equal
					intv := v.Right.(string)
					cd := "(" + v.Left.Field + " <> " + intv + ")"
					strs = append(strs, cd)
				} else if v.Op == "less" { // 小于
					intv := v.Right.(string)
					cd := "(" + v.Left.Field + " < " + intv + ")"
					strs = append(strs, cd)
				} else if v.Op == "less_or_equal" { // 小于或等于
					intv := v.Right.(string)
					cd := "(" + v.Left.Field + " <= " + intv + ")"
					strs = append(strs, cd)
				} else if v.Op == "greater" { //大于
					intv := v.Right.(string)
					cd := "(" + v.Left.Field + " > " + intv + ")"
					strs = append(strs, cd)
				} else if v.Op == "greater_or_equal" { //大于或等于
					intv := v.Right.(string)
					cd := "(" + v.Left.Field + " >= " + intv + ")"
					strs = append(strs, cd)
				} else if v.Op == "between" { //属于范围
					// intv := v.Right.(string)
					s := reflect.ValueOf(v.Right)
					if s.Len() != 2 {
						continue
					}
					start := s.Index(0).Interface().(string)
					end := s.Index(1).Interface().(string)

					cd := "(" + v.Left.Field + " between " + start + " and " + end + ")"
					strs = append(strs, cd)
				} else if v.Op == "not_between" { //不属于范围
					s := reflect.ValueOf(v.Right)
					if s.Len() != 2 {
						continue
					}
					start := s.Index(0).Interface().(string)
					end := s.Index(1).Interface().(string)

					cd := "(" + v.Left.Field + " not between " + start + " and " + end + ")"
					strs = append(strs, cd)
				} else if v.Op == "is_empty" { //is_empty
					cd := "(" + v.Left.Field + " is null)"
					strs = append(strs, cd)
				}
			} else if v.Left.Type == "select" || v.Left.Type == "select2" { //select
				strv := v.Right.(string)

				if v.Op == "select_equals" {
					cd := "(" + v.Left.Field + " = '" + strv + "')"
					strs = append(strs, cd)
				} else if v.Op == "select_not_equals" { //Lt
					cd := "(" + v.Left.Field + " <> '" + strv + "')"
					strs = append(strs, cd)
				} else if v.Op == "select_any_in" { //select_any_in
					cd := "(" + v.Left.Field + " like '%" + strv + "%')"
					strs = append(strs, cd)
				} else if v.Op == "select_not_any_in" { //Lt
					cd := "(" + v.Left.Field + " not like '%" + strv + "%')"
					strs = append(strs, cd)
				}
			} else if v.Left.Type == "datetime" || v.Left.Type == "time" || v.Left.Type == "date" { //datetime
				// intv := v.Right.(float64)
				if v.Op == "equal" { // equal
					intv := v.Right.(string)
					cd := "(" + v.Left.Field + " = '" + intv + "')"
					strs = append(strs, cd)
				} else if v.Op == "not_equal" { //not_equal
					intv := v.Right.(string)
					cd := "(" + v.Left.Field + " <> '" + intv + "')"
					strs = append(strs, cd)
				} else if v.Op == "less" { // 小于
					intv := v.Right.(string)
					cd := "(" + v.Left.Field + " < '" + intv + "')"
					strs = append(strs, cd)
				} else if v.Op == "less_or_equal" { // 小于或等于
					intv := v.Right.(string)
					cd := "(" + v.Left.Field + " <= '" + intv + "')"
					strs = append(strs, cd)
				} else if v.Op == "greater" { //大于
					intv := v.Right.(string)
					cd := "(" + v.Left.Field + " > '" + intv + "')"
					strs = append(strs, cd)
				} else if v.Op == "greater_or_equal" { //大于或等于
					intv := v.Right.(string)
					cd := "(" + v.Left.Field + " >= '" + intv + "')"
					strs = append(strs, cd)
				} else if v.Op == "between" { //属于范围
					// intv := v.Right.(string)
					s := reflect.ValueOf(v.Right)
					if s.Len() != 2 {
						continue
					}
					start := s.Index(0).Interface().(string)
					end := s.Index(1).Interface().(string)

					cd := "(" + v.Left.Field + " between '" + start + "' and '" + end + "')"
					strs = append(strs, cd)
				} else if v.Op == "not_between" { //不属于范围
					s := reflect.ValueOf(v.Right)
					if s.Len() != 2 {
						continue
					}
					start := s.Index(0).Interface().(string)
					end := s.Index(1).Interface().(string)

					cd := "(" + v.Left.Field + " not between '" + start + "' and '" + end + "')"
					strs = append(strs, cd)
				} else if v.Op == "is_empty" { //is_empty
					cd := "(" + v.Left.Field + " is null)"
					strs = append(strs, cd)
				}
			}
		}

	}

	return strs
}

/*

func generate() (interface{}, bool) {
	//s := []string{"123", "345", "abc"}
	//s := 123
	s := "mmm"
	return s, true
}
func test() {
	origin, ok := generate()
	if ok {
		switch reflect.TypeOf(origin).Kind() {
		case reflect.Slice, reflect.Array:
			s := reflect.ValueOf(origin)
			for i := 0; i < s.Len(); i++ {
				fmt.Println(s.Index(i))
			}
		case reflect.String:
			s := reflect.ValueOf(origin)
			fmt.Println(s.String(), "I am a string type variable.")
		case reflect.Int:
			s := reflect.ValueOf(origin)
			t := s.Int()
			fmt.Println(t, " I am a int type variable.")
		}
	}
}

*/
