#!/usr/bin/env bats

@test "ct: log create int" {
	CONFIG_FILE="${BATS_TMPDIR}/ct${RANDOM}.json"
	DB_FILE="${BATS_TMPDIR}/ct${RANDOM}.db"
	run ct init --config-file "${CONFIG_FILE}" --db-file "${DB_FILE}"
	printf '%s\n' 'output: ' "${output}" >&2
	[ $status -eq 0 ]

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
    [ $(echo "${output}" | jq -r '.configs[] | select(.metric_id==1) | select(.opt=="data_type") | .val') == "int" ]

    rm -f "${CONFIG_FILE}" "${DB_FILE}"
}

@test "ct: log int with timestamp" {
	CONFIG_FILE="${BATS_TMPDIR}/ct${RANDOM}.json"
	DB_FILE="${BATS_TMPDIR}/ct${RANDOM}.db"
	run ct init --config-file "${CONFIG_FILE}" --db-file "${DB_FILE}"
	printf '%s\n' 'output: ' "${output}" >&2
	[ $status -eq 0 ]

	run ct metric create --config-file "${CONFIG_FILE}" --metric-name test --data-type int
	printf '%s\n' 'output: ' "${output}" >&2
	[ $status -eq 0 ]

	run ct log create --config-file "${CONFIG_FILE}" --metric-name test --metric-value 1 --timestamp 2020-01-01
	printf '%s\n' 'output: ' "${output}" >&2
	[ $status -eq 0 ]

	run ct dump --config-file "${CONFIG_FILE}"
	printf '%s\n' 'output: ' "${output}" >&2
	[ $status -eq 0 ]
    [ $(echo "${output}" | jq -r .logs[0].value) == "1" ]
    [ $(echo "${output}" | jq -r .logs[0].timestamp) == "2020-01-01T00:00:00Z" ]

    rm -f "${CONFIG_FILE}" "${DB_FILE}"
}

@test "ct: log float" {
	CONFIG_FILE="${BATS_TMPDIR}/ct${RANDOM}.json"
	DB_FILE="${BATS_TMPDIR}/ct${RANDOM}.db"
	run ct init --config-file "${CONFIG_FILE}" --db-file "${DB_FILE}"
	printf '%s\n' 'output: ' "${output}" >&2
	[ $status -eq 0 ]

	run ct metric create --config-file "${CONFIG_FILE}" --metric-name test --data-type float
	printf '%s\n' 'output: ' "${output}" >&2
	[ $status -eq 0 ]

	run ct log create --config-file "${CONFIG_FILE}" --metric-name test --metric-value 1.1
	printf '%s\n' 'output: ' "${output}" >&2
	[ $status -eq 0 ]

	run ct dump --config-file "${CONFIG_FILE}"
	printf '%s\n' 'output: ' "${output}" >&2
	[ $status -eq 0 ]
    [ $(echo "${output}" | jq -r .metrics[0].name) == "test" ]
    [ $(echo "${output}" | jq -r .logs[0].value) == "1.1" ]
    [ $(echo "${output}" | jq -r '.configs[] | select(.metric_id==1) | select(.opt=="data_type") | .val') == "float" ]

    rm -f "${CONFIG_FILE}" "${DB_FILE}"
}

@test "ct: log float with timestamp" {
	CONFIG_FILE="${BATS_TMPDIR}/ct${RANDOM}.json"
	DB_FILE="${BATS_TMPDIR}/ct${RANDOM}.db"
	run ct init --config-file "${CONFIG_FILE}" --db-file "${DB_FILE}"
	printf '%s\n' 'output: ' "${output}" >&2
	[ $status -eq 0 ]

	run ct metric create --config-file "${CONFIG_FILE}" --metric-name test --data-type float
	printf '%s\n' 'output: ' "${output}" >&2
	[ $status -eq 0 ]

	run ct log create --config-file "${CONFIG_FILE}" --metric-name test --metric-value 1.1 --timestamp 2020-01-01
	printf '%s\n' 'output: ' "${output}" >&2
	[ $status -eq 0 ]

	run ct dump --config-file "${CONFIG_FILE}"
	printf '%s\n' 'output: ' "${output}" >&2
	[ $status -eq 0 ]
    [ $(echo "${output}" | jq -r .logs[0].value) == "1.1" ]
    [ $(echo "${output}" | jq -r .logs[0].timestamp) == "2020-01-01T00:00:00Z" ]

    rm -f "${CONFIG_FILE}" "${DB_FILE}"
}

@test "ct: log create bool" {
	CONFIG_FILE="${BATS_TMPDIR}/ct${RANDOM}.json"
	DB_FILE="${BATS_TMPDIR}/ct${RANDOM}.db"
	run ct init --config-file "${CONFIG_FILE}" --db-file "${DB_FILE}"
	printf '%s\n' 'output: ' "${output}" >&2
	[ $status -eq 0 ]

	run ct metric create --config-file "${CONFIG_FILE}" --metric-name test --data-type bool
	printf '%s\n' 'output: ' "${output}" >&2
	[ $status -eq 0 ]

	run ct log create --config-file "${CONFIG_FILE}" --metric-name test --metric-value 0
	printf '%s\n' 'output: ' "${output}" >&2
	[ $status -eq 0 ]

	run ct dump --config-file "${CONFIG_FILE}"
	printf '%s\n' 'output: ' "${output}" >&2
	[ $status -eq 0 ]
    [ $(echo "${output}" | jq -r .logs[0].value) == "0" ]
    [ $(echo "${output}" | jq -r .metrics[0].name) == "test" ]
    [ $(echo "${output}" | jq -r '.configs[] | select(.metric_id==1) | select(.opt=="data_type") | .val') == "bool" ]

    rm -f "${CONFIG_FILE}" "${DB_FILE}"
}

@test "ct: log create bool with timestamp" {
	CONFIG_FILE="${BATS_TMPDIR}/ct${RANDOM}.json"
	DB_FILE="${BATS_TMPDIR}/ct${RANDOM}.db"
	run ct init --config-file "${CONFIG_FILE}" --db-file "${DB_FILE}"
	printf '%s\n' 'output: ' "${output}" >&2
	[ $status -eq 0 ]

	run ct metric create --config-file "${CONFIG_FILE}" --metric-name test --data-type bool
	printf '%s\n' 'output: ' "${output}" >&2
	[ $status -eq 0 ]

	run ct log create --config-file "${CONFIG_FILE}" --metric-name test --metric-value 0 --timestamp 2020-01-01
	printf '%s\n' 'output: ' "${output}" >&2
	[ $status -eq 0 ]

	run ct log create --config-file "${CONFIG_FILE}" --metric-name test --metric-value 1 --timestamp 2020-01-02
	printf '%s\n' 'output: ' "${output}" >&2
	[ $status -eq 0 ]

	run ct log create --config-file "${CONFIG_FILE}" --metric-name test --metric-value yes --timestamp 2020-01-03
	printf '%s\n' 'output: ' "${output}" >&2
	[ $status -eq 0 ]

	run ct log create --config-file "${CONFIG_FILE}" --metric-name test --metric-value no --timestamp 2020-01-04
	printf '%s\n' 'output: ' "${output}" >&2
	[ $status -eq 0 ]

	run ct dump --config-file "${CONFIG_FILE}"
	printf '%s\n' 'output: ' "${output}" >&2
	[ $status -eq 0 ]
    [ $(echo "${output}" | jq -r .logs[0].value) == "0" ]
    [ $(echo "${output}" | jq -r .logs[0].timestamp) == "2020-01-01T00:00:00Z" ]

    rm -f "${CONFIG_FILE}" "${DB_FILE}"
}


@test "ct: log create wrong type" {
	CONFIG_FILE="${BATS_TMPDIR}/ct${RANDOM}.json"
	DB_FILE="${BATS_TMPDIR}/ct${RANDOM}.db"
	run ct init --config-file "${CONFIG_FILE}" --db-file "${DB_FILE}"
	printf '%s\n' 'output: ' "${output}" >&2

	run ct metric create --config-file "${CONFIG_FILE}" --metric-name test --data-type float
	printf '%s\n' 'output: ' "${output}" >&2
	[ $status -eq 0 ]

	[ $status -eq 0 ]
	run ct log create --config-file "${CONFIG_FILE}" --metric-name test --metric-value 2.1 --timestamp 2020-01-01
	printf '%s\n' 'output: ' "${output}" >&2
	[ $status -eq 0 ]

	run ct log create --config-file "${CONFIG_FILE}" --metric-name test --metric-value foo --timestamp 2020-01-02
	printf '%s\n' 'output: ' "${output}" >&2
	[ $status -eq 1 ]

    rm -f "${CONFIG_FILE}" "${DB_FILE}"
}


@test "ct: log exceed daily frequency" {
	CONFIG_FILE="${BATS_TMPDIR}/ct${RANDOM}.json"
	DB_FILE="${BATS_TMPDIR}/ct${RANDOM}.db"
	run ct init --config-file "${CONFIG_FILE}" --db-file "${DB_FILE}"
	printf '%s\n' 'output: ' "${output}" >&2
	[ $status -eq 0 ]

	run ct metric create --config-file "${CONFIG_FILE}" --metric-name test --data-type float
	printf '%s\n' 'output: ' "${output}" >&2
	[ $status -eq 0 ]

	run ct log create --config-file "${CONFIG_FILE}" --metric-name test --metric-value 1 --timestamp 2020-01-01
	printf '%s\n' 'output: ' "${output}" >&2
	[ $status -eq 0 ]

	run ct log create --config-file "${CONFIG_FILE}" --metric-name test --metric-value 2.1 --timestamp 2020-01-01
	printf '%s\n' 'output: ' "${output}" >&2
	[ $status -eq 1 ]

    rm -f "${CONFIG_FILE}" "${DB_FILE}"
}

@test "ct: log create update" {
	CONFIG_FILE="${BATS_TMPDIR}/ct${RANDOM}.json"
	DB_FILE="${BATS_TMPDIR}/ct${RANDOM}.db"
	run ct init --config-file "${CONFIG_FILE}" --db-file "${DB_FILE}"
	printf '%s\n' 'output: ' "${output}" >&2
	[ $status -eq 0 ]

	run ct metric create --config-file "${CONFIG_FILE}" --metric-name test --data-type float
	printf '%s\n' 'output: ' "${output}" >&2
	[ $status -eq 0 ]

	run ct log create --config-file "${CONFIG_FILE}" --metric-name test --metric-value 1 --timestamp 2020-01-01
	printf '%s\n' 'output: ' "${output}" >&2
	[ $status -eq 0 ]

	run ct log create --config-file "${CONFIG_FILE}" --metric-name test --update --metric-value 2.1 --timestamp 2020-01-01
	printf '%s\n' 'output: ' "${output}" >&2
	[ $status -eq 0 ]

	run ct dump --config-file "${CONFIG_FILE}"
	printf '%s\n' 'output: ' "${output}" >&2
	[ $status -eq 0 ]
    [ $(echo "${output}" | jq -r .logs[0].value) == "2.1" ]

    rm -f "${CONFIG_FILE}" "${DB_FILE}"
}

@test "ct: log quiet" {
	CONFIG_FILE="${BATS_TMPDIR}/ct${RANDOM}.json"
	DB_FILE="${BATS_TMPDIR}/ct${RANDOM}.db"
	run ct init --config-file "${CONFIG_FILE}" --db-file "${DB_FILE}"
	printf '%s\n' 'output: ' "${output}" >&2
	[ $status -eq 0 ]

	run ct log --config-file "${CONFIG_FILE}" --metric test --value 1 --timestamp 2020-01-01
	printf '%s\n' 'output: ' "${output}" >&2
	[ $status -eq 0 ]

	run ct log --config-file "${CONFIG_FILE}" --metric test --value 1 --timestamp 2020-01-01 --quiet
	printf '%s\n' 'output: ' "${output}" >&2
	[ $status -eq 0 ]

	run ct log --config-file "${CONFIG_FILE}" --metric test --value 1 --timestamp 2020-01-01
	printf '%s\n' 'output: ' "${output}" >&2
	[ $status -eq 1 ]

    rm -f "${CONFIG_FILE}" "${DB_FILE}"
}

@test "ct: log stdin" {
	CONFIG_FILE="${BATS_TMPDIR}/ct${RANDOM}.json"
	DB_FILE="${BATS_TMPDIR}/ct${RANDOM}.db"
	run ct init --config-file "${CONFIG_FILE}" --db-file "${DB_FILE}"
	printf '%s\n' 'output: ' "${output}" >&2
	[ $status -eq 0 ]

	run ct log --config-file "${CONFIG_FILE}" --metric test --value 1
	printf '%s\n' 'output: ' "${output}" >&2
	[ $status -eq 0 ]

    echo 2.1 > ${BATS_TMPDIR}/$$.txt
	run ct log --config-file "${CONFIG_FILE}" --metric test --edit < ${BATS_TMPDIR}/$$.txt 
	printf '%s\n' 'output: ' "${output}" >&2
	[ $status -eq 0 ]

	run ct dump --config-file "${CONFIG_FILE}"
	printf '%s\n' 'output: ' "${output}" >&2
	[ $status -eq 0 ]
    [ $(echo "${output}" | jq -r .logs[0].value) == "2.1" ]

    rm -f "${CONFIG_FILE}" "${DB_FILE}"
}

@test "ct: log stdin exceed daily frequency" {
	CONFIG_FILE="${BATS_TMPDIR}/ct${RANDOM}.json"
	DB_FILE="${BATS_TMPDIR}/ct${RANDOM}.db"
	run ct init --config-file "${CONFIG_FILE}" --db-file "${DB_FILE}"
	printf '%s\n' 'output: ' "${output}" >&2
	[ $status -eq 0 ]

    echo 2.1 > ${BATS_TMPDIR}/$$.txt
	run ct log --config-file "${CONFIG_FILE}" --metric test --timestamp 2020-01-01 < ${BATS_TMPDIR}/$$.txt 
	printf '%s\n' 'output: ' "${output}" >&2
	[ $status -eq 0 ]

	run ct log --config-file "${CONFIG_FILE}" --metric test --timestamp 2020-01-01 < ${BATS_TMPDIR}/$$.txt 
	printf '%s\n' 'output: ' "${output}" >&2
	[ $status -eq 1 ]

	run ct log --config-file "${CONFIG_FILE}" --metric test --timestamp 2020-01-01 --quiet < ${BATS_TMPDIR}/$$.txt 
	printf '%s\n' 'output: ' "${output}" >&2
	[ $status -eq 0 ]

    rm -f "${CONFIG_FILE}" "${DB_FILE}"
}

