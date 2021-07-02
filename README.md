## Trainer

Trainer is a configurable, extensible mock state machine for
microservice ecosystems that consist of multiple microservices
that coordinate amongst each other to accomplish a goal.  While
it is extremely flexible and can be used for nearly any microservice
ecosystem, it is particularly useful for an ecosystem that must contact
one or more external services as a component of processing its request.

It has two main functions.
* Trainer can kick off tests and provide the result of that test by contacting appropriate microservices.
* Trainer can mock any external service that the microservices will contact
and return the expected information. This allows for tests with completely controlled inputs.

### Rationale

Unit tests are relatively easy to create.  Integration tests are far more difficult.
Tests that can exercise a running ecosystem are even more difficult, and this tool addresses that need.

## Theory of operation

Trainer is a series of plans, one of which can be running at any time.
Each plan consists of a series of transactions, each one of which represents either a set
of actions that need to be performed, or an API transaction, in
which case a series of actions are defined based upon whether the
transaction was successful. As trainer moves through a plan,
specific actions are requested from the API, and responses are then
collected from the API to gauge whether or not the specified
plan succeeded.

## Limitations

- The service must be restarted to reload the configuration file.  There is
  currently no concept whatsoever of dynamic plan loading or reloading.
- Trainer can only use configuration files.  There is no other backend currently
  supported.  Pull requests welcome.
- Trainer is poorly suited for stress-testing.
- While there is a limited ability for parallel processing (a callback
  can be split and other processes can occur while a callback is open),
  generally trainer will iterate through one state at a time.  This means
  there are some structural limitations that could only be fixed with a
  pretty extensive rewrite.

## Starting

Trainer is designed as a PCF (Pivotal Cloud Foundry) based microservice.
Practically, this means that it accepts a PORT environment variable to figure
out which port it needs to listen on.  It would be a "todo" to add support for
different types of cloud environments, or make it environment agnostic.  Trainer is
unaware of the environment it's running under other than that environment
variable, so tweaks should be simple.

The following command line arguments are supported (If one is required but has a default, you don't have to specify):
| Name | Env | Default | Required | Description |
| - | - | - | - | - |
| apiport | PORT | 8080 | Yes | The port to listen on |
| apihost | APILISTENHOST | localhost | Yes | The host to listen on (try 0.0.0.0) |
| apiuser | APIAUTHUSERNAME | | Yes | The username to be used to auth to the API |
| apipass | APIAUTHPASSWORD | | Yes | The password to be used to auth to the API |
| loglevel | LOGLEVEL | WARNING | Yes | The Loglevel (TRACE, DEBUG, INFO, WARNING, ERROR) |
| configfile | CONFIGFILE | config.yml | Yes | The config file |
| testmode | | | No | For development purposes |
| testurl | | | No | For development purposes |

Trainer uses kingpin, so see the [Go kingpin documentation](https://pkg.go.dev/gopkg.in/alecthomas/kingpin.v2?utm_source=godoc) for further details on command line parsing.

## Using

A plan can be thought of as a test.  A plan is kicked off through
the API and runs to completion.  What completion *means* is a fully configurable thing.
It could mean that you run through to a failure or success transaction which sets a variable.
It could mean that it connects to an API somewhere and sends the information there. Whatever you
think a good outcome (or a bad outcome) is, trainer can probably accommodate.

To get started, write a plan and add it to config.yml.  Where you put the plans is fully configurable.
There are some test plans in the "data" directory of the source code, please feel free to
use those for inspiration.

Once these plans are written and the service is successfully started,
then you can run a test.  Go ahead and launch the plan (see API below).  The plan should succeed, fail, or hang.

You can hit the "status" endpoint at any point to see how your tests are going.  That can also contain some valuable
debugging info, and even more if you have configured your plans to set variables or log.

Trainer is, generally, not plug and play.  If your microservices or other services connect to an external
service, there is really no way for trainer to mock and/or proxy that service at the present time.  So the ecosystem
has to be aware that it's running a test.  You can configure trainer to send specific headers that your services can
recognize, or to send special payloads. Bottom line, your service has
to know when to contact trainer instead of its own external services.

Once you have that set up, you gain the ability for "canary testing", in that you can submit tests to your
services at the same time the services are serving production workloads.  Trainer was designed with
this use case in mind, and works well for that purpose.

## TODOS

Enhancement ideas:

- Stored config somewhere other than yaml files
- Web based UI for configuration and monitoring
- Multiple in progress tests (somehow!)
- Proxying of external services

## API

### Launch

```
/launch/<plan>
```

Set the current plan to the specified plan.

The plan must exist or an error will be returned.

The plan will be implicitly reset.

<!--
If you specify a "planincludes" array in the configuration, you may add plans as individual
files,

```
planincludes:
  - filename1
  - filename2
```

These will be appended AFTER the plans specified in the config file.
Practically, this will make no difference.

In each plan, if you add an externalvars field pointing to a file,
pointing to a key/value pair yaml file (nothing more complicated than
that, please) you can load in external variables. This is good for things
like the authorization headers, etc., which can be templated.

See the test configs for an example.
-->

### Reset

```
/remove
```

The current plan will be removed. This will completely reset all state,
including variables, etc. Don't do this until you are sure you don't
need that output anymore.

### Status

```
/status
```

Get the current running state of the application. This contains the
state history of the current run. It also contains the entire
state structure, the disposition, the kitchen sink, and a king sized
waterbed.

Meaning, it dumps a _lot_ of info, but if you take some time to
understand what it's telling you, it's very useful for monitoring,
control, and troubleshooting.

### Config

```
/config
```

Dumps the current running config. Useful to know that it's running
the right config if you have automatic deployment/CI/CD.

## Configuration

Configuration is done through a yaml file, which can be configured
on the command line. All other files necessary are relative to the
base configuration file. Subdirectories may be used for other
necessary files, just use the path relative to the config file.
For example, if the config file is in

```
/home/trainer
```

Then the config file would be in

```
/home/trainer/config.xml
```

To create a plan specific directory just

```
mkdir /home/trainer/plan1
```

And create a file underneath. Then reference it in an appropriate
entry:

```
file: plan1/file
```

The yaml file has a very specific syntax, described here.

### Action

An action has the following syntax:

```
type: <type>
  args:
    <arg1>: <arg1value>
    <arg2>: <arg2value>
    ...
```

There are different kinds of actions, and all of the different kinds
of actions take different args. The actions are:

#### Advance

###### Purpose

This action advances to the specified transaction. Any further
actions to process in the current transaction are skipped.

###### Args

| Arg | Description                   |
| --- | ----------------------------- |
| txn | The transaction to advance to |

#### Callback

###### Purpose

This action calls back to a specific URL.
It can also be used to query a specific URL (not as a callback).

###### Args

| Arg                 | Type                 | Description                                       |
| ------------------- | -------------------- | ------------------------------------------------- |
| url                 | template (see below) | The URL to call.                                  |
| payload_contenttype | json/yaml            | the content type of the payload                   |
| payload             | file                 | The filename of the payload                       |
| auth_header         | string               | The auth header. Basic <hash>                     |
| method              | POST/GET             | the method to use when sending                    |
| response_type       | string               | the expected response type (json, yaml, string)   |
| save                | map                  | the variables to save from the json (see below)   |
| save_response       | variable             | the variable name to save the full response into  |
| save_response_map   | variable             | if set, copy the json decoded response into a map |
| ignore_failure      | boolean              | if true, keep going even if the callback fails.   |
| headers             | map                  | arbitrary headers.  keys and values must be strings. |

args that are used by a particular action are ignored.

Any arg preceded by an underscore (for example, "\_context"), is reserved
and should not be used by user configuration.

Note that this is for a callback that returns without making any interstitiary calls to trainer.
If you have such a need, use the split callback functionality below.

If response_type is "string", don't attempt to use save_response_map.  A map is not generated with a string.
Unsure what this will do, but it might panic, or just do nothing.

#### cb_split/cb_finish

These two actions create a split callback.

The cb_split action creates, but does not complete, a callback.  The callback is started - with the same
arguments as a regular callback - but it is held in a state of "stasis".  What this means is that other actions
can be run while this callback runs in the background.

This is designed so that the process the callback starts can send responses back to trainer in order to complete
whichever actions it needs to take in order to complete.

Once it has done what you are expecting, then use the cb_finish action.  This gathers the response from the
callback and finishes the execution.

cb_split cannot be run with a pending callback.  cb_finish cannot be run with no callback.  You should finish
any pending callbacks, even if actions in between fail (set your failure variable, advance to cb_finish, and then
take action based upon your failure variable).  Don't skip past the finish because behavior then is not defined.

#### Conditional

###### Purpose

This action tests a condition and branches to different transactions
based upon the result of the condition.

###### Args

| Arg                    | Type   | Description                                              |
| ---------------------- | ------ | -------------------------------------------------------- |
| term:variable          | string | The variable to compare against                          |
| term:conditional       | string | The type of conditional, see below                       |
| term:conditional_value | int    | The value to compare the variable against                |
| term:conditional_var   | string | The variable to compare the variable against             |
| advance_true           | string | the transaction to advance to if the comparison succeeds |
| advance_false          | string | the transaction to advance to if the comparison fails    |

###### Conditional types

| Type | Description                                             |
| ---- | ------------------------------------------------------- |
| eq   | Match if the variable is equal to the conditional value |
| ne   | Match if not equal                                      |
| gt   | Match if greater than                                   |
| ge   | Match if greater than or equal                          |
| lt   | Match if less than                                      |
| le   | Match if less than or equal                             |

All comparisons are done via Go rules. This means that orderable
types can be ordered (gt, ge, lt, le) and comparable types can be

Don't count on any other types being comparable.

If conditional_var is set conditional_value is ignored.

#### Set

###### Purpose

Set variable values

###### Args

| Arg      | Description                       |
| -------- | --------------------------------- |
| variable | the variable to set               |
| value    | the value to set the variable to  |
| source   | the source to set the variable to |

A variable can be set to values of any type, but it must match the
type the variable was declared with. For example, setting a boolean
variable to a string might not work very well. Setting an int to a
float may work but will have unintended consequences.

If source is set, it will copy the value of the source variable to the
destination variable.

#### Log

###### Purpose

Log something to the log.

###### Args

| Arg | Description |
| ----| ------------ |
| value | What to log |
| loglevel | The loglevel, one of TRACE, DEBUG, INFO, WARNING, ERROR, CRITICAL |

#### Match

###### Purpose

Match a request against a file.

When provided a json file and a json input, the matching occurs
based upon the parsed json and NOT the actual text string. This
means that it can be in any order and still match. For example,
if you have the response:

```
{
  "1": {
    "2": "3",
    "4": "5",
    "6": "7"
  }
}
```

And the match file:

```
{
  "1": {
    "4": "5"
  }
}
```

The other fields will be ignored. ONLY "4": "5" matching is
sufficient for the match to succeed. In other words, if you want
to match on something, you have to provide it as part of the match
file.

Note also that when something is provided to match, it can be of
any complexity, but it has to _exactly_ match the structure of the
response JSON.

###### Args

| Arg             | Type   | Description                                        |
| --------------- | ------ | -------------------------------------------------- |
| match_file      | file   | The file containing the json to match              |
| match_file_type | string   | The type of the data in the match file (json, yaml, string)             |
| advance_true    | string | transaction to advance to if the match succeeds    |
| advance_false   | string | transaction to advance to if the match fails       |
| variable        | string | the variable name containing the response to match |
| response_type   | string   | the type of data contained in match_compare_var (json, yaml, string)   |

#### Math

###### Purpose

Perform a math operation on a variable

Note that all operations are floating point.

The result of the operation is stored in the supplied variable.

The math operations supported are add, subtract, multiply, divide,
and any _one or two operand_ math operation imported by the math library.

###### Args

| Arg      | Type   | Description                                    |
| -------- | ------ | ---------------------------------------------- |
| action   | string | A math operation                               |
| value    | float  | the value on the right side of the operation   |
| variable | string | the variable on the left side of the operation |

#### Wait

###### Purpose

Wait a given number of seconds before proceeding.
Use with caution as this will hang the running test until complete.

###### Args

| Arg      | Description                                                       |
| -------- | ----------------------------------------------------------------- |
| duration | The number of seconds to wait in seconds (floating point allowed) |

Please note that there is a resolution of somewhere around 200ms,
as this is the interval the internal ticker uses.

#### Url

###### Purpose

This is an internal action which is autogenerated in specific cases.
It is generated when a particular transaction has
a "url" field. It is always at the _end_ of the init_action array,
and behaves just as an ordinary action does.

If you wish to specifically include this action, you may do so, but in that
case, do NOT specify a "url" field inside the transaction you are including
this in. If you do so, the behaviors are undefined. Also, do not
include this in the actions list of an on_expected or on_unexpected
clause. If you do so, the behavior is undefined, and is almost
guaranteed to not do what you expect.

Note that the "url" field in the transaction exists to provide for
backwards compatibility with existing tests. For
new tests, you should specify this action specifically.

If a satisfy group is specified, you may optionally include an on_expected
argument. This takes the same format as on_expected in the transaction root.
If it is not specified, then the transaction on_expected is used.

You may not specify an on_unexpected, as when used by a satisfy_group, this
concept makes no sense for an individual action.

###### Note

| Arg              | Description                               |
| ---------------- | ----------------------------------------- |
| url              | the url to be waited for                  |
| save_body        | the variable to save the body into        |
| save_body_as_map | save the body as a map into this variable |
| data             | the file containing the expected data     |
| data_type        | the type of the data ("json" or "yaml")   |

###### Note

When a request is received while a test is running, it is made available
to a url call for processing. There is one request per url call.
The url call will wait until it receives a request, or if one is already
waiting, it will process the waiting call.

Any calls will block until processed by a url action.

### Satisfy Groups

There are situations, in specific kinds of actions, where one might want to perform an
"or" operation. Meaning, you could have two actions of the same type, and want to choose
one to run at runtime. So the concept of satisfy groups were added.1

In order to use one, add a satisfy*group argument at the *root* of any action. It is a string.
If there is no satisfy group, the behavior is the same as it would have been previously. If
a satisfy*group is created, *and* if two or more actions are a member of the same satisfy group,
then a "satisfy" step will be performed before the action is executed. Whichever action is satisfied
as defined by the action, that is the action that will execute.

Currently, the only action that can use this functionality is URL.

### Variables

There is sometimes a need to carry data across transactions inside
a plan. Because of this need, we have thoughtfully provided
"variables" in order to fulfill this need. Define these variables
at the beginning of a plan as follows:

```
variables:
  variable1: value1
  variable2: value2
```

These may be declared, and they may be initially set to a value
if needed. Variables may be any type supported by JSON/Go, which includes
strings, ints, floats, booleans, etc., but the type they are initially
set to may define how they can be used. For example, a boolean can't
be used in a math operation. You probably also can't save a response
into an int. So these are powerful, but use them carefully.

If a variable is not declared, it will be created automatically under
most circumstances. The exception is []interface{} maps: if an
attempt to access or set one with an out of bounds index is made, the array
will not be resized and the access will fail.

There are some cases where an array MUST be declared. This is when the
variable is used for other purposes, such as with stop_var. When
in doubt, just declare the variable with a sane default and see if that
solves the issue.

Inside certain actions, there are "save" directives. The "save"
(or save_response) directives direct the action to save specific
information into a variable to be accessed latter (for example, from
the "match" or "callback" actions). These can also be referenced in
the "url" argument to a callback action via templating (see below).

There may be any number of variables containing any amount of information,
but it should be clear that at the moment the places in which they may
be used are limited.

In certain cases, a variable can also be an interface map (in Go parlance,
a map[string]interface{}). Such variables can be accessed in any place
variables are used using "dot notation". Say, for example, that you have
a variable structure that is like this:

```
variable1: string
variable2: map[string]interface
  sub1: string
  sub2: map[string]interface
    deepersub1: string
  sub3: []interface{}
    [1]: string
    [2]: string
```

You can thus access the variables this way:

```
variable1
variable2.sub1
variable2.sub2.deepersub1
variable2.sub3[1]
variable2.sub3[2]
```

In nearly all cases, int, float32, and float64 are convertible, though
care must be taken as some loss of accuracy is possible when converting
from float to int. This is particularly troublesome when you have an
int variable and are trying to do float operations to it. This probably
won't work.

### Bases

At the root of a config, a map of "bases" may be set. These are
designed to only be accessed inside callback url variables. There
may be global bases, or they may also be set inside plans. If you
set these inside plans, they will override the global bases created
inside the configuration root, with one exception.

If a "TestURL" is specified on the command line or environment
variable when starting the application, it will be added automatically
into the global bases (whether or not it exists) and also any plan bases,
_even if there is a plan bases that overrides the config bases_.

This allows one to use the test URL even if no other bases are provided.

The purpose of these is to be used in the following way:

```
action: callback
  url: <<index .Bases "testurl">>/something
```

So we are then able to connect to the appropriate URL, either for
testing, or if using another URL, for connecting to multiple APIs
without having to hardcode them into the callback action themselves.

### Templates

As mentioned, in a very specific circumstance (right now limited to
the callback url argument) it is possible to use templates to
substitute variables in. There are two different sets that can be
substituted:

```
url: <<index .Bases "blah">>/something
url: http://something.com/api/v1/dosomething/<<index .Variables "blah">>
```

This is useful for being able to send information into an API that
was gathered from an earlier call.

You can also introduce these template calls into files as well! Every
file, when loaded, is run through the templater to substitute variables
in.

### Transactions

A transaction looks a bit like this:

```
txn:
  init_actions:
    - <action>  (see above)
    - <action>
    ...
  url: <a url from which to wait for responses)
  save_body: <An optional variable to save data to as a string>
  save_body_as_map: <An optional variable to save data to as a map>
  data: <the data to expect from the url>
  datatype: <the datatype of the data>
  on_expected:
    response: <the file containing the expected response>
    response_contenttype: <the type of data contained in said response>
    response_code: <which code are we expecting?
    action:
     - <action>
     - <action>
     ...
  on_unexpected:
    ...

```

Notes:

- init_actions declares the actions to be run inside a transaction.
  Historically, a transaction that did not include init_actions
  was a valid transaction, and was treated as a transaction with only
  a "url" action.  This is now deprecated, and every transaction should
  have an init_actions array.  You may see the word "standalone" used for
  a transaction that has an init_actions field but no url field.  This is
  a deprecated term.  Currently, all transactions should be considered
  "standalone".  The ability to have no init_actions will be removed in
  a further release.  (This is why it's called init_actions, it was originally
  designed to run before the transaction started.  That's no longer a thing)
- data is optional. If it is specified, then it is compared with
  the contents of the file. If datatype is set to json or yaml,
  it will compare the data itself instead of a text-based comparison.
- on_expected is run when the data matches and the url matches.
  It sends back the appropriate response code.
- If there is no advance action or no action with an implicit
  advance, then the plan will stall and will require a reset.
  Don't design your plans to do that unless you run a dispose action
  first.
- For further details on "url", see the "url" action above. If you
  specify a url action directly, do not specify "data", "data_type",
  "save_body", or "save_body_as_map", as they will be unused.

If you add a `txninclude` option to a plan, you may specify a
file, similarly to how plans are specified by the `planinclude`
option. It is important to note that the transactions that are inside
the yaml file are loaded _first_. If you want to specify a transaction
that was loded from a txninclude directive as the initial transaction
in a plan, you must include a `start_transaction` directive.
This allows you to specify which plan you want to be the initial one.
If it is not specified, it will be the first transaction loaded, which
will either be the first transaction in the yaml file, or if no yaml
transactions are specified, the first transaction file specified in
the txninclude array.
