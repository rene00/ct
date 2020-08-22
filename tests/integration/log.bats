#!/usr/bin/env bats

@test "ct: log" {
	CONFIG_FILE="${BATS_TMPDIR}/ct.json"
	run ct init --config-file "${CONFIG_FILE}"
	printf '%s\n' 'output: ' "${output}" >&2
	[ $status -eq 0 ]

	run ct log --config-file "${CONFIG_FILE}" --metric test --value 1
	printf '%s\n' 'output: ' "${output}" >&2
	[ $status -eq 0 ]

	rm -f "${CONFIG_FILE}" "${BATS_TMPDIR}/ct.db"
}

@test "ct: log default frequency" {
	CONFIG_FILE="${BATS_TMPDIR}/ct.json"
	run ct init --config-file "${CONFIG_FILE}"
	printf '%s\n' 'output: ' "${output}" >&2
	[ $status -eq 0 ]

	run ct log --config-file "${CONFIG_FILE}" --metric test --value 1
	printf '%s\n' 'output: ' "${output}" >&2
	[ $status -eq 0 ]

	run ct dump --config-file "${CONFIG_FILE}" 
    [ $status -eq 0 ]
	printf '%s\n' 'output: ' "${output}" >&2
    [ $(echo "${output}" | jq -r .[].metric_config.frequency) == "daily" ]

	rm -f "${CONFIG_FILE}" "${BATS_TMPDIR}/ct.db"
}

@test "ct: log with timestamp" {
	CONFIG_FILE="${BATS_TMPDIR}/ct.json"
	run ct init --config-file "${CONFIG_FILE}"
	printf '%s\n' 'output: ' "${output}" >&2
	[ $status -eq 0 ]

	run ct log --config-file "${CONFIG_FILE}" --metric test --value 1 --timestamp 2020-01-01
	printf '%s\n' 'output: ' "${output}" >&2
	[ $status -eq 0 ]

	run ct dump --config-file "${CONFIG_FILE}"
	printf '%s\n' 'output: ' "${output}" >&2
	[ $status -eq 0 ]
    [ $(echo "${output}" | jq -r .[].metric_data[].value) == "1" ]
    [ $(echo "${output}" | jq -r .[].metric_data[].timestamp) == "2020-01-01T00:00:00Z" ]

	rm -f "${CONFIG_FILE}" "${BATS_TMPDIR}/ct.db"
}


@test "ct: log wrong type" {
	CONFIG_FILE="${BATS_TMPDIR}/ct.json"
	run ct init --config-file "${CONFIG_FILE}"
	printf '%s\n' 'output: ' "${output}" >&2

	[ $status -eq 0 ]
	run ct log --config-file "${CONFIG_FILE}" --metric test --value 2.1 --timestamp 2020-01-01
	printf '%s\n' 'output: ' "${output}" >&2
	[ $status -eq 0 ]

	run ct log --config-file "${CONFIG_FILE}" --metric test --value foo --timestamp 2020-01-02
	printf '%s\n' 'output: ' "${output}" >&2
	[ $status -eq 1 ]

	rm -f "${CONFIG_FILE}" "${BATS_TMPDIR}/ct.db"
}


@test "ct: log with yes-no" {
	CONFIG_FILE="${BATS_TMPDIR}/ct.json"
	run ct init --config-file "${CONFIG_FILE}"
	printf '%s\n' 'output: ' "${output}" >&2
	[ $status -eq 0 ]

	run ct configure --config-file "${CONFIG_FILE}" --metric test --data-type bool
	printf '%s\n' 'output: ' "${output}" >&2
	[ $status -eq 0 ]

	run ct log --config-file "${CONFIG_FILE}" --metric test --value true
	printf '%s\n' 'output: ' "${output}" >&2
	[ $status -eq 0 ]

	rm -f "${CONFIG_FILE}" "${BATS_TMPDIR}/ct.db"
}

@test "ct: log exceed daily frequency" {
	CONFIG_FILE="${BATS_TMPDIR}/ct.json"
	run ct init --config-file "${CONFIG_FILE}"
	printf '%s\n' 'output: ' "${output}" >&2
	[ $status -eq 0 ]

	run ct log --config-file "${CONFIG_FILE}" --metric test --value 1 --timestamp 2020-01-01
	printf '%s\n' 'output: ' "${output}" >&2
	[ $status -eq 0 ]

	run ct configure --config-file "${CONFIG_FILE}" --metric test --frequency daily
	printf '%s\n' 'output: ' "${output}" >&2
	[ $status -eq 0 ]

	run ct log --config-file "${CONFIG_FILE}" --metric test --value 2.1 --timestamp 2020-01-01
	printf '%s\n' 'output: ' "${output}" >&2
	[ $status -eq 1 ]

	rm -f "${CONFIG_FILE}" "${BATS_TMPDIR}/ct.db"
}

@test "ct: log exceed daily frequency with data-type" {
	CONFIG_FILE="${BATS_TMPDIR}/ct.json"
	run ct init --config-file "${CONFIG_FILE}"
	printf '%s\n' 'output: ' "${output}" >&2
	[ $status -eq 0 ]

	run ct log --config-file "${CONFIG_FILE}" --metric test --value 1 --timestamp 2020-01-01
	printf '%s\n' 'output: ' "${output}" >&2
	[ $status -eq 0 ]

	run ct configure --config-file "${CONFIG_FILE}" --metric test --frequency daily --data-type float
	printf '%s\n' 'output: ' "${output}" >&2
	[ $status -eq 0 ]

	run ct log --config-file "${CONFIG_FILE}" --metric test --value 2.1 --timestamp 2020-01-02
	printf '%s\n' 'output: ' "${output}" >&2
	[ $status -eq 0 ]

	run ct log --config-file "${CONFIG_FILE}" --metric test --value 2.1 --timestamp 2020-01-01
	printf '%s\n' 'output: ' "${output}" >&2
	[ $status -eq 1 ]
    [ ${lines[0]} == "Error: Already logged metric within frequency" ]

	rm -f "${CONFIG_FILE}" "${BATS_TMPDIR}/ct.db"
}
