# Study Rules Evaluator

Study rules evaluator is a library and a command line tool to evaluate study rules against prepared survey submissions and test the participant state changes.
For now only participant flags is handled and only one single participant (virtual one) is considered.

## 


## Scenario

A Scenario is a test suite for study rules.

It's composed of :

- `state`: An initial participant state (set of flags, possibly empty)
- `time` : An optional starting time time (format YYYY-MM-DD HH:mm:ss) 
- `label`: An optional label 
- `submits`: List of surveys submissions

```json
{
    "time": "2024-11-04 12:00:00",
        "label": "Scenario label (shown in results)",
        "state": {
            "flags": {}
        },
        "submits": [

        ]
}

```

Each Survey Submission, has :

- One survey response in `file`  or embedded using `data` entry
- A set of assertions 
- An optional time to override (or set) the submission time 

File path are relative to the scenario file itself.

```
{
    "file":"response.json"
    "data": {}
    "time": {
        "days": 1
        "weeks": 0,
    },
    "asserts": [
        "assertions expression..."
    ]
}
```

### Time shift

For each submission, you can define time using 2 ways

Fixed time value

```json
"time": {
    "fixed": "2024-11-04 12:00:00"
}
```

Relative time shift from the last submission time

```json
"time": {
    "duration": "1h",
    "days": 1,
    "weeks: 0
}
```

`duration` uses golang `time.Duration` format (1h, 12m, 60s, ...), but maximum unit is hour
2 shortcuts are available `days` and `weeks` adding duration for respectively days (86400 seconds) and weeks (7 days).
All values are added to the time shift, if you set days and weeks or days and duration, the shift will be added.
Default values for each is 0.

### Assertions 
An assertions is represented by an expression to make test on participant state after the rules are applied on the survey response submission.
The language used in golang expression language `Expr` from https://expr-lang.org/ 
assertion expression must return a boolean value

When an assertion is evaluated, variables are availbles:
- `state` : the participant state after submissions are applied
- `previousState` : the participant state before the submissions are applied

For example, the following tests for existence of entry 'bg1' in flags after submission.
````
    'bg1' in state.Flags
```

Example:

```json
{
        "time": "2024-11-04 12:00:00",
        "label": "Scenario label (shown in results)",
        "state": {
            "flags": {}
        },
        "submits": [
            { 
                "file": "../surveys/vaccQ10.json",
                "asserts": [
                ]
            },
            { 
                "file": "../surveys/intake.json",
                "asserts": [
                    "'bg2' in state.Flags"
                ]
            }
        ]
    }
```