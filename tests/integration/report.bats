#!/usr/bin/env bats

@test "ct: report all" {
	CONFIG_FILE="${BATS_TMPDIR}/ct.json"
	run ct init --config-file "${CONFIG_FILE}"
	printf '%s\n' 'output: ' "${output}" >&2
	[ $status -eq 0 ]

	run ct log --config-file "${CONFIG_FILE}" --metric test1 --value 1 --timestamp 2020-01-01
	printf '%s\n' 'output: ' "${output}" >&2
	[ $status -eq 0 ]

	run ct log --config-file "${CONFIG_FILE}" --metric test2 --value 1 --timestamp 2020-01-01
	printf '%s\n' 'output: ' "${output}" >&2
	[ $status -eq 0 ]

	run ct report all --config-file "${CONFIG_FILE}" --metrics test1
	printf '%s\n' 'output: ' "${output}" >&2
	[ $status -eq 0 ]
	echo "${output}" | grep -q "test1"
	[ $status -eq 0 ]

	run ct report all --config-file "${CONFIG_FILE}" --metrics test2
	printf '%s\n' 'output: ' "${output}" >&2
	[ $status -eq 0 ]
	echo "${output}" | grep -q "test2"
	[ $status -eq 0 ]

	run ct report all --config-file "${CONFIG_FILE}" --metrics test1,test2
	printf '%s\n' 'output: ' "${output}" >&2
	[ $status -eq 0 ]
	echo "${output}" | grep -q "test1"
	echo "${output}" | grep -q "test2"
	[ $status -eq 0 ]

	rm -f "${CONFIG_FILE}" "${BATS_TMPDIR}/ct.db"
}

@test "ct: report monthly-average" {
	CONFIG_FILE="${BATS_TMPDIR}/ct.json"
	run ct init --config-file "${CONFIG_FILE}"
	printf '%s\n' 'output: ' "${output}" >&2
	[ $status -eq 0 ]

	run ct log --config-file "${CONFIG_FILE}" --metric test1 --value 1 --timestamp 2020-01-01
	printf '%s\n' 'output: ' "${output}" >&2
	[ $status -eq 0 ]

	run ct log --config-file "${CONFIG_FILE}" --metric test2 --value 1 --timestamp 2020-02-01
	printf '%s\n' 'output: ' "${output}" >&2
	[ $status -eq 0 ]

	run ct report monthly-average --config-file "${CONFIG_FILE}" --metrics test1
	printf '%s\n' 'output: ' "${output}" >&2
	[ $status -eq 0 ]
	echo "${output}" | grep "2020-01" | grep -q "test1"
	[ $status -eq 0 ]

	run ct report monthly-average --config-file "${CONFIG_FILE}" --metrics test2
	printf '%s\n' 'output: ' "${output}" >&2
	[ $status -eq 0 ]
	echo "${output}" | grep "2020-02" | grep -q "test2"
	[ $status -eq 0 ]

	run ct report monthly-average --config-file "${CONFIG_FILE}" --metrics test1,test2
	printf '%s\n' 'output: ' "${output}" >&2
	[ $status -eq 0 ]
	echo "${output}" | grep "2020-01" | grep -q "test1"
	echo "${output}" | grep "2020-02" | grep -q "test2"
	[ $status -eq 0 ]

	rm -f "${CONFIG_FILE}" "${BATS_TMPDIR}/ct.db"
}

@test "ct: report streak" {
	CONFIG_FILE="${BATS_TMPDIR}/ct.json"
	run ct init --config-file "${CONFIG_FILE}"
	printf '%s\n' 'output: ' "${output}" >&2
	[ $status -eq 0 ]

	run ct log --config-file "${CONFIG_FILE}" --metric test1 --value 1 --timestamp 2020-01-01
	printf '%s\n' 'output: ' "${output}" >&2
	[ $status -eq 0 ]

	run ct configure --config-file "${CONFIG_FILE}" --metric test1 --data-type bool
	printf '%s\n' 'output: ' "${output}" >&2
	[ $status -eq 0 ]

	run ct log --config-file "${CONFIG_FILE}" --metric test1 --value 1 --timestamp 2020-01-02
	printf '%s\n' 'output: ' "${output}" >&2
	[ $status -eq 0 ]

	run ct log --config-file "${CONFIG_FILE}" --metric test2 --value 1 --timestamp 2020-02-01
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

	rm -f "${CONFIG_FILE}" "${BATS_TMPDIR}/ct.db"
}
