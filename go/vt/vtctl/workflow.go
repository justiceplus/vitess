/*
Copyright 2019 The Vitess Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package vtctl

import (
	"context"
	"fmt"

	"github.com/spf13/pflag"

	"vitess.io/vitess/go/vt/workflow"
	"vitess.io/vitess/go/vt/wrangler"
)

// This file contains the workflows command group for vtctl.

const workflowsGroupName = "Workflows"

var (
	// WorkflowManager contains our manager. It needs to be set or else all
	// commands will be disabled.
	WorkflowManager *workflow.Manager
)

func init() {
	addCommandGroup(workflowsGroupName)

	addCommand(workflowsGroupName, command{
		name:   "WorkflowCreate",
		method: commandWorkflowCreate,
		params: "[--skip_start] <factoryName> [parameters...]",
		help:   "Creates the workflow with the provided parameters. The workflow is also started, unless -skip_start is specified.",
	})
	addCommand(workflowsGroupName, command{
		name:   "WorkflowStart",
		method: commandWorkflowStart,
		params: "<uuid>",
		help:   "Starts the workflow.",
	})
	addCommand(workflowsGroupName, command{
		name:   "WorkflowStop",
		method: commandWorkflowStop,
		params: "<uuid>",
		help:   "Stops the workflow.",
	})
	addCommand(workflowsGroupName, command{
		name:   "WorkflowDelete",
		method: commandWorkflowDelete,
		params: "<uuid>",
		help:   "Deletes the finished or not started workflow.",
	})
	addCommand(workflowsGroupName, command{
		name:   "WorkflowWait",
		method: commandWorkflowWait,
		params: "<uuid>",
		help:   "Waits for the workflow to finish.",
	})

	addCommand(workflowsGroupName, command{
		name:   "WorkflowTree",
		method: commandWorkflowTree,
		params: "",
		help:   "Displays a JSON representation of the workflow tree.",
	})
	addCommand(workflowsGroupName, command{
		name:   "WorkflowAction",
		method: commandWorkflowAction,
		params: "<path> <name>",
		help:   "Sends the provided action name on the specified path.",
	})
}

func commandWorkflowCreate(ctx context.Context, wr *wrangler.Wrangler, subFlags *pflag.FlagSet, args []string) error {
	if WorkflowManager == nil {
		return fmt.Errorf("no workflow.Manager registered")
	}

	skipStart := subFlags.Bool("skip_start", false, "If set, the workflow will not be started.")
	if err := subFlags.Parse(args); err != nil {
		return err
	}
	if subFlags.NArg() < 1 {
		return fmt.Errorf("the <factoryName> argument is required for the WorkflowCreate command")
	}
	factoryName := subFlags.Arg(0)

	uuid, err := WorkflowManager.Create(ctx, factoryName, subFlags.Args()[1:])
	if err != nil {
		return err
	}
	wr.Logger().Printf("uuid: %v\n", uuid)

	if !*skipStart {
		return WorkflowManager.Start(ctx, uuid)
	}
	return nil
}

func commandWorkflowStart(ctx context.Context, wr *wrangler.Wrangler, subFlags *pflag.FlagSet, args []string) error {
	if WorkflowManager == nil {
		return fmt.Errorf("no workflow.Manager registered")
	}

	if err := subFlags.Parse(args); err != nil {
		return err
	}
	if subFlags.NArg() != 1 {
		return fmt.Errorf("the <uuid> argument is required for the WorkflowStart command")
	}
	uuid := subFlags.Arg(0)
	return WorkflowManager.Start(ctx, uuid)
}

func commandWorkflowStop(ctx context.Context, wr *wrangler.Wrangler, subFlags *pflag.FlagSet, args []string) error {
	if WorkflowManager == nil {
		return fmt.Errorf("no workflow.Manager registered")
	}

	if err := subFlags.Parse(args); err != nil {
		return err
	}
	if subFlags.NArg() != 1 {
		return fmt.Errorf("the <uuid> argument is required for the WorkflowStop command")
	}
	uuid := subFlags.Arg(0)
	return WorkflowManager.Stop(ctx, uuid)
}

func commandWorkflowDelete(ctx context.Context, wr *wrangler.Wrangler, subFlags *pflag.FlagSet, args []string) error {
	if WorkflowManager == nil {
		return fmt.Errorf("no workflow.Manager registered")
	}

	if err := subFlags.Parse(args); err != nil {
		return err
	}
	if subFlags.NArg() != 1 {
		return fmt.Errorf("the <uuid> argument is required for the WorkflowDelete command")
	}
	uuid := subFlags.Arg(0)
	return WorkflowManager.Delete(ctx, uuid)
}

func commandWorkflowWait(ctx context.Context, wr *wrangler.Wrangler, subFlags *pflag.FlagSet, args []string) error {
	if WorkflowManager == nil {
		return fmt.Errorf("no workflow.Manager registered")
	}

	if err := subFlags.Parse(args); err != nil {
		return err
	}
	if subFlags.NArg() != 1 {
		return fmt.Errorf("the <uuid> argument is required for the WorkflowWait command")
	}
	uuid := subFlags.Arg(0)
	return WorkflowManager.Wait(ctx, uuid)
}

func commandWorkflowTree(ctx context.Context, wr *wrangler.Wrangler, subFlags *pflag.FlagSet, args []string) error {
	if WorkflowManager == nil {
		return fmt.Errorf("no workflow.Manager registered")
	}

	if err := subFlags.Parse(args); err != nil {
		return err
	}
	if subFlags.NArg() != 0 {
		return fmt.Errorf("the WorkflowTree command takes no parameter")
	}

	tree, err := WorkflowManager.NodeManager().GetFullTree()
	if err != nil {
		return err
	}
	wr.Logger().Printf("%v\n", string(tree))
	return nil
}

func commandWorkflowAction(ctx context.Context, wr *wrangler.Wrangler, subFlags *pflag.FlagSet, args []string) error {
	if WorkflowManager == nil {
		return fmt.Errorf("no workflow.Manager registered")
	}

	if err := subFlags.Parse(args); err != nil {
		return err
	}
	if subFlags.NArg() != 2 {
		return fmt.Errorf("the <path> and <name> arguments are required for the WorkflowAction command")
	}
	ap := &workflow.ActionParameters{
		Path: subFlags.Arg(0),
		Name: subFlags.Arg(1),
	}

	return WorkflowManager.NodeManager().Action(ctx, ap)
}
