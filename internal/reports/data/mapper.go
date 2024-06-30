package data

import "github.com/hop-/golog"

func GetDataTypeOf(name string) any {
	golog.Debug("data type", name)
	switch name {
	case "RiskAssessmentData":
		return new(RiskAssessmentData)
	case "TestData":
		return new(TestData)
	case "SnapshotData":
		return new(SnapshotData)
	case "DiffData":
		return new(DiffData)
	default:
		return new(GeneralData)
	}
}
