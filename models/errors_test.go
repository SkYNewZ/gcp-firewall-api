package models

import (
	"fmt"
	"net/http"
	"testing"

	"google.golang.org/api/googleapi"
)

type TestCase struct {
	Title    string
	Expected interface{}
	Got      interface{}
}

func TestGoogleApplicationError(t *testing.T) {
	expectedError := &googleapi.Error{
		Code:    1234,
		Message: "dummy",
	}

	// Execute function
	testedError := NewGoogleApplicationError(expectedError)

	suite := []TestCase{
		TestCase{
			Title:    "Error should not be nil",
			Expected: false,
			Got:      testedError == nil,
		},
		TestCase{
			Title:    "Error code should be identical",
			Expected: expectedError.Code,
			Got:      testedError.Code,
		},
		TestCase{
			Title:    "Error message shoud be identical",
			Expected: fmt.Sprintf("Google error: %s", expectedError.Message),
			Got:      testedError.Message,
		},
		TestCase{
			Title:    "Error() method should return error as JSON",
			Expected: fmt.Sprintf(`{"code":%d,"message":"Google error: %s"}`, expectedError.Code, expectedError.Message),
			Got:      testedError.Error(),
		},
	}

	// Launch test
	for _, suiteCase := range suite {
		t.Run(suiteCase.Title, func(t *testing.T) {
			if suiteCase.Expected != suiteCase.Got {
				t.Errorf("Got '%v' want '%v'", suiteCase.Got, suiteCase.Expected)
			}
		})
	}
}

func TestBadTokenError(t *testing.T) {
	expectedError1 := ApplicationError{
		Code:    http.StatusBadRequest,
		Message: "Invalid Bearer token. Please make sure you are using 'gcloud auth print-identity-token'",
	}

	expectedError2 := ApplicationError{
		Code:    http.StatusBadRequest,
		Message: "Invalid Bearer token. Please make sure you are using 'gcloud auth print-identity-token': foo",
	}

	// Execute function
	testedError1 := NewBadTokenError()
	testedError2 := NewBadTokenError("foo")

	suite := []TestCase{
		TestCase{
			Title:    "Error should not be nil",
			Expected: false,
			Got:      testedError1 == nil,
		},
		TestCase{
			Title:    "Error should not be nil",
			Expected: false,
			Got:      testedError2 == nil,
		},
		TestCase{
			Title:    "Error code should be identical",
			Expected: expectedError1.Code,
			Got:      testedError1.Code,
		},
		TestCase{
			Title:    "Error code should be identical",
			Expected: expectedError2.Code,
			Got:      testedError2.Code,
		},
		TestCase{
			Title:    "Error message shoud be identical",
			Expected: expectedError1.Message,
			Got:      testedError1.Message,
		},
		TestCase{
			Title:    "Error message shoud be identical",
			Expected: expectedError2.Message,
			Got:      testedError2.Message,
		},
		TestCase{
			Title:    "Error() method should return error as JSON",
			Expected: fmt.Sprintf(`{"code":%d,"message":"%s"}`, expectedError1.Code, expectedError1.Message),
			Got:      testedError1.Error(),
		},
		TestCase{
			Title:    "Error() method should return error as JSON",
			Expected: fmt.Sprintf(`{"code":%d,"message":"%s"}`, expectedError2.Code, expectedError2.Message),
			Got:      testedError2.Error(),
		},
	}

	// Launch test
	for _, suiteCase := range suite {
		t.Run(suiteCase.Title, func(t *testing.T) {
			if suiteCase.Expected != suiteCase.Got {
				t.Errorf("Got '%v' want '%v'", suiteCase.Got, suiteCase.Expected)
			}
		})
	}
}

func TestForbiddenError(t *testing.T) {
	expectedError1 := ApplicationError{
		Code:    http.StatusForbidden,
		Message: http.StatusText(http.StatusForbidden),
	}

	expectedError2 := ApplicationError{
		Code:    http.StatusForbidden,
		Message: "foo",
	}

	// Execute function
	testedError1 := NewForbiddenError()
	testedError2 := NewForbiddenError("foo")

	suite := []TestCase{
		TestCase{
			Title:    "Error should not be nil",
			Expected: false,
			Got:      testedError1 == nil,
		},
		TestCase{
			Title:    "Error should not be nil",
			Expected: false,
			Got:      testedError2 == nil,
		},
		TestCase{
			Title:    "Error code should be identical",
			Expected: expectedError1.Code,
			Got:      testedError1.Code,
		},
		TestCase{
			Title:    "Error code should be identical",
			Expected: expectedError2.Code,
			Got:      testedError2.Code,
		},
		TestCase{
			Title:    "Error message shoud be identical",
			Expected: expectedError1.Message,
			Got:      testedError1.Message,
		},
		TestCase{
			Title:    "Error message shoud be identical",
			Expected: expectedError2.Message,
			Got:      testedError2.Message,
		},
		TestCase{
			Title:    "Error() method should return error as JSON",
			Expected: fmt.Sprintf(`{"code":%d,"message":"%s"}`, expectedError1.Code, expectedError1.Message),
			Got:      testedError1.Error(),
		},
		TestCase{
			Title:    "Error() method should return error as JSON",
			Expected: fmt.Sprintf(`{"code":%d,"message":"%s"}`, expectedError2.Code, expectedError2.Message),
			Got:      testedError2.Error(),
		},
	}

	// Launch test
	for _, suiteCase := range suite {
		t.Run(suiteCase.Title, func(t *testing.T) {
			if suiteCase.Expected != suiteCase.Got {
				t.Errorf("Got '%v' want '%v'", suiteCase.Got, suiteCase.Expected)
			}
		})
	}
}

func TestNewInternalError(t *testing.T) {
	expectedError := ApplicationError{
		Code:    http.StatusInternalServerError,
		Message: http.StatusText(http.StatusInternalServerError),
	}

	// Execute function
	testedError := NewInternalError()

	suite := []TestCase{
		TestCase{
			Title:    "Error should not be nil",
			Expected: false,
			Got:      testedError == nil,
		},
		TestCase{
			Title:    "Error code should be identical",
			Expected: expectedError.Code,
			Got:      testedError.Code,
		},
		TestCase{
			Title:    "Error message shoud be identical",
			Expected: expectedError.Message,
			Got:      testedError.Message,
		},
		TestCase{
			Title:    "Error() method should return error as JSON",
			Expected: fmt.Sprintf(`{"code":%d,"message":"%s"}`, expectedError.Code, expectedError.Message),
			Got:      testedError.Error(),
		},
	}

	// Launch test
	for _, suiteCase := range suite {
		t.Run(suiteCase.Title, func(t *testing.T) {
			if suiteCase.Expected != suiteCase.Got {
				t.Errorf("Got '%v' want '%v'", suiteCase.Got, suiteCase.Expected)
			}
		})
	}
}

func TestNewMethodNotAllowedError(t *testing.T) {
	expectedError := ApplicationError{
		Code:    http.StatusMethodNotAllowed,
		Message: http.StatusText(http.StatusMethodNotAllowed),
	}

	// Execute function
	testedError := NewMethodNotAllowedError()

	suite := []TestCase{
		TestCase{
			Title:    "Error should not be nil",
			Expected: false,
			Got:      testedError == nil,
		},
		TestCase{
			Title:    "Error code should be identical",
			Expected: expectedError.Code,
			Got:      testedError.Code,
		},
		TestCase{
			Title:    "Error message shoud be identical",
			Expected: expectedError.Message,
			Got:      testedError.Message,
		},
		TestCase{
			Title:    "Error() method should return error as JSON",
			Expected: fmt.Sprintf(`{"code":%d,"message":"%s"}`, expectedError.Code, expectedError.Message),
			Got:      testedError.Error(),
		},
	}

	// Launch test
	for _, suiteCase := range suite {
		t.Run(suiteCase.Title, func(t *testing.T) {
			if suiteCase.Expected != suiteCase.Got {
				t.Errorf("Got '%v' want '%v'", suiteCase.Got, suiteCase.Expected)
			}
		})
	}
}

func TestNewNotFoundError(t *testing.T) {
	expectedError := ApplicationError{
		Code:    http.StatusNotFound,
		Message: http.StatusText(http.StatusNotFound),
	}

	// Execute function
	testedError := NewNotFoundError()

	suite := []TestCase{
		TestCase{
			Title:    "Error should not be nil",
			Expected: false,
			Got:      testedError == nil,
		},
		TestCase{
			Title:    "Error code should be identical",
			Expected: expectedError.Code,
			Got:      testedError.Code,
		},
		TestCase{
			Title:    "Error message shoud be identical",
			Expected: expectedError.Message,
			Got:      testedError.Message,
		},
		TestCase{
			Title:    "Error() method should return error as JSON",
			Expected: fmt.Sprintf(`{"code":%d,"message":"%s"}`, expectedError.Code, expectedError.Message),
			Got:      testedError.Error(),
		},
	}

	// Launch test
	for _, suiteCase := range suite {
		t.Run(suiteCase.Title, func(t *testing.T) {
			if suiteCase.Expected != suiteCase.Got {
				t.Errorf("Got '%v' want '%v'", suiteCase.Got, suiteCase.Expected)
			}
		})
	}
}

func TestNewBadRequestError(t *testing.T) {
	expectedError1 := ApplicationError{
		Code:    http.StatusBadRequest,
		Message: http.StatusText(http.StatusBadRequest),
	}

	expectedError2 := ApplicationError{
		Code:    http.StatusBadRequest,
		Message: "foo",
	}

	// Execute function
	testedError1 := NewBadRequestError()
	testedError2 := NewBadRequestError("foo")

	suite := []TestCase{
		TestCase{
			Title:    "Error should not be nil",
			Expected: false,
			Got:      testedError1 == nil,
		},
		TestCase{
			Title:    "Error should not be nil",
			Expected: false,
			Got:      testedError2 == nil,
		},
		TestCase{
			Title:    "Error code should be identical",
			Expected: expectedError1.Code,
			Got:      testedError1.Code,
		},
		TestCase{
			Title:    "Error code should be identical",
			Expected: expectedError2.Code,
			Got:      testedError2.Code,
		},
		TestCase{
			Title:    "Error message shoud be identical",
			Expected: expectedError1.Message,
			Got:      testedError1.Message,
		},
		TestCase{
			Title:    "Error message shoud be identical",
			Expected: expectedError2.Message,
			Got:      testedError2.Message,
		},
		TestCase{
			Title:    "Error() method should return error as JSON",
			Expected: fmt.Sprintf(`{"code":%d,"message":"%s"}`, expectedError1.Code, expectedError1.Message),
			Got:      testedError1.Error(),
		},
		TestCase{
			Title:    "Error() method should return error as JSON",
			Expected: fmt.Sprintf(`{"code":%d,"message":"%s"}`, expectedError2.Code, expectedError2.Message),
			Got:      testedError2.Error(),
		},
	}

	// Launch test
	for _, suiteCase := range suite {
		t.Run(suiteCase.Title, func(t *testing.T) {
			if suiteCase.Expected != suiteCase.Got {
				t.Errorf("Got '%v' want '%v'", suiteCase.Got, suiteCase.Expected)
			}
		})
	}
}
