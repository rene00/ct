#!/usr/bin/env bats

@test "ct: dump" {
	CONFIG_FILE="${BATS_TMPDIR}/ct.json"
	run ct init --config-file "${CONFIG_FILE}"
	printf '%s\n' 'output: ' "${output}" >&2
	[ $status -eq 0 ]

	run ct dump --config-file "${CONFIG_FILE}"
	printf '%s\n' 'output: ' "${output}" >&2
	[ $status -eq 0 ]
    [ $(echo "${output}" | jq -r .) == "[]" ]
}
