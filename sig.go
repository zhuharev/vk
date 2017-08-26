package vk

import (
	"crypto/md5"
	"fmt"
	"net/url"
	"sort"
)

type Value struct {
	Key   string
	Value string
}

type Values []Value

func NewValuesFromParam(param url.Values) Values {
	vs := Values{}
	for k, v := range param {
		if len(v) > 0 {
			vs = append(vs, Value{k, v[0]})
		}
	}
	return vs
}

func (v Values) Less(i, j int) bool {
	return v[i].Key < v[j].Key
}

func (v Values) Len() int      { return len(v) }
func (v Values) Swap(i, j int) { v[i], v[j] = v[j], v[i] }

func (v Values) Sig(method, secret string) string {
	res := ""
	sort.Sort(v)
	for _, v := range v {
		res += v.Key + "=" + v.Value
	}

	str := "/" + method + "?" + res

	return fmt.Sprintf("%x", md5.Sum([]byte(str)))
}

func Sig(param url.Values, method, secret string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte("/"+method+"?"+param.Encode())))
}
