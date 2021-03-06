package aws

import (
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/aws"
	"../configuration"
)

type NodeInfo struct {
	InternalIp  string
	InternalDns string
	ExternalIp  string
	ExternalDns string
	Zone string
}

func MasterNodes(config *configuration.InputVars) []NodeInfo {
	return loadNodesOfType("master", config.ProjectId)
}

func InfraNodes(config *configuration.InputVars) []NodeInfo {
	return loadNodesOfType("infra", config.ProjectId)
}

func AppNodes(config *configuration.InputVars) []NodeInfo {
	return loadNodesOfType("app", config.ProjectId)
}

func BastionNode(config *configuration.InputVars) NodeInfo {
	bastion := loadNodesOfType("bastion", config.ProjectId)
	if len(bastion) < 1 {
		panic("No bastion host found")
	}

	return bastion[0]
}

func loadNodesOfType(type_ string, projectid string) []NodeInfo {
	params := &ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			{
				Name: aws.String("tag:Type"),
				Values: []*string{aws.String(type_)},
			},
			{
				Name: aws.String("tag:ProjectId"),
				Values: []*string{aws.String(projectid)},
			},
			{
				Name: aws.String("instance-state-name"),
				Values: []*string{aws.String("running")},
			},
		},
	}

	result, err := Client.DescribeInstances(params)

	if err != nil {
		panic(err)
	}

	return extractNodes(result)
}

func extractNodes(result *ec2.DescribeInstancesOutput) []NodeInfo {
	numInstances := len(result.Reservations)
	var nodes []NodeInfo

	for i := 0; i < numInstances; i++ {
		current := result.Reservations[i]

		if len(current.Instances) < 1 || len(current.Instances[0].NetworkInterfaces) < 1 {
			continue
		}

		state := *current.Instances[0].State.Name
		if state != "running" {
			continue
		}

		for _,instance := range current.Instances {
			nodes = append(nodes, extractNodeInfo(instance))
		}
	}

	return nodes
}

func extractNodeInfo(current *ec2.Instance) NodeInfo {
	var node NodeInfo

	node.InternalIp = *current.PrivateIpAddress
	node.InternalDns = *current.PrivateDnsName

	if current.PublicIpAddress != nil {
		node.ExternalIp = *current.PublicIpAddress
	}

	if current.PublicDnsName != nil {
		node.ExternalDns = *current.PublicDnsName
	}

	if current.Placement != nil {
		node.Zone = *current.Placement.AvailabilityZone
	}

	return node
}