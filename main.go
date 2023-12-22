package main

import (
	"fmt"

	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/ec2"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		// Creating security group for EC2 instances
		securityGroupArgs := &ec2.SecurityGroupArgs{
			// Ingress rules are used to control the incoming traffic to your instance
			Ingress: ec2.SecurityGroupIngressArray{
				// Allow all incoming traffic to port 22
				&ec2.SecurityGroupIngressArgs{
					Protocol:   pulumi.String("tcp"),
					FromPort:   pulumi.Int(22),
					ToPort:     pulumi.Int(22),
					CidrBlocks: pulumi.StringArray{pulumi.String("0.0.0.0/0")},
				},
				&ec2.SecurityGroupIngressArgs{
					Protocol:   pulumi.String("tcp"),
					FromPort:   pulumi.Int(8080),
					ToPort:     pulumi.Int(8080),
					CidrBlocks: pulumi.StringArray{pulumi.String("0.0.0.0/0")},
				},
			},
			Egress: ec2.SecurityGroupEgressArray{
				// Open all outgoing traffic to all ports
				&ec2.SecurityGroupEgressArgs{
					// All traffic allowed
					Protocol:   pulumi.String("-1"),
					FromPort:   pulumi.Int(0),
					ToPort:     pulumi.Int(0),
					CidrBlocks: pulumi.StringArray{pulumi.String("0.0.0.0/0")},
				},
			},
		}

		securitGroup, err := ec2.NewSecurityGroup(ctx, "jenkins-sg", securityGroupArgs)
		if err != nil {
			return err
		}

		// Create new Key Pair for EC2 instance ssh access
		keyPair, err := ec2.NewKeyPair(ctx, "local-ssh", &ec2.KeyPairArgs{
			PublicKey: pulumi.String("ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQDHSuUD5v6iQEOBo5epY61p1uOZQX7/zlkcdAEZK+3pNU04+uAm0uGUu6CIJlGDM7DOQ1W0jZFYork5utVexgeKViPUS7E39I4EN9YttD5KOFQyEnhkbLQ450a6FNTky+XPU2v54ywb0dPrqKErWYq5agL97UT7NdRb8s6Ov6lmoSCCYvCXcet48j75XPTLrq75dcneS3f/kfj7NrGQYdo+IiDeciaydRGhExX6qryK6I3RF1KYaelUZNjAOAW+vikvqPT40ea+ix5rllSTg1vXyT7daHIzH27Z7sDRi7v2q1JakN6Py2dpEIGs/UkB+LAIC//LXgI8xUiubPWi+MdcfkbEUOyoKT8nB+5Po7WLvfO9afkCW7nf8kk5+gBLfdflHD1WWOe+KIbSWUtkHqxantyGJybuy4gdDsOhzKZvmE2Jb85ebbhOXTYWIaKbaS89D4UW5D+SxmdEslMNO314Fl0TjPCqx+wNbEPcYUhGIdL54mZ3wyQ23BsYygRHjac= generalpy@GENERALPY-MAIN"),
		})
		if err != nil {
			return err
		}

		jenkinsServer, err := ec2.NewInstance(ctx, "jenkins-server", &ec2.InstanceArgs{
			// Free tier machine t2.micro
			InstanceType:        pulumi.String("t2.micro"),
			VpcSecurityGroupIds: pulumi.StringArray{securitGroup.ID()},
			// Id of os image to use
			Ami:     pulumi.String("ami-0a0f1259dd1c90938"),
			KeyName: keyPair.KeyName,
		})
		if err != nil {
			return err
		}

		fmt.Println("Jenkins server public IP: ", jenkinsServer.PublicIp)
		fmt.Println("Jenkins server public DNS: ", jenkinsServer.PublicDns)

		ctx.Export("publicIp", jenkinsServer.PublicIp)
		ctx.Export("publicDns", jenkinsServer.PublicDns)

		return nil

	})

}
