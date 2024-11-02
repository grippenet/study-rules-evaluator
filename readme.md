# Study Rules Evaluator

Study rules evaluator is a library and a command line tool to evaluate study rules against prepared survey submissions and test the participant state changes.
For now only participant flags is handled and only one single participant (virtual one) is considered.

## Scenario

A Scenario is a test suite for study rules.

It's composed of :

- An initial participant state (set of flags, possibly empty)
- A set of survey Submissions

Each Submission, has :

- One survey response (embeded in scenario or from a file)
- A set of assertions
- An optional time to override (or set) the submission time 

An assertions is represented by an expression to make test on participant state after the rules are applied on the survey response submission.
The language used in golang expression language `Expr` from https://expr-lang.org/ 
assertion expression must return a boolean value

When an assertion is evaluated, variables are availbles:
- `state` : the participant state after submission is applied