package main

import (
	"archive/zip"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

var (
	reSlideMasterPath = regexp.MustCompile(`^ppt/slideMasters/slideMaster(\d+)\.xml$`)
	reSlideLayoutPath = regexp.MustCompile(`^ppt/slideLayouts/slideLayout(\d+)\.xml$`)
	reCTMaster        = regexp.MustCompile(`/ppt/slideMasters/slideMaster(\d+)\.xml`)
	reCTLayout        = regexp.MustCompile(`/ppt/slideLayouts/slideLayout(\d+)\.xml`)
	reRel             = regexp.MustCompile(`<Relationship[^>]*Id="([^"]+)"[^>]*Type="([^"]+)"[^>]*Target="([^"]+)"[^>]*/?>`)
	reMasterRID       = regexp.MustCompile(`<p:sldMasterId[^>]*r:id="([^"]+)"`)
)

const (
	relTypeSlideMaster = "http://schemas.openxmlformats.org/officeDocument/2006/relationships/slideMaster"
	relTypeSlideLayout = "http://schemas.openxmlformats.org/officeDocument/2006/relationships/slideLayout"
)

type pkgInfo struct {
	path string

	entrySet map[string]struct{}

	masters []int
	layouts []int

	ctMasterOverrides map[int]struct{}
	ctLayoutOverrides map[int]struct{}

	presentationMasterRIDs []string
	presentationRelsByID   map[string]relationship

	slideLayoutTargets map[string]string // slide rel part -> layout target
	layoutToMaster     map[string]string // layout part -> master target
}

type relationship struct {
	typeURI string
	target  string
}

func main() {
	baseline := flag.String("baseline", "examples/output/pp_multi_master_reference.pptx", "PowerPoint-authored baseline .pptx")
	candidate := flag.String("candidate", "examples/output/36_multi_master_smoke.pptx", "candidate .pptx to validate")
	flag.Parse()

	baseInfo, err := inspectPackage(*baseline)
	if err != nil {
		fmt.Fprintf(os.Stderr, "baseline inspect failed: %v\n", err)
		os.Exit(1)
	}
	candInfo, err := inspectPackage(*candidate)
	if err != nil {
		fmt.Fprintf(os.Stderr, "candidate inspect failed: %v\n", err)
		os.Exit(1)
	}

	errs := validateAgainstBaseline(baseInfo, candInfo)
	printSummary(baseInfo, candInfo)

	if len(errs) > 0 {
		fmt.Println("RESULT=FAIL")
		for _, e := range errs {
			fmt.Printf("- %s\n", e)
		}
		os.Exit(1)
	}
	fmt.Println("RESULT=PASS")
}

func inspectPackage(path string) (*pkgInfo, error) {
	zr, err := zip.OpenReader(path)
	if err != nil {
		return nil, fmt.Errorf("open zip %q: %w", path, err)
	}
	defer func() {
		_ = zr.Close()
	}()

	info := &pkgInfo{
		path:                 path,
		entrySet:             make(map[string]struct{}, len(zr.File)),
		ctMasterOverrides:    map[int]struct{}{},
		ctLayoutOverrides:    map[int]struct{}{},
		presentationRelsByID: map[string]relationship{},
		slideLayoutTargets:   map[string]string{},
		layoutToMaster:       map[string]string{},
	}

	for _, f := range zr.File {
		name := f.Name
		info.entrySet[name] = struct{}{}
		if m := reSlideMasterPath.FindStringSubmatch(name); len(m) == 2 {
			idx, _ := strconv.Atoi(m[1])
			info.masters = append(info.masters, idx)
		}
		if m := reSlideLayoutPath.FindStringSubmatch(name); len(m) == 2 {
			idx, _ := strconv.Atoi(m[1])
			info.layouts = append(info.layouts, idx)
		}
	}
	sort.Ints(info.masters)
	sort.Ints(info.layouts)

	ct, err := readZipEntry(zr.File, "[Content_Types].xml")
	if err != nil {
		return nil, err
	}
	for _, m := range reCTMaster.FindAllStringSubmatch(ct, -1) {
		idx, _ := strconv.Atoi(m[1])
		info.ctMasterOverrides[idx] = struct{}{}
	}
	for _, m := range reCTLayout.FindAllStringSubmatch(ct, -1) {
		idx, _ := strconv.Atoi(m[1])
		info.ctLayoutOverrides[idx] = struct{}{}
	}

	presXML, err := readZipEntry(zr.File, "ppt/presentation.xml")
	if err != nil {
		return nil, err
	}
	for _, m := range reMasterRID.FindAllStringSubmatch(presXML, -1) {
		info.presentationMasterRIDs = append(info.presentationMasterRIDs, m[1])
	}

	presRels, err := readZipEntry(zr.File, "ppt/_rels/presentation.xml.rels")
	if err != nil {
		return nil, err
	}
	for _, m := range reRel.FindAllStringSubmatch(presRels, -1) {
		info.presentationRelsByID[m[1]] = relationship{typeURI: m[2], target: m[3]}
	}

	for entry := range info.entrySet {
		if !strings.HasPrefix(entry, "ppt/slides/_rels/slide") || !strings.HasSuffix(entry, ".xml.rels") {
			continue
		}
		txt, readErr := readZipEntry(zr.File, entry)
		if readErr != nil {
			return nil, readErr
		}
		for _, m := range reRel.FindAllStringSubmatch(txt, -1) {
			if m[2] == relTypeSlideLayout {
				info.slideLayoutTargets[entry] = m[3]
				break
			}
		}
	}

	for entry := range info.entrySet {
		if !strings.HasPrefix(entry, "ppt/slideLayouts/_rels/slideLayout") || !strings.HasSuffix(entry, ".xml.rels") {
			continue
		}
		txt, readErr := readZipEntry(zr.File, entry)
		if readErr != nil {
			return nil, readErr
		}
		for _, m := range reRel.FindAllStringSubmatch(txt, -1) {
			if m[2] == relTypeSlideMaster {
				layoutName := strings.TrimSuffix(filepath.Base(entry), ".rels")
				info.layoutToMaster[layoutName] = m[3]
				break
			}
		}
	}

	return info, nil
}

func validateAgainstBaseline(base, cand *pkgInfo) []string {
	errs := make([]string, 0)

	if len(base.masters) < 2 {
		errs = append(errs, "baseline does not look multi-master (expected >=2 masters)")
	}
	if len(cand.masters) < 2 {
		errs = append(errs, fmt.Sprintf("candidate has only %d master(s); expected multi-master (>=2)", len(cand.masters)))
	}

	for _, idx := range cand.masters {
		if _, ok := cand.ctMasterOverrides[idx]; !ok {
			errs = append(errs, fmt.Sprintf("missing [Content_Types] override for /ppt/slideMasters/slideMaster%d.xml", idx))
		}
	}
	for _, idx := range cand.layouts {
		if _, ok := cand.ctLayoutOverrides[idx]; !ok {
			errs = append(errs, fmt.Sprintf("missing [Content_Types] override for /ppt/slideLayouts/slideLayout%d.xml", idx))
		}
	}

	for _, rid := range cand.presentationMasterRIDs {
		rel, ok := cand.presentationRelsByID[rid]
		if !ok {
			errs = append(errs, fmt.Sprintf("presentation.xml references master r:id=%s but relation is missing", rid))
			continue
		}
		if rel.typeURI != relTypeSlideMaster {
			errs = append(errs, fmt.Sprintf("presentation.xml master r:id=%s points to non-slideMaster rel type %s", rid, rel.typeURI))
		}
		part := "ppt/" + strings.TrimPrefix(rel.target, "../")
		if _, exists := cand.entrySet[part]; !exists {
			errs = append(errs, fmt.Sprintf("presentation master rel target missing part: %s", part))
		}
	}

	masterUsage := map[string]int{}
	for slideRel, target := range cand.slideLayoutTargets {
		layoutPart := cleanPartTarget("ppt/slides", target)
		if _, exists := cand.entrySet[layoutPart]; !exists {
			errs = append(errs, fmt.Sprintf("%s points to missing layout part %s", slideRel, layoutPart))
			continue
		}
		layoutName := filepath.Base(layoutPart)
		masterTarget, ok := cand.layoutToMaster[layoutName]
		if !ok {
			errs = append(errs, fmt.Sprintf("layout rel missing for %s", layoutName))
			continue
		}
		masterPart := cleanPartTarget("ppt/slideLayouts", masterTarget)
		if _, exists := cand.entrySet[masterPart]; !exists {
			errs = append(errs, fmt.Sprintf("layout %s points to missing master part %s", layoutName, masterPart))
			continue
		}
		masterUsage[filepath.Base(masterPart)]++
	}

	if len(cand.masters) > 1 && len(masterUsage) < 2 {
		errs = append(errs, "slides are not distributed across multiple master families")
	}

	return dedupe(errs)
}

func cleanPartTarget(baseDir, target string) string {
	joined := filepath.ToSlash(filepath.Clean(filepath.Join(baseDir, target)))
	return strings.TrimPrefix(joined, "./")
}

func readZipEntry(files []*zip.File, name string) (string, error) {
	for _, f := range files {
		if f.Name != name {
			continue
		}
		rc, err := f.Open()
		if err != nil {
			return "", fmt.Errorf("open zip entry %s: %w", name, err)
		}
		defer func() {
			_ = rc.Close()
		}()
		buf, err := io.ReadAll(rc)
		if err != nil {
			return "", fmt.Errorf("read zip entry %s: %w", name, err)
		}
		return string(buf), nil
	}
	return "", fmt.Errorf("missing zip entry: %s", name)
}

func printSummary(base, cand *pkgInfo) {
	fmt.Printf("BASELINE=%s masters=%d layouts=%d\n", base.path, len(base.masters), len(base.layouts))
	fmt.Printf("CANDIDATE=%s masters=%d layouts=%d\n", cand.path, len(cand.masters), len(cand.layouts))
}

func dedupe(in []string) []string {
	seen := map[string]struct{}{}
	out := make([]string, 0, len(in))
	for _, msg := range in {
		if _, ok := seen[msg]; ok {
			continue
		}
		seen[msg] = struct{}{}
		out = append(out, msg)
	}
	sort.Strings(out)
	return out
}
