package tokenparser

import (
	"fmt"
	"testing"
)

func checkField(t *testing.T, name string, test string, expected string) {
	if test != expected {
		t.Errorf("%s field mismatch: got '%s', expected '%s'\n", name, test, expected)
	}
}

func TestParse(t *testing.T) {
	testString := `[28/Jun/2013 12:54:48] example.com "GET /some/url HTTP/1.0" "Mozilla" "_session=3829834;" 28.314` + "\n"
	fmt.Println("Non strict parsing")
	fmt.Println("test pattern:", testString)

	parser := New()
	var date, time, vhost, method string

	parser.Skip('[')
	parser.UpTo(' ', &date)
	parser.Skip(' ')
	parser.UpTo(']', &time)
	parser.SkipMultiple(2)
	parser.UpTo(' ', &vhost)
	parser.Skip(' ')
	parser.SkipAny()
	parser.UpTo(' ', &method)
	parsed := parser.ParseString(testString)

	if !parsed {
		t.Error("Parsing test string returns FALSE")
	}
	checkField(t, "date", date, "28/Jun/2013")
	checkField(t, "time", time, "12:54:48")
	checkField(t, "vhost", vhost, "example.com")
	checkField(t, "method", method, "GET")

}

func TestParseStrict(t *testing.T) {
	testString := `[28/Jun/2013 12:54:48] example.com  "GET /some/url HTTP/1.0" "Mozilla" "_session=3829834;" 28.314` + "\n"
	fmt.Println("Strict parsing")
	fmt.Println("test pattern:", testString)

	parser := New()
	var date, time, vhost, method, url, browser, cookies, reqtime string

	parser.Skip('[')
	parser.UpTo(' ', &date)
	parser.Skip(' ')
	parser.UpTo(']', &time)
	parser.SkipAny()
	parser.SkipAny()
	parser.UpTo(' ', &vhost)
	parser.SkipAll(' ')
	parser.SkipAny()
	parser.UpTo(' ', &method)
	parser.Skip(' ')
	parser.UpTo(' ', &url)
	parser.Skip(' ')
	parser.SkipTo(' ')
	parser.Skip(' ')
	parser.Skip('"')
	parser.UpTo('"', &browser)
	parser.Skip('"')
	parser.SkipTo('"')
	parser.Skip('"')
	parser.UpTo('"', &cookies)
	parser.Skip('"')
	parser.Skip(' ')
	parser.UpTo('\n', &reqtime)
	parser.Skip('\n')
	parser.Strict = true
	parsed := parser.ParseString(testString)

	if !parsed {
		t.Error("Parsing test string returns FALSE")
	}
	checkField(t, "date", date, "28/Jun/2013")
	checkField(t, "time", time, "12:54:48")
	checkField(t, "vhost", vhost, "example.com")
	checkField(t, "method", method, "GET")
	checkField(t, "url", url, "/some/url")
	checkField(t, "browser", browser, "Mozilla")
	checkField(t, "cookies", cookies, "_session=3829834;")
	checkField(t, "reqtime", reqtime, "28.314")
}

func TestSearchString(t *testing.T) {
	testString := `[23/Apr/2014 00:00:48] postfix/cleanup[29385] EBDBE28A0129: message-id=<20140423200004.EBDBE28A0129@mx.example.com>` + "\n"
	fmt.Println("SearchString testing")
	fmt.Println("test pattern:", testString)

	parser := New()
	var method, pid string

	parser.SearchString("postfix")
	parser.SkipMultiple(8)
	parser.UpTo('[', &method)
	parser.Skip('[')
	parser.UpTo(']', &pid)
	parser.SkipTo('\n')
	parser.Skip('\n')
	parser.Strict = true
	parsed := parser.ParseString(testString)

	if !parsed {
		t.Error("Parsing test string returns FALSE")
	}

	checkField(t, "method", method, "cleanup")
	checkField(t, "pid", pid, "29385")
}

func BenchmarkStrictParse(b *testing.B) {
	testString := `[28/Jun/2013 12:54:48] example.com "GET /some/url HTTP/1.0" "Mozilla" "_session=3829834;" 28.314` + "\n"
	parser := New()
	var date, tm, vhost, method, url, browser, cookies, reqtime string

	parser.Skip('[')
	parser.UpTo(' ', &date)
	parser.Skip(' ')
	parser.UpTo(']', &tm)
	parser.SkipMultiple(2)
	parser.UpTo(' ', &vhost)
	parser.SkipMultiple(2)
	parser.UpTo(' ', &method)
	parser.SkipAny()
	parser.UpTo(' ', &url)
	parser.Skip(' ')
	parser.SkipTo(' ')
	parser.SkipMultiple(2)
	parser.UpTo('"', &browser)
	parser.Skip('"')
	parser.SkipTo('"')
	parser.Skip('"')
	parser.UpTo('"', &cookies)
	parser.SkipMultiple(2)
	parser.UpTo('\n', &reqtime)
	parser.Skip('\n')
	parser.Strict = true
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		pd := parser.ParseString(testString)
		if !pd {
			b.Error("Parser return FALSE")
		}
	}
}
