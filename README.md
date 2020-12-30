# cfn-diff

`cfn-diff` is a diff tool for AWS CloudFormation resources.

## Installation

**go get**

```sh
$ go get github.com/kmrok/cfn-diff/cmd/cfn-diff
```

**build**

```sh
$ git clone https://github.com/kmrok/cfn-diff.git
$ make install
```

## Getting Started

Export environment variables for your AWS credentials:

```
$ export AWS_ACCESS_KEY_ID=XXXXXX
$ export AWS_SECRET_ACCESS_KEY=XXXXXX
$ export AWS_REGION=XXXXXX
```

Add `.cfndiff.yml` ( or `cfndiff.yml` ) file to your repository:

```yaml
# .cfndiff.yml

run:
  # Set to true if the system terminates abnormally when a difference is detected. (Optional)
  enable_ci_mode: true

# Describe the name of the stack for which you want to detect drift. (Optional)
# Regular expressions can be used.
# Default: All stacks
stack_with_drift_detection:
  - stg-.*

# Describe the name of the stack you do not want to detect drift. (Optional)
# Regular expressions can be used.
stack_without_drift_detection:
  - stg-aurora
  - stg-.*-ecs-task-def
  - stg-alb

# Describe the correspondence between the stack name and the template file. (Optional)
stack_template_maps:
  - stack_name: stg-ecr
    template_name: ./ecr.yaml
  - stack_name: stg-ecs-task-def
    template_name: ./ecs-task-def.yaml
```

## Usage

### Drift detection

```sh
$ cfn-diff drift
$ cfn-diff drift -f config.yml
```

### Changeset difference detection

```sh
$ cfn-diff changeset
$ cfn-diff changeset -f config.yml
```
