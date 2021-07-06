Go language plugin for ThoughtWorks [Gauge](http://getgauge.io).

## Getting started in 3 steps
* Install [Gauge](http://gauge.org) by following instructions [here](https://docs.gauge.org/getting_started/installing-gauge.html) and gauge go plugin as `gauge install go`
* Initialize a golang gauge project in subfolder within `GOPATH/src` directory as: 
```sh
cd $GOPATH/src
mkdir gaugeproject
cd gaugeproject
gauge init go
```
* Run specs: `gauge run specs`

## Installation
* Install [Gauge](http://gauge.org). Follow instructions [here](https://docs.gauge.org/getting_started/installing-gauge.html).
* Install Gauge-Go language plugin as: `gauge install go`

### Build from Source
* Checkout the repository
* Run commands:
```sh
go run build/make.go
go run build/make.go --install
```
This installs the gauge-go plugin to its default location based on OS.

#### Run Unit tests

`go test ./...`

#### Generate gauge-go distributables

`go run build/make.go --distro`

This generates the gauge-go plugin as zip distributable. In order to generate it for all platforms, add `--all-platforms` flag in above command.

#### Example project
A sample project illustrating Gauge features using Golang & selenium webdriver can be found [here](https://github.com/getgauge-contrib/gauge-example-go).

## Usage

If you are new to Gauge, please read the [Gauge documentation](https://docs.gauge.org/index.html) to know about how Gauge works.

**Initialize:**

To initialize a project with gauge-go, in an empty directory run:
```sh
$ gauge init go
```
**Note: Create your project in `$GOPATH/src`.**

**Run specs:**

```sh
$ gauge run specs/
```

**Run specs with remote debugging enabled:**

```sh
$ GAUGE_DEBUG_OPTS=40000 gauge run specs/
```
## Methods

### Step implementation

**`var _ = gauge.Step(<step-text>, fn)`**

Use the `gauge.Step()` method to implement your steps. For example:

```go
import (
    "github.com/getgauge-contrib/gauge-go/gauge"
	. "github.com/getgauge-contrib/gauge-go/testsuit"
)    

var _ = gauge.Step("Vowels in English language are <vowels>.", func (vowelsGiven string) {
    if vowelsGiven != "aeiou" {
        T.Fail(fmt.Errorf("want: %s, got: %s", "aeiou", vowelsGiven))
    }
})
```
You can use `T.Fail` or `T.Errorf` functions to mark a step as failure, which are under package `github.com/getgauge-contrib/gauge-go/testsuit`. You can also use any assertion libraries like [testify](https://github.com/stretchr/testify), by passing `testsuit.T` as argument to assertions.

E.g: `assert.Equal(testsuit.T, actualCount, expectedCount, "got: %d, want: %d", actualCount, expectedCount)`

### Execution Hooks

gauge-go supports tagged [execution hooks](https://docs.gauge.org/writing-specifications.html#execution-hooks). These methods are available for each type of hook:

"Before" hooks:

- **`BeforeSuite(fn, tags, operator)`** - Executed before the test suite begins
- **`BeforeSpec(fn, tags, operator)`** - Executed before each specification
- **`BeforeScenario(fn, tags, operator)`** - Executed before each scenario
- **`BeforeStep(fn, tags, operator)`**- Execute before each step

"After" hooks:

- **`AfterSuite(fn, tags, operator)`** - Executed after the test suite begins
- **`AfterSpec(fn, tags, operator)`** - Executed after each specification
- **`AfterScenario(fn, tags, operator)`** - Executed after each scenario
- **`AfterStep(fn, tags, operator)`**- Execute after each step

Here's an example of a hook that is executed before each scenario:

```go
var _ = gauge.BeforeScenario(function () {
  assert.equal(vowels.join(""), "aeiou")
}, []string{"foo", "bar"}, testsuit.AND)
```

- *`tags`*
An array of strings for the tags for which to execute the current hook. These are only useful at specification or scenario level. If tags is empty, the provided hook is executed on each occurrence of the hook.
- *`operator`*: Valid values: `"AND"`, `"OR"`.
This controls whether the current hook is executed when all of the tags match (in case of `"AND"`), or if any of the tags match (in case of `OR`).

### Data Stores

Step implementations can share custom data across scenarios, specifications and suites using data stores.

There are 3 different types of data stores based on the lifecycle of when it gets cleared. These are present in the `github.com/getgauge-contrib/gauge-go/gauge` package of gauge-go plugin.

#### Scenario store

This data store keeps values added to it in the lifecycle of the scenario execution. Values are cleared after every scenario executes.

**Store a value:**

```go
GetScenarioStore()[key] = value
```

**Retrieve a value:**

```go
value := GetScenarioStore()[key]
```

#### Specification store

This data store keeps values added to it in the lifecycle of the specification execution. Values are cleared after every specification executes.

**Store a value:**

```go
GetSpecStore()[key] = value
```

**Retrieve a value:**

```go
value := GetSpecStore()[key]
```

#### Suite store

This data store keeps values added to it in the lifecycle of the entire suite's execution. Values are cleared after entire suite executes.

**Store a value:**

```go
GetSuiteStore()[key] = value
```

**Retrieve a value:**

```go
value := GetSuiteStore()[key]
```

**Note:** Suite Store is not advised to be used when executing specs in parallel. The values are not retained between parallel streams of execution.

### Custom Messages

Custom messages/data can be added to execution reports using the below API from the step implementations. The API is under the package `github.com/getgauge-contrib/gauge-go/gauge` which should be imported.

These messages will appear under steps in the execution reports.

```go
gauge.WriteMessage("my custom message")
gauge.WriteMessage("Say %s to %s %d times", "hello", "Gauge", 10) //prints: Say hello to Gauge 10 times
```

### Custom Screenshot

You can specify a custom function to grab a screenshot on step failure. By default, gauge-go takes screenshot of the current screen using the `gauge_screenshot` binary.

This custom function should be set on the `gauge.CustomScreenshotFn` property in test implementation code and it should return a base64 encoded byte array of the image data that gauge-go will use as image content on failure.

```go
import "github.com/getgauge-contrib/gauge-go/gauge"

func init() {
	gauge.CustomScreenshotFn = func() []byte {
		screenshot, err := driver.Screenshot()
        if err != nil {
            return nil
        }
		return screenshot
	}
}
```

## License

![GNU Public License version 3.0](http://www.gnu.org/graphics/gplv3-127x51.png)
Gauge-Go is released under [GNU Public License version 3.0](http://www.gnu.org/licenses/gpl-3.0.txt)
