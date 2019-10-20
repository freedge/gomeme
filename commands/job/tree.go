package job

import (
	"fmt"
	"strings"

	"github.com/freedge/gomeme/client"
	"github.com/freedge/gomeme/types"

	"github.com/freedge/gomeme/commands"
)

type node struct {
	jobid         string
	outgoingEdges []*node
	incomingEdges []*node
	status        *types.Status
	analysed      bool
}

type treenode struct {
	thejob *node
	shift  int
}

// jobsStatusCommand retrieve a list of jobs
type jobTreeCommand struct {
	jobsStatusCommand
	nodes map[string]*node
	tree  []treenode
}

func (cmd *jobTreeCommand) addNode(job *types.Status) *node {
	_, found := cmd.nodes[job.JobId]
	if !found {
		cmd.nodes[job.JobId] = &node{job.JobId, make([]*node, 0), make([]*node, 0), job, false}
	}
	return cmd.nodes[job.JobId]
}

func (cmd *jobTreeCommand) addEdge(fromNode *node, to *types.Status) {
	toNode := cmd.addNode(to)
	toNode.incomingEdges = append(toNode.incomingEdges, fromNode)
	fromNode.outgoingEdges = append(fromNode.outgoingEdges, toNode)
}

func (cmd *jobTreeCommand) Run() (i interface{}, err error) {
	// retrieve a list of jobs
	_, err = cmd.GetJobs()
	if err != nil {
		return
	}

	if len(cmd.reply.Statuses) > 42 {
		err = fmt.Errorf("there are too much (%d) jobs selected", len(cmd.reply.Statuses))
		return
	}

	cmd.nodes = make(map[string]*node, 0)

	// retrieve the dependencies between each job. It does not
	// explore jobs outside the initial list
	for it, job := range cmd.reply.Statuses {
		fromNode := cmd.addNode(&cmd.reply.Statuses[it])

		reply := &types.JobsStatusReply{}
		err = client.Call("GET", JOBS_STATUS, nil, map[string]string{
			"neighborhood": "skip",
			"direction":    "depend",
			"jobid":        job.JobId,
		}, &reply)
		if err != nil {
			return
		}
		for it2, subjob := range reply.Statuses {
			// only consider subjob
			if subjob.JobId == job.JobId {
				continue
			}
			// add
			cmd.addEdge(fromNode, &reply.Statuses[it2])
		}
	}

	// do some magic here to simplify the graph?

	// build a tree out of this
	cmd.tree = make([]treenode, 0)
	i = &cmd.tree
	for _, anode := range cmd.nodes {
		// visit starting from ancestor
		if len(anode.incomingEdges) > 0 {
			continue
		}
		cmd.visit(0, anode)
	}
	return
}

func (cmd *jobTreeCommand) visit(shift int, anode *node) {
	anode.analysed = true
	cmd.tree = append(cmd.tree, treenode{anode, shift})

	for _, subnode := range anode.outgoingEdges {
		if !subnode.analysed {
			cmd.visit(shift+1, subnode)
		}
	}
}

func (cmd *jobTreeCommand) PrettyPrint(interface{}) error {
	for _, atreenode := range cmd.tree {
		c := strings.Repeat(" ", 2*atreenode.shift)
		c += atreenode.thejob.status.JobId
		fmt.Printf("%-30.30s %-15.15s %8.8s\n", c, atreenode.thejob.status.Name, atreenode.thejob.status.Status)
	}
	return nil
}

func init() {
	commands.Register("job.tree", &jobTreeCommand{})
}
