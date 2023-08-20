package runtime

//import (
//  "fmt"
//  "time"

//  "github.com/hofstadter-io/hof/flow/flow"
//)

//type FlowEnricher func(*Runtime, *flow.Flow) error

//func (R *Runtime) EnrichFlows(flows []string, enrich FlowEnricher) error {
//  start := time.Now()
//  defer func() {
//    end := time.Now()
//    R.Stats.Add("enrich/flow", end.Sub(start))
//  }()

//  if R.Flags.Verbosity > 1 {
//    fmt.Println("Runtime.Flow: ", flows)
//    for _, node := range R.Nodes {
//      node.Print()
//    }
//  }

//  // Find only the datamodel nodes
//  // TODO, dedup any references
//  fs := []*flow.Flow{}
//  for _, node := range R.Nodes {
//    // check for Chat root
//    if node.Hof.Chat.Root {

//      if !keepFilter(node, generators) {
//        continue
//      }
//      upgrade := func(n *hof.Node[gen.Generator]) *gen.Generator {
//        v := gen.NewGenerator(n)
//        return v
//      }
//      u := hof.Upgrade[any, gen.Generator](node, upgrade, nil)
//      // we'd like this line in upgrade, but...
//      // how do we make T a Node[T] type (or ensure that it has a hof)
//      // u.T.Hof = u.Hof
//      gen := u.T
//      gen.Node = u
//      gens = append(gens, gen)
//    }
//  }

//  R.Workflows = cs

//  for _, c := range R.Workflows {
//    err := enrich(R, c)
//    if err != nil {
//      return err
//    }
//  }


//  return nil
//}
