#!/usr/bin/env bats

@test "ct: init" {
	CONFIG_FILE="${BATS_TMPDIR}/ct.json"
	run ./ct init --config-file "${CONFIG_FILE}"
	printf '%s\n' 'output: ' "${output}" >&2
	[ $status -eq 0 ]
	[ "$(jq -r .ct.db_file ${CONFIG_FILE})" == "${BATS_TMPDIR}/ct.db" ]
	rm -f "${CONFIG_FILE}" "${BATS_TMPDIR}/ct.db"
}

@test "ct: log" {
	CONFIG_FILE="${BATS_TMPDIR}/ct.json"
	run ./ct init --config-file "${CONFIG_FILE}"
	printf '%s\n' 'output: ' "${output}" >&2
	[ $status -eq 0 ]

	run ./ct log --config-file "${CONFIG_FILE}" --metric test --value 1
	printf '%s\n' 'output: ' "${output}" >&2
	[ $status -eq 0 ]

	run ./ct log --config-file "${CONFIG_FILE}" --metric test --value 2
	printf '%s\n' 'output: ' "${output}" >&2
	[ $status -eq 0 ]

	rm -f "${CONFIG_FILE}" "${BATS_TMPDIR}/ct.db"
}

@test "ct: log with timestamp" {
	CONFIG_FILE="${BATS_TMPDIR}/ct.json"
	run ./ct init --config-file "${CONFIG_FILE}"
	printf '%s\n' 'output: ' "${output}" >&2
	[ $status -eq 0 ]

	run ./ct log --config-file "${CONFIG_FILE}" --metric test --value 1 --timestamp 2020-01-01
	printf '%s\n' 'output: ' "${output}" >&2
	[ $status -eq 0 ]

	run ./ct log --config-file "${CONFIG_FILE}" --metric test --value 2 --timestamp 2020-01-02
	printf '%s\n' 'output: ' "${output}" >&2
	[ $status -eq 0 ]

	OUTPUT_FILE="${BATS_TMPDIR}/output.txt"
	run ./ct report --config-file "${CONFIG_FILE}" --metric test
	printf '%s\n' 'output: ' "${output}" >&2
	[ $status -eq 0 ]
	echo "${output}" > "${OUTPUT_FILE}"
	run grep -q "2020-01-01" "${OUTPUT_FILE}" && grep -q "2020-01-02" "${OUTPUT_FILE}"
	[ $status -eq 0 ]

	rm -f "${CONFIG_FILE}" "${BATS_TMPDIR}/ct.db"
}

@test "ct: report" {
	CONFIG_FILE="${BATS_TMPDIR}/ct.json"
	run ./ct init --config-file "${CONFIG_FILE}"
	printf '%s\n' 'output: ' "${output}" >&2
	[ $status -eq 0 ]

	run ./ct log --config-file "${CONFIG_FILE}" --metric test --value 1
	printf '%s\n' 'output: ' "${output}" >&2
	[ $status -eq 0 ]

	run ./ct report --config-file "${CONFIG_FILE}" --metric test
	printf '%s\n' 'output: ' "${output}" >&2
	[ $status -eq 0 ]
	echo "$lines[0]" | grep -q "test 1"
	[ $status -eq 0 ]

	rm -f "${CONFIG_FILE}" "${BATS_TMPDIR}/ct.db"
}

@test "ct: configure" {
	CONFIG_FILE="${BATS_TMPDIR}/ct.json"
	run ./ct init --config-file "${CONFIG_FILE}"
	printf '%s\n' 'output: ' "${output}" >&2
	[ $status -eq 0 ]

	run ./ct configure --config-file "${CONFIG_FILE}" --metric test --frequency daily
	printf '%s\n' 'output: ' "${output}" >&2
	[ $status -eq 0 ]

	run ./ct dump --config-file "${CONFIG_FILE}" 
    [ $status -eq 0 ]
    [ $(echo "${output}" | jq -r .[].metric_name) == "test" ]
    [ $(echo "${output}" | jq -r .[].metric_config.frequency) == "daily" ]

	run ./ct log --config-file "${CONFIG_FILE}" --metric test --value 1
	printf '%s\n' 'output: ' "${output}" >&2
	[ $status -eq 0 ]

	run ./ct log --config-file "${CONFIG_FILE}" --metric test --value 2
	printf '%s\n' 'output: ' "${output}" >&2
	[ $status -eq 1 ]

	run ./ct configure --config-file "${CONFIG_FILE}" --metric test --frequency notSupportedFrequency
	printf '%s\n' 'output: ' "${output}" >&2
	[ $status -eq 1 ]

	run ./ct configure --config-file "${CONFIG_FILE}" --metric test --data-type int
	printf '%s\n' 'output: ' "${output}" >&2
	[ $status -eq 0 ]

	run ./ct configure --config-file "${CONFIG_FILE}" --metric test --data-type float
	printf '%s\n' 'output: ' "${output}" >&2
	[ $status -eq 0 ]

	rm -f "${CONFIG_FILE}" "${BATS_TMPDIR}/ct.db"
}

