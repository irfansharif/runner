load("@bazel_gazelle//:def.bzl", "gazelle")
load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

gazelle(
    name = "gazelle",
    prefix = "github.com/irfansharif/runner",
)

go_library(
    name = "runner",
    srcs = [
        "runner.go",
        "runner.s",
        "runtime.go",
    ],
    importpath = "github.com/irfansharif/runner",
    visibility = ["//visibility:public"],
)

go_test(
    name = "runner_test",
    srcs = [
        "runner_test.go",
        "runtime_test.go",
    ],
    embed = [":runner"],
    deps = ["@com_github_stretchr_testify//assert"],
)
