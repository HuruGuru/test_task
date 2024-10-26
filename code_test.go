package main

import (
	"errors"
	"testing"
)

var (
	logError_flag        bool
	logInfo_flag         bool
	writeToDatabase_flag bool
	writeToDatabaseError error
)

func mockLogError(msg string) {
	logError_flag = true
}

func mockLogInfo(msg string) {
	logInfo_flag = true
}

func mockWriteToDatabase(doc Document) error {
	writeToDatabase_flag = true
	return writeToDatabaseError
}

func TestProcessDocument(t *testing.T) {
	logError = mockLogError
	logInfo = mockLogInfo
	writeToDatabase = mockWriteToDatabase

	tests := []struct {
		name        string
		input       []byte
		expectedRes bool
		expectedErr error
		dbErr       error
	}{
		{
			name:        "Valid doc",
			input:       []byte(`{"header": "Test header", "line_items":["test_item_1", "test_item_2", "test_item_3"]}`),
			expectedRes: true,
			expectedErr: nil,
			dbErr:       nil,
		},
		{
			name:        "Invalid doc",
			input:       []byte(`{header:"Test header", "line_items":["test_item_1", "test_item_2", test_item_3""]}`),
			expectedRes: false,
			expectedErr: errors.New("invalid character 'h' looking for beginning of object key string"),
			dbErr:       nil,
		},
		{
			name:        "Missing header",
			input:       []byte(`{"line_items":["test_item_1", "test_item_2", "test_item_3"]}`),
			expectedRes: false,
			expectedErr: errors.New("validation error"),
			dbErr:       nil,
		},
		{
			name:        "Missing line items",
			input:       []byte(`{"header": "Test header", "line_items":[]}`),
			expectedRes: false,
			expectedErr: errors.New("validation error"),
			dbErr:       nil,
		},
		{
			name:        "Database write rror",
			input:       []byte(`{"header": "Test header", "line_items":["test_item_1", "test_item_2", "test_item_3"]}`),
			expectedRes: false,
			expectedErr: errors.New("database error"),
			dbErr:       errors.New("database error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logError_flag = false
			logInfo_flag = false
			writeToDatabase_flag = false
			writeToDatabaseError = tt.dbErr

			result, err := ProcessDocument(tt.input)
			if result != tt.expectedRes {
				t.Errorf("Expected result %v, got %v", tt.expectedRes, result)
			}

			if err == nil && tt.expectedErr != nil {
				t.Errorf("Expected error %v, but got nil", tt.expectedErr)
			} else if err != nil && tt.expectedErr == nil {
				t.Errorf("Expected no error, but got %v", err)
			} else if err != nil && tt.expectedErr != nil && err.Error() != tt.expectedErr.Error() {
				t.Errorf("Expected error %v, got %v", tt.expectedErr, err)
			}
		})
	}
}
