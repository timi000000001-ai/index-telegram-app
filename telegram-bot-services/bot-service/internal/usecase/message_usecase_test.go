package usecase

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/mock"
)

// MockSearchUsecase is a mock type for the SearchUsecase type
type MockSearchUsecase struct {
	mock.Mock
}

// Search provides a mock function with given fields: query, page, pageSize, filter
func (_m *MockSearchUsecase) Search(query string, page int, pageSize int, filter string) ([]byte, error) {
	args := _m.Called(query, page, pageSize, filter)
	// Handle the case where the first argument is nil
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]byte), args.Error(1)
}

func (_m *MockSearchUsecase) IndexToMeilisearch(data map[string]interface{}) error {
	args := _m.Called(data)
	return args.Error(0)
}

// MockStorageRepository is a mock type for the StorageRepository type
type MockStorageRepository struct {
	mock.Mock
}

// SaveAndIndex provides a mock function with given fields: data
func (_m *MockStorageRepository) SaveAndIndex(data map[string]interface{}) error {
	args := _m.Called(data)
	return args.Error(0)
}

// SaveToPocketBase provides a mock function with given fields: data
func (_m *MockStorageRepository) SaveToPocketBase(data map[string]interface{}) error {
	args := _m.Called(data)
	return args.Error(0)
}

// IndexToMeilisearch provides a mock function with given fields: data
func (_m *MockStorageRepository) IndexToMeilisearch(data map[string]interface{}) error {
	args := _m.Called(data)
	return args.Error(0)
}

func TestHandleCallbackLogic(t *testing.T) {
	// Test cases
	testCases := []struct {
		name          string
		callbackData  string
		messageText   string
		mockSearch    bool
		mockReturn    []byte
		mockError     error
		expectedText  string
		expectedError error
	}{
		{
			name:         "Next Page Success",
			callbackData: "next_all",
			messageText:  "<b>🔍 关键字: test</b> (第 1 页 / 共 3 页)\n\n",
			mockSearch:   true,
			mockReturn:   []byte(`{"hits":[{"MESSAGE_ID": 456.0, "text":"result 2"}],"query":"test","totalPages":3,"page":2}`),
			mockError:    nil,
			expectedText: "<b>🔍 关键字: test</b> (第 2 页 / 共 3 页)\n\n<b>6. 💬 消息</b> from 未知\n<blockquote>result 2</blockquote>\n",
		},
		{
			name:         "Prev Page Success",
			callbackData: "prev_all",
			messageText:  "<b>🔍 关键字: test</b> (第 2 页 / 共 3 页)",
			mockSearch:   true,
			mockReturn:   []byte(`{"hits":[{"MESSAGE_ID": 123.0, "text": "result 1"}],"totalPages":3,"page":1}`),
			mockError:    nil,
			expectedText:  "<b>🔍 关键字: test</b> (第 1 页 / 共 3 页)\n\n<b>1. 💬 消息</b> from 未知\n<blockquote>result 1</blockquote>\n",
		},
		{
			name:          "Already on First Page",
			callbackData:  "prev_all",
			messageText:   "<b>🔍 关键字: test</b> (第 1 页 / 共 3 页)\n\n",
			mockSearch:    false,
			expectedError: errors.New("已经是第一页了"),
		},
		{
			name:          "Already on Last Page",
			callbackData:  "next_all",
			messageText:   "<b>🔍 关键字: test</b> (第 3 页 / 共 3 页)\n\n",
			mockSearch:    false,
			expectedError: errors.New("已经是最后一页了"),
		},
		{
			name:         "Filter Change",
			callbackData: "filter_group",
			messageText:  "<b>🔍 关键字: test</b> (第 1 页 / 共 1 页)",
			mockSearch:   true,
			mockReturn:   []byte(`{"hits":[{"TITLE": "Group A", "chat_username": "group_a", "TYPE": "group"}],"totalPages":1,"page":1}`),
			mockError:    nil,
			expectedText:  "<b>🔍 关键字: test</b> (第 1 页 / 共 1 页)\n\n<b>1. <a href=\"https://t.me/group_a\">Group A</a></b> 👥\n\n",
		},
		{
			name:          "Current Page No-op",
			callbackData:  "current",
			messageText:   "<b>🔍 关键字: test</b> (第 1 页 / 共 3 页)\n\n",
			mockSearch:    false,
			expectedText:  "",
			expectedError: nil,
		},
		{
			name:          "Unknown Action",
			callbackData:  "unknown_action",
			messageText:   "<b>🔍 关键字: test</b> (第 1 页 / 共 3 页)\n\n",
			mockSearch:    false,
			expectedError: errors.New("未知操作"),
		},
		{
			name:          "Cannot Parse Query",
			callbackData:  "next_all",
			messageText:   "Invalid message text",
			mockSearch:    false,
			expectedError: errors.New("无法解析查询关键字"),
		},
		{
			name:          "Cannot Parse Page",
			callbackData:  "next_all",
			messageText:   "<b>🔍 关键字: test</b> (第 a 页 / 共 b 页)\n\n",
			mockSearch:    false,
			expectedError: errors.New("无法解析页码"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockSearchUsecase := new(MockSearchUsecase)

			if tc.mockSearch {
				mockSearchUsecase.On("Search", "test", mock.AnythingOfType("int"), mock.AnythingOfType("int"), mock.AnythingOfType("string")).Return(tc.mockReturn, tc.mockError).Once()
			}

			

			mockSearchUsecase.AssertExpectations(t)
		})
	}
}