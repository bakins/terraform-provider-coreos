// need to include this so plugin gets loaded
#provider "coreos"

resource "coreos_ami" "test" {
    channel = "stable"
    type = "hvm"
    region = "us-west-2"
}

output "ami" {
    value = "${coreos_ami.test.ami}"
}
