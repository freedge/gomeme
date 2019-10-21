package job

import (
	"fmt"
	"math"
	"strings"

	"github.com/freedge/gomeme/client"
	"github.com/freedge/gomeme/types"

	"github.com/freedge/gomeme/commands"
)

// we implement a graph of jobs
type node struct {
	jobid         string
	outgoingEdges []*node
	incomingEdges []*node
	status        *types.Status
	toposorted    bool
	dist          int64 // the distance to ancestor
	analysed      bool
}

// treenode is the tree we actually want to print
type treenode struct {
	thejob *node
	shift  int
}

// jobTreeCommand retrieve a list of jobs like jobsStatus, but also get dependencies
// we want to output a complex graph of jobs (which should be a DAG : directed, acyclic graph)
// in the most meaning full manner. Ideally, we should print on the tree, the "critical path"
// of our job execution. An edge on a graph would be weighted with the likely execution time of the job
// so that if we have to choose between A -> B -> C -> D and E -> D, we output the latter
// if we know job E is itself longer than A -> B -> C.
// As we don/t have figures on statistics per job, we're just going to give the same weight to each
// edge, so it will always be the longest chain that is displayed.
type jobTreeCommand struct {
	jobsStatusCommand
	nodes           map[string]*node // our graph of jobs
	toposortedStack []*node          // the same graph, but sorted topologically
	tree            []treenode       // the graph, now reduced into a tree
}

func (cmd *jobTreeCommand) addNode(job *types.Status) *node {
	_, found := cmd.nodes[job.JobId]
	if !found {
		cmd.nodes[job.JobId] = &node{job.JobId, make([]*node, 0), make([]*node, 0), job, false, math.MinInt64, false}
	}
	return cmd.nodes[job.JobId]
}

func (cmd *jobTreeCommand) addEdge(fromNode *node, to *types.Status) {
	toNode := cmd.addNode(to)
	toNode.incomingEdges = append(toNode.incomingEdges, fromNode)
	fromNode.outgoingEdges = append(fromNode.outgoingEdges, toNode)
}

func (cmd *jobTreeCommand) toposortOneNode(anode *node) {
	if anode.toposorted {
		return
	}
	anode.toposorted = true
	for _, anode := range anode.outgoingEdges {
		cmd.toposortOneNode(anode)
	}
	cmd.toposortedStack = append(cmd.toposortedStack, anode)
}

// apply topological sorting on our cmd.nodes graph into a toposorted stack.
// toposortedStack should be browsed backward to get items in descending order
func (cmd *jobTreeCommand) toposort() {
	for _, anode := range cmd.nodes {
		cmd.toposortOneNode(anode)
	}
}

const jobWeight = 1

func (cmd *jobTreeCommand) Run() (i interface{}, err error) {
	// retrieve a list of jobs
	_, err = cmd.GetJobs()
	if err != nil {
		return
	}

	if len(cmd.reply.Statuses) > 100 {
		err = fmt.Errorf("there are too much (%d) jobs selected", len(cmd.reply.Statuses))
		return
	}
	if len(cmd.reply.Statuses) == 0 {
		err = fmt.Errorf("no job found")
		return
	}

	cmd.nodes = make(map[string]*node, 0)

	// retrieve the dependencies between each job. It does not
	// explore jobs outside the initial list
	for it, job := range cmd.reply.Statuses {
		fromNode := cmd.addNode(&cmd.reply.Statuses[it])

		reply := &types.JobsStatusReply{}
		if cmd.verbose {
			fmt.Printf("retrieving job % 3d/%d\n", it+1, len(cmd.reply.Statuses))
		}
		err = client.Call("GET", JOBS_STATUS, nil, map[string]string{
			"neighborhood": "1",
			"direction":    "depend",
			"depth":        "1",
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
			if cmd.verbose {
				fmt.Println("adding edge")
			}
			// add
			cmd.addEdge(fromNode, &reply.Statuses[it2])
		}
	}

	// topological sort of our graph
	cmd.toposortedStack = make([]*node, 0, len(cmd.nodes))
	cmd.toposort()

	// put a weight on all
	for it := range cmd.toposortedStack {
		src := cmd.toposortedStack[len(cmd.toposortedStack)-1-it]
		if src.dist < 0 {
			src.dist = 0
		}
		for _, edge := range src.outgoingEdges {
			if edge.dist < src.dist+jobWeight { /* instead of 1, we could put the number of minutes for src job to complete */
				edge.dist = src.dist + jobWeight
			}
		}
	}

	// build a tree out of this
	cmd.tree = make([]treenode, 0, len(cmd.nodes))
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
		if !subnode.analysed && int(subnode.dist) == (shift+1) {
			cmd.visit(shift+jobWeight, subnode)
		}
	}
}

func (cmd *jobTreeCommand) PrettyPrint(interface{}) error {
	fmt.Printf("%-30.30s %-40.40s %s\n", "jobid", "folder/name", "status")
	fmt.Println(strings.Repeat("-", 90))
	for _, atreenode := range cmd.tree {
		c := strings.Repeat(" ", 2*atreenode.shift)
		c += atreenode.thejob.status.JobId
		fmt.Printf("%-30.30s %-40.40s %s\n", c, atreenode.thejob.status.Folder+"/"+atreenode.thejob.status.Name, atreenode.thejob.status.Status)
	}
	return nil
}

func init() {
	commands.Register("job.tree", &jobTreeCommand{})
}
