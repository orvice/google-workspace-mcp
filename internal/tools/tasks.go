package tools

import (
	"context"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"go.orx.me/mcp/google-workspace/internal/utils"
	"google.golang.org/api/tasks/v1"
)

// ListTaskListsInput defines input for list_task_lists tool
type ListTaskListsInput struct {
	Email string `json:"email" jsonschema:"required,description=Email address to access Google Tasks"`
}

// ListTaskListsOutput defines output for list_task_lists tool
type ListTaskListsOutput struct {
	TaskLists string `json:"taskLists" jsonschema:"description=List of task lists"`
}

// ListTasksInput defines input for list_tasks tool
type ListTasksInput struct {
	Email      string `json:"email" jsonschema:"required,description=Email address to access Google Tasks"`
	TaskListID string `json:"taskListId" jsonschema:"required,description=Task list identifier (use @default for the default task list)"`
}

// ListTasksOutput defines output for list_tasks tool
type ListTasksOutput struct {
	Tasks string `json:"tasks" jsonschema:"description=List of tasks"`
}

// CreateTaskInput defines input for create_task tool
type CreateTaskInput struct {
	Email      string `json:"email" jsonschema:"required,description=Email address to access Google Tasks"`
	TaskListID string `json:"taskListId" jsonschema:"required,description=Task list identifier (use @default for the default task list)"`
	Title      string `json:"title" jsonschema:"required,description=Title of the task"`
	Notes      string `json:"notes,omitempty" jsonschema:"description=Notes describing the task"`
	Due        string `json:"due,omitempty" jsonschema:"description=Due date in RFC3339 format (e.g. 2024-12-31T00:00:00Z)"`
}

// CreateTaskOutput defines output for create_task tool
type CreateTaskOutput struct {
	Result string `json:"result" jsonschema:"description=Result of task creation"`
}


// UpdateTaskInput defines input for update_task tool
type UpdateTaskInput struct {
	Email      string `json:"email" jsonschema:"required,description=Email address to access Google Tasks"`
	TaskListID string `json:"taskListId" jsonschema:"required,description=Task list identifier"`
	TaskID     string `json:"taskId" jsonschema:"required,description=Task identifier"`
	Title      string `json:"title,omitempty" jsonschema:"description=New title of the task"`
	Notes      string `json:"notes,omitempty" jsonschema:"description=New notes for the task"`
	Status     string `json:"status,omitempty" jsonschema:"description=Task status: needsAction or completed"`
	Due        string `json:"due,omitempty" jsonschema:"description=Due date in RFC3339 format"`
}

// UpdateTaskOutput defines output for update_task tool
type UpdateTaskOutput struct {
	Result string `json:"result" jsonschema:"description=Result of task update"`
}

// DeleteTaskInput defines input for delete_task tool
type DeleteTaskInput struct {
	Email      string `json:"email" jsonschema:"required,description=Email address to access Google Tasks"`
	TaskListID string `json:"taskListId" jsonschema:"required,description=Task list identifier"`
	TaskID     string `json:"taskId" jsonschema:"required,description=Task identifier"`
}

// DeleteTaskOutput defines output for delete_task tool
type DeleteTaskOutput struct {
	Result string `json:"result" jsonschema:"description=Result of task deletion"`
}

// CompleteTaskInput defines input for complete_task tool
type CompleteTaskInput struct {
	Email      string `json:"email" jsonschema:"required,description=Email address to access Google Tasks"`
	TaskListID string `json:"taskListId" jsonschema:"required,description=Task list identifier"`
	TaskID     string `json:"taskId" jsonschema:"required,description=Task identifier"`
}

// CompleteTaskOutput defines output for complete_task tool
type CompleteTaskOutput struct {
	Result string `json:"result" jsonschema:"description=Result of marking task as completed"`
}

// ListTaskLists handles the list_task_lists tool call
func ListTaskLists(ctx context.Context, req *mcp.CallToolRequest, input ListTaskListsInput) (*mcp.CallToolResult, ListTaskListsOutput, error) {
	srv, err := utils.NewTasksClient(input.Email)
	if err != nil {
		return nil, ListTaskListsOutput{}, err
	}

	taskLists, err := srv.Tasklists.List().MaxResults(100).Do()
	if err != nil {
		return nil, ListTaskListsOutput{}, fmt.Errorf("failed to list task lists: %w", err)
	}

	var resp string
	if len(taskLists.Items) == 0 {
		resp = "No task lists found."
	} else {
		resp = "Task lists:\n"
		for _, tl := range taskLists.Items {
			resp += fmt.Sprintf("- %s (ID: %s)\n", tl.Title, tl.Id)
		}
	}

	return nil, ListTaskListsOutput{TaskLists: resp}, nil
}


// ListTasks handles the list_tasks tool call
func ListTasks(ctx context.Context, req *mcp.CallToolRequest, input ListTasksInput) (*mcp.CallToolResult, ListTasksOutput, error) {
	srv, err := utils.NewTasksClient(input.Email)
	if err != nil {
		return nil, ListTasksOutput{}, err
	}

	taskList, err := srv.Tasks.List(input.TaskListID).MaxResults(100).Do()
	if err != nil {
		return nil, ListTasksOutput{}, fmt.Errorf("failed to list tasks: %w", err)
	}

	var resp string
	if len(taskList.Items) == 0 {
		resp = "No tasks found."
	} else {
		resp = "Tasks:\n"
		for _, t := range taskList.Items {
			status := "[ ]"
			if t.Status == "completed" {
				status = "[x]"
			}
			due := ""
			if t.Due != "" {
				due = fmt.Sprintf(" (Due: %s)", t.Due)
			}
			resp += fmt.Sprintf("%s %s%s (ID: %s)\n", status, t.Title, due, t.Id)
		}
	}

	return nil, ListTasksOutput{Tasks: resp}, nil
}

// CreateTask handles the create_task tool call
func CreateTask(ctx context.Context, req *mcp.CallToolRequest, input CreateTaskInput) (*mcp.CallToolResult, CreateTaskOutput, error) {
	srv, err := utils.NewTasksClient(input.Email)
	if err != nil {
		return nil, CreateTaskOutput{}, err
	}

	task := &tasks.Task{
		Title: input.Title,
		Notes: input.Notes,
		Due:   input.Due,
	}

	createdTask, err := srv.Tasks.Insert(input.TaskListID, task).Do()
	if err != nil {
		return nil, CreateTaskOutput{}, fmt.Errorf("failed to create task: %w", err)
	}

	resp := fmt.Sprintf("Task created successfully:\nTitle: %s\nID: %s", createdTask.Title, createdTask.Id)
	if createdTask.Due != "" {
		resp += fmt.Sprintf("\nDue: %s", createdTask.Due)
	}

	return nil, CreateTaskOutput{Result: resp}, nil
}

// UpdateTask handles the update_task tool call
func UpdateTask(ctx context.Context, req *mcp.CallToolRequest, input UpdateTaskInput) (*mcp.CallToolResult, UpdateTaskOutput, error) {
	srv, err := utils.NewTasksClient(input.Email)
	if err != nil {
		return nil, UpdateTaskOutput{}, err
	}

	// Get existing task first
	existingTask, err := srv.Tasks.Get(input.TaskListID, input.TaskID).Do()
	if err != nil {
		return nil, UpdateTaskOutput{}, fmt.Errorf("failed to get task: %w", err)
	}

	// Update fields if provided
	if input.Title != "" {
		existingTask.Title = input.Title
	}
	if input.Notes != "" {
		existingTask.Notes = input.Notes
	}
	if input.Status != "" {
		existingTask.Status = input.Status
	}
	if input.Due != "" {
		existingTask.Due = input.Due
	}

	updatedTask, err := srv.Tasks.Update(input.TaskListID, input.TaskID, existingTask).Do()
	if err != nil {
		return nil, UpdateTaskOutput{}, fmt.Errorf("failed to update task: %w", err)
	}

	resp := fmt.Sprintf("Task updated successfully:\nTitle: %s\nStatus: %s\nID: %s", updatedTask.Title, updatedTask.Status, updatedTask.Id)

	return nil, UpdateTaskOutput{Result: resp}, nil
}


// DeleteTask handles the delete_task tool call
func DeleteTask(ctx context.Context, req *mcp.CallToolRequest, input DeleteTaskInput) (*mcp.CallToolResult, DeleteTaskOutput, error) {
	srv, err := utils.NewTasksClient(input.Email)
	if err != nil {
		return nil, DeleteTaskOutput{}, err
	}

	err = srv.Tasks.Delete(input.TaskListID, input.TaskID).Do()
	if err != nil {
		return nil, DeleteTaskOutput{}, fmt.Errorf("failed to delete task: %w", err)
	}

	return nil, DeleteTaskOutput{Result: "Task deleted successfully"}, nil
}

// CompleteTask handles the complete_task tool call
func CompleteTask(ctx context.Context, req *mcp.CallToolRequest, input CompleteTaskInput) (*mcp.CallToolResult, CompleteTaskOutput, error) {
	srv, err := utils.NewTasksClient(input.Email)
	if err != nil {
		return nil, CompleteTaskOutput{}, err
	}

	// Get existing task first
	existingTask, err := srv.Tasks.Get(input.TaskListID, input.TaskID).Do()
	if err != nil {
		return nil, CompleteTaskOutput{}, fmt.Errorf("failed to get task: %w", err)
	}

	existingTask.Status = "completed"

	updatedTask, err := srv.Tasks.Update(input.TaskListID, input.TaskID, existingTask).Do()
	if err != nil {
		return nil, CompleteTaskOutput{}, fmt.Errorf("failed to complete task: %w", err)
	}

	resp := fmt.Sprintf("Task marked as completed:\nTitle: %s\nID: %s", updatedTask.Title, updatedTask.Id)

	return nil, CompleteTaskOutput{Result: resp}, nil
}

// RegisterTasksTools registers all Tasks-related tools with the MCP server
func RegisterTasksTools(server *mcp.Server) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "list_task_lists",
		Description: "List all Google Tasks task lists for a user",
	}, ListTaskLists)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "list_tasks",
		Description: "List all tasks in a specific task list",
	}, ListTasks)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "create_task",
		Description: "Create a new task in a task list",
	}, CreateTask)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "update_task",
		Description: "Update an existing task",
	}, UpdateTask)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "delete_task",
		Description: "Delete a task from a task list",
	}, DeleteTask)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "complete_task",
		Description: "Mark a task as completed",
	}, CompleteTask)
}
