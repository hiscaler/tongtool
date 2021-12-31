package cache

import "testing"

func TestGenerateKey(t *testing.T) {
	type User struct {
		ID   int
		Name string
	}
	type testCase struct {
		Number int
		Values []interface{}
		Key    string
	}

	testCases := []testCase{
		{1, []interface{}{1, 2, 3}, "123"},
		{2, []interface{}{0, -1, 2, 3}, "0-123"},
		{3, []interface{}{1.1, 2.12, 3.123}, "1.102.123.12"},
		{4, []interface{}{1.1, 2.12, 3.123}, "1.102.123.12"},
		{5, []interface{}{"a", "b", "c"}, "abc"},
		{6, []interface{}{"a", "b", "c", 1, 2, 3}, "abc123"},
		{7, []interface{}{true, true, false, false}, "truetruefalsefalse"},
		{8, []interface{}{[]int{1, 2, 3}}, "123"},
		{9, []interface{}{[...]int{1, 2, 3, 4}}, "1234"},
		{10, []interface{}{struct {
			Username string
			Age      int
		}{}}, "Age0Username"},
		{11, []interface{}{struct {
			Username string
			Age      int
		}{"John", 12}}, "Age12UsernameJohn"},
		{12, []interface{}{User{
			ID:   1,
			Name: "John",
		}}, "User:ID1NameJohn"},
	}
	for _, tc := range testCases {
		key := GenerateKey(tc.Values...)
		if key != tc.Key {
			t.Errorf("%d: 期待：%s 实际：%s", tc.Number, tc.Key, key)
		}
	}
}

func BenchmarkGenerateKey(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GenerateKey([]interface{}{struct {
			Username string
			Age      int
		}{"John", 12}})
	}
}
