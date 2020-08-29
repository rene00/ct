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


