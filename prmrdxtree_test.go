package prmrdxtree

import (
	"encoding/json"
	"fmt"
	"regexp"
	"testing"
)

func TestPrmRdxTree(t *testing.T) {
	nds := &Nds[int]{}
	type Params []string
	paths := []struct {
		Key string
		Params
		Val int
	}{
		{"ab:/cde:/fgh", Params{".AB.", ".CD."}, 1},
		{"ab:/cDe:/fgh", Params{".EF.", ".GH."}, 2},
		{"ab:cDe:fGh", Params{".IJ.", ".KL."}, 3},
		{"ab:cde:f", Params{".MN.", ".OP."}, 4},
		{"ab:cde:fghij", Params{".QR.", ".ST."}, 5},
	}
	for _, path := range paths {
		nds.Cnfg(path.Key, path.Val)
	}
	bs, _ := json.MarshalIndent(nds, "", "  ")
	println(string(regexp.MustCompile(
		`\{\},*\s+`,
	).ReplaceAll(bs, nil)))

	for _, path := range paths {
		pathPrms := replace(
			path.Key, ':', path.Params,
		)
		prms, val := nds.Find(pathPrms)
		fmt.Printf("%q %q\n", prms.Val[:prms.Len], path.Params)
		if path.Val != val {
			t.Error("val error")
		}
		for i, prm := range path.Params {
			if prm != prms.Val[i] {
				t.Error("params error")
			}
		}
		fmt.Printf(
			"key: %s\nparams: %q\nval: %d\n%s: %d\n\n",
			path.Key,
			path.Params,
			path.Val,
			pathPrms, val,
		)
	}
}

func replace(str string, chr byte, slcstr []string) string {
	l := 0
	for i := 0; i < len(slcstr); i++ {
		l += len(slcstr[i])
	}
	ret := make([]byte, len(str)+l)
	n := 0
	j := 0
	for i := 0; i < len(str); i++ {
		if str[i] == chr {
			n += copy(ret[n:], slcstr[j])
			j++
			if j == len(slcstr) {
				n += copy(ret[n:], str[i+1:])
				return string(ret[:n])
			}
			continue
		}
		ret[n] = str[i]
		n++
	}
	return string(ret[:n])
}
