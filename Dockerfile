# E2E テストの実行用イメージ。
# 例: docker run --rm -e E2E_ENV=stg e2e
FROM golang:1.26

WORKDIR /e2e

RUN go install gotest.tools/gotestsum@v1.12.0

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# 依存とテストを事前にビルドしておき、実行時間を短くする。
RUN go test -tags e2e -count=1 -run '^$' ./tests/...

ENV E2E_ENV=local
ENTRYPOINT ["gotestsum", "--junitfile", "report.xml", "--format", "testname", "--", "-tags", "e2e", "-count=1", "./tests/..."]
