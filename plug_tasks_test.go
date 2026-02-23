package obsidian_test

import (
	"testing"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/testutil"

	obsidian "github.com/powerman/goldmark-obsidian"
	"github.com/powerman/goldmark-obsidian/obsast"
)

func TestPlugTasks(t *testing.T) {
	markdown := goldmark.New(
		goldmark.WithRendererOptions(
			html.WithUnsafe(),
		),
		goldmark.WithExtensions(
			obsidian.NewPlugTasks(
				obsidian.WithPlugTasksListClass(""),
				obsidian.WithPlugTasksListItemNotCheckedClass(""),
				obsidian.WithPlugTasksListItemCheckedClass(""),
				obsidian.WithPlugTasksCheckboxClass(""),
			),
		),
	)
	testutil.DoTestCaseFile(markdown, "testdata/plug_tasks.txt", t, testutil.ParseCliCaseArg()...)
}

func TestPlugTasks_Default(t *testing.T) {
	markdown := goldmark.New(
		goldmark.WithRendererOptions(
			html.WithXHTML(),
			html.WithUnsafe(),
		),
		goldmark.WithExtensions(
			obsidian.NewPlugTasks(),
		),
	)
	testutil.DoTestCaseFile(markdown, "testdata/plug_tasks_default.txt", t, testutil.ParseCliCaseArg()...)
}

func TestPlugTasks_Options(t *testing.T) {
	markdown := goldmark.New(
		goldmark.WithRendererOptions(
			html.WithUnsafe(),
		),
		goldmark.WithExtensions(
			obsidian.NewPlugTasks(
				obsidian.WithPlugTasksStatusType('X', obsast.PlugTasksStatusTypeDone),
				obsidian.WithPlugTasksStatusTypes(map[rune]obsast.PlugTasksStatusType{
					'~': obsast.PlugTasksStatusTypeInProgress,
					'b': obsast.PlugTasksStatusTypeNonTask,
				}),
				obsidian.WithPlugTasksListClass(""),
				obsidian.WithPlugTasksListItemNotCheckedClass("tasks-item not-checked"),
				obsidian.WithPlugTasksListItemCheckedClass("tasks-item"),
				obsidian.WithPlugTasksListItemStatusAttr("data-tasks"),
				obsidian.WithPlugTasksCheckboxClass("tasks-checkbox"),
			),
		),
	)
	testutil.DoTestCaseFile(markdown, "testdata/plug_tasks_options.txt", t, testutil.ParseCliCaseArg()...)
}
