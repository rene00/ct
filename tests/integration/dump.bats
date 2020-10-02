#!/usr/bin/env bats

@test "ct: dump" {
	CONFIG_FILE="${BATS_TMPDIR}/ct${RANDOM}.json"
	DB_FILE="${BATS_TMPDIR}/ct${RANDOM}.db"
	run ct init --config-file "${CONFIG_FILE}" --db-file "${DB_FILE}"
	printf '%s\n' 'output: ' "${output}" >&2
	[ $status -eq 0 ]

	run ct dump --config-file "${CONFIG_FILE}"
	printf '%s\n' 'output: ' "${output}" >&2
	[ $status -eq 0 ]
    [ $(echo "${output}" | jq -r .metrics) == "null" ]
    [ $(echo "${output}" | jq -r .configs) == "null" ]
    [ $(echo "${output}" | jq -r .logs) == "null" ]

	run ct metric create --config-file "${CONFIG_FILE}" --metric-name test --data-type int
	printf '%s\n' 'output: ' "${output}" >&2
	[ $status -eq 0 ]

	run ct log create --config-file "${CONFIG_FILE}" --metric-name test --metric-value 1
	printf '%s\n' 'output: ' "${output}" >&2
	[ $status -eq 0 ]

	run ct dump --config-file "${CONFIG_FILE}"
	printf '%s\n' 'output: ' "${output}" >&2
	[ $status -eq 0 ]
    [ $(echo "${output}" | jq -r .metrics[0].name) == "test" ]
    [ $(echo "${output}" | jq -r .logs[0].value) == "1" ]

    rm -f "${CONFIG_FILE}" "${DB_FILE}"
}
