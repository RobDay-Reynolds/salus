# salus
Monitor processes and services

# Minimum Requirements
Currently Salus works with:
- Golang >= 1.8

# Example config to run checksd:
- Start checks daemon:
```bash
go run main/checksd.go -c <(cat <<EOF
{
  "checkStatusPath": "/tmp/summary.json",
  "checksPollTime": 1000000000,
  "checks": [
    {
      "CheckProperties": {
        "Address": "www.google.com",
        "Timeout": 5000000000
      },
      "Type": "icmp"
    }
  ]
}
EOF
)
```
- The checks daemon will listen on port 8080. i.e. `curl -v localhost:8080`
