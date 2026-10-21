package names

import (
	"errors"
	"net/url"
	"testing"
)

func must[T any](value T, err error) T {
	if err != nil {
		panic(err)
	}
	return value
}

func checkEq[T comparable](t *testing.T, value T, expect T) {
	t.Helper()
	if value != expect {
		t.Errorf("unexpected: %v != %v", value, expect)
	}
}

func checkPanic[T any](t *testing.T, f func() T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("function did not panic")
		}
	}()

	f()
}

func checkOk[T any](t *testing.T, parse func(T) (Project, error), value T, expect Project) {
	t.Helper()
	actual, err := parse(value)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	} else if actual != expect {
		t.Errorf("unexpected result: %v, expected %v", actual, expect)
	}
}

func checkErr[T any](t *testing.T, parse func(T) (Project, error), value T, expect error) {
	t.Helper()
	_, err := parse(value)
	if !errors.Is(err, expect) {
		t.Errorf("unexpected error: %v, expected %v", err, expect)
	}
}

func TestParseProjectName(t *testing.T) {
	checkOk(t, ParseProjectName, ".index", Project{".index"})
	checkOk(t, ParseProjectName, "foo", Project{"foo"})
	checkOk(t, ParseProjectName, "foo@10", Project{"foo@10"})

	checkErr(t, ParseProjectName, ".foo", ErrInvalidProjectName)
	checkErr(t, ParseProjectName, "meow/miau", ErrInvalidProjectName)
	checkErr(t, ParseProjectName, "meow\\miau", ErrInvalidProjectName)
	checkErr(t, ParseProjectName, "meow%miau", ErrInvalidProjectName)
	checkErr(t, ParseProjectName, "meow*", ErrInvalidProjectName)
	checkErr(t, ParseProjectName, "foo@10@11", ErrInvalidProjectName)
	checkErr(t, ParseProjectName, "foo@awa", ErrInvalidProjectName)
	checkErr(t, ParseProjectName, "foo@0", ErrInvalidProjectName)
	checkErr(t, ParseProjectName, "foo@-10", ErrInvalidProjectName)
}

func TestParseProjectURL(t *testing.T) {
	checkOk(t, ParseProjectURL, must(url.Parse("")), Project{".index"})
	checkOk(t, ParseProjectURL, must(url.Parse("/")), Project{".index"})
	checkOk(t, ParseProjectURL, must(url.Parse("/foo")), Project{"foo"})
	checkOk(t, ParseProjectURL, must(url.Parse("/foo/")), Project{"foo"})

	checkErr(t, ParseProjectURL, must(url.Parse("/foo/bar")), ErrInvalidProjectURL)
	checkErr(t, ParseProjectURL, must(url.Parse("/foo*")), ErrInvalidProjectURL)
	checkErr(t, ParseProjectURL, must(url.Parse("/.index")), ErrInvalidProjectURL)
}

func TestProjectProperties(t *testing.T) {
	checkEq(t, Project{}.IsValid(), false)
	checkEq(t, Project{"foo"}.IsValid(), true)

	checkPanic(t, func() Project { return Project{}.MustValid() })
	checkEq(t, Project{"foo"}.MustValid(), Project{"foo"})

	checkPanic(t, func() string { return Project{}.String() })
	checkEq(t, Project{"foo"}.String(), "foo")
	checkEq(t, Project{".index"}.String(), ".index")

	checkPanic(t, func() bool { return Project{}.IsIndex() })
	checkEq(t, Project{"foo"}.IsIndex(), false)
	checkEq(t, Project{".index"}.IsIndex(), true)

	checkPanic(t, func() string { x, _, _ := Project{}.GetPullRequest(); return x })

	baseName, prNumber, isPreview := Project{"foo"}.GetPullRequest()
	checkEq(t, baseName, "foo")
	checkEq(t, prNumber, 0)
	checkEq(t, isPreview, false)

	baseName, prNumber, isPreview = Project{"foo@10"}.GetPullRequest()
	checkEq(t, baseName, "foo")
	checkEq(t, prNumber, 10)
	checkEq(t, isPreview, true)
}
