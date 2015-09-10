package html2text

import (
	"testing"
)

func TestText(t *testing.T) {
	testCases := []struct {
		input    string
		expected string
	}{
		{
			`<li>
	<a href="/new" data-ga-click="Header, create new repository, icon:repo"><span class="octicon octicon-repo"></span> New repository</a>
</li>`,
			`* /new`,
		},
		{
			`hi

			<br>

	hello <a href="https://google.com">google</a>
	<br><br>
	test<p>List:</p>

	<ul>
		<li><a href="foo">Foo</a></li>
		<li><a href="http://www.microshwhat.com/bar/soapy">Barsoap</a></li>
        <li>Baz</li>
	</ul>
`,
			`hi
hello https://google.com

test

List:

* foo
* http://www.microshwhat.com/bar/soapy
* Baz`,
		},
		// Malformed input html.
		{
			`hi

			hello <a href="https://google.com">google</a>

			test<p>List:</p>

			<ul>
				<li><a href="foo">Foo</a>
				<li><a href="/
		bar/baz">Bar</a>
		        <li>Baz</li>
			</ul>
		`,
			`hi hello https://google.com test

List:

* foo
* /\n\t\tbar/baz
* Baz`,
		},
		// preformatted
		{
			`<div><pre>
this

should

stay</pre></div>`,
			`this

should

stay`,
		},
		// code
		{
			`<pre><code>
this

	should

stay</code></pre>`,
			`this

	should

stay`,
		},
		// list
		{
			`<ul>
			<li>one</li>
			<li>two</li>
			<li>three</li>
			</ul>

			<ol>
			<li>one</li>
			<li>two</li>
			<li>three</li>
			</ol>`,
			`* one
* two
* three
1. one
2. two
3. three`,
		},
		// new lines
		{
			`<p>hello</p>


			<br>

			<p>hello</p>`,
			`hello

hello`,
		},
	}

	for _, testCase := range testCases {
		text, err := FromString(testCase.input)
		if err != nil {
			t.Error(err)
		}
		if testCase.expected != text {
			t.Errorf("Input did not match expression\nExpected:\n>>>>\n%s\n<<<<\n\nGenerated:\n>>>>\n%s\n<<<<\n\n", testCase.expected, text)
		}
	}
}
