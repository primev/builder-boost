package commons

//type Mode byte
//
//const (
//	UNDEFINED_mo Mode = iota
//	BUILDER
//	SEARCHER
//)

type MsgMode byte

const (
	UNDEFINED_sm MsgMode = iota
	ToBuilders
	ToSearchers
)
