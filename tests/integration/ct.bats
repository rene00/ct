#!/usr/bin/env bats

@test "ct: init" {
	CONFIG_FILE="${BATS_TMPDIR}/ct.json"
	run ct init --config-file "${CONFIG_FILE}"
	printf '%s\n' 'output: ' "${output}" >&2
	[ $status -eq 0 ]
	[ "$(jq -r .ct.db_file ${CONFIG_FILE})" == "${BATS_TMPDIR}/ct.db" ]
	rm -f "${CONFIG_FILE}" "${BATS_TMPDIR}/ct.db"
}

@test "ct: db migrate" {
	CONFIG_FILE="${BATS_TMPDIR}/ct.json"
	run ct init --config-file "${CONFIG_FILE}"
	printf '%s\n' 'output: ' "${output}" >&2
	[ $status -eq 0 ]

	run ct db migrate --config-file "${CONFIG_FILE}" --run
	printf '%s\n' 'output: ' "${output}" >&2
	[ $status -eq 0 ]
}


@test "ct: configure" {
	CONFIG_FILE="${BATS_TMPDIR}/ct.json"
	run ct init --config-file "${CONFIG_FILE}"
	printf '%s\n' 'output: ' "${output}" >&2
	[ $status -eq 0 ]

	run ct configure --config-file "${CONFIG_FILE}" --metric test
	printf '%s\n' 'output: ' "${output}" >&2
	[ $status -eq 0 ]

	run ct configure --config-file "${CONFIG_FILE}" --metric test --value-text "foo" --data-type "float"
	printf '%s\n' 'output: ' "${output}" >&2
	[ $status -eq 0 ]

	run ct dump --config-file "${CONFIG_FILE}" 
    [ $status -eq 0 ]
	printf '%s\n' 'output: ' "${output}" >&2
    [ $(echo "${output}" | jq -r .[].metric_name) == "test" ]
    [ $(echo "${output}" | jq -r .[].metric_config.value_text) == "foo" ]
    [ $(echo "${output}" | jq -r .[].metric_config.data_type) == "float" ]

	run ct log --config-file "${CONFIG_FILE}" --metric test --value 1 --timestamp 2019-01-01
	printf '%s\n' 'output: ' "${output}" >&2
	[ $status -eq 0 ]

	run ct log --config-file "${CONFIG_FILE}" --metric test --value 2.1 --timestamp 2019-01-02
	printf '%s\n' 'output: ' "${output}" >&2
	[ $status -eq 0 ]

	run ct log --config-file "${CONFIG_FILE}" --metric test --value 2 --timestamp 20-01-02
	printf '%s\n' 'output: ' "${output}" >&2
	[ $status -eq 1 ]

	run ct configure --config-file "${CONFIG_FILE}" --metric test --data-type notSupportedDataType
	printf '%s\n' 'output: ' "${output}" >&2
	[ $status -eq 1 ]

	run ct configure --config-file "${CONFIG_FILE}" --metric test --data-type int
	printf '%s\n' 'output: ' "${output}" >&2
	[ $status -eq 0 ]

	run ct log --config-file "${CONFIG_FILE}" --metric test --value 2.0 --timestamp 2019-01-03
	printf '%s\n' 'output: ' "${output}" >&2
	[ $status -eq 1 ]

	run ct log --config-file "${CONFIG_FILE}" --metric test --value 2 --timestamp 2019-01-03
	printf '%s\n' 'output: ' "${output}" >&2
	[ $status -eq 0 ]

	run ct configure --config-file "${CONFIG_FILE}" --metric test --data-type float
	printf '%s\n' 'output: ' "${output}" >&2
	[ $status -eq 0 ]

	rm -f "${CONFIG_FILE}" "${BATS_TMPDIR}/ct.db"
}

