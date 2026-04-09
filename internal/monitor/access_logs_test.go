package monitor

import "testing"

func TestParseAccessLogLines(t *testing.T) {
	lines := []string{
		"INFO[0016] [591505516 0ms] inbound/http[http-in]: inbound connection from 127.0.0.1:42370",
		"INFO[0016] [591505516 0ms] inbound/http[http-in]: inbound connection to cp.cloudflare.com:80",
		"INFO[0016] outbound/pool[proxy-pool]: → cp.cloudflare.com:80 ⇒ BWG_LA_HY2 [tcp]",
		"ERROR[0017] [591505516 165ms] inbound/http[http-in]: process connection from 127.0.0.1:42370: read http request: use of closed network connection",
	}

	logs := parseAccessLogLines(lines, 200)
	if len(logs) != 1 {
		t.Fatalf("expected 1 merged log, got %d", len(logs))
	}

	got := logs[0]
	if got.SourceIP != "127.0.0.1" {
		t.Fatalf("unexpected source ip: %q", got.SourceIP)
	}
	if got.Target != "cp.cloudflare.com:80" {
		t.Fatalf("unexpected target: %q", got.Target)
	}
	if got.OutboundNode != "BWG_LA_HY2" {
		t.Fatalf("unexpected outbound node: %q", got.OutboundNode)
	}
	if got.Status != "error" {
		t.Fatalf("unexpected status: %q", got.Status)
	}
	if got.Error == "" {
		t.Fatal("expected error message")
	}
}

func TestParseAccessLogLinesLimit(t *testing.T) {
	lines := []string{
		"INFO[0016] [1 0ms] inbound/http[http-in]: inbound connection from 127.0.0.1:42370",
		"INFO[0016] [1 0ms] inbound/http[http-in]: inbound connection to one.example:80",
		"INFO[0016] [2 0ms] inbound/http[http-in]: inbound connection from 127.0.0.1:42371",
		"INFO[0016] [2 0ms] inbound/http[http-in]: inbound connection to two.example:80",
		"INFO[0016] [3 0ms] inbound/http[http-in]: inbound connection from 127.0.0.1:42372",
		"INFO[0016] [3 0ms] inbound/http[http-in]: inbound connection to three.example:80",
	}
	logs := parseAccessLogLines(lines, 2)
	if len(logs) != 2 {
		t.Fatalf("expected 2 logs, got %d", len(logs))
	}
	if logs[0].Target != "two.example:80" || logs[1].Target != "three.example:80" {
		t.Fatalf("unexpected targets after limit: %#v", logs)
	}
}
