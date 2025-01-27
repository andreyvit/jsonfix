package jsonfix

import (
	"encoding/json"
	"strings"
	"testing"
)

func Test(t *testing.T) {
	cases := []struct {
		name     string
		input    string
		expected string
		invalid  bool
	}{
		{
			name:     "trailing_comma_in_object",
			input:    `{"foo": 1, "bar": 2,}`,
			expected: `{"foo": 1, "bar": 2}`,
		},
		{
			name: "trailing_comma_in_multiline_object",
			input: `{
				"foo": 1,
				"bar": 2,
			}`,
			expected: `{
				"foo": 1,
				"bar": 2
			}`,
		},
		{
			name: "trailing_comma_in_multiline_array",
			input: `[
				{"foo": 1},
				{"bar": 2},
			]`,
			expected: `[
				{"foo": 1},
				{"bar": 2}
			]`,
		},
		{
			name:     "comments_in_multiline_object",
			input:    "// header comment\n{\n\t// one\n\t\"foo\": 1,\n\t// fubar\n\t// two\n\t\"bar\": 2,\n}\n// trailing comments\n// are allowed too\n",
			expected: "\n{\n\t\n\t\"foo\": 1,\n\t\n\t\n\t\"bar\": 2\n}\n\n\n",
		},
		{
			name: "comments_ignored_inside_strings",
			input: `
			{
				"foo": "http://example.com/", // real comment
				"bar": "http://example.com\\",// third comment
				"boz": "\"http://example.com/",// another comment
			}`,
			expected: `
			{
				"foo": "http://example.com/", <NOTHING>
				"bar": "http://example.com\\",
				"boz": "\"http://example.com/"
			}`,
		},
		{
			name: "true_false_null_in_object_values",
			input: `{
				"foo": true,
				"bar": false,
				"baz": null
			}`,
			expected: `{
				"foo": true,
				"bar": false,
				"baz": null
			}`,
		},
		{
			name:     "true_false_null_in_array_values_1",
			input:    `[true, false, null]`,
			expected: `[true, false, null]`,
		},
		{
			name:     "true_false_null_in_array_values_2",
			input:    `[true,false,null]`,
			expected: `[true,false,null]`,
		},
		{
			name:     "bare_keys",
			input:    `{foo: 1, bar: 2}`,
			expected: `{"foo": 1, "bar": 2}`,
		},
		{
			name:     "bare_keys_with_special_chars_1",
			input:    `{foo@bar: 1}`,
			expected: `{"foo@bar": 1}`,
		},
		{
			name:     "bare_keys_with_special_chars_2",
			input:    `{foo$$$bar: 1}`,
			expected: `{"foo$$$bar": 1}`,
		},
		{
			name:     "mixed_bare_and_quoted_keys",
			input:    `{foo: 1, "bar": 2, baz: 3}`,
			expected: `{"foo": 1, "bar": 2, "baz": 3}`,
		},
		{
			name: "bare_keys_with_spaces",
			input: `{
				foo: 1,
				bar: 2,
				baz: 3
			}`,
			expected: `{
				"foo": 1,
				"bar": 2,
				"baz": 3
			}`,
		},
		{
			name:     "bare_keys_with_spaces_before_colon",
			input:    `{foo   : 1}`,
			expected: `{"foo"   : 1}`,
		},
		{
			name:     "bare_true_false_null_in_object_values",
			input:    `{foo: true, bar: false, baz: null}`,
			expected: `{"foo": true, "bar": false, "baz": null}`,
		},
		{
			name:     "bare_word_in_object_value_should_remain_unchanged",
			input:    `{"foo": bar}`,
			expected: `{"foo": bar}`,
			invalid:  true,
		},
		{
			name:     "bare_word_in_array_value_should_remain_unchanged",
			input:    `[foo,bar,baz]`,
			expected: `[foo,bar,baz]`,
			invalid:  true,
		},
		{
			name:     "nested_objects",
			input:    `{foo: {bar: 1, baz: {qux: 2}}}`,
			expected: `{"foo": {"bar": 1, "baz": {"qux": 2}}}`,
		},
		{
			name:     "object_in_array",
			input:    `[{foo: 1}, {bar: 2}]`,
			expected: `[{"foo": 1}, {"bar": 2}]`,
		},
		{
			name:     "array_in_object",
			input:    `{foo: [1, {bar: 2}, 3]}`,
			expected: `{"foo": [1, {"bar": 2}, 3]}`,
		},
		{
			name:     "deeply_nested_mixed",
			input:    `{foo: [{bar: {baz: [1, {qux: 2}]}}]}`,
			expected: `{"foo": [{"bar": {"baz": [1, {"qux": 2}]}}]}`,
		},
		{
			name: "complex_nested_with_comments",
			input: `{
				// header
				foo: {
					bar: [
						{baz: 1}, // first
						{qux: 2}, // second
					],
					// nested object
					obj: {
						key: "value",
						arr: [1,2,3,],
					},
				},
			}`,
			expected: `{
				<NOTHING>
				"foo": {
					"bar": [
						{"baz": 1}, <NOTHING>
						{"qux": 2} <NOTHING>
					],
					<NOTHING>
					"obj": {
						"key": "value",
						"arr": [1,2,3]
					}
				}
			}`,
		},
		{
			name:     "unicode_in_bare_keys",
			input:    `{привет: 1, 你好: 2, γεια: 3}`,
			expected: `{"привет": 1, "你好": 2, "γεια": 3}`,
		},
		{
			name:     "empty_objects_and_arrays",
			input:    `{foo: [], bar: {}}`,
			expected: `{"foo": [], "bar": {}}`,
		},
		{
			name:     "mixed_whitespace",
			input:    "{foo:\t1,\nbar:  2,\r\nbaz:\n3}",
			expected: "{\"foo\":\t1,\n\"bar\":  2,\r\n\"baz\":\n3}",
		},
		{
			name: "special_values_in_keys",
			input: `{
				true: 1,
				false: 2,
				null: 3,
				42: 4,
				3.14: 5,
				1e10: 6
			}`,
			expected: `{
				"true": 1,
				"false": 2,
				"null": 3,
				"42": 4,
				"3.14": 5,
				"1e10": 6
			}`,
		},
		{
			name: "numbers",
			input: `{
				int: 42,
				negative: -17,
				float: 3.14159,
				exp: 1.23e-4,
				bigexp: 1.23E+10
			}`,
			expected: `{
				"int": 42,
				"negative": -17,
				"float": 3.14159,
				"exp": 1.23e-4,
				"bigexp": 1.23E+10
			}`,
		},
		{
			name: "escaped_characters",
			input: `{
				"escaped//comment": "not//comment",
				foo: "contains\"quote",
				bar: "back\\slash",
				baz: "tab\t\r\n"
			}`,
			expected: `{
				"escaped//comment": "not//comment",
				"foo": "contains\"quote",
				"bar": "back\\slash",
				"baz": "tab\t\r\n"
			}`,
		},
		{
			name:     "multiple_commas",
			input:    `[1,, 2]`,
			expected: `[1,, 2]`,
			invalid:  true,
		},
		{
			name:     "comment_characters_in_string",
			input:    `{"//": "//not//a//comment"}`,
			expected: `{"//": "//not//a//comment"}`,
		},
		{
			name: "deep_nesting_with_all_features",
			input: `{
				array: [
					{ nested: { deeper: { evenDeeper: true } } },
					{ unicode你好: ["mixed", 42, null, false] },
					// comment in array
					{ trailing: "comma", }, // another comment
					{} // empty
				],
				empty: {
					// nothing here
				}
			}`,
			expected: `{
				"array": [
					{ "nested": { "deeper": { "evenDeeper": true } } },
					{ "unicode你好": ["mixed", 42, null, false] },
					<NOTHING>
					{ "trailing": "comma" }, <NOTHING>
					{} <NOTHING>
				],
				"empty": {
					<NOTHING>
				}
			}`,
		},
		{
			name:     "bare_word_starting_with_true",
			input:    `{trueish: 1}`,
			expected: `{"trueish": 1}`,
		},
		{
			name:     "bare_word_starting_with_false",
			input:    `{falsehood: 1}`,
			expected: `{"falsehood": 1}`,
		},
		{
			name:     "bare_word_starting_with_null",
			input:    `{nullable: 1}`,
			expected: `{"nullable": 1}`,
		},
		{
			name: "non_string_keys",
			input: `{
				042: 1,
				+1: 2,
				.5: 3,
				1.: 4,
				1.2.3: 5,
				1e: 6,
				--1: 7,
				1-2: 8,
				true: 9,
				false: 10,
				null: 11
			}`,
			expected: `{
				"042": 1,
				"+1": 2,
				".5": 3,
				"1.": 4,
				"1.2.3": 5,
				"1e": 6,
				"--1": 7,
				"1-2": 8,
				"true": 9,
				"false": 10,
				"null": 11
			}`,
		},
		{
			name: "block_comments",
			input: `{
				/* this is a comment */
				"foo": 1,
				"bar": /* inline comment */ 2,
				/* multi
				   line
				   comment */ "baz": 3,
				"qux": 4 /* trailing comment */
			}`,
			expected: `{
				<NOTHING>
				"foo": 1,
				"bar":  2,
				<NOTHING>
<NOTHING>
<NOTHING> "baz": 3,
				"qux": 4 <NOTHING>
			}`,
		},
		{
			name: "block_comments_in_strings",
			input: `{
				"foo": "/* not a comment */",
				"bar": "text with /* in the middle",
				"baz": "ends with /*"
			}`,
			expected: `{
				"foo": "/* not a comment */",
				"bar": "text with /* in the middle",
				"baz": "ends with /*"
			}`,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			c.expected = strings.ReplaceAll(c.expected, "<NOTHING>", "")
			fixed := string(Bytes([]byte(c.input)))
			t.Logf("Fixed JSON = %s", fixed)
			if fixed != c.expected {
				t.Errorf("invalid fixed JSON\nwanted %q\ngot    %q", c.expected, fixed)
			}
			if c.invalid {
				return
			}
			var v interface{}
			if err := json.Unmarshal([]byte(fixed), &v); err != nil {
				t.Errorf("** %v", err)
			}
		})
	}
}
