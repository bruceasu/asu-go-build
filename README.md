# asu-go-build

`asu-go-build` is a simple cross-compilation helper for Go projects.
It builds the project for the current platform and a set of predefined target platforms (Windows, Linux, etc.), making it easy to prepare distributable binaries without complex configuration.

---

## Features

* **Windows “lazy package”**

  * Run the provided executable directly; no need to remember long `go build` commands.
* **Linux/macOS support**

  * Use the `go-build.sh` shell script to achieve the same multi-platform build.
* **Automatic local build**

  * Builds for the current `GOOS/GOARCH` using `go env`.
* **Cross-compilation**

  * Preconfigured to build for:

    * `windows/amd64`
    * `windows/386`
    * `linux/amd64`
    * (optional: `darwin/amd64`, can be uncommented in source)
* **Customizable project name**

  * Provided via `-project` flag, or inferred from `go.mod` / current directory.

---

## Installation

Requirements:

* Go 1.18+ (must be available in PATH)

Clone and build:

```bash
git clone https://github.com/yourname/asu-go-build.git
cd asu-go-build
go build -o asu-go-build main.go
```

On Windows, the output will be `asu-go-build.exe`.

---

## Usage

```
asu-go-build -project <name> -output <dir>
```

### Options

| Option     | Description                                           | Default     |
| ---------- | ----------------------------------------------------- | ----------- |
| `-project` | Project name; if not provided, inferred automatically | current dir |
| `-output`  | Output directory for binaries                         | `bin`       |

---

## Examples

### Windows (lazy mode)

Simply run:

```powershell
asu-go-build.exe -project myapp -output bin
```

Outputs:

```
bin/myapp-windows-amd64.exe
bin/myapp-windows-386.exe
bin/myapp-linux-amd64
```

### Linux/macOS

Use the included script:

```bash
./go-build.sh -project myapp -output bin
```

---

## How It Works

1. Builds the local executable using the detected `GOOS`/`GOARCH`.
2. Cross-compiles for preconfigured target platforms.
3. Produces binaries named `<project>-<os>-<arch>[.exe]` under the output directory.

---

## Customization

* To add/remove target platforms, edit the `platforms` slice in `main.go`.
* macOS (`darwin/amd64` or `darwin/arm64`) can be enabled by uncommenting lines in the source.

---

## License

```
Apache License 2.0

Copyright 2025 BruceAsu

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

   http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
```


