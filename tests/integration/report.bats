#!/usr/bin/env bats

@test "ct: report daily" {
	CONFIG_FILE="${BATS_TMPDIR}/ct${RANDOM}.json"
	DB_FILE="${BATS_TMPDIR}/ct${RANDOM}.db"
	run ct init --config-file "${CONFIG_FILE}" --db-file "${DB_FILE}"
	printf '%s\n' 'output: ' "${output}" >&2
	[ $status -eq 0 ]

	run ct metric create --config-file "${CONFIG_FILE}" --metric-name test1 --data-type int
	printf '%s\n' 'output: ' "${output}" >&2
	[ $status -eq 0 ]

	run ct metric create --config-file "${CONFIG_FILE}" --metric-name test2 --data-type int
	printf '%s\n' 'output: ' "${output}" >&2
	[ $status -eq 0 ]

	run ct log create --config-file "${CONFIG_FILE}" --metric-name test1 --metric-value 1 --timestamp 2020-01-01
	printf '%s\n' 'output: ' "${output}" >&2
	[ $status -eq 0 ]

	run ct log create --config-file "${CONFIG_FILE}" --metric-name test2 --metric-value 1 --timestamp 2020-01-01
	printf '%s\n' 'output: ' "${output}" >&2
	[ $status -eq 0 ]

	run ct report daily --config-file "${CONFIG_FILE}" --metric test1
	printf '%s\n' 'output: ' "${output}" >&2
	[ $status -eq 0 ]

	run ct report daily --config-file "${CONFIG_FILE}" --metric test2
	printf '%s\n' 'output: ' "${output}" >&2
	[ $status -eq 0 ]

    rm -f "${CONFIG_FILE}" "${DB_FILE}"
}

@test "ct: report monthly" {
	CONFIG_FILE="${BATS_TMPDIR}/ct${RANDOM}.json"
	DB_FILE="${BATS_TMPDIR}/ct${RANDOM}.db"
	run ct init --config-file "${CONFIG_FILE}" --db-file "${DB_FILE}"
	printf '%s\n' 'output: ' "${output}" >&2
	[ $status -eq 0 ]

	run ct metric create --config-file "${CONFIG_FILE}" --metric-name test1 --data-type int --metric-type gauge
	printf '%s\n' 'output: ' "${output}" >&2
	[ $status -eq 0 ]

	run ct metric create --config-file "${CONFIG_FILE}" --metric-name test2 --data-type int --metric-type counter
	printf '%s\n' 'output: ' "${output}" >&2
	[ $status -eq 0 ]

	run ct log create --config-file "${CONFIG_FILE}" --metric-name test1 --metric-value 1 --timestamp 2020-01-01
	printf '%s\n' 'output: ' "${output}" >&2
	[ $status -eq 0 ]

	run ct log create --config-file "${CONFIG_FILE}" --metric-name test2 --metric-value 1 --timestamp 2020-02-01
	printf '%s\n' 'output: ' "${output}" >&2
	[ $status -eq 0 ]

	run ct report monthly --config-file "${CONFIG_FILE}" --metric test1
	printf '%s\n' 'output: ' "${output}" >&2
	[ $status -eq 0 ]

	run ct report monthly --config-file "${CONFIG_FILE}" --metric test2
	printf '%s\n' 'output: ' "${output}" >&2
	[ $status -eq 0 ]

    rm -f "${CONFIG_FILE}" "${DB_FILE}"
}

@test "ct: report streak" {
	CONFIG_FILE="${BATS_TMPDIR}/ct${RANDOM}.json"
	DB_FILE="${BATS_TMPDIR}/ct${RANDOM}.db"
	run ct init --config-file "${CONFIG_FILE}" --db-file "${DB_FILE}"
	printf '%s\n' 'output: ' "${output}" >&2
	[ $status -eq 0 ]

	run ct metric create --config-file "${CONFIG_FILE}" --metric-name test1 --data-type int
	printf '%s\n' 'output: ' "${output}" >&2
	[ $status -eq 0 ]

	run ct metric create --config-file "${CONFIG_FILE}" --metric-name test2 --data-type int
	printf '%s\n' 'output: ' "${output}" >&2
	[ $status -eq 0 ]

	run ct log create --config-file "${CONFIG_FILE}" --metric-name test1 --metric-value 1 --timestamp 2020-01-01
	printf '%s\n' 'output: ' "${output}" >&2
	[ $status -eq 0 ]

	run ct configure --config-file "${CONFIG_FILE}" --metric test1 --data-type bool
	printf '%s\n' 'output: ' "${output}" >&2
	[ $status -eq 0 ]

	run ct log create --config-file "${CONFIG_FILE}" --metric-name test1 --metric-value 1 --timestamp 2020-01-02
	printf '%s\n' 'output: ' "${output}" >&2
	[ $status -eq 0 ]

	run ct log create --config-file "${CONFIG_FILE}" --metric-name test2 --metric-value 1 --timestamp 2020-02-01
	printf '%s\n' 'output: ' "${output}" >&2
	[ $status -eq 0 ]

	run ct configure --config-file "${CONFIG_FILE}" --metric test2 --data-type bool
	printf '%s\n' 'output: ' "${output}" >&2
	[ $status -eq 0 ]

	run ct report streak --config-file "${CONFIG_FILE}" --metric test1
	printf '%s\n' 'output: ' "${output}" >&2
	[ $status -eq 0 ]
	echo "${output}" | grep "test1" | awk '{print $4}' | grep 2
	[ $status -eq 0 ]

	run ct report streak --config-file "${CONFIG_FILE}" --metric test2
	printf '%s\n' 'output: ' "${output}" >&2
	[ $status -eq 0 ]
	echo "${output}" | grep "test2" | awk '{print $4}' | grep 1
	[ $status -eq 0 ]

    rm -f "${CONFIG_FILE}" "${DB_FILE}"
}
