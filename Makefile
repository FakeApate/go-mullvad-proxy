AMI_SCHEMA := mullvad_ami_schema.json
RELAY_SCHEMA := mullvad_relay_schema.json
OUT_RELAY := endpoints/relay.go
OUT_AMI := endpoints/ami.go

.PHONY: generate-jsonschema build

generate-jsonschema:
	go-jsonschema --help >/dev/null 2>&1 || { printf "go-jsonschema not found.\nTo install: go install github.com/atombender/go-jsonschema@latest\nMake sure to add go to the \$$PATH\n"; exit 1; }
	go-jsonschema -p endpoints -o ${OUT_RELAY} ${RELAY_SCHEMA}
	go-jsonschema -p endpoints -o ${OUT_AMI} ${AMI_SCHEMA}