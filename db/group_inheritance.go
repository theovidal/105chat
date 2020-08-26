package db

// GroupInheritance defines a relation between a child group and a parent group that'll copy the
// permissions of the child
type GroupInheritance struct {
	// Identifier of the inheritance, Twitter snowflake
	ID uint `json:"-" gorm:"primary_key,AUTO_INCREMENT"`
	// Identifier of the parent group, Twitter snowflake
	ParentGroupID uint `json:"parent_group_id"`
	// Identifier of the child group, Twitter snowflake
	ChildGroupID uint `json:"child_group_id"`
}

// FindGroupInheritances returns all the inheritances of a group
func FindGroupInheritances(group uint) (inheritances []GroupInheritance) {
	Database.Where("parent_group_id = ?", group).Find(&inheritances)
	return
}

// AppendGroupInheritances appends all the inheritances of a group to its structure
func AppendGroupInheritances(group *Group) {
	for _, inheritance := range FindGroupInheritances(group.ID) {
		group.Inheritances = append(group.Inheritances, inheritance.ChildGroupID)
	}
}
