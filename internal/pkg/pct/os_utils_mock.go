package pct

import (
	"io/fs"
	"os"
)

type mockStatResponse struct {
	expectedName string
	mockError    error
}

type fileInfo = fs.FileInfo

var fileInfoMock fileInfo

var mockStatResponses []mockStatResponse
var mockIsNotExistResponse bool

type osUtilHelpersImplMock struct{}

func (osUtilHelpersImplMock) Stat(name string) (os.FileInfo, error) {
	mockIsNotExistResponse = false
	for _, mockStatResponse := range mockStatResponses {
		if name == mockStatResponse.expectedName {
			if mockStatResponse.mockError != nil {
				mockIsNotExistResponse = true
			}
			return fileInfoMock, mockStatResponse.mockError
		}
	}
	return nil, nil
}

func (osUtilHelpersImplMock) IsNotExist(err error) bool {
	return mockIsNotExistResponse
}
