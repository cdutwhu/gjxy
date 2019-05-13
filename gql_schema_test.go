package gjxy

import "testing"

func TestGQLBuild(t *testing.T) {
	s := Str(`
type StaffPersonal {
	-RefId: String
	LocalId: String
	StateProvinceId: String
	OtherIdList: OtherIdList
}
	
type OtherIdList {
	OtherId: OtherId
}

type OtherIdList1 {
	OtherId1: OtherId1
}

type OtherIdList2 {
	OtherId2: OtherId2
}
	
type OtherId {
	-Type: String
}`)

	s = Str(SchemaMake(s, "OtherIdList2", "OtherId2", "String"))

	// s := Str("")
	// s = Str(s.GQLBuild("StaffPersonal", "RefId", "String"))
	// s = Str(s.GQLBuild("StaffPersonal", "LocalId", "String"))
	// s = Str(s.GQLBuild("Recent", "SchoolLocalId", "String"))
	// s = Str(s.GQLBuild("Recent", "LocalCampusId", "String"))
	// s = Str(s.GQLBuild("StaffPersonal", "StateProvinceId", "String"))
	// s = Str(s.GQLBuild("NAPLANClassListType", "ClassCode", "[String]"))
	// s = Str(s.GQLBuild("StaffPersonal", "OtherIdList", "OtherIdList"))

	fPln(s)
}
