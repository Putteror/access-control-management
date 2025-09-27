package common

var ACCESS_RECORD_TYPE = []string{
	"in",
	"out",
}

var ACCESS_RECORD_RESULT = []string{
	"success",
	"failed",
	"unknown",
}

func ValidateAccessRecordType(accessRecordType string) bool {
	for _, v := range ACCESS_RECORD_TYPE {
		if v == accessRecordType {
			return true
		}
	}
	return false
}

func ValidateAccessRecordResult(accessRecordResult string) bool {
	for _, v := range ACCESS_RECORD_RESULT {
		if v == accessRecordResult {
			return true
		}
	}
	return false
}
