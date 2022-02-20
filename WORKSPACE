workspace(name = "hello")

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
    "go_local_sdk",
    "go_rules_dependencies",
)

# TODO(irfansharif): Point to mirrored, public URL with modified runtime. For
# now this points to local checkout of go source tree.
go_local_sdk(
    name = "go_sdk",
    path = "/Users/irfansharif/Software/src/github.com/irfansharif/runner/modules/go",
)

go_rules_dependencies()

# Load gazelle dependencies.
load(
    "@bazel_gazelle//:deps.bzl",
    "gazelle_dependencies",
    "go_repository",
)

gazelle_dependencies()
