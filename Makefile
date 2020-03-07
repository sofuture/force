.PHONY: all
all:
	go install -mod=vendor github.com/gravitational/force/tool/force

.PHONY: oneshot
oneshot:
	cd examples/oneshot && force


.PHONY: vendor
vendor:
	go mod vendor

.PHONY: tidy
tidy:
	go mod tidy

.PHONY: install-docs
install-docs:
	pip install mkdocs==1.0.4
	pip install git+https://github.com/simonrenger/markdown-include-lines.git

.PHONY: serve-docs
serve-docs:
	mkdocs serve

##

.PHONY: 1-simple
1-simple:
	$(MAKE) all
	cd examples/demo/1-simple && force

.PHONY: 2-watch
2-watch:
	$(MAKE) all
	force -d examples/demo/2-watch/g.force

.PHONY: 3-github
3-github:
	$(MAKE) all
	force examples/demo/3-github/g.force --setup=./examples/github/setup.force

.PHONY: 4-docker
4-docker:
	$(MAKE) all
	force examples/demo/4-docker/g.force --setup=./examples/github/setup.force
##

.PHONY: ssh
ssh:
	$(MAKE) all
	cd examples/ssh && force -d ssh.force

.PHONY: aws
aws:
	$(MAKE) all
	cd examples/aws && force -d aws.force

.PHONY: flows
flows:
	$(MAKE) all
	cd examples/flows && force -d ci.force

.PHONY: github
github:
	$(MAKE) all
	cd examples/github && force -d ci.force

.PHONY: github-branches
github-branches:
	$(MAKE) all
	cd examples/github-branches && force -d ci.force --setup=../github/setup.force

.PHONY: slack
slack:
	$(MAKE) all
	cd examples/slack && force -d ci.force --setup=../github/setup.force

.PHONY: buildbox
buildbox:
	$(MAKE) all
	cd examples/teleport/buildbox && force --setup=../../github/setup.force

.PHONY: teleport
teleport:
	$(MAKE) all
	cd examples/teleport && force -d teleport.force

.PHONY: teleport-reload
teleport-reload:
	$(MAKE) all
	cd examples/teleport && force -d reload.force

.PHONY: teleport-apply
teleport-apply:
	$(MAKE) all
	cd examples/teleport && force -d apply.force --setup=setup-local.force


.PHONY: kube
kube:
	$(MAKE) all
	cd examples/kube && force -d kube.force --setup=./setup.force

.PHONY: vars
vars:
	$(MAKE) all
	cd examples/vars && force -d

.PHONY: reload
reload:
	$(MAKE) all
	cd examples/reload && force -d reload.force

.PHONY: conditionals
conditionals:
	$(MAKE) all
	cd examples/conditionals && force -d conditionals.force

.PHONY: hello
hello:
	$(MAKE) all
	cd examples/hello && force -d

.PHONY: hello-lambda
hello-lambda:
	$(MAKE) all
	cd examples/hello-lambda && force -d


.PHONY: inception
inception:
	$(MAKE) all
	cd inception && force -d inception.force --setup=./setup.force

.PHONY: mkdocs
mkdocs:
	$(MAKE) all
	cd mkdocs && force -d mkdocs.force --setup=../examples/github/setup.force


.PHONY: kbuild
kbuild:
	$(MAKE) all
	cd examples/kbuild && force kbuild.force


.PHONY: marshal
marshal:
	$(MAKE) all
	cd examples/marshal && force marshal.force


.PHONY: sloccount
sloccount:
	find . -path ./vendor -prune -o -name "*.go" -print0 | xargs -0 wc -l
