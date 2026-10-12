package prmrdxtree

import (
	"errors"
)

type (
	Sgmn[T any] struct {
		Nd   string  `json:"Node,omitempty"`
		Sgmn string  `json:"Segment,omitempty"`
		Val  T       `json:"Value,omitempty"`
		Nds  *Nds[T] `json:"Nodes,omitempty"`
	}
	Nds[T any] [128]Sgmn[T]
)

func prms[T any](chr byte, sgmn string, val T) Sgmn[T] {
	var (
		nds  *Nds[T]
		val_ T
	)
	j := len(sgmn)
	for i := j - 1; i > -1; i-- {
		if sgmn[i] == ':' {
			nds = &Nds[T]{':': Sgmn[T]{
				":", sgmn[i+1 : j], val, nds,
			}}
			val = val_
			j = i
		}
	}
	return Sgmn[T]{
		string(chr), sgmn[:j], val, nds,
	}
}

const prmsMax = 32

func (nds *Nds[T]) Cnfg(path string, val T) error {
	pathLen := len(path)
	prmsI := 0
	for _, chr := range path {
		if chr < 32 || chr > 127 {
			return errors.New("invalid char")
		}
		if chr == ':' {
			prmsI++
			if prmsI > prmsMax {
				return errors.New("prmsNum > 32")
			}
		}
	}
	var val_ T
	i := 0
	for i < pathLen {
		chr := path[i]
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
					var nds_ Nds[T]
					nds_[sj] = Sgmn[T]{
						string(sj),
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
					var nds_ Nds[T]
					nds_[sgmn[j]] = Sgmn[T]{
						string(sgmn[j]),
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
				println("4", path)
				return nil
			} else {
				if nds[chr].Nds == nil {
					//i < pathLen && j == sgmnLen
					chr_ := path[i]
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
			sgmn := nds[':'].Sgmn
			sgmnLen := len(sgmn)
			i_ := i
			for i < pathLen {
				chr := path[i]
				if sgmnLen != 0 && chr == sgmn[0] {
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
			i_ = i + sgmnLen
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
