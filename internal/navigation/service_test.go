package navigation

import "testing"

func TestHomeIncludesTeacherWorkflowLinks(t *testing.T) {
	home := NewServiceWithFixtures().Home()
	if home.Title != "学校教师资源导航" {
		t.Fatalf("unexpected home title: %q", home.Title)
	}
	if len(home.Primary) != 4 {
		t.Fatalf("expected four primary links, got %d", len(home.Primary))
	}
	if len(home.RightPanel) != 3 {
		t.Fatalf("expected three right panel links, got %d", len(home.RightPanel))
	}
}

func TestImportBatchMakesLinksAvailableByGroup(t *testing.T) {
	service := NewServiceWithFixtures()
	result := service.ImportBatch([]ImportLink{
		{ID: "math", Group: "数学教研", Title: "数学资源", URL: "https://resources.example.test/math", SortOrder: 1},
		{ID: "invalid", Group: "数学教研", Title: "", URL: "https://resources.example.test/invalid", SortOrder: 2},
	})
	if result.Imported != 1 || result.Rejected != 1 {
		t.Fatalf("unexpected import result: %+v", result)
	}
	links := service.List("数学教研")
	if len(links) != 1 || links[0].ID != "math" {
		t.Fatalf("unexpected group links: %+v", links)
	}
}

func TestImportedLinksAreOrderedBySortThenTitle(t *testing.T) {
	service := NewServiceWithFixtures()
	service.ImportBatch([]ImportLink{
		{ID: "c", Group: "同分组", Title: "语文", URL: "https://resources.example.test/c", SortOrder: 5},
		{ID: "a", Group: "同分组", Title: "地理", URL: "https://resources.example.test/a", SortOrder: 5},
		{ID: "b", Group: "同分组", Title: "历史", URL: "https://resources.example.test/b", SortOrder: 5},
	})
	links := service.List("同分组")
	expected := []string{"地理", "历史", "语文"}
	if len(links) != len(expected) {
		t.Fatalf("expected %d links, got %d", len(expected), len(links))
	}
	for index, title := range expected {
		if links[index].Title != title {
			t.Fatalf("expected title %q at position %d, got %q", title, index, links[index].Title)
		}
	}
}

func TestResourceStatusReportsActiveAndInactiveLinks(t *testing.T) {
	service := NewServiceWithFixtures()
	service.ImportBatch([]ImportLink{
		{ID: "inactive", Group: "状态检查", Title: "停用资源", URL: "https://resources.example.test/inactive", Status: StatusInactive, SortOrder: 1},
		{ID: "active", Group: "状态检查", Title: "启用资源", URL: "https://resources.example.test/active", Status: StatusActive, SortOrder: 2},
	})
	results := service.Validate("状态检查")
	if len(results) != 2 {
		t.Fatalf("expected two status results, got %d", len(results))
	}
	if results[0].Valid || results[0].Reason != "resource is inactive" {
		t.Fatalf("unexpected inactive result: %+v", results[0])
	}
	if !results[1].Valid || results[1].Reason != "" {
		t.Fatalf("unexpected active result: %+v", results[1])
	}
}
