package valid

import (
	"testing"
)

func TestIsURL(t *testing.T) {
	t.Parallel()

	var tests = []struct {
		param    string
		expected bool
	}{
		{"http://foo.bar#com", true},                        // with hash tag
		{"http://foobar.com", true},                         // simple http
		{"https://foobar.com", true},                        // simple https
		{"foobar.com", true},                                // with no scheme
		{"http://foobar.coffee/", true},                     // with longer than 3 char tld
		{"http://foobar.中文网/", true},                        // with unicode tld
		{"http://foobar.org/", true},                        // allow orgs
		{"http://foobar.ORG", true},                         // allow uppercase
		{"http://foobar.org:8080/", true},                   // allow with port spec
		{"ftp://foobar.ru/", true},                          // ftp schemes are valid
		{"ftp.foo.bar", true},                               // with no scheme
		{"http://user:pass@www.foobar.com/", true},          // url with auth valid
		{"http://user:pass@www.foobar.com/path/file", true}, // url with auth and file path
		{"http://127.0.0.1/", true},                         // ip allowed
		{"http://duckduckgo.com/?q=%2F", true},              // query string
		{"http://localhost:3000/", true},                    // single hostname
		{"http://foobar.com/?foo=bar#baz=qux", true},        // query string and hashtag
		{"http://foobar.com?foo=bar", true},                 // query string
		{"http://www.xn--froschgrn-x9a.net/", true},         // -- chars
		{"http://foobar.com/a-", true},                      // - chars
		{"http://foobar.پاکستان/", true},                    // unicode tld
		{"http://foobar.c_o_m", false},                      // - in tld
		{"http://_foobar.com", false},                       // hostname starting with _
		{"http://foo_bar.com", true},                        // _ in hostname
		{"http://user:pass@foo_bar_bar.bar_foo.com", true},  // auth with _ in hostnames
		{"", true},                                        // allowing because we shoudl use "required" validator to check this
		{"xyz://foobar.com", false},                       // invalid scheme
		{".com", false},                                   // just a tld
		{"rtmp://foobar.com", false},                      // rtmp scheme
		{"http://foobar.com#baz=qux", true},               // hashtag
		{"http://foobar.com/t$-_.+!*\\'(),", true},        // allowed chars
		{"http://www.foobar.com/~foobar", true},           // tilde in path
		{"http://www.-foobar.com/", false},                // - as start of hostname
		{"http://www.foo---bar.com/", false},              // multiple ---
		{"http://r6---snnvoxuioq6.googlevideo.com", true}, // wtf ?
		{"mailto:someone@example.com", true},              // mailto scheme
		{"irc://irc.server.org/channel", false},           // irc scheme
		{"irc://#channel@network", true},                  // irc scheme
		{"/abs/test/dir", false},                          // just a path
		{"./rel/test/dir", false},                         // just a path
		{"http://foo^bar.org", false},                     // ^ in host
		{"http://foo&*bar.org", false},                    // * in host
		{"http://foo&bar.org", false},                     // & in  host
		{"http://foo bar.org", false},                     // <space> in host
		{"http://www.foo.bar.org", true},                  // 4 domain tld
		{"http://www.foo.co.uk", true},                    // 4 domain tld
		{"foo", false},                                    // single word
		{"http://.foo.com", false},                        // missing subdomain
		{"http://,foo.com", false},                        // domain starting with ,
		{",foo.com", false},                               // domain starting with ,
		{"http://myservice.:9093/", true},                 // domain ending with .
		// according to issues #62 #66
		{"https://pbs.twimg.com/profile_images/560826135676588032/j8fWrmYY_normal.jpeg", true},
		// according to #125
		{"http://prometheus-alertmanager.service.q:9093", true},
		{"aio1_alertmanager_container-63376c45:9093", true},
		{"https://www.logn-123-123.url.with.sigle.letter.d:12345/url/path/foo?bar=zzz#user", true},
		{"http://me.example.com", true},
		{"http://www.me.example.com", true},
		{"https://farm6.static.flickr.com", true},
		{"https://zh.wikipedia.org/wiki/Wikipedia:%E9%A6%96%E9%A1%B5", true},
		// According to #87
		{"http://hyphenated-host-name.example.co.in", true},
		{"http://cant-end-with-hyphen-.example.com", false},
		{"http://-cant-start-with-hyphen.example.com", false},
		{"http://www.domain-can-have-dashes.com", true},
		{"http://m.abcd.com/test.html", true},
		{"http://m.abcd.com/a/b/c/d/test.html?args=a&b=c", true},
		{"http://[::1]:9093", true},
		{"http://[::1]:909388", false},
		{"1200::AB00:1234::2552:7777:1313", false},
		{"http://[2001:db8:a0b:12f0::1]/index.html", true},
		{"http://[1200:0000:AB00:1234:0000:2552:7777:1313]", true},
		{"http://user:pass@[::1]:9093/a/b/c/?a=v#abc", true},
		{"https://127.0.0.1/a/b/c?a=v&c=11d", true},
		{"https://foo_bar.example.com", true},
		{"http://foo_bar.example.com", true},
		{"http://foo_bar_fizz_buzz.example.com", true},
		{"http://_cant_start_with_underescore", false},
		{"http://cant_end_with_underescore_", false},
		{"foo_bar.example.com", true},
		{"foo_bar_fizz_buzz.example.com", true},
		{"http://hello_world.example.com", true},
		// According to #212
		{"foo_bar-fizz-buzz:1313", true},
		{"foo_bar-fizz-buzz:13:13", false},
		{"foo_bar-fizz-buzz://1313", false},
	}
	for _, test := range tests {
		err := IsURL(test.param, nil, nil)
		if (err == nil) != test.expected {
			t.Errorf("Expected IsURL(%q) to be %v, got %v", test.param, test.expected, (err == nil))
			if err != nil {
				t.Errorf("IsURL(%q): %s", test.param, err.Error())
			}
		}
	}
}

func TestContains(t *testing.T) {
	t.Parallel()

	t.Run("missing params", func(t *testing.T) {
		if err := Contains("narf", nil, []string{}); err == nil {
			t.Error("expected range params error")
		}
	})

	t.Run("has match", func(t *testing.T) {
		if err := Contains("narf", nil, []string{"blah", "blop", "narf"}); err != nil {
			t.Error(err)
		}
	})

	t.Run("no match", func(t *testing.T) {
		if err := Contains("narfty", nil, []string{"blah", "blop", "narf"}); err == nil {
			t.Error("expected no match error")
		}
	})

}

func TestFloatRange(t *testing.T) {
	t.Parallel()

	t.Run("wrong number of params", func(t *testing.T) {
		if err := Range(1, nil, []string{}); err == nil {
			t.Error("expected range params error")
		}

		if err := Range(1, nil, []string{"1", "2", "3"}); err == nil {
			t.Error("expected range params error")
		}
	})

	t.Run("unsupported range type", func(t *testing.T) {
		if err := Range("1", nil, []string{"1", "2"}); err == nil {
			t.Error("expected unsupported range type")
		}
	})

	t.Run("float in range", func(t *testing.T) {
		if err := Range(5.2, nil, []string{"5.1", "6"}); err != nil {
			t.Error(err)
		}

		if err := Range(5.1, nil, []string{"5.1", "6"}); err != nil {
			t.Error(err)
		}

		if err := Range(6.0, nil, []string{"5.1", "6"}); err != nil {
			t.Error(err)
		}
	})

	t.Run("float not in range", func(t *testing.T) {
		if err := Range(5.0, nil, []string{"5.1", "6"}); err == nil {
			t.Error("expected a less than min error")
		}

		if err := Range(6.1, nil, []string{"5.1", "6"}); err == nil {
			t.Error("expected a greater than max error")
		}
	})
}

func TestIntRange(t *testing.T) {
	t.Parallel()

	t.Run("wrong number of params", func(t *testing.T) {
		if err := Range(1, nil, []string{}); err == nil {
			t.Error("expected range params error")
		}

		if err := Range(1, nil, []string{"1", "2", "3"}); err == nil {
			t.Error("expected range params error")
		}
	})

	t.Run("unsupported range type", func(t *testing.T) {
		if err := Range("1", nil, []string{"1", "2"}); err == nil {
			t.Error("expected unsupported range type")
		}
	})

	t.Run("in in range", func(t *testing.T) {
		if err := Range(5, nil, []string{"4", "6"}); err != nil {
			t.Error(err)
		}

		if err := Range(5, nil, []string{"5", "6"}); err != nil {
			t.Error(err)
		}

		if err := Range(6, nil, []string{"5", "6"}); err != nil {
			t.Error(err)
		}
	})

	t.Run("int not in range", func(t *testing.T) {
		if err := Range(4, nil, []string{"5", "6"}); err == nil {
			t.Error("expected a less than min error")
		}

	})
}

func TestIsWebhookURL(t *testing.T) {
	t.Parallel()

	var tests = []struct {
		param    string
		expected bool
	}{
		{"https://webhook.site/c7935814-fdc5-4c31-a04e-3b123a671228", true},            // http real domain
		{"https://mike:moose@webhook.site/c7935814-fdc5-4c31-a04e-3b123a671228", true}, // http real domain with auth
		{"http://mike:moose@webhook.site/c7935814-fdc5-4c31-a04e-3b123a671228", true},  // http real domain with auth
		{"http://mike:moose@127.0.0.1/c7935814-fdc5-4c31-a04e-3b123a671228", false},    // http with auth and reserved ip of 127.0.0.1
		{"http://mike:moose@localhost/c7935814-fdc5-4c31-a04e-3b123a671228", false},    // http with auth and reserved ip of 127.0.0.1 via hostname
	}
	for _, test := range tests {
		err := IsWebhookURL(test.param, nil, nil)
		if (err == nil) != test.expected {
			t.Errorf("Expected IsWebhookURL(%q) to be %v, got %v", test.param, test.expected, (err == nil))
			if err != nil {
				t.Errorf("IsWebhookURL(%q): %s", test.param, err.Error())
			}
		}
	}
}
