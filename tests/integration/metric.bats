#!/usr/bin/env bats

@test "ct: metric create" {
	CONFIG_FILE="${BATS_TMPDIR}/ct${RANDOM}.json"
	DB_FILE="${BATS_TMPDIR}/ct${RANDOM}.db"
	run ct init --config-file "${CONFIG_FILE}" --db-file "${DB_FILE}"

	run ct metric create --config-file "${CONFIG_FILE}" --metric-name test
	printf '%s\n' 'output: ' "${output}" >&2
	[ $status -eq 0 ]

	run ct dump --config-file "${CONFIG_FILE}"
	printf '%s\n' 'output: ' "${output}" >&2
	[ $status -eq 0 ]
    [ $(echo "${output}" | jq -r .metrics[0].name) == "test" ]

    rm -f "${CONFIG_FILE}" "${DB_FILE}"
}

@test "ct: metric create --value-text" {
	CONFIG_FILE="${BATS_TMPDIR}/ct${RANDOM}.json"
	DB_FILE="${BATS_TMPDIR}/ct${RANDOM}.db"
	run ct init --config-file "${CONFIG_FILE}" --db-file "${DB_FILE}"

	run ct metric create --config-file "${CONFIG_FILE}" --metric-name test --value-text "is this a test?"
	printf '%s\n' 'output: ' "${output}" >&2
	[ $status -eq 0 ]

	run ct dump --config-file "${CONFIG_FILE}"
	printf '%s\n' 'output: ' "${output}" >&2
	[ $status -eq 0 ]
    [ $(echo "${output}" | jq -r .configs[0].val) == "is this a test?" ]

    rm -f "${CONFIG_FILE}" "${DB_FILE}"
}

@test "ct: metric create --data-type" {
	CONFIG_FILE="${BATS_TMPDIR}/ct${RANDOM}.json"
	DB_FILE="${BATS_TMPDIR}/ct${RANDOM}.db"
	run ct init --config-file "${CONFIG_FILE}" --db-file "${DB_FILE}"

	run ct metric create --config-file "${CONFIG_FILE}" --metric-name test --data-type bool
	printf '%s\n' 'output: ' "${output}" >&2
	[ $status -eq 0 ]

	run ct dump --config-file "${CONFIG_FILE}"
	printf '%s\n' 'output: ' "${output}" >&2
	[ $status -eq 0 ]
    [ $(echo "${output}" | jq -r .configs[0].val) == "bool" ]

	run ct metric delete --config-file "${CONFIG_FILE}" --metric-name test
	printf '%s\n' 'output: ' "${output}" >&2
	[ $status -eq 0 ]

	run ct metric create --config-file "${CONFIG_FILE}" --metric-name test --data-type float
	printf '%s\n' 'output: ' "${output}" >&2
	[ $status -eq 0 ]

	run ct dump --config-file "${CONFIG_FILE}"
	printf '%s\n' 'output: ' "${output}" >&2
	[ $status -eq 0 ]
    [ $(echo "${output}" | jq -r .configs[0].val) == "float" ]

	run ct metric delete --config-file "${CONFIG_FILE}" --metric-name test
	printf '%s\n' 'output: ' "${output}" >&2
	[ $status -eq 0 ]

	run ct metric create --config-file "${CONFIG_FILE}" --metric-name test --data-type int
	printf '%s\n' 'output: ' "${output}" >&2
	[ $status -eq 0 ]

	run ct dump --config-file "${CONFIG_FILE}"
	printf '%s\n' 'output: ' "${output}" >&2
	[ $status -eq 0 ]
    [ $(echo "${output}" | jq -r .configs[0].val) == "int" ]

	run ct metric create --config-file "${CONFIG_FILE}" --metric-name test --data-type notSupported
	printf '%s\n' 'output: ' "${output}" >&2
	[ $status -eq 1 ]

    rm -f "${CONFIG_FILE}" "${DB_FILE}"
}

@test "ct: metric delete" {
	CONFIG_FILE="${BATS_TMPDIR}/ct${RANDOM}.json"
	DB_FILE="${BATS_TMPDIR}/ct${RANDOM}.db"
	run ct init --config-file "${CONFIG_FILE}" --db-file "${DB_FILE}"

	run ct metric create --config-file "${CONFIG_FILE}" --metric-name test
	printf '%s\n' 'output: ' "${output}" >&2
	[ $status -eq 0 ]

	run ct metric delete --config-file "${CONFIG_FILE}" --metric-name test
	printf '%s\n' 'output: ' "${output}" >&2
	[ $status -eq 0 ]

	run ct dump --config-file "${CONFIG_FILE}"
	printf '%s\n' 'output: ' "${output}" >&2
	[ $status -eq 0 ]
    [ $(echo "${output}" | jq -r .metrics) == "null" ]

    rm -f "${CONFIG_FILE}" "${DB_FILE}"
}

@test "ct: metric list" {
	CONFIG_FILE="${BATS_TMPDIR}/ct${RANDOM}.json"
	DB_FILE="${BATS_TMPDIR}/ct${RANDOM}.db"
	run ct init --config-file "${CONFIG_FILE}" --db-file "${DB_FILE}"

	run ct metric create --config-file "${CONFIG_FILE}" --metric-name test
	printf '%s\n' 'output: ' "${output}" >&2
	[ $status -eq 0 ]

	run ct metric list --config-file "${CONFIG_FILE}"
	printf '%s\n' 'output: ' "${output}" >&2
	[ $status -eq 0 ]

    rm -f "${CONFIG_FILE}" "${DB_FILE}"
}
