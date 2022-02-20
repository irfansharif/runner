# Copyright 2022 Irfan Sharif.
# 
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
# 
#     http://www.apache.org/licenses/LICENSE-2.0
# 
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or
# implied. See the License for the specific language governing
# permissions and limitations under the License.

test:
	bazel test //:all --test_arg='-test.v' --test_output=all

bench:
	bazel test //:all --nocache_test_results \
		--test_arg='-test.v' --test_output=all --test_arg='-test.run=-' \
		--test_arg='-test.bench=.' --test_arg='-test.benchtime=10000x'

generate: FORCE
	bazel run //:gazelle -- update-repos \
		-from_file=go.mod -prune=true \
		-build_file_proto_mode=disable_global \
		-to_macro=DEPS.bzl%go_deps &> /dev/null
	bazel run //:gazelle &> /dev/null

go:
	git submodule update --init --recursive
	cd modules/go/src && ./make.bash
	
FORCE: ;
