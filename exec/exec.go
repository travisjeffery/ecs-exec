package exec

import (
	"os"
	"os/exec"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	ec2pkg "github.com/aws/aws-sdk-go/service/ec2"
	ecspkg "github.com/aws/aws-sdk-go/service/ecs"
	"github.com/segmentio/go-log"
)

// Exec the cmd on the container instances for the cluster.
func Exec(cluster *string, cmd *[]string) {
	session := session.New(&aws.Config{Region: aws.String(os.Getenv("AWS_REGION"))})
	ecs := ecspkg.New(session)
	ec2 := ec2pkg.New(session)

	containers := containerInstances(ecs, cluster)
	ids := instanceIds(ecs, containers)
	ips := privateIps(ec2, ids)

	var wg sync.WaitGroup
	for _, ip := range ips {
		wg.Add(1)
		go func(ip string) {
			defer wg.Done()
			l := log.New(os.Stderr, log.INFO, ip)
			args := []string{ip}
			args = append(args, *cmd...)
			cmd := exec.Command("ssh", args...)
			cmd.Stdout = l
			cmd.Stderr = l
			err := cmd.Run()
			if err != nil {
				l.Error("failed: %s", err)
			}
		}(*ip)
	}
	wg.Wait()
}

func privateIps(ec2 *ec2pkg.EC2, ids []*string) []*string {
	output, err := ec2.DescribeInstances(&ec2pkg.DescribeInstancesInput{
		InstanceIds: ids,
	})
	log.Check(err)
	var ips []*string
	for _, r := range output.Reservations {
		for _, i := range r.Instances {
			if i.PrivateIpAddress != nil {
				ips = append(ips, i.PrivateIpAddress)
			}
		}
	}
	return ips
}

func instanceIds(e *ecspkg.ECS, containerInstances []*string) []*string {
	output, err := e.DescribeContainerInstances(&ecspkg.DescribeContainerInstancesInput{
		ContainerInstances: containerInstances,
	})
	log.Check(err)
	var ids []*string
	for _, i := range output.ContainerInstances {
		ids = append(ids, i.Ec2InstanceId)
	}
	return ids
}

func containerInstances(e *ecspkg.ECS, cluster *string) []*string {
	instances := []*string{}
	next := aws.String("")
	for next != nil {
		input := ecspkg.ListContainerInstancesInput{
			Cluster: cluster,
		}
		if next != nil && *next != "" {
			input.NextToken = next
		}
		output, err := e.ListContainerInstances(&input)
		log.Check(err)
		instances = append(instances, output.ContainerInstanceArns...)
		next = output.NextToken
	}
	return instances
}
