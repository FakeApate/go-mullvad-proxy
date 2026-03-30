AMI_SCHEMA := mullvad_ami_schema.json
RELAY_SCHEMA := mullvad_relay_schema.json
OUT_RELAY := relay.go
OUT_AMI := ami.go

.PHONY: generate-jsonschema build

generate-jsonschema:
	go-jsonschema --help >/dev/null 2>&1 || { printf "go-jsonschema not found.\nTo install: go install github.com/atombender/go-jsonschema@latest\nMake sure to add go to the \$$PATH\n"; exit 1; }
	go-jsonschema -p mullvad -o ${OUT_RELAY} ${RELAY_SCHEMA}
	go-jsonschema -p mullvad -o ${OUT_AMI} ${AMI_SCHEMA}

build:
	go build .