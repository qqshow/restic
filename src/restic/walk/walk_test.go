package walk_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"restic"
	"restic/archiver"
	"restic/pipe"
	"restic/repository"
	. "restic/test"
	"restic/walk"
)

func TestWalkTree(t *testing.T) {
	repo, cleanup := repository.TestRepository(t)
	defer cleanup()

	dirs, err := filepath.Glob(TestWalkerPath)
	OK(t, err)

	// archive a few files
	arch := archiver.New(repo)
	sn, _, err := arch.Snapshot(nil, dirs, nil)
	OK(t, err)

	// flush repo, write all packs
	OK(t, repo.Flush())

	done := make(chan struct{})

	// start tree walker
	treeJobs := make(chan walk.TreeJob)
	go walk.Tree(repo, *sn.Tree, done, treeJobs)

	// start filesystem walker
	fsJobs := make(chan pipe.Job)
	resCh := make(chan pipe.Result, 1)

	f := func(string, os.FileInfo) bool {
		return true
	}
	go pipe.Walk(dirs, f, done, fsJobs, resCh)

	for {
		// receive fs job
		fsJob, fsChOpen := <-fsJobs
		Assert(t, !fsChOpen || fsJob != nil,
			"received nil job from filesystem: %v %v", fsJob, fsChOpen)
		if fsJob != nil {
			OK(t, fsJob.Error())
		}

		var path string
		fsEntries := 1
		switch j := fsJob.(type) {
		case pipe.Dir:
			path = j.Path()
			fsEntries = len(j.Entries)
		case pipe.Entry:
			path = j.Path()
		}

		// receive tree job
		treeJob, treeChOpen := <-treeJobs
		treeEntries := 1

		OK(t, treeJob.Error)

		if treeJob.Tree != nil {
			treeEntries = len(treeJob.Tree.Nodes)
		}

		Assert(t, fsChOpen == treeChOpen,
			"one channel closed too early: fsChOpen %v, treeChOpen %v",
			fsChOpen, treeChOpen)

		if !fsChOpen || !treeChOpen {
			break
		}

		Assert(t, filepath.Base(path) == filepath.Base(treeJob.Path),
			"paths do not match: %q != %q", filepath.Base(path), filepath.Base(treeJob.Path))

		Assert(t, fsEntries == treeEntries,
			"wrong number of entries: %v != %v", fsEntries, treeEntries)
	}
}

type delayRepo struct {
	repo  restic.Repository
	delay time.Duration
}

func (d delayRepo) LoadTree(id restic.ID) (*restic.Tree, error) {
	time.Sleep(d.delay)
	return d.repo.LoadTree(id)
}

var repoFixture = filepath.Join("testdata", "walktree-test-repo.tar.gz")

var walktreeTestItems = []string{
	"testdata/0/0/0/0",
	"testdata/0/0/0/1",
	"testdata/0/0/0/10",
	"testdata/0/0/0/100",
	"testdata/0/0/0/101",
	"testdata/0/0/0/102",
	"testdata/0/0/0/103",
	"testdata/0/0/0/104",
	"testdata/0/0/0/105",
	"testdata/0/0/0/106",
	"testdata/0/0/0/107",
	"testdata/0/0/0/108",
	"testdata/0/0/0/109",
	"testdata/0/0/0/11",
	"testdata/0/0/0/110",
	"testdata/0/0/0/111",
	"testdata/0/0/0/112",
	"testdata/0/0/0/113",
	"testdata/0/0/0/114",
	"testdata/0/0/0/115",
	"testdata/0/0/0/116",
	"testdata/0/0/0/117",
	"testdata/0/0/0/118",
	"testdata/0/0/0/119",
	"testdata/0/0/0/12",
	"testdata/0/0/0/120",
	"testdata/0/0/0/121",
	"testdata/0/0/0/122",
	"testdata/0/0/0/123",
	"testdata/0/0/0/124",
	"testdata/0/0/0/125",
	"testdata/0/0/0/126",
	"testdata/0/0/0/127",
	"testdata/0/0/0/13",
	"testdata/0/0/0/14",
	"testdata/0/0/0/15",
	"testdata/0/0/0/16",
	"testdata/0/0/0/17",
	"testdata/0/0/0/18",
	"testdata/0/0/0/19",
	"testdata/0/0/0/2",
	"testdata/0/0/0/20",
	"testdata/0/0/0/21",
	"testdata/0/0/0/22",
	"testdata/0/0/0/23",
	"testdata/0/0/0/24",
	"testdata/0/0/0/25",
	"testdata/0/0/0/26",
	"testdata/0/0/0/27",
	"testdata/0/0/0/28",
	"testdata/0/0/0/29",
	"testdata/0/0/0/3",
	"testdata/0/0/0/30",
	"testdata/0/0/0/31",
	"testdata/0/0/0/32",
	"testdata/0/0/0/33",
	"testdata/0/0/0/34",
	"testdata/0/0/0/35",
	"testdata/0/0/0/36",
	"testdata/0/0/0/37",
	"testdata/0/0/0/38",
	"testdata/0/0/0/39",
	"testdata/0/0/0/4",
	"testdata/0/0/0/40",
	"testdata/0/0/0/41",
	"testdata/0/0/0/42",
	"testdata/0/0/0/43",
	"testdata/0/0/0/44",
	"testdata/0/0/0/45",
	"testdata/0/0/0/46",
	"testdata/0/0/0/47",
	"testdata/0/0/0/48",
	"testdata/0/0/0/49",
	"testdata/0/0/0/5",
	"testdata/0/0/0/50",
	"testdata/0/0/0/51",
	"testdata/0/0/0/52",
	"testdata/0/0/0/53",
	"testdata/0/0/0/54",
	"testdata/0/0/0/55",
	"testdata/0/0/0/56",
	"testdata/0/0/0/57",
	"testdata/0/0/0/58",
	"testdata/0/0/0/59",
	"testdata/0/0/0/6",
	"testdata/0/0/0/60",
	"testdata/0/0/0/61",
	"testdata/0/0/0/62",
	"testdata/0/0/0/63",
	"testdata/0/0/0/64",
	"testdata/0/0/0/65",
	"testdata/0/0/0/66",
	"testdata/0/0/0/67",
	"testdata/0/0/0/68",
	"testdata/0/0/0/69",
	"testdata/0/0/0/7",
	"testdata/0/0/0/70",
	"testdata/0/0/0/71",
	"testdata/0/0/0/72",
	"testdata/0/0/0/73",
	"testdata/0/0/0/74",
	"testdata/0/0/0/75",
	"testdata/0/0/0/76",
	"testdata/0/0/0/77",
	"testdata/0/0/0/78",
	"testdata/0/0/0/79",
	"testdata/0/0/0/8",
	"testdata/0/0/0/80",
	"testdata/0/0/0/81",
	"testdata/0/0/0/82",
	"testdata/0/0/0/83",
	"testdata/0/0/0/84",
	"testdata/0/0/0/85",
	"testdata/0/0/0/86",
	"testdata/0/0/0/87",
	"testdata/0/0/0/88",
	"testdata/0/0/0/89",
	"testdata/0/0/0/9",
	"testdata/0/0/0/90",
	"testdata/0/0/0/91",
	"testdata/0/0/0/92",
	"testdata/0/0/0/93",
	"testdata/0/0/0/94",
	"testdata/0/0/0/95",
	"testdata/0/0/0/96",
	"testdata/0/0/0/97",
	"testdata/0/0/0/98",
	"testdata/0/0/0/99",
	"testdata/0/0/0",
	"testdata/0/0/1/0",
	"testdata/0/0/1/1",
	"testdata/0/0/1/10",
	"testdata/0/0/1/100",
	"testdata/0/0/1/101",
	"testdata/0/0/1/102",
	"testdata/0/0/1/103",
	"testdata/0/0/1/104",
	"testdata/0/0/1/105",
	"testdata/0/0/1/106",
	"testdata/0/0/1/107",
	"testdata/0/0/1/108",
	"testdata/0/0/1/109",
	"testdata/0/0/1/11",
	"testdata/0/0/1/110",
	"testdata/0/0/1/111",
	"testdata/0/0/1/112",
	"testdata/0/0/1/113",
	"testdata/0/0/1/114",
	"testdata/0/0/1/115",
	"testdata/0/0/1/116",
	"testdata/0/0/1/117",
	"testdata/0/0/1/118",
	"testdata/0/0/1/119",
	"testdata/0/0/1/12",
	"testdata/0/0/1/120",
	"testdata/0/0/1/121",
	"testdata/0/0/1/122",
	"testdata/0/0/1/123",
	"testdata/0/0/1/124",
	"testdata/0/0/1/125",
	"testdata/0/0/1/126",
	"testdata/0/0/1/127",
	"testdata/0/0/1/13",
	"testdata/0/0/1/14",
	"testdata/0/0/1/15",
	"testdata/0/0/1/16",
	"testdata/0/0/1/17",
	"testdata/0/0/1/18",
	"testdata/0/0/1/19",
	"testdata/0/0/1/2",
	"testdata/0/0/1/20",
	"testdata/0/0/1/21",
	"testdata/0/0/1/22",
	"testdata/0/0/1/23",
	"testdata/0/0/1/24",
	"testdata/0/0/1/25",
	"testdata/0/0/1/26",
	"testdata/0/0/1/27",
	"testdata/0/0/1/28",
	"testdata/0/0/1/29",
	"testdata/0/0/1/3",
	"testdata/0/0/1/30",
	"testdata/0/0/1/31",
	"testdata/0/0/1/32",
	"testdata/0/0/1/33",
	"testdata/0/0/1/34",
	"testdata/0/0/1/35",
	"testdata/0/0/1/36",
	"testdata/0/0/1/37",
	"testdata/0/0/1/38",
	"testdata/0/0/1/39",
	"testdata/0/0/1/4",
	"testdata/0/0/1/40",
	"testdata/0/0/1/41",
	"testdata/0/0/1/42",
	"testdata/0/0/1/43",
	"testdata/0/0/1/44",
	"testdata/0/0/1/45",
	"testdata/0/0/1/46",
	"testdata/0/0/1/47",
	"testdata/0/0/1/48",
	"testdata/0/0/1/49",
	"testdata/0/0/1/5",
	"testdata/0/0/1/50",
	"testdata/0/0/1/51",
	"testdata/0/0/1/52",
	"testdata/0/0/1/53",
	"testdata/0/0/1/54",
	"testdata/0/0/1/55",
	"testdata/0/0/1/56",
	"testdata/0/0/1/57",
	"testdata/0/0/1/58",
	"testdata/0/0/1/59",
	"testdata/0/0/1/6",
	"testdata/0/0/1/60",
	"testdata/0/0/1/61",
	"testdata/0/0/1/62",
	"testdata/0/0/1/63",
	"testdata/0/0/1/64",
	"testdata/0/0/1/65",
	"testdata/0/0/1/66",
	"testdata/0/0/1/67",
	"testdata/0/0/1/68",
	"testdata/0/0/1/69",
	"testdata/0/0/1/7",
	"testdata/0/0/1/70",
	"testdata/0/0/1/71",
	"testdata/0/0/1/72",
	"testdata/0/0/1/73",
	"testdata/0/0/1/74",
	"testdata/0/0/1/75",
	"testdata/0/0/1/76",
	"testdata/0/0/1/77",
	"testdata/0/0/1/78",
	"testdata/0/0/1/79",
	"testdata/0/0/1/8",
	"testdata/0/0/1/80",
	"testdata/0/0/1/81",
	"testdata/0/0/1/82",
	"testdata/0/0/1/83",
	"testdata/0/0/1/84",
	"testdata/0/0/1/85",
	"testdata/0/0/1/86",
	"testdata/0/0/1/87",
	"testdata/0/0/1/88",
	"testdata/0/0/1/89",
	"testdata/0/0/1/9",
	"testdata/0/0/1/90",
	"testdata/0/0/1/91",
	"testdata/0/0/1/92",
	"testdata/0/0/1/93",
	"testdata/0/0/1/94",
	"testdata/0/0/1/95",
	"testdata/0/0/1/96",
	"testdata/0/0/1/97",
	"testdata/0/0/1/98",
	"testdata/0/0/1/99",
	"testdata/0/0/1",
	"testdata/0/0/2/0",
	"testdata/0/0/2/1",
	"testdata/0/0/2/10",
	"testdata/0/0/2/100",
	"testdata/0/0/2/101",
	"testdata/0/0/2/102",
	"testdata/0/0/2/103",
	"testdata/0/0/2/104",
	"testdata/0/0/2/105",
	"testdata/0/0/2/106",
	"testdata/0/0/2/107",
	"testdata/0/0/2/108",
	"testdata/0/0/2/109",
	"testdata/0/0/2/11",
	"testdata/0/0/2/110",
	"testdata/0/0/2/111",
	"testdata/0/0/2/112",
	"testdata/0/0/2/113",
	"testdata/0/0/2/114",
	"testdata/0/0/2/115",
	"testdata/0/0/2/116",
	"testdata/0/0/2/117",
	"testdata/0/0/2/118",
	"testdata/0/0/2/119",
	"testdata/0/0/2/12",
	"testdata/0/0/2/120",
	"testdata/0/0/2/121",
	"testdata/0/0/2/122",
	"testdata/0/0/2/123",
	"testdata/0/0/2/124",
	"testdata/0/0/2/125",
	"testdata/0/0/2/126",
	"testdata/0/0/2/127",
	"testdata/0/0/2/13",
	"testdata/0/0/2/14",
	"testdata/0/0/2/15",
	"testdata/0/0/2/16",
	"testdata/0/0/2/17",
	"testdata/0/0/2/18",
	"testdata/0/0/2/19",
	"testdata/0/0/2/2",
	"testdata/0/0/2/20",
	"testdata/0/0/2/21",
	"testdata/0/0/2/22",
	"testdata/0/0/2/23",
	"testdata/0/0/2/24",
	"testdata/0/0/2/25",
	"testdata/0/0/2/26",
	"testdata/0/0/2/27",
	"testdata/0/0/2/28",
	"testdata/0/0/2/29",
	"testdata/0/0/2/3",
	"testdata/0/0/2/30",
	"testdata/0/0/2/31",
	"testdata/0/0/2/32",
	"testdata/0/0/2/33",
	"testdata/0/0/2/34",
	"testdata/0/0/2/35",
	"testdata/0/0/2/36",
	"testdata/0/0/2/37",
	"testdata/0/0/2/38",
	"testdata/0/0/2/39",
	"testdata/0/0/2/4",
	"testdata/0/0/2/40",
	"testdata/0/0/2/41",
	"testdata/0/0/2/42",
	"testdata/0/0/2/43",
	"testdata/0/0/2/44",
	"testdata/0/0/2/45",
	"testdata/0/0/2/46",
	"testdata/0/0/2/47",
	"testdata/0/0/2/48",
	"testdata/0/0/2/49",
	"testdata/0/0/2/5",
	"testdata/0/0/2/50",
	"testdata/0/0/2/51",
	"testdata/0/0/2/52",
	"testdata/0/0/2/53",
	"testdata/0/0/2/54",
	"testdata/0/0/2/55",
	"testdata/0/0/2/56",
	"testdata/0/0/2/57",
	"testdata/0/0/2/58",
	"testdata/0/0/2/59",
	"testdata/0/0/2/6",
	"testdata/0/0/2/60",
	"testdata/0/0/2/61",
	"testdata/0/0/2/62",
	"testdata/0/0/2/63",
	"testdata/0/0/2/64",
	"testdata/0/0/2/65",
	"testdata/0/0/2/66",
	"testdata/0/0/2/67",
	"testdata/0/0/2/68",
	"testdata/0/0/2/69",
	"testdata/0/0/2/7",
	"testdata/0/0/2/70",
	"testdata/0/0/2/71",
	"testdata/0/0/2/72",
	"testdata/0/0/2/73",
	"testdata/0/0/2/74",
	"testdata/0/0/2/75",
	"testdata/0/0/2/76",
	"testdata/0/0/2/77",
	"testdata/0/0/2/78",
	"testdata/0/0/2/79",
	"testdata/0/0/2/8",
	"testdata/0/0/2/80",
	"testdata/0/0/2/81",
	"testdata/0/0/2/82",
	"testdata/0/0/2/83",
	"testdata/0/0/2/84",
	"testdata/0/0/2/85",
	"testdata/0/0/2/86",
	"testdata/0/0/2/87",
	"testdata/0/0/2/88",
	"testdata/0/0/2/89",
	"testdata/0/0/2/9",
	"testdata/0/0/2/90",
	"testdata/0/0/2/91",
	"testdata/0/0/2/92",
	"testdata/0/0/2/93",
	"testdata/0/0/2/94",
	"testdata/0/0/2/95",
	"testdata/0/0/2/96",
	"testdata/0/0/2/97",
	"testdata/0/0/2/98",
	"testdata/0/0/2/99",
	"testdata/0/0/2",
	"testdata/0/0/3/0",
	"testdata/0/0/3/1",
	"testdata/0/0/3/10",
	"testdata/0/0/3/100",
	"testdata/0/0/3/101",
	"testdata/0/0/3/102",
	"testdata/0/0/3/103",
	"testdata/0/0/3/104",
	"testdata/0/0/3/105",
	"testdata/0/0/3/106",
	"testdata/0/0/3/107",
	"testdata/0/0/3/108",
	"testdata/0/0/3/109",
	"testdata/0/0/3/11",
	"testdata/0/0/3/110",
	"testdata/0/0/3/111",
	"testdata/0/0/3/112",
	"testdata/0/0/3/113",
	"testdata/0/0/3/114",
	"testdata/0/0/3/115",
	"testdata/0/0/3/116",
	"testdata/0/0/3/117",
	"testdata/0/0/3/118",
	"testdata/0/0/3/119",
	"testdata/0/0/3/12",
	"testdata/0/0/3/120",
	"testdata/0/0/3/121",
	"testdata/0/0/3/122",
	"testdata/0/0/3/123",
	"testdata/0/0/3/124",
	"testdata/0/0/3/125",
	"testdata/0/0/3/126",
	"testdata/0/0/3/127",
	"testdata/0/0/3/13",
	"testdata/0/0/3/14",
	"testdata/0/0/3/15",
	"testdata/0/0/3/16",
	"testdata/0/0/3/17",
	"testdata/0/0/3/18",
	"testdata/0/0/3/19",
	"testdata/0/0/3/2",
	"testdata/0/0/3/20",
	"testdata/0/0/3/21",
	"testdata/0/0/3/22",
	"testdata/0/0/3/23",
	"testdata/0/0/3/24",
	"testdata/0/0/3/25",
	"testdata/0/0/3/26",
	"testdata/0/0/3/27",
	"testdata/0/0/3/28",
	"testdata/0/0/3/29",
	"testdata/0/0/3/3",
	"testdata/0/0/3/30",
	"testdata/0/0/3/31",
	"testdata/0/0/3/32",
	"testdata/0/0/3/33",
	"testdata/0/0/3/34",
	"testdata/0/0/3/35",
	"testdata/0/0/3/36",
	"testdata/0/0/3/37",
	"testdata/0/0/3/38",
	"testdata/0/0/3/39",
	"testdata/0/0/3/4",
	"testdata/0/0/3/40",
	"testdata/0/0/3/41",
	"testdata/0/0/3/42",
	"testdata/0/0/3/43",
	"testdata/0/0/3/44",
	"testdata/0/0/3/45",
	"testdata/0/0/3/46",
	"testdata/0/0/3/47",
	"testdata/0/0/3/48",
	"testdata/0/0/3/49",
	"testdata/0/0/3/5",
	"testdata/0/0/3/50",
	"testdata/0/0/3/51",
	"testdata/0/0/3/52",
	"testdata/0/0/3/53",
	"testdata/0/0/3/54",
	"testdata/0/0/3/55",
	"testdata/0/0/3/56",
	"testdata/0/0/3/57",
	"testdata/0/0/3/58",
	"testdata/0/0/3/59",
	"testdata/0/0/3/6",
	"testdata/0/0/3/60",
	"testdata/0/0/3/61",
	"testdata/0/0/3/62",
	"testdata/0/0/3/63",
	"testdata/0/0/3/64",
	"testdata/0/0/3/65",
	"testdata/0/0/3/66",
	"testdata/0/0/3/67",
	"testdata/0/0/3/68",
	"testdata/0/0/3/69",
	"testdata/0/0/3/7",
	"testdata/0/0/3/70",
	"testdata/0/0/3/71",
	"testdata/0/0/3/72",
	"testdata/0/0/3/73",
	"testdata/0/0/3/74",
	"testdata/0/0/3/75",
	"testdata/0/0/3/76",
	"testdata/0/0/3/77",
	"testdata/0/0/3/78",
	"testdata/0/0/3/79",
	"testdata/0/0/3/8",
	"testdata/0/0/3/80",
	"testdata/0/0/3/81",
	"testdata/0/0/3/82",
	"testdata/0/0/3/83",
	"testdata/0/0/3/84",
	"testdata/0/0/3/85",
	"testdata/0/0/3/86",
	"testdata/0/0/3/87",
	"testdata/0/0/3/88",
	"testdata/0/0/3/89",
	"testdata/0/0/3/9",
	"testdata/0/0/3/90",
	"testdata/0/0/3/91",
	"testdata/0/0/3/92",
	"testdata/0/0/3/93",
	"testdata/0/0/3/94",
	"testdata/0/0/3/95",
	"testdata/0/0/3/96",
	"testdata/0/0/3/97",
	"testdata/0/0/3/98",
	"testdata/0/0/3/99",
	"testdata/0/0/3",
	"testdata/0/0/4/0",
	"testdata/0/0/4/1",
	"testdata/0/0/4/10",
	"testdata/0/0/4/100",
	"testdata/0/0/4/101",
	"testdata/0/0/4/102",
	"testdata/0/0/4/103",
	"testdata/0/0/4/104",
	"testdata/0/0/4/105",
	"testdata/0/0/4/106",
	"testdata/0/0/4/107",
	"testdata/0/0/4/108",
	"testdata/0/0/4/109",
	"testdata/0/0/4/11",
	"testdata/0/0/4/110",
	"testdata/0/0/4/111",
	"testdata/0/0/4/112",
	"testdata/0/0/4/113",
	"testdata/0/0/4/114",
	"testdata/0/0/4/115",
	"testdata/0/0/4/116",
	"testdata/0/0/4/117",
	"testdata/0/0/4/118",
	"testdata/0/0/4/119",
	"testdata/0/0/4/12",
	"testdata/0/0/4/120",
	"testdata/0/0/4/121",
	"testdata/0/0/4/122",
	"testdata/0/0/4/123",
	"testdata/0/0/4/124",
	"testdata/0/0/4/125",
	"testdata/0/0/4/126",
	"testdata/0/0/4/127",
	"testdata/0/0/4/13",
	"testdata/0/0/4/14",
	"testdata/0/0/4/15",
	"testdata/0/0/4/16",
	"testdata/0/0/4/17",
	"testdata/0/0/4/18",
	"testdata/0/0/4/19",
	"testdata/0/0/4/2",
	"testdata/0/0/4/20",
	"testdata/0/0/4/21",
	"testdata/0/0/4/22",
	"testdata/0/0/4/23",
	"testdata/0/0/4/24",
	"testdata/0/0/4/25",
	"testdata/0/0/4/26",
	"testdata/0/0/4/27",
	"testdata/0/0/4/28",
	"testdata/0/0/4/29",
	"testdata/0/0/4/3",
	"testdata/0/0/4/30",
	"testdata/0/0/4/31",
	"testdata/0/0/4/32",
	"testdata/0/0/4/33",
	"testdata/0/0/4/34",
	"testdata/0/0/4/35",
	"testdata/0/0/4/36",
	"testdata/0/0/4/37",
	"testdata/0/0/4/38",
	"testdata/0/0/4/39",
	"testdata/0/0/4/4",
	"testdata/0/0/4/40",
	"testdata/0/0/4/41",
	"testdata/0/0/4/42",
	"testdata/0/0/4/43",
	"testdata/0/0/4/44",
	"testdata/0/0/4/45",
	"testdata/0/0/4/46",
	"testdata/0/0/4/47",
	"testdata/0/0/4/48",
	"testdata/0/0/4/49",
	"testdata/0/0/4/5",
	"testdata/0/0/4/50",
	"testdata/0/0/4/51",
	"testdata/0/0/4/52",
	"testdata/0/0/4/53",
	"testdata/0/0/4/54",
	"testdata/0/0/4/55",
	"testdata/0/0/4/56",
	"testdata/0/0/4/57",
	"testdata/0/0/4/58",
	"testdata/0/0/4/59",
	"testdata/0/0/4/6",
	"testdata/0/0/4/60",
	"testdata/0/0/4/61",
	"testdata/0/0/4/62",
	"testdata/0/0/4/63",
	"testdata/0/0/4/64",
	"testdata/0/0/4/65",
	"testdata/0/0/4/66",
	"testdata/0/0/4/67",
	"testdata/0/0/4/68",
	"testdata/0/0/4/69",
	"testdata/0/0/4/7",
	"testdata/0/0/4/70",
	"testdata/0/0/4/71",
	"testdata/0/0/4/72",
	"testdata/0/0/4/73",
	"testdata/0/0/4/74",
	"testdata/0/0/4/75",
	"testdata/0/0/4/76",
	"testdata/0/0/4/77",
	"testdata/0/0/4/78",
	"testdata/0/0/4/79",
	"testdata/0/0/4/8",
	"testdata/0/0/4/80",
	"testdata/0/0/4/81",
	"testdata/0/0/4/82",
	"testdata/0/0/4/83",
	"testdata/0/0/4/84",
	"testdata/0/0/4/85",
	"testdata/0/0/4/86",
	"testdata/0/0/4/87",
	"testdata/0/0/4/88",
	"testdata/0/0/4/89",
	"testdata/0/0/4/9",
	"testdata/0/0/4/90",
	"testdata/0/0/4/91",
	"testdata/0/0/4/92",
	"testdata/0/0/4/93",
	"testdata/0/0/4/94",
	"testdata/0/0/4/95",
	"testdata/0/0/4/96",
	"testdata/0/0/4/97",
	"testdata/0/0/4/98",
	"testdata/0/0/4/99",
	"testdata/0/0/4",
	"testdata/0/0/5/0",
	"testdata/0/0/5/1",
	"testdata/0/0/5/10",
	"testdata/0/0/5/100",
	"testdata/0/0/5/101",
	"testdata/0/0/5/102",
	"testdata/0/0/5/103",
	"testdata/0/0/5/104",
	"testdata/0/0/5/105",
	"testdata/0/0/5/106",
	"testdata/0/0/5/107",
	"testdata/0/0/5/108",
	"testdata/0/0/5/109",
	"testdata/0/0/5/11",
	"testdata/0/0/5/110",
	"testdata/0/0/5/111",
	"testdata/0/0/5/112",
	"testdata/0/0/5/113",
	"testdata/0/0/5/114",
	"testdata/0/0/5/115",
	"testdata/0/0/5/116",
	"testdata/0/0/5/117",
	"testdata/0/0/5/118",
	"testdata/0/0/5/119",
	"testdata/0/0/5/12",
	"testdata/0/0/5/120",
	"testdata/0/0/5/121",
	"testdata/0/0/5/122",
	"testdata/0/0/5/123",
	"testdata/0/0/5/124",
	"testdata/0/0/5/125",
	"testdata/0/0/5/126",
	"testdata/0/0/5/127",
	"testdata/0/0/5/13",
	"testdata/0/0/5/14",
	"testdata/0/0/5/15",
	"testdata/0/0/5/16",
	"testdata/0/0/5/17",
	"testdata/0/0/5/18",
	"testdata/0/0/5/19",
	"testdata/0/0/5/2",
	"testdata/0/0/5/20",
	"testdata/0/0/5/21",
	"testdata/0/0/5/22",
	"testdata/0/0/5/23",
	"testdata/0/0/5/24",
	"testdata/0/0/5/25",
	"testdata/0/0/5/26",
	"testdata/0/0/5/27",
	"testdata/0/0/5/28",
	"testdata/0/0/5/29",
	"testdata/0/0/5/3",
	"testdata/0/0/5/30",
	"testdata/0/0/5/31",
	"testdata/0/0/5/32",
	"testdata/0/0/5/33",
	"testdata/0/0/5/34",
	"testdata/0/0/5/35",
	"testdata/0/0/5/36",
	"testdata/0/0/5/37",
	"testdata/0/0/5/38",
	"testdata/0/0/5/39",
	"testdata/0/0/5/4",
	"testdata/0/0/5/40",
	"testdata/0/0/5/41",
	"testdata/0/0/5/42",
	"testdata/0/0/5/43",
	"testdata/0/0/5/44",
	"testdata/0/0/5/45",
	"testdata/0/0/5/46",
	"testdata/0/0/5/47",
	"testdata/0/0/5/48",
	"testdata/0/0/5/49",
	"testdata/0/0/5/5",
	"testdata/0/0/5/50",
	"testdata/0/0/5/51",
	"testdata/0/0/5/52",
	"testdata/0/0/5/53",
	"testdata/0/0/5/54",
	"testdata/0/0/5/55",
	"testdata/0/0/5/56",
	"testdata/0/0/5/57",
	"testdata/0/0/5/58",
	"testdata/0/0/5/59",
	"testdata/0/0/5/6",
	"testdata/0/0/5/60",
	"testdata/0/0/5/61",
	"testdata/0/0/5/62",
	"testdata/0/0/5/63",
	"testdata/0/0/5/64",
	"testdata/0/0/5/65",
	"testdata/0/0/5/66",
	"testdata/0/0/5/67",
	"testdata/0/0/5/68",
	"testdata/0/0/5/69",
	"testdata/0/0/5/7",
	"testdata/0/0/5/70",
	"testdata/0/0/5/71",
	"testdata/0/0/5/72",
	"testdata/0/0/5/73",
	"testdata/0/0/5/74",
	"testdata/0/0/5/75",
	"testdata/0/0/5/76",
	"testdata/0/0/5/77",
	"testdata/0/0/5/78",
	"testdata/0/0/5/79",
	"testdata/0/0/5/8",
	"testdata/0/0/5/80",
	"testdata/0/0/5/81",
	"testdata/0/0/5/82",
	"testdata/0/0/5/83",
	"testdata/0/0/5/84",
	"testdata/0/0/5/85",
	"testdata/0/0/5/86",
	"testdata/0/0/5/87",
	"testdata/0/0/5/88",
	"testdata/0/0/5/89",
	"testdata/0/0/5/9",
	"testdata/0/0/5/90",
	"testdata/0/0/5/91",
	"testdata/0/0/5/92",
	"testdata/0/0/5/93",
	"testdata/0/0/5/94",
	"testdata/0/0/5/95",
	"testdata/0/0/5/96",
	"testdata/0/0/5/97",
	"testdata/0/0/5/98",
	"testdata/0/0/5/99",
	"testdata/0/0/5",
	"testdata/0/0/6/0",
	"testdata/0/0/6/1",
	"testdata/0/0/6/10",
	"testdata/0/0/6/100",
	"testdata/0/0/6/101",
	"testdata/0/0/6/102",
	"testdata/0/0/6/103",
	"testdata/0/0/6/104",
	"testdata/0/0/6/105",
	"testdata/0/0/6/106",
	"testdata/0/0/6/107",
	"testdata/0/0/6/108",
	"testdata/0/0/6/109",
	"testdata/0/0/6/11",
	"testdata/0/0/6/110",
	"testdata/0/0/6/111",
	"testdata/0/0/6/112",
	"testdata/0/0/6/113",
	"testdata/0/0/6/114",
	"testdata/0/0/6/115",
	"testdata/0/0/6/116",
	"testdata/0/0/6/117",
	"testdata/0/0/6/118",
	"testdata/0/0/6/119",
	"testdata/0/0/6/12",
	"testdata/0/0/6/120",
	"testdata/0/0/6/121",
	"testdata/0/0/6/122",
	"testdata/0/0/6/123",
	"testdata/0/0/6/124",
	"testdata/0/0/6/125",
	"testdata/0/0/6/126",
	"testdata/0/0/6/127",
	"testdata/0/0/6/13",
	"testdata/0/0/6/14",
	"testdata/0/0/6/15",
	"testdata/0/0/6/16",
	"testdata/0/0/6/17",
	"testdata/0/0/6/18",
	"testdata/0/0/6/19",
	"testdata/0/0/6/2",
	"testdata/0/0/6/20",
	"testdata/0/0/6/21",
	"testdata/0/0/6/22",
	"testdata/0/0/6/23",
	"testdata/0/0/6/24",
	"testdata/0/0/6/25",
	"testdata/0/0/6/26",
	"testdata/0/0/6/27",
	"testdata/0/0/6/28",
	"testdata/0/0/6/29",
	"testdata/0/0/6/3",
	"testdata/0/0/6/30",
	"testdata/0/0/6/31",
	"testdata/0/0/6/32",
	"testdata/0/0/6/33",
	"testdata/0/0/6/34",
	"testdata/0/0/6/35",
	"testdata/0/0/6/36",
	"testdata/0/0/6/37",
	"testdata/0/0/6/38",
	"testdata/0/0/6/39",
	"testdata/0/0/6/4",
	"testdata/0/0/6/40",
	"testdata/0/0/6/41",
	"testdata/0/0/6/42",
	"testdata/0/0/6/43",
	"testdata/0/0/6/44",
	"testdata/0/0/6/45",
	"testdata/0/0/6/46",
	"testdata/0/0/6/47",
	"testdata/0/0/6/48",
	"testdata/0/0/6/49",
	"testdata/0/0/6/5",
	"testdata/0/0/6/50",
	"testdata/0/0/6/51",
	"testdata/0/0/6/52",
	"testdata/0/0/6/53",
	"testdata/0/0/6/54",
	"testdata/0/0/6/55",
	"testdata/0/0/6/56",
	"testdata/0/0/6/57",
	"testdata/0/0/6/58",
	"testdata/0/0/6/59",
	"testdata/0/0/6/6",
	"testdata/0/0/6/60",
	"testdata/0/0/6/61",
	"testdata/0/0/6/62",
	"testdata/0/0/6/63",
	"testdata/0/0/6/64",
	"testdata/0/0/6/65",
	"testdata/0/0/6/66",
	"testdata/0/0/6/67",
	"testdata/0/0/6/68",
	"testdata/0/0/6/69",
	"testdata/0/0/6/7",
	"testdata/0/0/6/70",
	"testdata/0/0/6/71",
	"testdata/0/0/6/72",
	"testdata/0/0/6/73",
	"testdata/0/0/6/74",
	"testdata/0/0/6/75",
	"testdata/0/0/6/76",
	"testdata/0/0/6/77",
	"testdata/0/0/6/78",
	"testdata/0/0/6/79",
	"testdata/0/0/6/8",
	"testdata/0/0/6/80",
	"testdata/0/0/6/81",
	"testdata/0/0/6/82",
	"testdata/0/0/6/83",
	"testdata/0/0/6/84",
	"testdata/0/0/6/85",
	"testdata/0/0/6/86",
	"testdata/0/0/6/87",
	"testdata/0/0/6/88",
	"testdata/0/0/6/89",
	"testdata/0/0/6/9",
	"testdata/0/0/6/90",
	"testdata/0/0/6/91",
	"testdata/0/0/6/92",
	"testdata/0/0/6/93",
	"testdata/0/0/6/94",
	"testdata/0/0/6/95",
	"testdata/0/0/6/96",
	"testdata/0/0/6/97",
	"testdata/0/0/6/98",
	"testdata/0/0/6/99",
	"testdata/0/0/6",
	"testdata/0/0/7/0",
	"testdata/0/0/7/1",
	"testdata/0/0/7/10",
	"testdata/0/0/7/100",
	"testdata/0/0/7/101",
	"testdata/0/0/7/102",
	"testdata/0/0/7/103",
	"testdata/0/0/7/104",
	"testdata/0/0/7/105",
	"testdata/0/0/7/106",
	"testdata/0/0/7/107",
	"testdata/0/0/7/108",
	"testdata/0/0/7/109",
	"testdata/0/0/7/11",
	"testdata/0/0/7/110",
	"testdata/0/0/7/111",
	"testdata/0/0/7/112",
	"testdata/0/0/7/113",
	"testdata/0/0/7/114",
	"testdata/0/0/7/115",
	"testdata/0/0/7/116",
	"testdata/0/0/7/117",
	"testdata/0/0/7/118",
	"testdata/0/0/7/119",
	"testdata/0/0/7/12",
	"testdata/0/0/7/120",
	"testdata/0/0/7/121",
	"testdata/0/0/7/122",
	"testdata/0/0/7/123",
	"testdata/0/0/7/124",
	"testdata/0/0/7/125",
	"testdata/0/0/7/126",
	"testdata/0/0/7/127",
	"testdata/0/0/7/13",
	"testdata/0/0/7/14",
	"testdata/0/0/7/15",
	"testdata/0/0/7/16",
	"testdata/0/0/7/17",
	"testdata/0/0/7/18",
	"testdata/0/0/7/19",
	"testdata/0/0/7/2",
	"testdata/0/0/7/20",
	"testdata/0/0/7/21",
	"testdata/0/0/7/22",
	"testdata/0/0/7/23",
	"testdata/0/0/7/24",
	"testdata/0/0/7/25",
	"testdata/0/0/7/26",
	"testdata/0/0/7/27",
	"testdata/0/0/7/28",
	"testdata/0/0/7/29",
	"testdata/0/0/7/3",
	"testdata/0/0/7/30",
	"testdata/0/0/7/31",
	"testdata/0/0/7/32",
	"testdata/0/0/7/33",
	"testdata/0/0/7/34",
	"testdata/0/0/7/35",
	"testdata/0/0/7/36",
	"testdata/0/0/7/37",
	"testdata/0/0/7/38",
	"testdata/0/0/7/39",
	"testdata/0/0/7/4",
	"testdata/0/0/7/40",
	"testdata/0/0/7/41",
	"testdata/0/0/7/42",
	"testdata/0/0/7/43",
	"testdata/0/0/7/44",
	"testdata/0/0/7/45",
	"testdata/0/0/7/46",
	"testdata/0/0/7/47",
	"testdata/0/0/7/48",
	"testdata/0/0/7/49",
	"testdata/0/0/7/5",
	"testdata/0/0/7/50",
	"testdata/0/0/7/51",
	"testdata/0/0/7/52",
	"testdata/0/0/7/53",
	"testdata/0/0/7/54",
	"testdata/0/0/7/55",
	"testdata/0/0/7/56",
	"testdata/0/0/7/57",
	"testdata/0/0/7/58",
	"testdata/0/0/7/59",
	"testdata/0/0/7/6",
	"testdata/0/0/7/60",
	"testdata/0/0/7/61",
	"testdata/0/0/7/62",
	"testdata/0/0/7/63",
	"testdata/0/0/7/64",
	"testdata/0/0/7/65",
	"testdata/0/0/7/66",
	"testdata/0/0/7/67",
	"testdata/0/0/7/68",
	"testdata/0/0/7/69",
	"testdata/0/0/7/7",
	"testdata/0/0/7/70",
	"testdata/0/0/7/71",
	"testdata/0/0/7/72",
	"testdata/0/0/7/73",
	"testdata/0/0/7/74",
	"testdata/0/0/7/75",
	"testdata/0/0/7/76",
	"testdata/0/0/7/77",
	"testdata/0/0/7/78",
	"testdata/0/0/7/79",
	"testdata/0/0/7/8",
	"testdata/0/0/7/80",
	"testdata/0/0/7/81",
	"testdata/0/0/7/82",
	"testdata/0/0/7/83",
	"testdata/0/0/7/84",
	"testdata/0/0/7/85",
	"testdata/0/0/7/86",
	"testdata/0/0/7/87",
	"testdata/0/0/7/88",
	"testdata/0/0/7/89",
	"testdata/0/0/7/9",
	"testdata/0/0/7/90",
	"testdata/0/0/7/91",
	"testdata/0/0/7/92",
	"testdata/0/0/7/93",
	"testdata/0/0/7/94",
	"testdata/0/0/7/95",
	"testdata/0/0/7/96",
	"testdata/0/0/7/97",
	"testdata/0/0/7/98",
	"testdata/0/0/7/99",
	"testdata/0/0/7",
	"testdata/0/0/8/0",
	"testdata/0/0/8/1",
	"testdata/0/0/8/10",
	"testdata/0/0/8/100",
	"testdata/0/0/8/101",
	"testdata/0/0/8/102",
	"testdata/0/0/8/103",
	"testdata/0/0/8/104",
	"testdata/0/0/8/105",
	"testdata/0/0/8/106",
	"testdata/0/0/8/107",
	"testdata/0/0/8/108",
	"testdata/0/0/8/109",
	"testdata/0/0/8/11",
	"testdata/0/0/8/110",
	"testdata/0/0/8/111",
	"testdata/0/0/8/112",
	"testdata/0/0/8/113",
	"testdata/0/0/8/114",
	"testdata/0/0/8/115",
	"testdata/0/0/8/116",
	"testdata/0/0/8/117",
	"testdata/0/0/8/118",
	"testdata/0/0/8/119",
	"testdata/0/0/8/12",
	"testdata/0/0/8/120",
	"testdata/0/0/8/121",
	"testdata/0/0/8/122",
	"testdata/0/0/8/123",
	"testdata/0/0/8/124",
	"testdata/0/0/8/125",
	"testdata/0/0/8/126",
	"testdata/0/0/8/127",
	"testdata/0/0/8/13",
	"testdata/0/0/8/14",
	"testdata/0/0/8/15",
	"testdata/0/0/8/16",
	"testdata/0/0/8/17",
	"testdata/0/0/8/18",
	"testdata/0/0/8/19",
	"testdata/0/0/8/2",
	"testdata/0/0/8/20",
	"testdata/0/0/8/21",
	"testdata/0/0/8/22",
	"testdata/0/0/8/23",
	"testdata/0/0/8/24",
	"testdata/0/0/8/25",
	"testdata/0/0/8/26",
	"testdata/0/0/8/27",
	"testdata/0/0/8/28",
	"testdata/0/0/8/29",
	"testdata/0/0/8/3",
	"testdata/0/0/8/30",
	"testdata/0/0/8/31",
	"testdata/0/0/8/32",
	"testdata/0/0/8/33",
	"testdata/0/0/8/34",
	"testdata/0/0/8/35",
	"testdata/0/0/8/36",
	"testdata/0/0/8/37",
	"testdata/0/0/8/38",
	"testdata/0/0/8/39",
	"testdata/0/0/8/4",
	"testdata/0/0/8/40",
	"testdata/0/0/8/41",
	"testdata/0/0/8/42",
	"testdata/0/0/8/43",
	"testdata/0/0/8/44",
	"testdata/0/0/8/45",
	"testdata/0/0/8/46",
	"testdata/0/0/8/47",
	"testdata/0/0/8/48",
	"testdata/0/0/8/49",
	"testdata/0/0/8/5",
	"testdata/0/0/8/50",
	"testdata/0/0/8/51",
	"testdata/0/0/8/52",
	"testdata/0/0/8/53",
	"testdata/0/0/8/54",
	"testdata/0/0/8/55",
	"testdata/0/0/8/56",
	"testdata/0/0/8/57",
	"testdata/0/0/8/58",
	"testdata/0/0/8/59",
	"testdata/0/0/8/6",
	"testdata/0/0/8/60",
	"testdata/0/0/8/61",
	"testdata/0/0/8/62",
	"testdata/0/0/8/63",
	"testdata/0/0/8/64",
	"testdata/0/0/8/65",
	"testdata/0/0/8/66",
	"testdata/0/0/8/67",
	"testdata/0/0/8/68",
	"testdata/0/0/8/69",
	"testdata/0/0/8/7",
	"testdata/0/0/8/70",
	"testdata/0/0/8/71",
	"testdata/0/0/8/72",
	"testdata/0/0/8/73",
	"testdata/0/0/8/74",
	"testdata/0/0/8/75",
	"testdata/0/0/8/76",
	"testdata/0/0/8/77",
	"testdata/0/0/8/78",
	"testdata/0/0/8/79",
	"testdata/0/0/8/8",
	"testdata/0/0/8/80",
	"testdata/0/0/8/81",
	"testdata/0/0/8/82",
	"testdata/0/0/8/83",
	"testdata/0/0/8/84",
	"testdata/0/0/8/85",
	"testdata/0/0/8/86",
	"testdata/0/0/8/87",
	"testdata/0/0/8/88",
	"testdata/0/0/8/89",
	"testdata/0/0/8/9",
	"testdata/0/0/8/90",
	"testdata/0/0/8/91",
	"testdata/0/0/8/92",
	"testdata/0/0/8/93",
	"testdata/0/0/8/94",
	"testdata/0/0/8/95",
	"testdata/0/0/8/96",
	"testdata/0/0/8/97",
	"testdata/0/0/8/98",
	"testdata/0/0/8/99",
	"testdata/0/0/8",
	"testdata/0/0/9/0",
	"testdata/0/0/9/1",
	"testdata/0/0/9/10",
	"testdata/0/0/9/11",
	"testdata/0/0/9/12",
	"testdata/0/0/9/13",
	"testdata/0/0/9/14",
	"testdata/0/0/9/15",
	"testdata/0/0/9/16",
	"testdata/0/0/9/17",
	"testdata/0/0/9/18",
	"testdata/0/0/9/19",
	"testdata/0/0/9/2",
	"testdata/0/0/9/20",
	"testdata/0/0/9/21",
	"testdata/0/0/9/22",
	"testdata/0/0/9/23",
	"testdata/0/0/9/24",
	"testdata/0/0/9/25",
	"testdata/0/0/9/26",
	"testdata/0/0/9/27",
	"testdata/0/0/9/28",
	"testdata/0/0/9/29",
	"testdata/0/0/9/3",
	"testdata/0/0/9/30",
	"testdata/0/0/9/31",
	"testdata/0/0/9/32",
	"testdata/0/0/9/33",
	"testdata/0/0/9/34",
	"testdata/0/0/9/35",
	"testdata/0/0/9/36",
	"testdata/0/0/9/37",
	"testdata/0/0/9/38",
	"testdata/0/0/9/39",
	"testdata/0/0/9/4",
	"testdata/0/0/9/40",
	"testdata/0/0/9/41",
	"testdata/0/0/9/42",
	"testdata/0/0/9/43",
	"testdata/0/0/9/44",
	"testdata/0/0/9/45",
	"testdata/0/0/9/46",
	"testdata/0/0/9/47",
	"testdata/0/0/9/48",
	"testdata/0/0/9/49",
	"testdata/0/0/9/5",
	"testdata/0/0/9/50",
	"testdata/0/0/9/51",
	"testdata/0/0/9/52",
	"testdata/0/0/9/53",
	"testdata/0/0/9/54",
	"testdata/0/0/9/55",
	"testdata/0/0/9/56",
	"testdata/0/0/9/57",
	"testdata/0/0/9/58",
	"testdata/0/0/9/59",
	"testdata/0/0/9/6",
	"testdata/0/0/9/60",
	"testdata/0/0/9/61",
	"testdata/0/0/9/62",
	"testdata/0/0/9/63",
	"testdata/0/0/9/64",
	"testdata/0/0/9/65",
	"testdata/0/0/9/66",
	"testdata/0/0/9/67",
	"testdata/0/0/9/68",
	"testdata/0/0/9/7",
	"testdata/0/0/9/8",
	"testdata/0/0/9/9",
	"testdata/0/0/9",
	"testdata/0/0",
	"testdata/0",
	"testdata",
	"",
}

func TestDelayedWalkTree(t *testing.T) {
	WithTestEnvironment(t, repoFixture, func(repodir string) {
		repo := OpenLocalRepo(t, repodir)
		OK(t, repo.LoadIndex())

		root, err := restic.ParseID("937a2f64f736c64ee700c6ab06f840c68c94799c288146a0e81e07f4c94254da")
		OK(t, err)

		dr := delayRepo{repo, 100 * time.Millisecond}

		// start tree walker
		treeJobs := make(chan walk.TreeJob)
		go walk.Tree(dr, root, nil, treeJobs)

		i := 0
		for job := range treeJobs {
			expectedPath := filepath.Join(strings.Split(walktreeTestItems[i], "/")...)
			if job.Path != expectedPath {
				t.Fatalf("expected path %q (%v), got %q", walktreeTestItems[i], i, job.Path)
			}
			i++
		}

		if i != len(walktreeTestItems) {
			t.Fatalf("got %d items, expected %v", i, len(walktreeTestItems))
		}
	})
}

func BenchmarkDelayedWalkTree(t *testing.B) {
	WithTestEnvironment(t, repoFixture, func(repodir string) {
		repo := OpenLocalRepo(t, repodir)
		OK(t, repo.LoadIndex())

		root, err := restic.ParseID("937a2f64f736c64ee700c6ab06f840c68c94799c288146a0e81e07f4c94254da")
		OK(t, err)

		dr := delayRepo{repo, 10 * time.Millisecond}

		t.ResetTimer()

		for i := 0; i < t.N; i++ {
			// start tree walker
			treeJobs := make(chan walk.TreeJob)
			go walk.Tree(dr, root, nil, treeJobs)

			for _ = range treeJobs {
			}
		}
	})
}
