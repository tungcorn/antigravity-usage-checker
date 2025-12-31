package discovery

import (
	"testing"
)

// TestParseCommandLine tests Unix command line parsing from ps aux output.
func TestParseCommandLine(t *testing.T) {
	tests := []struct {
		name        string
		cmdLine     string
		wantPort    int
		wantToken   string
		wantErr     bool
	}{
		{
			name:      "Standard format with spaces",
			cmdLine:   "user 12345 0.5 1.2 12345 67890 ? Ssl Dec29 10:00 /path/to/language_server --extension_server_port 42100 --csrf_token abc123def456ghi789xyz12",
			wantPort:  42100,
			wantToken: "abc123def456ghi789xyz12",
			wantErr:   false,
		},
		{
			name:      "Format with equals sign",
			cmdLine:   "/opt/antigravity/language_server --extension_server_port=55000 --csrf_token=my-csrf-token-12345678901234567890",
			wantPort:  55000,
			wantToken: "my-csrf-token-12345678901234567890",
			wantErr:   false,
		},
		{
			name:      "macOS style output",
			cmdLine:   "tung             98765   0.0  0.3  4567890  12345   ??  S    10:30AM   0:05.23 /Applications/Antigravity.app/Contents/MacOS/language_server --extension_server_port 8080 --csrf_token abcdefghij1234567890abcdefghij12",
			wantPort:  8080,
			wantToken: "abcdefghij1234567890abcdefghij12",
			wantErr:   false,
		},
		{
			name:      "With PID in line",
			cmdLine:   "language_server --extension_server_port 9999 --csrf_token token123456789012345678901234 PID=12345",
			wantPort:  9999,
			wantToken: "token123456789012345678901234",
			wantErr:   false,
		},
		{
			name:    "Missing port and token",
			cmdLine: "some random process without relevant args",
			wantErr: true,
		},
		{
			name:      "Port only (should still work with alternative token pattern)",
			cmdLine:   "language_server --extension_server_port 5000 --some_flag abcdefghijklmnopqrstuvwxyz123456",
			wantPort:  5000,
			wantToken: "abcdefghijklmnopqrstuvwxyz123456",
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			info, err := parseCommandLine(tt.cmdLine)
			
			if tt.wantErr {
				if err == nil {
					t.Errorf("parseCommandLine() expected error, got nil")
				}
				return
			}
			
			if err != nil {
				t.Errorf("parseCommandLine() unexpected error: %v", err)
				return
			}
			
			if info.HTTPPort != tt.wantPort {
				t.Errorf("HTTPPort = %d, want %d", info.HTTPPort, tt.wantPort)
			}
			
			if info.CSRFToken != tt.wantToken {
				t.Errorf("CSRFToken = %s, want %s", info.CSRFToken, tt.wantToken)
			}
		})
	}
}

// TestParseProcessInfoJSON tests Windows PowerShell JSON parsing.
func TestParseProcessInfoJSON(t *testing.T) {
	tests := []struct {
		name      string
		jsonInput string
		wantPID   int
		wantPort  int
		wantToken string
		wantErr   bool
	}{
		{
			name:      "Valid JSON",
			jsonInput: `{"ProcessId":12345,"CommandLine":"C:\\path\\language_server.exe --extension_server_port 42100 --csrf_token abc123def456"}`,
			wantPID:   12345,
			wantPort:  42100,
			wantToken: "abc123def456",
			wantErr:   false,
		},
		{
			name:      "JSON with escaped quotes",
			jsonInput: `{"ProcessId":99999,"CommandLine":"\"C:\\Program Files\\Antigravity\\language_server.exe\" --extension_server_port 55555 --csrf_token my-long-csrf-token-here"}`,
			wantPID:   99999,
			wantPort:  55555,
			wantToken: "my-long-csrf-token-here",
			wantErr:   false,
		},
		{
			name:      "Empty JSON",
			jsonInput: "",
			wantErr:   true,
		},
		{
			name:      "Null JSON",
			jsonInput: "null",
			wantErr:   true,
		},
		{
			name:      "Invalid JSON",
			jsonInput: "not valid json",
			wantErr:   true,
		},
		{
			name:      "Missing ProcessId",
			jsonInput: `{"ProcessId":0,"CommandLine":""}`,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			info, err := parseProcessInfoJSON(tt.jsonInput)
			
			if tt.wantErr {
				if err == nil {
					t.Errorf("parseProcessInfoJSON() expected error, got nil")
				}
				return
			}
			
			if err != nil {
				t.Errorf("parseProcessInfoJSON() unexpected error: %v", err)
				return
			}
			
			if info.PID != tt.wantPID {
				t.Errorf("PID = %d, want %d", info.PID, tt.wantPID)
			}
			
			if info.HTTPPort != tt.wantPort {
				t.Errorf("HTTPPort = %d, want %d", info.HTTPPort, tt.wantPort)
			}
			
			if info.CSRFToken != tt.wantToken {
				t.Errorf("CSRFToken = %s, want %s", info.CSRFToken, tt.wantToken)
			}
		})
	}
}

// TestParseListeningPorts tests netstat output parsing.
func TestParseListeningPorts(t *testing.T) {
	tests := []struct {
		name       string
		netstatOut string
		pid        int
		wantPorts  []int
	}{
		{
			name: "Windows netstat output",
			netstatOut: `Active Connections

  Proto  Local Address          Foreign Address        State           PID
  TCP    127.0.0.1:42100        0.0.0.0:0              LISTENING       12345
  TCP    127.0.0.1:42101        0.0.0.0:0              LISTENING       12345
  TCP    127.0.0.1:8080         0.0.0.0:0              LISTENING       99999
`,
			pid:       12345,
			wantPorts: []int{42100, 42101},
		},
		{
			name: "No matching PID",
			netstatOut: `Active Connections

  Proto  Local Address          Foreign Address        State           PID
  TCP    127.0.0.1:8080         0.0.0.0:0              LISTENING       99999
`,
			pid:       12345,
			wantPorts: []int{},
		},
		{
			name:       "Empty output",
			netstatOut: "",
			pid:        12345,
			wantPorts:  []int{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ports, err := parseListeningPorts(tt.netstatOut, tt.pid)
			if err != nil {
				t.Errorf("parseListeningPorts() unexpected error: %v", err)
				return
			}
			
			if len(ports) != len(tt.wantPorts) {
				t.Errorf("got %d ports, want %d", len(ports), len(tt.wantPorts))
				return
			}
			
			for i, port := range ports {
				if port != tt.wantPorts[i] {
					t.Errorf("port[%d] = %d, want %d", i, port, tt.wantPorts[i])
				}
			}
		})
	}
}
