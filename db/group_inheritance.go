package db

type GroupInheritance struct {
	ParentGroupID uint `json:"parent_group_id"`
	ChildGroupID  uint `json:"child_group_id"`
}

func FindGroupInheritances(group uint) (inheritances []GroupInheritance) {
	Database.Where("parent_group_id = ?", group).Find(&inheritances)
	return
}
