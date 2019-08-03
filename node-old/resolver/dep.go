package resolver

type dep struct {
	name   string
	url    string
	sha256 string // of downloader url

	// the paths that the compiler should add to the include path search list
	includePaths []string

	// the paths that the source code will use to pull in this code
	// relative to the include path
	includes []string

	// the source files that btool will compile when pullin in this dep
	// relative to the root of the download
	sources []string

	// the header files for this dependency
	// relative to the root of the download
	headers []string
}

var deps = []*dep{
	&dep{
		name:   "googletest",
		url:    "https://github.com/google/googletest/archive/release-1.8.1.zip",
		sha256: "927827c183d01734cc5cfef85e0ff3f5a92ffe6188e0d18e909c5efebf28a0c7",

		includePaths: []string{
			"googletest-release-1.8.1/googletest/include",
			"googletest-release-1.8.1/googlemock/include",
		},
		includes: []string{
			"gtest/gtest.h",
			"gmock/gmock.h",
		},
		sources: []string{
			"googletest-release-1.8.1/googletest/src/gtest-all.cc",
			"googletest-release-1.8.1/googlemock/src/gmock-all.cc",
			"googletest-release-1.8.1/googlemock/src/gmock_main.cc",
		},
		headers: []string{
			"googletest-release-1.8.1/googletest/include/gtest/gtest.h",
			"googletest-release-1.8.1/googlemock/include/gmock/gmock.h",
		},
	},
}

func findDep(i string) *dep {
	for _, dep := range deps {
		for _, include := range dep.includes {
			if include == i {
				return dep
			}
		}
	}
	return nil
}
