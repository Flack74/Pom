package cmd

import (
	"fmt"
	"pom/config"
	"strings"

	"github.com/spf13/cobra"
)

var planCmd = &cobra.Command{
	Use:   "plan",
	Short: "📝 Plan and manage tasks for Pomodoro sessions",
	Long: `📝 Plan and Track Your Tasks

Organize your work by planning tasks for your Pomodoro sessions:
  • Create and manage task list
  • Track time spent on each task
  • Link tasks to Pomodoro sessions
  • View task completion statistics
  • Organize work by projects

Examples:
  pom plan add "Write documentation"    Add a new task
  pom plan list                        List all tasks
  pom plan done task-id                Mark task as complete
  pom plan delete task-id              Remove a task
  pom start -t task-id                 Start session for task`,
	Run: func(cmd *cobra.Command, args []string) {
		// ... existing code ...
	},
}

var addTaskCmd = &cobra.Command{
	Use:   "add [title]",
	Short: "Add a new task",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		title := strings.Join(args, " ")
		description, _ := cmd.Flags().GetString("description")
		tags, _ := cmd.Flags().GetStringSlice("tags")

		if err := config.AddTask(title, description, tags); err != nil {
			fmt.Printf("Error adding task: %v\n", err)
			return
		}

		fmt.Printf("Task added: %s\n", title)
	},
}

var listTasksCmd = &cobra.Command{
	Use:   "list",
	Short: "List all tasks",
	Run: func(cmd *cobra.Command, args []string) {
		showCompleted, _ := cmd.Flags().GetBool("all")
		if err := config.ListTasks(showCompleted); err != nil {
			fmt.Printf("Error listing tasks: %v\n", err)
		}
	},
}

var completeTaskCmd = &cobra.Command{
	Use:   "complete [task-id]",
	Short: "Mark a task as completed",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		taskID := args[0]
		if err := config.CompleteTask(taskID); err != nil {
			fmt.Printf("Error completing task: %v\n", err)
			return
		}

		fmt.Printf("Task %s marked as completed\n", taskID)
	},
}

func init() {
	addTaskCmd.Flags().String("description", "", "Task description")
	addTaskCmd.Flags().StringSlice("tags", []string{}, "Task tags (comma-separated)")

	listTasksCmd.Flags().Bool("all", false, "Show completed tasks")

	planCmd.AddCommand(addTaskCmd)
	planCmd.AddCommand(listTasksCmd)
	planCmd.AddCommand(completeTaskCmd)
	rootCmd.AddCommand(planCmd)
}
