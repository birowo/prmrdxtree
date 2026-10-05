package prmrdxtree

import (
	"errors"
)

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

func prms[T any](chr byte, sgmn string, val T) Sgmn[T] {
	var (
		nds  *Nds[T]
		val_ T
	)
	x := ""
	j := len(sgmn)
	for i := j - 1; i > -1; i-- {
		if sgmn[i] == ':' {
			nds = &Nds[T]{':': Sgmn[T]{
				":", x, sgmn[i+1 : j], val, nds,
			}}
			val = val_
			j = i
		}
		x = string(sgmn[i])
	}
	return Sgmn[T]{
		string(chr), x, sgmn[:j], val, nds,
	}
}

const prmsMax = 32

func (nds *Nds[T]) Cnfg(path string, val T) error {
	pathLen := len(path)
	var val_ T
	prmsI := 0
	i := 0
	for i < pathLen {
		chr := path[i]
		if chr < 32 || chr > 127 {
			return errors.New("invalid char")
		}
		if chr == ':' {
			prmsI++
			if prmsI > prmsMax {
				return errors.New("prmsNum > 32")
			}
		}
		i++
		if nds[chr].Nd == "" {
			nds[chr] = prms(chr, path[i:], val)
			println("1", path)
			return nil
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
							return errors.New("prmsNum > 32")
						}
					}
					//j < sgmnLen && i < pathLen
					x := ""
					if (j + 1) < sgmnLen {
						x = string(sgmn[j+1])
					}
					var nds_ Nds[T]
					nds_[sj] = Sgmn[T]{
						string(sj),
						x,
						sgmn[j+1:],
						nds[chr].Val,
						nds[chr].Nds,
					}
					nds_[pi] = prms(
						pi, path[i+1:], val,
					)
					nds[chr].Sgmn = sgmn[:j]
					nds[chr].Val = val_
					nds[chr].Nds = &nds_
					println("2", path)
					return nil
				}
				j++
				i++
			}
			if i == pathLen {
				if j != sgmnLen {
					//i == pathLen && j < sgmnLen
					x := ""
					if (j + 1) < sgmnLen {
						x = string(sgmn[j+1])
					}
					var nds_ Nds[T]
					nds_[sgmn[j]] = Sgmn[T]{
						string(sgmn[j]),
						x,
						sgmn[j+1:],
						nds[chr].Val,
						nds[chr].Nds,
					}
					nds[chr].Sgmn = sgmn[:j]
					nds[chr].Val = val
					nds[chr].Nds = &nds_
					println("3", path)
					return nil
				}
				//i == pathLen && j == sgmnLen
				nds[chr].Val = val
				println("4", path, nds[chr].X)
				return nil
			} else {
				if nds[chr].Nds == nil {
					//i < pathLen && j == sgmnLen
					chr_ := path[i]
					if chr_ == ':' {
						prmsI++
						if prmsI > prmsMax {
							return errors.New("prmsNum > 32")
						}
					}
					var nds_ Nds[T]
					nds_[chr_] = prms(
						chr_, path[i+1:], val,
					)
					nds[chr].Nds = &nds_
					println("5", path)
					return nil
				}
			}
		}
		nds = nds[chr].Nds
	}
	return nil
}

type params struct {
	Val [32]string
	Len int
}

func (nds *Nds[T]) Find(path string) (params, T) {
	pathLen := len(path)
	var (
		prms params
		val  T
	)
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
				if string(chr) == nds[':'].X {
					break
				}
				if isNds_ && nds_[chr].Nd == string(chr) {
					break
				}
				i++
			}
			prms.Val[prms.Len] = path[i_:i]
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
