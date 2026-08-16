package navigation

import (
	"fmt"
	"net/url"
	"sort"
	"strings"

	"golang.org/x/text/collate"
	"golang.org/x/text/language"
)

type Service struct {
	store      *memoryStore
	primary    []HomeSection
	rightPanel []HomeSection
}

func NewService() *Service {
	return NewServiceWithFixtures()
}

func NewServiceWithFixtures() *Service {
	return &Service{
		store: newMemoryStore([]Link{
			{ID: "schedule", Group: "核心教学", Title: "课表", URL: "https://teacher.example.test/schedule", SortOrder: 10, Status: StatusActive, Description: "查看本周课程安排"},
			{ID: "grades", Group: "核心教学", Title: "成绩录入", URL: "https://teacher.example.test/grades", SortOrder: 20, Status: StatusActive, Description: "录入和提交学生成绩"},
			{ID: "calendar", Group: "行政支持", Title: "校历", URL: "https://teacher.example.test/calendar", SortOrder: 10, Status: StatusActive, Description: "查看学期校历和重要日期"},
			{ID: "leave", Group: "行政支持", Title: "请假系统", URL: "https://teacher.example.test/leave", SortOrder: 20, Status: StatusActive, Description: "提交和跟踪请假申请"},
			{ID: "lesson-plan", Group: "教师工具", Title: "教案模板", URL: "https://teacher.example.test/lesson-plans", SortOrder: 10, Status: StatusActive, Description: "下载常用教案模板"},
			{ID: "research", Group: "教师工具", Title: "教研平台", URL: "https://teacher.example.test/research", SortOrder: 20, Status: StatusActive, Description: "参与校内教研活动"},
			{ID: "repair", Group: "教师工具", Title: "设备报修", URL: "https://teacher.example.test/repairs", SortOrder: 30, Status: StatusActive, Description: "提交教室设备报修"},
		}),
		primary: []HomeSection{
			{ID: "schedule", Title: "课表", Description: "查看本周课程安排", URL: "/teacher/schedule"},
			{ID: "grades", Title: "成绩录入", Description: "录入和提交学生成绩", URL: "/teacher/grades"},
			{ID: "calendar", Title: "校历", Description: "查看学期校历和重要日期", URL: "/school/calendar"},
			{ID: "leave", Title: "请假系统", Description: "提交和跟踪请假申请", URL: "/teacher/leave"},
		},
		rightPanel: []HomeSection{
			{ID: "lesson-plan", Title: "教案模板", Description: "下载常用教案模板", URL: "https://teacher.example.test/lesson-plans"},
			{ID: "research", Title: "教研平台", Description: "参与校内教研活动", URL: "https://teacher.example.test/research"},
			{ID: "repair", Title: "设备报修", Description: "提交教室设备报修", URL: "https://teacher.example.test/repairs"},
		},
	}
}

func (s *Service) Home() HomePage {
	primary := append([]HomeSection(nil), s.primary...)
	right := append([]HomeSection(nil), s.rightPanel...)
	return HomePage{Title: "学校教师资源导航", Primary: primary, RightPanel: right, ResourceGroups: s.Groups()}
}

func (s *Service) Groups() []string {
	groups := s.store.groups()
	sort.Strings(groups)
	return groups
}

func (s *Service) List(group string) []Link {
	links := s.store.list(group)
	titleCollator := collate.New(language.SimplifiedChinese)
	sort.SliceStable(links, func(i, j int) bool {
		if links[i].SortOrder != links[j].SortOrder {
			return links[i].SortOrder < links[j].SortOrder
		}
		return titleCollator.CompareString(links[i].Title, links[j].Title) < 0
	})
	return links
}

func (s *Service) ImportBatch(inputs []ImportLink) ImportResult {
	result := ImportResult{Errors: make([]string, 0)}
	for index, input := range inputs {
		link, err := normalizeImport(input)
		if err != nil {
			result.Rejected++
			result.Errors = append(result.Errors, fmt.Sprintf("item %d: %s", index, err))
			continue
		}
		s.store.add(link)
		result.Imported++
	}
	return result
}

func (s *Service) Validate(group string) []ValidationResult {
	links := s.List(group)
	results := make([]ValidationResult, 0, len(links))
	for _, link := range links {
		result := ValidationResult{ID: link.ID, Title: link.Title, Status: link.Status, Valid: true}
		if link.Status != StatusActive {
			result.Valid = false
			result.Reason = "resource is inactive"
		} else if !validURL(link.URL) {
			result.Valid = false
			result.Reason = "resource URL is invalid"
		}
		results = append(results, result)
	}
	return results
}

func normalizeImport(input ImportLink) (Link, error) {
	if strings.TrimSpace(input.ID) == "" {
		return Link{}, fmt.Errorf("id is required")
	}
	if strings.TrimSpace(input.Group) == "" {
		return Link{}, fmt.Errorf("group is required")
	}
	if strings.TrimSpace(input.Title) == "" {
		return Link{}, fmt.Errorf("title is required")
	}
	if !validURL(input.URL) {
		return Link{}, fmt.Errorf("url must use http or https")
	}
	status := input.Status
	if status == "" {
		status = StatusActive
	}
	if status != StatusActive && status != StatusInactive {
		return Link{}, fmt.Errorf("status must be active or inactive")
	}
	return Link{
		ID:          strings.TrimSpace(input.ID),
		Group:       strings.TrimSpace(input.Group),
		Title:       strings.TrimSpace(input.Title),
		URL:         strings.TrimSpace(input.URL),
		SortOrder:   input.SortOrder,
		Status:      status,
		Description: strings.TrimSpace(input.Description),
	}, nil
}

func validURL(raw string) bool {
	parsed, err := url.ParseRequestURI(strings.TrimSpace(raw))
	if err != nil || parsed.Host == "" {
		return false
	}
	return parsed.Scheme == "http" || parsed.Scheme == "https"
}
