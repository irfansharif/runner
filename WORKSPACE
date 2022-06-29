workspace(name = "runner")

load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")

http_archive(
    name = "io_bazel_rules_go",
    sha256 = "5c4bd27429b1a307d51cd23d4677126aa6315fff608f0cd85c5bfb642a13b953",
    strip_prefix = "cockroachdb-rules_go-23b381c",
    urls = [
        "https://storage.googleapis.com/public-bazel-artifacts/bazel/cockroachdb-rules_go-v0.27.0-52-g23b381c.tar.gz",
    ],
)

# Load gazelle. This lets us auto-generate BUILD.bazel files throughout the
# repo.
http_archive(
    name = "bazel_gazelle",
    sha256 = "9fba095e4bebd8c6748154ca53c365862af47fa1651f7c0d25459e6ca5bb208f",
    strip_prefix = "bazelbuild-bazel-gazelle-3ea1d64",
    urls = [
        # v0.24.0
        "https://storage.googleapis.com/public-bazel-artifacts/bazel/bazelbuild-bazel-gazelle-v0.24.0-0-g3ea1d64.tar.gz",
    ],
)

# Load up go dependencies (the ones listed under go.mod).
load("//:DEPS.bzl", "go_deps")

# gazelle:repository_macro DEPS.bzl%go_deps
go_deps()

load(
    "@io_bazel_rules_go//go:deps.bzl",
    "go_download_sdk",
    "go_rules_dependencies",
)

go_download_sdk(
    name = "go_sdk",
    sdks = {
        "darwin_arm64": ("go1.17.6.darwin-arm64.tar.gz", "73b33986a3a4d7a4083b8457f9a62ffbc6803b31b2509c630f33327c5717e4fb"),
        "darwin_amd64": ("go1.17.6.darwin-amd64.tar.gz", "91f3f14e07acf4322d864e0c4ef70607d4b1e101883c53903568f11b1493becd"),
        "freebsd_amd64": ("go1.17.6.freebsd-amd64.tar.gz", "6d11cf238eaaad96d51c711f8af9d327e1809403ab1b1b597ac13d0b45dfac76"),
        "linux_amd64": ("go1.17.11.linux-amd64.tar.gz", "77da33c8b2699ffb6920439087964b059c6cbfb0101be3ac11b93c43b5ab2c5a"),
    },
    urls = ["https://github.com/irfansharif/utils/blob/master/goroutine-nanos/{}?raw=true"],
    version = "1.17.11",
)

go_rules_dependencies()

# Load gazelle dependencies.
load(
    "@bazel_gazelle//:deps.bzl",
    "gazelle_dependencies",
    "go_repository",
)

gazelle_dependencies()
