package mermaid

import (
	"strconv"
	"strings"
	"time"
)

// Durations Mermaid accepts on a task, in the units it spells them.
const (
	hoursPerDay  = 24
	daysPerWeek  = 7
	daysPerMonth = 30
)

// ganttDateLayouts are the date spellings understood without consulting the
// diagram's dateFormat line, which Mermaid writes in its own token language
// rather than Go's reference layout.
func ganttDateLayouts() []string {
	return []string{
		"2006-01-02",
		"2006/01/02",
		"02-01-2006",
		"01/02/2006",
		"2006-01-02 15:04",
		time.RFC3339,
	}
}

// resolveGanttSchedule fills in every task's Schedule, walking the tasks in
// order so "after <id>" can refer back to a task already placed.
//
// Bars used to be drawn at a fixed width and offset, so a 3-day task and a
// 15-day task looked identical and neither sat under its date on the axis.
func resolveGanttSchedule(gantt *GanttDiagram) {
	ends := make(map[string]time.Time)
	var previousEnd time.Time

	for si := range gantt.Sections {
		for ti := range gantt.Sections[si].Tasks {
			task := &gantt.Sections[si].Tasks[ti]
			start, ok := ganttTaskStart(task.Start, ends, previousEnd)
			if !ok {
				continue
			}
			end := ganttTaskEnd(start, task.Duration)
			task.Schedule = &GanttSchedule{
				Start:       start,
				End:         end,
				IsMilestone: strings.Contains(task.Status, "milestone") || !end.After(start),
			}
			if task.ID != "" {
				ends[task.ID] = end
			}
			previousEnd = end
		}
	}
}

// ganttTaskStart resolves a start field: an explicit date, "after <id>", or an
// empty field meaning "when the previous task ended".
func ganttTaskStart(field string, ends map[string]time.Time, previousEnd time.Time) (time.Time, bool) {
	field = strings.TrimSpace(field)
	if field == "" {
		if previousEnd.IsZero() {
			return time.Time{}, false
		}
		return previousEnd, true
	}
	if after, ok := strings.CutPrefix(strings.ToLower(field), "after "); ok {
		for id := range strings.FieldsSeq(after) {
			if end, known := ends[id]; known {
				return end, true
			}
		}
		if previousEnd.IsZero() {
			return time.Time{}, false
		}
		return previousEnd, true
	}
	return parseGanttDate(field)
}

func parseGanttDate(value string) (time.Time, bool) {
	for _, layout := range ganttDateLayouts() {
		if parsed, err := time.Parse(layout, value); err == nil {
			return parsed, true
		}
	}
	return time.Time{}, false
}

// ganttTaskEnd advances a start date by a duration such as "10d", "2w" or "8h".
// An unreadable or zero duration returns the start date, which marks the task
// as a milestone.
func ganttTaskEnd(start time.Time, duration string) time.Time {
	duration = strings.TrimSpace(strings.ToLower(duration))
	if duration == "" {
		return start
	}
	// A duration may also be spelled as the end date.
	if end, ok := parseGanttDate(duration); ok {
		return end
	}

	unit := duration[len(duration)-1:]
	amount, err := strconv.ParseFloat(strings.TrimSpace(duration[:len(duration)-1]), 64)
	if err != nil {
		return start
	}
	switch unit {
	case "h":
		return start.Add(time.Duration(amount * float64(time.Hour)))
	case "w":
		return start.AddDate(0, 0, int(amount*daysPerWeek))
	case "m":
		return start.AddDate(0, 0, int(amount*daysPerMonth))
	case "d":
		return start.AddDate(0, 0, int(amount))
	default:
		return start
	}
}

// ganttDateRange is the span the chart's x axis covers.
type ganttDateRange struct {
	start time.Time
	end   time.Time
	ok    bool
}

func ganttRangeOf(gantt *GanttDiagram) ganttDateRange {
	out := ganttDateRange{}
	for _, section := range gantt.Sections {
		for _, task := range section.Tasks {
			if task.Schedule == nil {
				continue
			}
			if !out.ok {
				out = ganttDateRange{start: task.Schedule.Start, end: task.Schedule.End, ok: true}
				continue
			}
			if task.Schedule.Start.Before(out.start) {
				out.start = task.Schedule.Start
			}
			if task.Schedule.End.After(out.end) {
				out.end = task.Schedule.End
			}
		}
	}
	// A chart whose tasks all fall on one day still needs a non-zero span to
	// divide by.
	if out.ok && !out.end.After(out.start) {
		out.end = out.start.AddDate(0, 0, 1)
	}
	return out
}

// fraction maps a date onto 0..1 across the chart's span.
func (r ganttDateRange) fraction(at time.Time) float64 {
	if !r.ok {
		return 0
	}
	span := r.end.Sub(r.start)
	if span <= 0 {
		return 0
	}
	value := at.Sub(r.start).Seconds() / span.Seconds()
	return min(max(value, 0), 1)
}

// tickDates returns up to count evenly spaced labels across the span.
func (r ganttDateRange) tickDates(count int) []time.Time {
	if !r.ok || count < 2 {
		return nil
	}
	span := r.end.Sub(r.start)
	ticks := make([]time.Time, 0, count)
	for i := range count {
		offset := time.Duration(float64(span) * float64(i) / float64(count-1))
		ticks = append(ticks, r.start.Add(offset))
	}
	return ticks
}
