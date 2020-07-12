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

	run ./ct dump --config-file "${CONFIG_FILE}"
	printf '%s\n' 'output: ' "${output}" >&2
	[ $status -eq 0 ]
    [ $(echo "${output}" | jq -r .[].metric_data[].value) == "1" ]
    [ $(echo "${output}" | jq -r .[].metric_data[].timestamp) == "2020-01-01T00:00:00Z" ]

	rm -f "${CONFIG_FILE}" "${BATS_TMPDIR}/ct.db"
}

@test "ct: report all" {
	CONFIG_FILE="${BATS_TMPDIR}/ct.json"
	run ./ct init --config-file "${CONFIG_FILE}"
	printf '%s\n' 'output: ' "${output}" >&2
	[ $status -eq 0 ]

	run ./ct log --config-file "${CONFIG_FILE}" --metric test1 --value 1 --timestamp 2020-01-01
	printf '%s\n' 'output: ' "${output}" >&2
	[ $status -eq 0 ]

	run ./ct log --config-file "${CONFIG_FILE}" --metric test2 --value 1 --timestamp 2020-01-01
	printf '%s\n' 'output: ' "${output}" >&2
	[ $status -eq 0 ]

	run ./ct report --config-file "${CONFIG_FILE}" --report-type=all --metrics test1
	printf '%s\n' 'output: ' "${output}" >&2
	[ $status -eq 0 ]
	echo "${output}" | grep -q "test1"
	[ $status -eq 0 ]

	run ./ct report --config-file "${CONFIG_FILE}" --report-type=all --metrics test2
	printf '%s\n' 'output: ' "${output}" >&2
	[ $status -eq 0 ]
	echo "${output}" | grep -q "test2"
	[ $status -eq 0 ]

	run ./ct report --config-file "${CONFIG_FILE}" --report-type=all --metrics test1,test2
	printf '%s\n' 'output: ' "${output}" >&2
	[ $status -eq 0 ]
	echo "${output}" | grep -q "test1"
	echo "${output}" | grep -q "test2"
	[ $status -eq 0 ]

	rm -f "${CONFIG_FILE}" "${BATS_TMPDIR}/ct.db"
}

@test "ct: report monthly-average" {
	CONFIG_FILE="${BATS_TMPDIR}/ct.json"
	run ./ct init --config-file "${CONFIG_FILE}"
	printf '%s\n' 'output: ' "${output}" >&2
	[ $status -eq 0 ]

	run ./ct log --config-file "${CONFIG_FILE}" --metric test1 --value 1 --timestamp 2020-01-01
	printf '%s\n' 'output: ' "${output}" >&2
	[ $status -eq 0 ]

	run ./ct log --config-file "${CONFIG_FILE}" --metric test2 --value 1 --timestamp 2020-02-01
	printf '%s\n' 'output: ' "${output}" >&2
	[ $status -eq 0 ]

	run ./ct report --config-file "${CONFIG_FILE}" --report-type=monthly-average --metrics test1
	printf '%s\n' 'output: ' "${output}" >&2
	[ $status -eq 0 ]
	echo "${output}" | grep "2020-01" | grep -q "test1"
	[ $status -eq 0 ]

	run ./ct report --config-file "${CONFIG_FILE}" --report-type=monthly-average --metrics test2
	printf '%s\n' 'output: ' "${output}" >&2
	[ $status -eq 0 ]
	echo "${output}" | grep "2020-02" | grep -q "test2"
	[ $status -eq 0 ]

	run ./ct report --config-file "${CONFIG_FILE}" --report-type=monthly-average --metrics test1,test2
	printf '%s\n' 'output: ' "${output}" >&2
	[ $status -eq 0 ]
	echo "${output}" | grep "2020-01" | grep -q "test1"
	echo "${output}" | grep "2020-02" | grep -q "test2"
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

	run ./ct configure --config-file "${CONFIG_FILE}" --metric test --value-text "foo"
	printf '%s\n' 'output: ' "${output}" >&2
	[ $status -eq 0 ]

	run ./ct dump --config-file "${CONFIG_FILE}" 
    [ $status -eq 0 ]
	printf '%s\n' 'output: ' "${output}" >&2
    [ $(echo "${output}" | jq -r .[].metric_name) == "test" ]
    [ $(echo "${output}" | jq -r .[].metric_config.frequency) == "daily" ]
    [ $(echo "${output}" | jq -r .[].metric_config.value_text) == "foo" ]

	run ./ct log --config-file "${CONFIG_FILE}" --metric test --value 1 --timestamp 2019-01-01
	printf '%s\n' 'output: ' "${output}" >&2
	[ $status -eq 0 ]

	run ./ct log --config-file "${CONFIG_FILE}" --metric test --value 2 --timestamp 20-01-02
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

