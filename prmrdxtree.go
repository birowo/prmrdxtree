package prmrdxtree

import (
	"encoding/json"
	"fmt"
	"log"
	"regexp"
)

func prms[T any](chr byte, sgmn string, val T) Sgmn[T] {
	var nds *Nds[T]
	x := ""
	j := len(sgmn)
	for i := j - 1; i > -1; i-- {
		if sgmn[i] == ':' {
			nds = &Nds[T]{':': Sgmn[T]{
				":", x, sgmn[i+1 : j], val, nds,
			}}
			var val_ T
			val = val_
			j = i
		}
		x = string(sgmn[i])
	}
	return Sgmn[T]{
		string(chr), x, sgmn[:j], val, nds,
	}
}

type (
	Sgmn[T any] struct {
		Nd   string  `json:"Nd,omitempty"`
		X    string  `json:"X,omitempty"`
		Sgmn string  `json:"Sgmn,omitempty"`
		Val  T       `json:"Val,omitempty"`
		Nds  *Nds[T] `json:"Nds,omitempty"`
	}
	Nds[T any] [128]Sgmn[T]
)

const prmsMax = 32

func (nds *Nds[T]) Cnfg(path string, val T) {
	pathLen := len(path)
	prmsI := 0
	i := 0
	for i < pathLen {
		chr := path[i]
		if chr == ':' {
			prmsI++
			if prmsI > prmsMax {
				log.Println("prmsLen > ", prmsMax)
				return
			}
		}
		i++
		if nds[chr].Nd == "" {
			nds[chr] = prms(
				chr, path[i:], val,
			)
			println("1", path)
			return
		} else {
			sgmn := nds[chr].Sgmn
			sgmnLen := len(sgmn)
			j := 0
			for j < sgmnLen && i < pathLen {
				pi := path[i]
				sj := sgmn[j]
				if sj != pi {
					if pi == ':' {
						prmsI++
						if prmsI > prmsMax {
							log.Println("prmsLen > ", prmsMax)
							return
						}
					}
					//j < sgmnLen && i < pathLen
					var nds_ Nds[T]
					if len(sgmn[j+1:]) != 0 {
						nds_[sj] = Sgmn[T]{
							string(sj),
							string(sgmn[j+1]),
							sgmn[j+1:],
							nds[chr].Val,
							nds[chr].Nds,
						}
					} else {
						log.Println("error conflict", path)
						return
					}
					nds[chr].Sgmn = sgmn[:j]
					var val_ T
					nds[chr].Val = val_
					nds[chr].Nds = &nds_
					nds_[pi] = prms(
						pi, path[i+1:], val,
					)
					println("2", path)
					return
				}
				j++
				i++
			}
			if i == pathLen {
				if j != sgmnLen {
					//i == pathLen && j < sgmnLen
					var nds_ Nds[T]
					if sgmn[j+1:] != "" {
						nds_[sgmn[j]] = Sgmn[T]{
							string(sgmn[j]),
							string(sgmn[j+1]),
							sgmn[j+1:],
							nds[chr].Val,
							nds[chr].Nds,
						}
					} else {
						log.Println("error conflict", path)
						return
					}
					nds[chr].Sgmn = sgmn[:j]
					nds[chr].Val = val
					nds[chr].Nds = &nds_
					println("3", path)
					return
				}
				//i == pathLen && j == sgmnLen
				nds[chr].Val = val
				println("4", path, nds[chr].X)
				return
			} else {
				if nds[chr].Nds == nil {
					//i < pathLen && j == sgmnLen
					var nds_ Nds[T]
					nds[chr].Nds = &nds_
					chr := path[i]
					if chr == ':' {
						prmsI++
						if prmsI > prmsMax {
							log.Println(
								"prmsLen > ", prmsMax,
							)
							return
						}
					}
					nds_[chr] = prms(
						chr, path[i+1:], val,
					)
					println("5", path)
					return
				}
			}
		}
		nds = nds[chr].Nds
	}
}

type params struct {
	Items [32]string
	Len   int
}

func (nds *Nds[T]) Find(path string) (params, T) {
	pathLen := len(path)
	var prms params
	var val T
	i := 0
	for i < pathLen {
		chr := path[i]
		if nds[chr].Nd == string(chr) {
			sgmn_ := nds[chr].Sgmn
			i++
			i_ := i + len(sgmn_)
			if i_ <= pathLen && sgmn_ == path[i:i_] {
				if i_ != pathLen && nds[chr].Nds != nil {
					nds = nds[chr].Nds
					i = i_
					continue
				}
				//println(path, nds[chr].Val)
				return prms, nds[chr].Val
			}
		} else if nds[':'].Nd == ":" {
			nds_ := nds[':'].Nds
			isNds_ := nds_ != nil
			i_ := i
			for i < pathLen {
				chr := path[i]
				if chr == nds[':'].X[0] {
					break
				}
				if isNds_ && nds_[chr].Nd == string(chr) {
					break
				}
				i++
			}
			prms.Items[prms.Len] = path[i_:i]
			prms.Len++
			//println("param:", path[i_:i])
			sgmn := nds[':'].Sgmn
			i_ = i + len(sgmn)
			if i_ <= pathLen && nds[':'].Sgmn == path[i:i_] {
				if isNds_ && i_ != pathLen {
					nds = nds_
					i = i_
					continue
				}
				//println(path, nds[':'].Val)
				return prms, nds[':'].Val
			}
		}
		return params{}, val
	}
	return params{}, val
}
func main() {
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
		fmt.Printf("%q %q\n", prms.Items[:prms.Len], path.Params)
		for i, prm := range path.Params {
			if prm != prms.Items[i] {
				log.Println("params error")
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
