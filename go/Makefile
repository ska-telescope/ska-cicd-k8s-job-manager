INSTALL_LOCATION ?= ${HOME}/.local/bin

.DEFAULT_GOAL := build
.PHONY: build

all: install

install: build
	@mkdir -p ${INSTALL_LOCATION}
	@mv kubectl-job_manager ${INSTALL_LOCATION}
	@echo "kubectl job-manager plugin installed."
	@echo "Ensure ${INSTALL_LOCATION} is in your system's PATH."

uninstall:
	@rm ${INSTALL_LOCATION}/kubectl-job_manager
	@echo "kubectl job-manager plugin uninstalled."

requirements:
	go mod tidy

build: requirements
	go build -o kubectl-job_manager *.go

clean:
	@rm kubectl-job_manager