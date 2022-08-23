package check

func VerifyStudentID(studentID string) bool {
	return len(studentID) == 8
}
