# CoreOS AMI provider for Terraform

This [Terraform](http://terraform.io) provider is for dynamically finding the latest CoreOS AMI.

## Status

Development/Testing

## Install

You can install in one of two ways:

```shell
$ go get github.com/bakins/terraform-provider-coreos
```

or

```shell
$ mkdir -p $GOPATH/src/github.com/bakins
$ cd $GOPATH/src/github.com/bakins
$ git clone https://github.com/bakins/terraform-provider-coreos.git
$ go build; cp terraform-provider-coreos /usr/local/bin
# or...
$ go install
```

Note: You may need to add something like the following to `~/.terraformrc` if you get an error about missing the coreos provider when running terraform:

```
providers {
  coreos = "/Users/bakins/go/bin/terraform-provider-coreos"
}
```

## Usage

Simple usage:

```
resource "coreos_ami" "test" {
    channel = "stable"
    type = "hvm"
    region = "us-west-2"
}

output "ami" {
    value = "${coreos_ami.test.ami}"
}
```

The resource `coreos_ami` has the following optional fields:

- `channel` - can be "stable", "beta", or "alpha". defaults to "stable".
- `type` - virtualization type: "pv" or "hvm". defaults to "pv".
- `region` - AWS region. defaults to "us-west-2"

The resulting AMI is availible in the `ami` output of the resource -- `coreos_ami.test.ami` in this example.

More realistic usage:

```
variable "aws_region" {
    description = "AWS Region"
    default = "us-west-2"
}

resource "coreos_ami" "nodes" {
    channel = "stable"
    type = "hvm"
    region = "${var. aws_region}"
}

resource "aws_instance" "mynode" {
    ami = "${coreos_ami.nodes.ami}"
    instance_type = "t2.medium"
...
}
```

