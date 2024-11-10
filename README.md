# Matrix-Task

![Build and Release](https://github.com/tanq16/matrix-task/actions/workflows/release.yml/badge.svg)

This project is a task-management app that uses the [Eisenhower Matrix](https://asana.com/resources/eisenhower-matrix) categorization. It is meant to be a dead simple app that can store tasks, mark them as complete, or delete them. All completed tasks automatically go to an Archive.

The intent behind making this app wasn't to create this app specifically, but to test the capabilities of an LLM. I used a paid LLM platform to create this app. Roughly everything was written by the LLM, but I had to fix a couple of minor things (10-12 lines at max) to get it working.

Generally speaking, first the application was built, tested, errors manually fixed, and containerization included.

Next, I used the LLM to generate code to use the `embed` package in Go to embed the static and template files. Overall, this repo is an experiment that also helped produce a task management app as a result.

To run the application, simply download the binary from releases and execute. I build both `x86-64` and `ARM64` variants.
