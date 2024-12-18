package entity

const UNSET_CTOR_PARAM_INDEX int = -1

type EntityField struct {
	name           string
	isIdField      bool
	ctorParanIndex int
	varAccessor    *VarAccessor
	isEnumNumber   bool
}
