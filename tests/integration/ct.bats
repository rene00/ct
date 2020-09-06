#!/usr/bin/env bats

@test "ct: init" {
	CONFIG_FILE="${BATS_TMPDIR}/ct${RANDOM}.json"
    DB_FILE="${BATS_TMPDIR}/ct${RANDOM}.db"
	run ct init --config-file "${CONFIG_FILE}" --db-file "${DB_FILE}"
	printf '%s\n' 'output: ' "${output}" >&2
	[ $status -eq 0 ]
	[ "$(jq -r .ct.db_file ${CONFIG_FILE})" == "${DB_FILE}" ]
}

@test "ct: db migrate" {
	CONFIG_FILE="${BATS_TMPDIR}/ct.json"
    DB_FILE="${BATS_TMPDIR}/ct${RANDOM}.db"
	run ct init --config-file "${CONFIG_FILE}" --db-file "${DB_FILE}"
	printf '%s\n' 'output: ' "${output}" >&2
	[ $status -eq 0 ]

	run ct db migrate --config-file "${CONFIG_FILE}" --run
	printf '%s\n' 'output: ' "${output}" >&2
	[ $status -eq 0 ]
}


